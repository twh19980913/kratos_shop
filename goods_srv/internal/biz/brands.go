package biz

import (
	"time"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type User struct {
	ID int32 
	CreatedAt time.Time 
	UpdatedAt time.Time 
	DeletedAt time.Time
	Name string 
	Logo string
}

type BrandFilterRequest struct{
    Pages int32
    PagePerNums int32
}

type BrandInfoResponse struct{
    Id int32
    Name string
    Logo string
  }

type BrandListResponse struct{
	Total int32
    Data []*BrandInfoResponse
}

type BrandRequest struct{
    Id int32
    Name string
    Logo string
  }

type BrandsRepo interface {
	BrandList(context.Context,*BrandFilterRequest) (*BrandListResponse,error)
	CreateBrand (context.Context,*BrandRequest) (*BrandInfoResponse,error)
	DeleteBrand (context.Context,*BrandRequest) error
	UpdateBrand (context.Context,*BrandRequest) error
}

// GreeterUsecase is a Greeter usecase.
type BrandsUsecase struct {
	repo BrandsRepo
	log  *log.Helper
}

// NewGreeterUsecase new a Greeter usecase.
func NewBrandsUsecase(repo BrandsRepo, logger log.Logger) *BrandsUsecase {
	return &BrandsUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (b BrandsUsecase)BrandList(ctx context.Context,req *BrandFilterRequest) (*BrandListResponse,error){
	result,err :=  b.repo.BrandList(ctx,req)
	if err != nil {
		b.log.Error(err)
	}
	for _,data := range result.Data{
		b.log.Info(data.Name)
	}
	return result,nil
}

func (b BrandsUsecase)CreateBrand(ctx context.Context,req *BrandRequest) (*BrandInfoResponse,error){
	b.log.Error("报错了")
	rsp,err :=  b.repo.CreateBrand(ctx,req)
	return rsp,err
}

func (b BrandsUsecase)UpdateBrand(ctx context.Context,req *BrandRequest) error{
	return b.repo.UpdateBrand(ctx,req)
}

func (b BrandsUsecase)DeleteBrand(ctx context.Context,req *BrandRequest) error{
	return b.repo.DeleteBrand(ctx,req)
}