package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type CategoryListRequest struct{
    Id int32
    Level int32
  }
  
  type CategoryInfoRequest struct{
    Id int32
    Name string
    ParentCategory int32
    Level int32
    IsTab bool
  }
  
  type DeleteCategoryRequest struct{
    Id int32
  }
  
  type QueryCategoryRequest struct{
    Id int32
    Name string
  }
  
  type CategoryInfoResponse struct{
    Id int32
    Name string
    ParentCategory int32
    Level int32
    IsTab bool
  }
  
  type CategoryListResponse struct{
    Total int32
    Data []*CategoryInfoResponse
    JsonData string
  }
  
  type SubCategoryListResponse  struct{
	Total int32
    Info *CategoryInfoResponse
    SubCategorys []*CategoryInfoResponse
  }

type CategoryRepo interface {
	GetAllCategorysList(context.Context) (*CategoryListResponse,error)
	GetSubCategory(context.Context,*CategoryListRequest) (*SubCategoryListResponse,error)
	CreateCategory(context.Context,*CategoryInfoRequest) (*CategoryInfoResponse,error)
	DeleteCategory(context.Context,*DeleteCategoryRequest) error
	UpdateCategory(context.Context,*CategoryInfoRequest) error
}

// GreeterUsecase is a Greeter usecase.
type CategoryUsecase struct {
	repo CategoryRepo
	log  *log.Helper
}

// NewGreeterUsecase new a Greeter usecase.
func NewCategoryUsecase(repo CategoryRepo, logger log.Logger) *CategoryUsecase {
	return &CategoryUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (c CategoryUsecase) GetAllCategorysList(ctx context.Context) (*CategoryListResponse,error){
	return c.repo.GetAllCategorysList(ctx)
}
func (c CategoryUsecase) GetSubCategory(ctx context.Context,req *CategoryListRequest) (*SubCategoryListResponse,error){
	return c.repo.GetSubCategory(ctx,req)
}
func (c CategoryUsecase)CreateCategory(ctx context.Context,req *CategoryInfoRequest) (*CategoryInfoResponse,error){
	return c.CreateCategory(ctx,req)
}
func (c CategoryUsecase)DeleteCategory(ctx context.Context,req *DeleteCategoryRequest) error{
	return c.DeleteCategory(ctx,req)
}
func (c CategoryUsecase)UpdateCategory(ctx context.Context,req *CategoryInfoRequest) error{
	return c.UpdateCategory(ctx,req)
}