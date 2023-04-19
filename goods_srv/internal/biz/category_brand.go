package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type CategoryBrandFilterRequest  struct{
	Pages int32
    PagePerNums int32
  }

  type CategoryBrandListResponse struct{
	Total int32
	Data []*CategoryBrandResponse
  }

  type CategoryBrandResponse struct{
	Id int32
	Brand *BrandInfoResponse
	Category *CategoryInfoResponse
  }


  type CategoryBrandRequest struct{
	Id int32
	CategoryId int32
	BrandId int32
  }

type GoodsCategoryBrandRepo interface {
	CategoryBrandList(context.Context,*CategoryBrandFilterRequest)(*CategoryBrandListResponse,error)
	GetCategoryBrandList(context.Context,*CategoryInfoRequest) (*BrandListResponse,error)
	CreateCategoryBrand(context.Context,*CategoryBrandRequest) (*CategoryBrandResponse,error)
	DeleteCategoryBrand(context.Context,*CategoryBrandRequest) error
	UpdateCategoryBrand(context.Context,*CategoryBrandRequest) error
}

// GreeterUsecase is a Greeter usecase.
type GoodsCategoryBrandUsecase struct {
	repo GoodsCategoryBrandRepo
	log  *log.Helper
}

// NewGreeterUsecase new a Greeter usecase.
func NewGoodsCategoryBrandUsecase(repo GoodsCategoryBrandRepo, logger log.Logger) *GoodsCategoryBrandUsecase {
	return &GoodsCategoryBrandUsecase{repo: repo, log: log.NewHelper(logger)}
}

func(g *GoodsCategoryBrandUsecase)CategoryBrandList(ctx context.Context,req*CategoryBrandFilterRequest)(*CategoryBrandListResponse,error){
	return g.repo.CategoryBrandList(ctx,req)
}
func(g *GoodsCategoryBrandUsecase)GetCategoryBrandList(ctx context.Context,req*CategoryInfoRequest) (*BrandListResponse,error){
	return g.GetCategoryBrandList(ctx,req)
}
func(g *GoodsCategoryBrandUsecase)CreateCategoryBrand(ctx context.Context,req*CategoryBrandRequest) (*CategoryBrandResponse,error){
	return g.CreateCategoryBrand(ctx,req)
}
func(g *GoodsCategoryBrandUsecase)DeleteCategoryBrand(ctx context.Context,req*CategoryBrandRequest) error{
	return g.DeleteCategoryBrand(ctx,req)
}
func(g *GoodsCategoryBrandUsecase)UpdateCategoryBrand(ctx context.Context,req*CategoryBrandRequest) error{
	return g.UpdateCategoryBrand(ctx,req)
}