package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type GoodsFilterRequest  struct{
	PriceMin int32
	PriceMax int32
	IsHot bool
	IsNew bool
	IsTab bool
	TopCategory int32
	Pages int32
	PagePerNums int32
	KeyWords string
	Brand int32
}

type GoodsInfoResponse struct{
	Id int32
	CategoryId int32
	Name string
	GoodsSn string
	ClickNum int32
	SoldNum int32
	FavNum int32
	MarketPrice float32
	ShopPrice float32
	GoodsBrief string
	GoodsDesc string
	ShipFree bool
	Images []string
	DescImages []string
	GoodsFrontImage string
	IsNew bool
	IsHot bool
	OnSale bool
	AddTime int64
	Category *CategoryBriefInfoResponse
	Brand *BrandInfoResponse
}

type CategoryBriefInfoResponse struct{
	Id int32
	Name string
}

type GoodsListResponse struct{
	Total int32
	Data []*GoodsInfoResponse
}

type BatchGoodsIdInfo struct{
	Id []int32
}

type DeleteGoodsInfo struct{
	Id int32
}

type GoodInfoRequest struct{
	Id int32
}

type CreateGoodsInfo struct{
	Id int32
	Name string
	GoodsSn string
	Stocks int32
	MarketPrice float32
	ShopPrice float32
	GoodsBrief string
	GoodsDesc string
	ShipFree bool
	Images []string
	DescImages []string
	GoodsFrontImage string
	IsNew bool
	IsHot bool
	OnSale bool
	CategoryId int32
	BrandId int32
}

type GoodsRepo interface {
	GoodsList(ctx context.Context,req *GoodsFilterRequest) (*GoodsListResponse,error)
	BatchGetGoods(ctx context.Context,req *BatchGoodsIdInfo)(*GoodsListResponse,error)
	GetGoodsDetail(ctx context.Context,req *GoodInfoRequest) (*GoodsInfoResponse,error)
	CreateGoods(ctx context.Context,req *CreateGoodsInfo)(*GoodsInfoResponse,error)
	DeleteGoods(ctx context.Context,req *DeleteGoodsInfo) error
	UpdateGoods(ctx context.Context,req *CreateGoodsInfo) error
}

// GreeterUsecase is a Greeter usecase.
type GoodsUsecase struct {
	repo GoodsRepo
	log  *log.Helper
}

// NewGreeterUsecase new a Greeter usecase.
func NewGoodsUsecase(repo GoodsRepo, logger log.Logger) *GoodsUsecase {
	return &GoodsUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (g *GoodsUsecase)GoodsList(ctx context.Context,req *GoodsFilterRequest) (*GoodsListResponse,error){
	return g.repo.GoodsList(ctx,req)
}

func (g *GoodsUsecase)BatchGetGoods(ctx context.Context,req *BatchGoodsIdInfo)(*GoodsListResponse,error){
	return g.repo.BatchGetGoods(ctx,req)
}
func (g *GoodsUsecase)GetGoodsDetail(ctx context.Context,req *GoodInfoRequest) (*GoodsInfoResponse,error){
	return g.repo.GetGoodsDetail(ctx,req)
}

func (g *GoodsUsecase)CreateGoods(ctx context.Context,req *CreateGoodsInfo)(*GoodsInfoResponse,error){
	return g.repo.CreateGoods(ctx,req)
}
func (g *GoodsUsecase)DeleteGoods(ctx context.Context,req *DeleteGoodsInfo) error{
	return g.repo.DeleteGoods(ctx,req)
}
func (g *GoodsUsecase)UpdateGoods(ctx context.Context,req *CreateGoodsInfo) error{
	return g.repo.UpdateGoods(ctx,req)
}