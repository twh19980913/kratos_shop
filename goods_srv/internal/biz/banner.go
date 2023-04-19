package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type BannerRequest struct{
    Id int32
    Index int32
    Image string
    Url string
  }

  type BannerResponse struct{
    Id int32
    Index int32
	Image string
    Url string
  }

  type BannerListResponse struct{
    Total int32
    Data []*BannerResponse
  }

type BannerRepo interface {
	BannerList(context.Context) (*BannerListResponse,error)
	CreateBanner(context.Context,*BannerRequest) (*BannerResponse,error)
	DeleteBanner(context.Context,*BannerRequest) error
	UpdateBanner(context.Context,*BannerRequest) error
}

// GreeterUsecase is a Greeter usecase.
type BannerUsecase struct {
	repo BannerRepo
	log  *log.Helper
}

// NewGreeterUsecase new a Greeter usecase.
func NewBannerUsecase(repo BannerRepo, logger log.Logger) *BannerUsecase {
	return &BannerUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (b *BannerUsecase)BannerList (ctx context.Context) (*BannerListResponse,error){
	return b.repo.BannerList(ctx)
}

func (b *BannerUsecase) CreateBanner(ctx context.Context,req *BannerRequest) (*BannerResponse,error){
	return b.repo.CreateBanner(ctx,req)
}

func  (b *BannerUsecase) DeleteBanner(ctx context.Context,req *BannerRequest) error{
	return b.repo.DeleteBanner(ctx,req)
}

func  (b *BannerUsecase) UpdateBanner(ctx context.Context,req *BannerRequest) error{
	return b.repo.UpdateBanner(ctx,req)
}