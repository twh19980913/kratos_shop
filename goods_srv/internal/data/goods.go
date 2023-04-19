package data

import (
	"context"
	"encoding/json"
	"fmt"
	"goods_srv/internal/biz"
	"strconv"

	errors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/olivere/elastic/v7"
)

// 类型，这个子段是否能为null，这个子段应该设置为可以为null还是设置为空
// 实际开发过程中，尽量设置为不为null
type Goods struct {
	BaseModel

	CategoryID int32 `gorm:"type:int;not null"`
	Category   Category
	BrandsID   int32 `gorm:"type:int;not null"`
	Brands     Brands

	OnSale   bool `gorm:"default:false;not null"`
	ShipFree bool `gorm:"default:false;not null"`
	IsNew    bool `gorm:"default:false;not null"`
	IsHot    bool `gorm:"default:false;not null"`

	Name            string   `gorm:"type:varchar(50);not null"`
	GoodsSn         string   `gorm:"type:varchar(50);not null"`
	ClickNum        int32    `gorm:"type:int;default:0;not null"`
	SoldNum         int32    `gorm:"type:int;default:0;not null"`
	FavNum          int32    `gorm:"type:int;default:0;not null"`
	MarketPrice     float32  `gorm:"not null"`
	ShopPrice       float32  `gorm:"not null"`
	GoodsBrief      string   `gorm:"type:varchar(100);not null"`
	Images          GormList `gorm:"type:varchar(1000);not null"`
	DescImages      GormList `gorm:"type:varchar(1000);not null"`
	GoodsFrontImage string   `gorm:"type:varchar(200);not null"`
}


type goodsRepo struct {
	data     *Data
	esClient *ESClient
	log      *log.Helper
}

// NewGreeterRepo .
func NewGoodsRepo(data *Data, logger log.Logger, esClient *ESClient) biz.GoodsRepo {
	return &goodsRepo{
		data:     data,
		log:      log.NewHelper(logger),
		esClient: esClient,
	}
}

func ModelToResponse(goods Goods) biz.GoodsInfoResponse {
	goodsInfoResponse := biz.GoodsInfoResponse{
		Id:              goods.ID,
		CategoryId:      goods.CategoryID,
		Name:            goods.Name,
		GoodsSn:         goods.GoodsSn,
		ClickNum:        goods.ClickNum,
		SoldNum:         goods.SoldNum,
		FavNum:          goods.FavNum,
		MarketPrice:     goods.MarketPrice,
		ShopPrice:       goods.ShopPrice,
		GoodsBrief:      goods.GoodsBrief,
		ShipFree:        goods.ShipFree,
		GoodsFrontImage: goods.GoodsFrontImage,
		IsNew:           goods.IsNew,
		IsHot:           goods.IsHot,
		OnSale:          goods.OnSale,
		DescImages:      goods.DescImages,
		Images:          goods.Images,
		Category: &biz.CategoryBriefInfoResponse{
			Id:   goods.Category.ID,
			Name: goods.Category.Name,
		},
		Brand: &biz.BrandInfoResponse{
			Id:   goods.Brands.ID,
			Name: goods.Brands.Name,
			Logo: goods.Brands.Logo,
		},
	}
	return goodsInfoResponse
}

func (s *goodsRepo) GoodsList(ctx context.Context, req *biz.GoodsFilterRequest) (*biz.GoodsListResponse, error) {
	goodsListResponse := &biz.GoodsListResponse{}

	// match bool 复合查询
	q := elastic.NewBoolQuery()

	// var goods []Goods
	localDB := s.data.db.Model(Goods{})
	if req.KeyWords != "" {
		// 搜索
		// localDB = localDB.Where("name LIKE ?","%" + req.KeyWords + "%")
		q = q.Must(elastic.NewMultiMatchQuery(req.KeyWords, "name", "goods_brief"))
	}
	if req.IsHot {
		// localDB = localDB.Where(Goods{IsHot: true})
		q = q.Filter(elastic.NewTermQuery("is_hot", req.IsHot))
	}
	if req.IsNew {
		// localDB = localDB.Where(Goods{IsNew: true})
		q = q.Filter(elastic.NewTermQuery("is_new", req.IsNew))
	}

	if req.PriceMin > 0 {
		// localDB = localDB.Where("shop_price >= ?",req.PriceMin)
		q = q.Filter(elastic.NewRangeQuery("shop_price").Gte(req.PriceMin))
	}
	if req.PriceMax > 0 {
		// localDB = localDB.Where("shop_price <= ?",req.PriceMax)
		q = q.Filter(elastic.NewRangeQuery("shop_price").Lte(req.PriceMax))
	}
	if req.Brand > 0 {
		// localDB = localDB.Where("brand_id = ?",req.Brand)
		q = q.Filter(elastic.NewTermQuery("brands_id", req.Brand))
	}

	// 通过category 查询商品
	var subQuery string
	categoryIds := make([]int32, 0)
	if req.TopCategory > 0 {
		var category Category
		if result := s.data.db.First(&category, req.TopCategory); result.RowsAffected == 0 {
			return nil, errors.New(404, "NotFound", "商品分类不存在")
		}

		if category.Level == 1 {
			subQuery = fmt.Sprintf("select id from category where parent_category_id in (select id from category WHERE parent_category_id=%d)", req.TopCategory)
		} else if category.Level == 2 {
			subQuery = fmt.Sprintf("select id from category WHERE parent_category_id=%d", req.TopCategory)
		} else if category.Level == 3 {
			subQuery = fmt.Sprintf("select id from category WHERE id=%d", req.TopCategory)
		}

		type Result struct {
			ID int32 `json:"id"`
		}
		var results []Result
		s.data.db.Model(Category{}).Raw(subQuery).Scan(&results)
		for _, re := range results {
			categoryIds = append(categoryIds, re.ID)
		}
		q = q.Filter(elastic.NewTermsQuery("category_id",categoryIds))
		// localDB = localDB.Where(fmt.Sprintf("category_id in (%s)",subQuery))
	}

	if req.Pages == 0{
		req.Pages  = 1
	}
	switch{
	case req.PagePerNums > 100:
		req.PagePerNums = 100
	case req.PagePerNums <= 0:
		req.PagePerNums = 10
	}

	result,err := s.esClient.esClient.Search().Index(EsGoods{}.GetIndexName()).Query(q).From(int(req.Pages)).Size(int(req.PagePerNums)).Do(context.Background())
	if err != nil {
		return nil,err
	}

	goodsIds := make([]int32,0)
	goodsListResponse.Total = int32(result.Hits.TotalHits.Value)
	for _,value := range result.Hits.Hits{
		goods := EsGoods{}
		_ = json.Unmarshal(value.Source,&goods)
		goodsIds = append(goodsIds, goods.ID)
	}
	// 拿到total
	// var count int64
	// localDB.Count(&count)
	// goodsListResponse.Total = int32(count)
	// 使用es的目的是搜索出来商品的id，通过id拿到具体的字段信息是通过mysql来完成的
	// 我们使用es是用来做搜索的，是否应该将所有的mysql字段全部在es中保存一份
	// es用来做搜索，这个时候我们一般只把搜索和过滤的字段信息保存到es中
	// es可以用来当作mysql使用，但实际上mysql和es之间是互补的关系
	// 一般mysql用来做存储使用 es做搜索使用
	// es想要提高性能，就要将es的内存设置的够大

	// 查询id在某个数组中的
	var goods []Goods
	re := localDB.Preload("Category").Preload("Brands").Find(&goods,goodsIds)
	if result.Error != nil {
		return nil, re.Error
	}

	for _, good := range goods {
		goodsInfoResponse := ModelToResponse(good)
		goodsListResponse.Data = append(goodsListResponse.Data, &goodsInfoResponse)
	}

	return goodsListResponse, nil
}

// 现在用户提交订单有多个商品，你得批量查询商品的信息把
func (s *goodsRepo) BatchGetGoods(ctx context.Context, req *biz.BatchGoodsIdInfo) (*biz.GoodsListResponse, error) {
	goodsListResponse := &biz.GoodsListResponse{}
	var goods []Goods

	result := s.data.db.Where(req.Id).Find(&goods)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	for _, good := range goods {
		fmt.Println(good.Name)
		goodsInfoResponse := ModelToResponse(good)
		goodsListResponse.Data = append(goodsListResponse.Data, &goodsInfoResponse)
	}
	goodsListResponse.Total = int32(result.RowsAffected)
	return goodsListResponse, nil
}

func (s *goodsRepo) GetGoodsDetail(ctx context.Context, req *biz.GoodInfoRequest) (*biz.GoodsInfoResponse, error) {
	var goods Goods

	if result := s.data.db.Preload("Category").Preload("Brands").First(&goods, req.Id); result.RowsAffected == 0 {
		return nil, errors.New(404, "NotFound", "商品不存在")
	}
	goodsInfoResponse := ModelToResponse(goods)
	return &goodsInfoResponse, nil
}

func (s *goodsRepo) CreateGoods(ctx context.Context, req *biz.CreateGoodsInfo) (*biz.GoodsInfoResponse, error) {
	var category Category
	if result := s.data.db.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, errors.New(404, "NotFound", "商品不存在")
	}

	var brand Brands
	if result := s.data.db.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, errors.New(404, "NotFound", "品牌不存在")
	}

	goods := Goods{
		Brands:          brand,
		BrandsID:        brand.ID,
		Category:        category,
		CategoryID:      category.ID,
		Name:            req.Name,
		GoodsSn:         req.GoodsSn,
		MarketPrice:     req.MarketPrice,
		ShopPrice:       req.ShopPrice,
		GoodsBrief:      req.GoodsBrief,
		ShipFree:        req.ShipFree,
		Images:          req.Images,
		DescImages:      req.DescImages,
		GoodsFrontImage: req.GoodsFrontImage,
		IsNew:           req.IsNew,
		IsHot:           req.IsHot,
		OnSale:          req.OnSale,
	}
	tx := s.data.db.Begin()
	result := tx.Save(&goods)
	if result.Error != nil {
		return nil, result.Error
	}

	esModel := EsGoods{
		ID: goods.ID,
		CategoryID: goods.CategoryID,
		BrandsID: goods.BrandsID,
		OnSale: goods.OnSale,
		ShipFree: goods.ShipFree,
		IsNew: goods.IsNew,
		IsHot: goods.IsHot,
		Name: goods.Name,
		ClickNum: goods.ClickNum,
		SoldNum: goods.SoldNum,
		FavNum: goods.FavNum,
		MarketPrice: goods.MarketPrice,
		GoodsBrief: goods.GoodsBrief,
		ShopPrice: goods.ShopPrice,
	}

	_, err := s.esClient.esClient.Index().Index(esModel.GetIndexName()).BodyJson(esModel).Id(strconv.Itoa(int(goods.ID))).Do(context.Background())
	if err != nil {
		tx.Rollback()
		return nil,err
	}
	tx.Commit()
	return &biz.GoodsInfoResponse{Id: goods.ID}, nil
}

func (s *goodsRepo) DeleteGoods(ctx context.Context, req *biz.DeleteGoodsInfo) error {
	if result := s.data.db.Delete(&Goods{BaseModel: BaseModel{ID: req.Id}}); result.Error != nil {
		return errors.New(404, "NotFound", "品牌不存在")
	}
	return nil
}

func (s *goodsRepo) UpdateGoods(ctx context.Context, req *biz.CreateGoodsInfo) error {
	var goods Goods
	if result := s.data.db.First(&goods, req.Id); result.RowsAffected == 0 {
		return errors.New(404, "NotFound", "商品不存在")
	}

	var category Category
	if result := s.data.db.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return errors.New(404, "NotFound", "商品不存在")
	}

	var brand Brands
	if result := s.data.db.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return errors.New(404, "NotFound", "品牌不存在")
	}

	goods.Brands = brand
	goods.BrandsID = brand.ID
	goods.Category = category
	goods.CategoryID = category.ID
	goods.Name = req.Name
	goods.GoodsSn = req.GoodsSn
	goods.MarketPrice = req.MarketPrice
	goods.ShopPrice = req.ShopPrice
	goods.GoodsBrief = req.GoodsBrief
	goods.ShipFree = req.ShipFree
	goods.Images = req.Images
	goods.DescImages = req.DescImages
	goods.GoodsFrontImage = req.GoodsFrontImage
	goods.IsNew = req.IsNew
	goods.IsHot = req.IsHot
	goods.OnSale = req.OnSale

	result := s.data.db.Save(&goods)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
