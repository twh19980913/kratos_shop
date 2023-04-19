package service

import (
	"context"

	pb "inventory_srv/api/helloworld/v1"
	"inventory_srv/internal/biz"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"google.golang.org/protobuf/types/known/emptypb"
)

type InventoryService struct {
	pb.UnimplementedInventoryServer
	iu *biz.InventoryUsecase
}

func NewInventoryService(iu *biz.InventoryUsecase) *InventoryService {
	return &InventoryService{iu: iu}
}

func(i *InventoryService)AutoReback (ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error){
	return i.iu.AutoReback(ctx,msgs...)
}

func (s *InventoryService) SetInv(ctx context.Context, req *pb.GoodsInvInfo) (*emptypb.Empty, error) {
	s.iu.SetInv(ctx,&biz.GoodsInvInfo{
		GoodsId: req.GoodsId,
		Num: req.Num,
	})
	return &emptypb.Empty{},nil
}
func (s *InventoryService) InvDetail(ctx context.Context, req *pb.GoodsInvInfo) (*pb.GoodsInvInfo, error) {
	rsp,_ := s.iu.InvDetail(ctx,&biz.GoodsInvInfo{
		GoodsId: req.GoodsId,
	})
	return &pb.GoodsInvInfo{
		GoodsId: rsp.GoodsId,
		Num: rsp.Num,
	},nil
}
func (s *InventoryService) Sell(ctx context.Context, req *pb.SellInfo) (*emptypb.Empty, error) {
	sellInfo := &biz.SellInfo{}
	for _,info := range req.GoodsInfo{
		sellInfo.GoodsInfo = append(sellInfo.GoodsInfo,&biz.GoodsInvInfo{
			GoodsId: info.GoodsId,
			Num: info.Num,
		})
	}
	sellInfo.OrderSn = req.OrderSn
	s.iu.Sell(ctx,sellInfo)
	return &emptypb.Empty{},nil
}
func (s *InventoryService) Reback(ctx context.Context, req *pb.SellInfo) (*emptypb.Empty, error) {
	sellInfo := &biz.SellInfo{}
	for _,info := range req.GoodsInfo{
		sellInfo.GoodsInfo = append(sellInfo.GoodsInfo,&biz.GoodsInvInfo{
			GoodsId: info.GoodsId,
			Num: info.Num,
		})
	}
	s.iu.Reback(ctx,sellInfo)
	return &emptypb.Empty{},nil
}
