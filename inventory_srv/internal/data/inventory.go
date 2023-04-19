package data

import (
	"context"
	"fmt"
	"inventory_srv/internal/biz"
	"sync"

	errors "github.com/go-kratos/kratos/v2/errors"
	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	// "github.com/apache/rocketmq-client-go/v2"
	"database/sql/driver"
	"encoding/json"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/go-kratos/kratos/v2/log"
)

type GoodsDetail struct {
	Goods int32
	Num int32
}
type GoodsDetailList []GoodsDetail

func (g GoodsDetailList) Value() (driver.Value, error){
	return json.Marshal(g)
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (g *GoodsDetailList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &g)
}

type Inventory struct{
	BaseModel
	Goods int32 `gorm:"type:int;index"`
	Stocks int32 `gorm:"type:int"`
	Version int32 `gorm:"type:int"` // 分布式锁乐观锁
}

type Delivery struct {
	Goods int32 `gorm:"type:int;index"`
	Nums int32 `gorm:"type:int"`
	OrderSn string `gorm:"type:varchar(200)"`
	Status string `gorm:"type:varchar(200)"` //1. 表示已扣减 2. 已归还 3. 失败
}

type StockSellDetail struct {
	OrderSn string `gorm:"type:varchar(200);index:idx_order_sn,unique;"`
	Status int32 `gorm:"type:varchar(200)"` //1 表示已扣减 2. 表示已归还
	Detail GoodsDetailList `gorm:"type:varchar(200)"`
}

func (StockSellDetail) TableName() string {
	return "stockselldetail"
}

type inventoryRepo struct {
	data *Data
	log  *log.Helper
}

// NewGreeterRepo .
func NewInventoryRepo(data *Data, logger log.Logger) biz.InventoryRepo {
	return &inventoryRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (i *inventoryRepo)SetInv(ctx context.Context,req *biz.GoodsInvInfo) error{
	// 设置库存 如果我要更新库存
	var inv Inventory
	i.data.db.Where(&Inventory{Goods: req.GoodsId}).First(&inv)
	inv.Goods = req.GoodsId
	inv.Stocks = req.Num

	i.data.db.Save(&inv)
	return nil
}

// TODO 基于redis的分布式锁实现
func (i *inventoryRepo) InvDetail(ctx context.Context,req *biz.GoodsInvInfo) (*biz.GoodsInvInfo,error){
	var inv Inventory
	if result := i.data.db.Where(&Inventory{Goods: req.GoodsId}).First(&inv);result.RowsAffected == 0{
		return nil,errors.New(404,"NotFound","没有库存信息")
	}
	return &biz.GoodsInvInfo{
		GoodsId: inv.Goods,
		Num: inv.Stocks,
	},nil
}
//TODO 基于redis分布式锁的实现方案
/**
redsync源码解读
setnx的作用
	将获取和设置值变成原子的操作
如果我的服务挂掉了 - 死锁
	a 设置了过期时间
	b 如果你设置了过期时间，那么如果过期时间到了我的业务还没有执行完怎么办？
		在过期之前刷新一下
		需要自己去启动协程来完成延时的工作
			延时的接口可能会带来负面影响 如果其中某一个服务一直没有执行完 2s就能执行完 但是一直申请延长锁
			导致别人永远获取不到锁，这个很要命
		421 2BjwsGxn7qNhykooMQSZzQ== 此时redis中的状态
分布式锁需要解决的问题 -- lua脚本去做
	互斥性 setnx
	死锁 过期时间
	安全性
		锁只能被持有该锁的用户删除，不能被其他用户删除
			TODO 当时设置的value值是多少只有当时的goroutine才能直到
			在删除的时候取出redis中的值和当前自己保存下来的值做对比
即使你这样实现了分布式锁但是还是会有问题的 redlock 红锁
五台redis实例
每一个库存实例setnx操作应该在多台服务器上进行
五台相同级别的服务器
应该拿到 5 / 2 + 1多数 谁先拿到三台谁成功
driftFactor 因为有始终漂移
m.expiry - now.Sub(start) - time.Duration(int64(float64(m.expiry)*m.driftFactor
**/
func (i *inventoryRepo) Sell(ctx context.Context,req *biz.SellInfo) error{
	// 扣减库存 本地事物[1:10,2:5,3:20]
	// 数据库基本应用场景，数据一致性
	// 并发情况之下 可能会出现超卖
	client := goredislib.NewClient(&goredislib.Options{
		Addr: "127.0.0.1:6379",
	})
	//新建redis连接池
	pool := goredis.NewPool(client)
	//
	rs := redsync.New(pool)
	tx := i.data.db.Begin()
	
	sellDetail := StockSellDetail{
		OrderSn: req.OrderSn,
		Status: 1,
	}

	var details []GoodsDetail
	for _,goodInfo := range req.GoodsInfo{
		details = append(details, GoodsDetail{
			Goods: goodInfo.GoodsId,
			Num: goodInfo.Num,
		})
		//m.Lock() // 获取锁 这把锁有问题么 假设有10w的并发 这里并不是请求的同一件商品
		var inv Inventory
		
		mutex := rs.NewMutex(fmt.Sprintf("goods_%d",goodInfo.GoodsId))
		if err := mutex.Lock(); err != nil {
			return errors.New(500,"Internal","获取redis分布式锁异常")
		}

		if result := i.data.db.Where(&Inventory{Goods: goodInfo.GoodsId}).First(&inv);result.RowsAffected == 0{
			tx.Rollback()
			return errors.New(22,"InvalidArgument","没有库存信息")
		}

		if inv.Stocks < goodInfo.Num{
			tx.Rollback()
			return errors.New(22,"InvalidArgument","库存不足")
		}
		// 扣减 会出现数据不一致的问题 - 锁 分布式锁
		inv.Stocks -= goodInfo.Num
		tx.Save(&inv)

		//m.Unlock() // 释放锁
		if ok, err := mutex.Unlock(); !ok || err != nil {
			return errors.New(500,"Internal","释放redis分布式锁异常")
		}
	}
	sellDetail.Detail = details
	// 写selldetail表
	if result := tx.Create(&sellDetail);result.RowsAffected == 0{
		tx.Rollback()
		return errors.New(500,"Internal","保存库存扣减历史失败")
	}
	tx.Commit() // 需要自己手动提交操作
	return nil
}


// TODO 基于mysql乐观锁实现
func (i *inventoryRepo) Sell4(ctx context.Context,req *biz.SellInfo) error{
	// 扣减库存 本地事物[1:10,2:5,3:20]
	// 数据库基本应用场景，数据一致性
	// 并发情况之下 可能会出现超卖
	tx := i.data.db.Begin()
	
	for _,goodInfo := range req.GoodsInfo{
		//m.Lock() // 获取锁 这把锁有问题么 假设有10w的并发 这里并不是请求的同一件商品
		var inv Inventory
		for {
			if result := i.data.db.Where(&Inventory{Goods: goodInfo.GoodsId}).First(&inv);result.RowsAffected == 0{
				tx.Rollback()
				return errors.New(22,"InvalidArgument","没有库存信息")
			}

			if inv.Stocks < goodInfo.Num{
				tx.Rollback()
				return errors.New(22,"InvalidArgument","库存不足")
			}
			// 扣减 会出现数据不一致的问题 - 锁 分布式锁
			inv.Stocks -= goodInfo.Num

			// update inventory set stocks = stocks - 1,version=version + 1 where 
			// goods=goods and version=version
			// 这种写法有瑕疵 
			// 零值会被gorm给忽略掉
			if result := tx.Model(&Inventory{}).Select("Stocks","Version").Where("goods = ? and version = ?",goodInfo.GoodsId,inv.Version).Updates(Inventory{Stocks: inv.Stocks,Version: inv.Version + 1});result.RowsAffected == 0{
				// 更新不成功 冲突了
				zap.S().Info("库存扣减失败")
			}else {
				break
			}
		}
		// tx.Save(&inv)

		//m.Unlock() // 释放锁
	}
	tx.Commit() // 需要自己手动提交操作
	return nil
}

// TODO 基于mysql的for update实现悲观锁
// 向mysql请求一把锁 for update 。使用for update的时候中注意，默认每个语句mysql都是默认
// 提交 需要关闭autocommit 如果在没有建立索引的子段上for update 则行锁升级为表锁
// 如果没有满足条件的结果 不会锁表
func (i *inventoryRepo) Sell3(ctx context.Context,req *biz.SellInfo) error{
	// 扣减库存 本地事物[1:10,2:5,3:20]
	// 数据库基本应用场景，数据一致性
	// 并发情况之下 可能会出现超卖
	tx := i.data.db.Begin()
	
	for _,goodInfo := range req.GoodsInfo{
		//m.Lock() // 获取锁 这把锁有问题么 假设有10w的并发 这里并不是请求的同一件商品
		var inv Inventory
		if result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(&Inventory{Goods: goodInfo.GoodsId}).First(&inv);result.RowsAffected == 0{
			tx.Rollback()
			return errors.New(22,"InvalidArgument","没有库存信息")
		}

		if inv.Stocks < goodInfo.Num{
			tx.Rollback()
			return errors.New(22,"InvalidArgument","库存不足")
		}
		// 扣减 会出现数据不一致的问题 - 锁 分布式锁
		inv.Stocks -= goodInfo.Num
		tx.Save(&inv)

		//m.Unlock() // 释放锁
	}
	tx.Commit() // 需要自己手动提交操作
	return nil
}

var m sync.Mutex
// TODO 基于互斥锁的本地实现
func (i *inventoryRepo) Sell2(ctx context.Context,req *biz.SellInfo) error{
	// 扣减库存 本地事物[1:10,2:5,3:20]
	// 数据库基本应用场景，数据一致性
	// 并发情况之下 可能会出现超卖
	tx := i.data.db.Begin()
	
	for _,goodInfo := range req.GoodsInfo{
		m.Lock() // 获取锁 这把锁有问题么 假设有10w的并发 这里并不是请求的同一件商品
		var inv Inventory
		if result := i.data.db.Where(&Inventory{Goods: goodInfo.GoodsId}).First(&inv);result.RowsAffected == 0{
			tx.Rollback()
			return errors.New(22,"InvalidArgument","没有库存信息")
		}

		if inv.Stocks < goodInfo.Num{
			tx.Rollback()
			return errors.New(22,"InvalidArgument","库存不足")
		}
		// 扣减 会出现数据不一致的问题 - 锁 分布式锁
		inv.Stocks -= goodInfo.Num
		tx.Save(&inv)

		m.Unlock() // 释放锁
	}
	tx.Commit() // 需要自己手动提交操作
	return nil
}

// TODO 原始Sell
func (i *inventoryRepo) Sell1(ctx context.Context,req *biz.SellInfo) error{
	// 扣减库存 本地事物[1:10,2:5,3:20]
	// 数据库基本应用场景，数据一致性
	// 并发情况之下 可能会出现超卖
	tx := i.data.db.Begin()
	for _,goodInfo := range req.GoodsInfo{
		var inv Inventory
		if result := i.data.db.Where(&Inventory{Goods: goodInfo.GoodsId}).First(&inv);result.RowsAffected == 0{
			tx.Rollback()
			return errors.New(22,"InvalidArgument","没有库存信息")
		}

		if inv.Stocks < goodInfo.Num{
			tx.Rollback()
			return errors.New(22,"InvalidArgument","库存不足")
		}
		// 扣减 会出现数据不一致的问题 - 锁 分布式锁
		inv.Stocks -= goodInfo.Num
		tx.Save(&inv)
	}
	tx.Commit() // 需要自己手动提交操作
	return nil
}

func (i *inventoryRepo) Reback(ctx context.Context,req *biz.SellInfo) error{
	// 库存归还 1、订单超时归还 2、订单创建失败 归还之前扣减的库存 3、手动归还
	tx := i.data.db.Begin()
	for _,goodInfo := range req.GoodsInfo{
		var inv Inventory
		if result := i.data.db.Where(&Inventory{Goods: goodInfo.GoodsId}).First(&inv);result.RowsAffected == 0{
			tx.Rollback()
			return errors.New(22,"InvalidArgument","没有库存信息")
		}

		// 归还 会出现数据不一致的问题 - 锁 分布式锁
		inv.Stocks += goodInfo.Num
		tx.Save(&inv)
	}
	tx.Commit() // 需要自己手动提交操作
	return nil
}

func (i *inventoryRepo)AutoReback (ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error){

	
	fmt.Println("我进来了")
	type OrderInfo struct{
		OrderSn string
	}
	//既然是归还库存，那么我应该具体的知道每件商品归还多少 但是有一个问题？ 重复归还
	//所以说这个接口应该确保幂等性，你不能因为消息的重复发送导致一个订单的库存归还多次
	//没有扣减的库存你别归还
	//如何确保这些都没有问题，新建一张表，这张表详细的记录了订单扣减细节，以及归还细节
	
	for m := range msgs{
		var orderInfo OrderInfo
		err := json.Unmarshal(msgs[m].Body,&orderInfo)
		if err != nil{
			zap.S().Errorf("解析json失败:%v\n",msgs[m].Body)
			return consumer.ConsumeSuccess,nil
		}

		// 去将inv的库存加回去，将sellDetail的status设置为2 要在事物中进行
		tx := i.data.db.Begin()
		var sellDetail StockSellDetail
		if result := tx.Model(&StockSellDetail{}).Where(StockSellDetail{OrderSn: orderInfo.OrderSn,Status: 1}).First(&sellDetail);result.RowsAffected == 0{
			return consumer.ConsumeSuccess,nil
		}
		// 如果查询到那么逐个归还
		for _,orderGood := range sellDetail.Detail{
			if result := tx.Model(&Inventory{}).Where(&Inventory{Goods: orderGood.Goods}).Update("stocks",gorm.Expr("stocks+?",orderGood.Num));result.RowsAffected == 0{
				tx.Rollback()
				return consumer.ConsumeRetryLater,nil
			}
		}
		sellDetail.Status = 2
		if result := tx.Model(&StockSellDetail{}).Where(&StockSellDetail{OrderSn: orderInfo.OrderSn}).Update("status",2);result.RowsAffected == 0{
			tx.Rollback()
			return consumer.ConsumeRetryLater,nil
		}
		tx.Commit()
		return consumer.ConsumeSuccess,nil
	}
	return consumer.ConsumeSuccess,nil
}
