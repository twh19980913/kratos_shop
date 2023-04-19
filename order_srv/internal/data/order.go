package data

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	v1 "order_srv/api/helloworld/v1"
	"order_srv/internal/biz"
	"time"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	errors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/apache/rocketmq-client-go/v2/consumer"
)

type ShoppingCart struct{
	BaseModel
	User int32 `gorm:"type:int;index"` //在购物车列表中我们需要查询当前用户的购物车记录
	Goods int32 `gorm:"type:int;index"` //加索引：我们需要查询时候， 1. 会影响插入性能 2. 会占用磁盘
	Nums int32 `gorm:"type:int"`
	Checked bool //是否选中
}

func (ShoppingCart) TableName() string {
	return "shoppingcart"
}

type OrderInfo struct{
	BaseModel

	User int32 `gorm:"type:int;index"`
	OrderSn string `gorm:"type:varchar(30);index"` //订单号，我们平台自己生成的订单号
	PayType string `gorm:"type:varchar(20) comment 'alipay(支付宝)， wechat(微信)'"`

	//status大家可以考虑使用iota来做
	Status string `gorm:"type:varchar(20)  comment 'PAYING(待支付), TRADE_SUCCESS(成功)， TRADE_CLOSED(超时关闭), WAIT_BUYER_PAY(交易创建), TRADE_FINISHED(交易结束)'"`
	TradeNo string `gorm:"type:varchar(100) comment '交易号'"` //交易号就是支付宝的订单号 查账
	OrderMount float32
	PayTime *time.Time `gorm:"type:datetime"`

	Address string `gorm:"type:varchar(100)"`
	SignerName string `gorm:"type:varchar(20)"`
	SingerMobile string `gorm:"type:varchar(11)"`
	Post string `gorm:"type:varchar(20)"`
}

func (OrderInfo) TableName() string {
	return "orderinfo"
}

type OrderGoods struct{
	BaseModel

	Order int32 `gorm:"type:int;index"`
	Goods int32 `gorm:"type:int;index"`

	//把商品的信息保存下来了 ， 字段冗余， 高并发系统中我们一般都不会遵循三范式  做镜像 记录
	GoodsName string `gorm:"type:varchar(100);index"`
	GoodsImage string `gorm:"type:varchar(200)"`
	GoodsPrice float32
	Nums int32 `gorm:"type:int"`
}

func (OrderGoods) TableName() string {
	return "ordergoods"
}

type orderRepo struct {
	data *Data
	goodsSrvClient v1.GoodsClient
	inventorySrvClient v1.InventoryClient
	log  *log.Helper
	orderListener *OrderListener
}

// NewGreeterRepo .
func NewOrderRepo(data *Data, logger log.Logger,goodsSrvClient v1.GoodsClient,inventorySrvClient v1.InventoryClient,orderListener *OrderListener) biz.OrderRepo {
	return &orderRepo{
		data: data,
		goodsSrvClient: goodsSrvClient,
		inventorySrvClient: inventorySrvClient,
		orderListener:orderListener,
		log:  log.NewHelper(logger),
	}
}


func (o *orderRepo)CartItemList(ctx context.Context,req *biz.UserInfo) (*biz.CartItemListResponse,error){
	// 获取用户的购物车列表
	var shopCarts []ShoppingCart
	var rsp biz.CartItemListResponse
	if result := o.data.db.Where(&ShoppingCart{User: req.Id}).Find(&shopCarts);result.Error != nil{
		return nil,result.Error
	}else {
		rsp.Total = int32(result.RowsAffected)
	}

	for _,shopCart := range shopCarts{
		rsp.Data = append(rsp.Data,&biz.ShopCartInfoResponse{
			Id: shopCart.ID,
			UserId: shopCart.User,
			GoodsId: shopCart.Goods,
			Nums: shopCart.Nums,
			Checked: shopCart.Checked,
		})
	}
	return &rsp,nil
}

func (o *orderRepo)CreateCartItem(ctx context.Context,req *biz.CartItemRequest) (*biz.ShopCartInfoResponse,error){
	// 将商品添加到购物车 1、 购物车中原本没有这件商品 - 新建一个记录
	// 2、这个商品之前添加到了购物车 - 合并
	var shopCart ShoppingCart
	if result := o.data.db.Where(&ShoppingCart{Goods: req.GoodsId,User: req.UserId}).First(&shopCart);result.RowsAffected == 1{ // 查询到了
		// 如果记录已经存在 则合并购物车记录 更新操作
		shopCart.Nums += req.Nums
	}else {
		// 添加操作
		shopCart.User = req.UserId
		shopCart.Goods = req.GoodsId
		shopCart.Nums = req.Nums
		shopCart.Checked = false
	}

	o.data.db.Save(&shopCart)
	return &biz.ShopCartInfoResponse{Id: shopCart.ID},nil
}

func (o *orderRepo)UpdateCartItem(ctx context.Context,req *biz.CartItemRequest) error{
	// 更新购物车记录，更新数量和选中状态
	var shopCart ShoppingCart
	if result := o.data.db.Where("goods = ? and user = ?",req.GoodsId,req.UserId).First(&shopCart);result.RowsAffected == 0{ // 查询到了
		return errors.New(404,"NotFound","购物车记录不存在")
	}

	shopCart.Checked = req.Checked
	if req.Nums > 0 {
		shopCart.Nums = req.Nums
	}
	o.data.db.Save(&shopCart)
	return nil
}

func (o *orderRepo)DeleteCartItem(ctx context.Context,req *biz.CartItemRequest) error{
	if result := o.data.db.Where("goods = ? and user = ?",req.GoodsId,req.UserId).Delete(&ShoppingCart{});result.RowsAffected == 0{ // 查询到了
		return errors.New(404,"NotFound","购物车记录不存在")
	}
	return nil
}

func (o *orderRepo)OrderList(ctx context.Context,req *biz.OrderFilterRequest) (*biz.OrderListResponse,error){
	var orders []OrderInfo
	var rsp biz.OrderListResponse

	var total int64
	// 是后台管理系统查询 还是电商系统查询
 	result := o.data.db.Where(&OrderInfo{User: req.UserId}).Count(&total)
	rsp.Total = int32(result.RowsAffected)

	// 分页
	o.data.db.Scopes(Paginate(int(req.Pages),int(req.PagePerNums))).Where(&OrderInfo{User: req.UserId}).Find(&orders)
	for _,order := range orders{
		rsp.Data = append(rsp.Data, &biz.OrderInfoResponse{
			Id: order.ID,
			UserId: order.User,
			OrderSn: order.OrderSn,
			PayType: order.PayType,
			Status: order.Status,
			Post: order.Post,
			Total: order.OrderMount,
			Address: order.Address,
			Name: order.SignerName,
			Mobile: order.SingerMobile,
			AddTime: order.CreatedAt.Format("2006-01-02 13:04:05"),
		})
	}
	return &rsp,nil
}

func (o *orderRepo)OrderDetail(ctx context.Context,req *biz.OrderRequest) (*biz.OrderInfoDetailResponse,error){
	var order OrderInfo
	var rsp biz.OrderInfoDetailResponse

	//这个订单的id是否是当前用户的订单 如果在web层用户传递过来一个id的订单，web层应该先查询一下订单是否是当前进程
	//在个人中心可以这样做，但是如果是后台管理系统 web层如果是后台管理系统，那么只传递order的id过来 如果是电商系统除了
	if result := o.data.db.Where(&OrderInfo{BaseModel: BaseModel{ID: req.Id},User: req.UserId}).First(&order);result.RowsAffected == 0{
		// 订单没有取到
		return nil,errors.New(404,"NotFound","订单不存在")
	}

	orderInfo := biz.OrderInfoResponse{}

	orderInfo.Id = order.ID
	orderInfo.UserId = order.User
	orderInfo.OrderSn = order.OrderSn
	orderInfo.PayType = order.PayType
	orderInfo.Status = order.Status
	orderInfo.Post = order.Post
	orderInfo.Total = order.OrderMount
	orderInfo.Address = order.Address
	orderInfo.Name = order.SignerName
	orderInfo.Mobile = order.SingerMobile

	rsp.OrderInfo = &orderInfo

	var orderGoods []OrderGoods
	if result := o.data.db.Where(&OrderGoods{Order: order.ID}).Find(&orderGoods);result.Error != nil{
		return nil,result.Error
	}
	for _,orderGood := range orderGoods{
		rsp.Goods = append(rsp.Goods, &biz.OrderItemResponse{
			GoodsId: orderGood.Goods,
			GoodsName: orderGood.GoodsName,
			GoodsPrice: orderGood.GoodsPrice,
			Nums: orderGood.Nums,
		})
	}

	return &rsp,nil
}

func GenerateOrderSn(userId int32) string {
	//订单号的生成规则
	/*
	年月日时分秒 + 用户id + 2位随机数
	*/
	now := time.Now()
	rand.Seed(time.Now().UnixNano())
	orderSn := fmt.Sprintf("%d%d%d%d%d%d%d%d",now.Year(),now.Month(),now.Day(),now.Hour(),
		now.Minute(),now.Nanosecond(),userId,rand.Intn(90) + 10)
	return orderSn
}

type OrderListener struct {
	data *Data
	goodsSrvClient v1.GoodsClient
	inventorySrvClient v1.InventoryClient
	Code int32
	Detail string
	ID int32
	OrderAmount float32
}

func NewOrderListener(data *Data,goodsSrvClient v1.GoodsClient,inventorySrvClient v1.InventoryClient) (*OrderListener, error) {
	return &OrderListener{data: data,goodsSrvClient: goodsSrvClient,inventorySrvClient: inventorySrvClient}, nil
}

func (o *OrderListener) ExecuteLocalTransaction (msg *primitive.Message) primitive.LocalTransactionState {
	var orderInfo OrderInfo
	json.Unmarshal(msg.Body,&orderInfo)
	fmt.Println("我进来了",orderInfo)
	var goodsIds []int32
	var shopCarts []ShoppingCart
	goodsNumsMap := make(map[int32]int32)
	fmt.Println("商品服务地址",o.goodsSrvClient)
	if result := o.data.db.Where(&ShoppingCart{User: orderInfo.User,Checked: true}).Find(&shopCarts);result.RowsAffected == 0{
		o.Code = 22
		o.Detail = "没有选中结算的商品"
		return primitive.RollbackMessageState
	}

	fmt.Println("查询到购物车信息")

	for _,shopCart := range shopCarts{
		goodsIds = append(goodsIds, shopCart.Goods)
		goodsNumsMap[shopCart.Goods] = shopCart.Nums
	}
	
	// 跨商品服务调用
	goods,err := o.goodsSrvClient.BatchGetGoods(context.Background(),&v1.BatchGoodsIdInfo{
		Id: goodsIds,
	})
	if err != nil {
		o.Code = 500
		o.Detail = "批量查询商品信息失败"
		return primitive.RollbackMessageState
	}

	fmt.Println("获取到选中的商品信息")
	
	var orderAmount float32
	var orderGoods []*OrderGoods
	var goodsInvInfo []*v1.GoodsInvInfo
	for _,good := range goods.Data{
		orderAmount += good.ShopPrice * float32(goodsNumsMap[good.Id]) // 通过商品的id 查询商品的购买数量
		orderGoods = append(orderGoods, &OrderGoods{
			Goods: good.Id,
			GoodsName:good.Name,
			GoodsImage: good.GoodsFrontImage,
			GoodsPrice:good.ShopPrice,
			Nums: goodsNumsMap[good.Id],
		})

		goodsInvInfo = append(goodsInvInfo, &v1.GoodsInvInfo{
			GoodsId: good.Id,
			Num: goodsNumsMap[good.Id],
		})
	}

	// 库存微服务进行库存扣减
	_,err = o.inventorySrvClient.Sell(context.Background(),&v1.SellInfo{
		GoodsInfo: goodsInvInfo,
		OrderSn: orderInfo.OrderSn,
	})
	if err != nil{
		o.Code = 500
		o.Detail = "库存扣减失败"
		return primitive.RollbackMessageState
	}

	tx := o.data.db.Begin()
	// 生成订单表
	orderInfo.OrderMount = orderAmount
	
	if result := tx.Save(&orderInfo);result.RowsAffected == 0{
		tx.Rollback()
		o.Code = 500
		o.Detail = "创建订单失败"
		return primitive.CommitMessageState
	}

	o.OrderAmount = orderAmount
	o.ID = orderInfo.ID

	for _,orderGood := range orderGoods{
		orderGood.Order = orderInfo.ID
	}

	// 批量插入orderGoods
	if result := tx.CreateInBatches(orderGoods,100);result.RowsAffected == 0{
		tx.Rollback()
		o.Code = 500
		o.Detail = "创建订单失败"
		return primitive.RollbackMessageState
	}

	// 删除选中的购物车数据
	if result := tx.Where(&ShoppingCart{User: orderInfo.User,Checked: true}).Delete(&ShoppingCart{});result.RowsAffected == 0{
		tx.Rollback()
		o.Code = 500
		o.Detail = "删除购物车记录失败"
		return primitive.RollbackMessageState
	}
	// 提交事物
	//==========================发送关闭消息
	var opt []producer.Option
	opt = append(opt, producer.WithGroupName("timeout_producer"))
	opt = append(opt, producer.WithNameServer([]string{"192.168.16.128:9876"}))
	p,err := rocketmq.NewProducer(opt...)
	if err != nil {
		fmt.Println("生成producer失败")
	}

	if err = p.Start();err != nil{
		fmt.Println("启动producer失败")
	}

	msg = primitive.NewMessage("order_timeout",msg.Body)
	msg.WithDelayTimeLevel(3)

	_,err = p.SendSync(context.Background(),msg)
	if err != nil {
		tx.Rollback()
		o.Code = 500
		o.Detail = "发送延迟消息成功"
		return primitive.CommitMessageState
	}

	//===============================发送关闭完成
	tx.Commit()
	o.Code = 200
	//本地执行逻辑无缘无故失败 代码异常 宕机
	return primitive.RollbackMessageState
}

func (o *OrderListener) CheckLocalTransaction (msg *primitive.MessageExt) primitive.LocalTransactionState {
	var orderInfo OrderInfo
	json.Unmarshal(msg.Body,&orderInfo)
	// 怎么检查之前的逻辑是否完成
	if result := o.data.db.Where(OrderInfo{OrderSn: orderInfo.OrderSn}).First(&orderInfo);result.RowsAffected == 0{
		return primitive.CommitMessageState // 并不能说明这里库存已经扣减了
	}
	//本地事务执行失败了 消息不要投递出去
	return primitive.RollbackMessageState
}

func (o *orderRepo)CreateOrder(ctx context.Context,req *biz.OrderRequest) (*biz.OrderInfoResponse,error){
	// 新建订单
	/*
		1、从购物车中获取到选中的商品
		2、商品的价格自己查询 - 访问商品服务(跨微服务调用)
		3、库存的扣减 - 访问库存服务去扣减库存(跨微服务)
		4、订单的基本信息表 - 订单的商品信息表
		5、从购物车中删除已购买的记录
	*/
	// ROCKET MQ发送消息
	orderListener := &OrderListener{data: o.data,goodsSrvClient: o.goodsSrvClient,inventorySrvClient: o.inventorySrvClient}
	
	p,err := rocketmq.NewTransactionProducer(
		orderListener,
		producer.WithNameServer([]string{"192.168.16.128:9876"}),
	)
	if err != nil {
		return nil,err
	}

	if err = p.Start();err != nil{
		return nil,err
	}
	order := OrderInfo{
		OrderSn: GenerateOrderSn(req.UserId),
		Address: req.Address,
		SignerName: req.Name,
		SingerMobile: req.Mobile,
		Post: req.Post,
		User: req.UserId,
	}

	jsonString,_ := json.Marshal(order)

	_,err = p.SendMessageInTransaction(context.Background(),primitive.NewMessage(
		"order_reback",jsonString))
	if err != nil {
		fmt.Printf("发送失败...%s\n",err)
		return nil, errors.New(500,"Internal","发送消息失败")
	}

	// 发送结束
	
	
	return &biz.OrderInfoResponse{Id: orderListener.ID,OrderSn: order.OrderSn,Total: orderListener.OrderAmount},nil
}

func (o *orderRepo)UpdateOrderStatus(ctx context.Context,req *biz.OrderStatus) error{
	// 先查询，再更新，实际上有两条sql语句执行，select 和 update语句
	if result := o.data.db.Model(&OrderInfo{}).Where("order_sn = ?",req.OrderSn).Update("status",req.Status);result.RowsAffected == 0{
		return errors.New(404,"NotFound","订单不存在")
	}
	return nil
}

func (o *orderRepo)OrderTimeout (ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error){
	for i := range msgs{
		var orderInfo OrderInfo
		_ = json.Unmarshal(msgs[i].Body,&orderInfo)
		
		fmt.Printf("获取到订单超时消息:%v\n",time.Now())
		// 查询订单支付状态，如果已经支付就什么都不做，如果未支付，归还库存
		var order OrderInfo
		if result := o.data.db.Model(&OrderInfo{}).Where(&OrderInfo{OrderSn: order.OrderSn}).First(&order);result.RowsAffected == 0{
			return consumer.ConsumeSuccess,nil
		}

		if order.Status != "TRADE_SUCCESS"{
			tx := o.data.db.Begin()
			// 修改订单的状态
			order.Status = "TRADE_CLOSED"
			tx.Save(&order)
			// 归还库存 我们可以模仿order中发送一个消息到order_reback中去
			p,err := rocketmq.NewProducer(producer.WithNameServer([]string{"192.168.16.128:9876"}))
			if err != nil {
				fmt.Println("生成producer失败")
			}

			if err = p.Start();err != nil{
				fmt.Println("启动producer失败")
			}

			_,err = p.SendSync(context.Background(),primitive.NewMessage("order_reback",msgs[i].Body))
			if err != nil {
				tx.Rollback()
				fmt.Printf("发送失败...%s\n",err)
				return consumer.ConsumeRetryLater,nil
			}
			
		}
	}
	return consumer.ConsumeSuccess,nil
}