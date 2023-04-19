package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

type GoodsInvInfo struct{
    GoodsId int32
    Num int32
}

type SellInfo struct{
    GoodsInfo []*GoodsInvInfo
	OrderSn string
}

// GreeterRepo is a Greater repo.
type InventoryRepo interface {
	SetInv(ctx context.Context,req *GoodsInvInfo) error
	InvDetail(ctx context.Context,req *GoodsInvInfo) (*GoodsInvInfo,error)
	Sell(ctx context.Context,req *SellInfo) error
	Reback(ctx context.Context,req *SellInfo) error
	AutoReback (ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error)
}

// GreeterUsecase is a Greeter usecase.
type InventoryUsecase struct {
	repo InventoryRepo
	log  *log.Helper
}

// NewGreeterUsecase new a Greeter usecase.
func NewInventoryUsecase(repo InventoryRepo, logger log.Logger) *InventoryUsecase {
	return &InventoryUsecase{repo: repo, log: log.NewHelper(logger)}
}

func(i *InventoryUsecase)SetInv(ctx context.Context,req *GoodsInvInfo) error{
	return i.repo.SetInv(ctx,req)
}
func(i *InventoryUsecase)InvDetail(ctx context.Context,req *GoodsInvInfo) (*GoodsInvInfo,error){
	return i.repo.InvDetail(ctx,req)
}
func(i *InventoryUsecase)Sell(ctx context.Context,req *SellInfo) error{
	return i.repo.Sell(ctx,req)
}
func(i *InventoryUsecase)Reback(ctx context.Context,req *SellInfo) error{
	return i.repo.Reback(ctx,req)
}

func(i *InventoryUsecase)AutoReback (ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error){
	return i.repo.AutoReback(ctx,msgs...)
}