package service

import (
	"context"

	pb "order_srv/api/helloworld/v1"
	"order_srv/internal/biz"
	"github.com/apache/rocketmq-client-go/v2/primitive"
  "github.com/apache/rocketmq-client-go/v2/consumer"
	"google.golang.org/protobuf/types/known/emptypb"
)

type OrderService struct {
	pb.UnimplementedOrderServer
	ou *biz.OrderUsecase
}

func NewOrderService(ou *biz.OrderUsecase) *OrderService {
	return &OrderService{ou: ou}
}

func(o *OrderService)OrderTimeout (ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error){
	return o.ou.OrderTimeout(ctx,msgs...)
}

func (s *OrderService) CartItemList(ctx context.Context, req *pb.UserInfo) (*pb.CartItemListResponse, error) {
	var cartItemListResponse pb.CartItemListResponse
	rsp,err := s.ou.CartItemList(ctx,&biz.UserInfo{
		Id: req.Id,
	})
	if err != nil {
		return nil,err
	}
	cartItemListResponse.Total = rsp.Total

	for _,shopCart := range rsp.Data{
		cartItemListResponse.Data = append(cartItemListResponse.Data, &pb.ShopCartInfoResponse{
			Id: shopCart.Id,
			UserId: shopCart.UserId,
			GoodsId: shopCart.GoodsId,
			Nums: shopCart.Nums,
			Checked: shopCart.Checked,
		})
	}

	return &cartItemListResponse,nil
}
func (s *OrderService) CreateCartItem(ctx context.Context, req *pb.CartItemRequest) (*pb.ShopCartInfoResponse, error) {
	rsp,err := s.ou.CreateCartItem(ctx,&biz.CartItemRequest{
		Id: req.Id,
		UserId: req.UserId,
		GoodsId: req.GoodsId,
		GoodsName: req.GoodsName,
		GoodsImage: req.GoodsImage,
		GoodsPrice: req.GoodsPrice,
		Nums: req.Nums,
		Checked: req.Checked,
	})
	if err != nil {
		return nil,err
	}

	return &pb.ShopCartInfoResponse{Id: rsp.Id},nil
}
func (s *OrderService) UpdateCartItem(ctx context.Context, req *pb.CartItemRequest) (*emptypb.Empty, error) {
	err := s.ou.UpdateCartItem(ctx,&biz.CartItemRequest{
		Id: req.Id,
		UserId: req.UserId,
		GoodsId: req.GoodsId,
		GoodsName: req.GoodsName,
		GoodsImage: req.GoodsImage,
		GoodsPrice: req.GoodsPrice,
		Nums: req.Nums,
		Checked: req.Checked,
	})
	if err != nil {
		return  &emptypb.Empty{},err
	}
	return &emptypb.Empty{},nil
}
func (s *OrderService) DeleteCartItem(ctx context.Context, req *pb.CartItemRequest) (*emptypb.Empty, error) {
	err := s.ou.DeleteCartItem(ctx,&biz.CartItemRequest{
		GoodsId: req.GoodsId,
		UserId: req.UserId,
	})
	if err != nil {
		return  &emptypb.Empty{},err
	}
	return &emptypb.Empty{},nil
}
func (s *OrderService) CreateOrder(ctx context.Context, req *pb.OrderRequest) (*pb.OrderInfoResponse, error) {
	rsp,err := s.ou.CreateOrder(ctx,&biz.OrderRequest{
		Id: req.Id,
		UserId: req.UserId,
		Address: req.Address,
		Name: req.Name,
		Mobile: req.Mobile,
		Post: req.Post,
	})
	if err != nil {
		return nil,err
	}

	return &pb.OrderInfoResponse{
		Id: rsp.Id,
		OrderSn: rsp.OrderSn,
		Total: rsp.Total,
	},nil
}
func (s *OrderService) OrderList(ctx context.Context, req *pb.OrderFilterRequest) (*pb.OrderListResponse, error) {
	var orderListResponse pb.OrderListResponse
	rsp,err := s.ou.OrderList(ctx,&biz.OrderFilterRequest{
		UserId: req.UserId,
		Pages: req.Pages,
		PagePerNums: req.PagePerNums,
	})
	if err != nil {
		return nil,err
	}
	orderListResponse.Total = rsp.Total

	for _,order := range rsp.Data{
		orderListResponse.Data = append(orderListResponse.Data, &pb.OrderInfoResponse{
			Id: order.Id,
			UserId: order.UserId,
			OrderSn: order.OrderSn,
			PayType: order.PayType,
			Status: order.Status,
			Post: order.Post,
			Total: order.Total,
			Address: order.Address,
			Name: order.Name,
			Mobile: order.Mobile,
			AddTime: order.AddTime,
		})
	}
	return &orderListResponse,nil
}
func (s *OrderService) OrderDetail(ctx context.Context, req *pb.OrderRequest) (*pb.OrderInfoDetailResponse, error) {
	var orderInfoDetailResponse pb.OrderInfoDetailResponse
	rsp,err := s.ou.OrderDetail(ctx,&biz.OrderRequest{
		Id: req.Id,
		UserId: req.UserId,
	})
	if err != nil {
		return nil,err
	}
	orderInfo := pb.OrderInfoResponse{}
	orderInfo.Id = rsp.OrderInfo.Id
	orderInfo.UserId = rsp.OrderInfo.UserId
	orderInfo.OrderSn = rsp.OrderInfo.OrderSn
	orderInfo.PayType = rsp.OrderInfo.PayType
	orderInfo.Status = rsp.OrderInfo.Status
	orderInfo.Post = rsp.OrderInfo.Post
	orderInfo.Total = rsp.OrderInfo.Total
	orderInfo.Address = rsp.OrderInfo.Address
	orderInfo.Name = rsp.OrderInfo.Name
	orderInfo.Mobile = rsp.OrderInfo.Mobile

	orderInfoDetailResponse.OrderInfo = &orderInfo

	for _,orderGood := range rsp.Goods{
		orderInfoDetailResponse.Goods = append(orderInfoDetailResponse.Goods, &pb.OrderItemResponse{
			GoodsId: orderGood.GoodsId,
			GoodsName: orderGood.GoodsName,
			GoodsPrice: orderGood.GoodsPrice,
			Nums: orderGood.Nums,
		})
	}

	return &orderInfoDetailResponse,nil
}
func (s *OrderService) UpdateOrderStatus(ctx context.Context, req *pb.OrderStatus) (*emptypb.Empty, error) {
	err := s.ou.UpdateOrderStatus(ctx,&biz.OrderStatus{
		OrderSn: req.OrderSn,
		Status: req.Status,
	})
	if err != nil {
		return nil,err
	}
	return &emptypb.Empty{},nil
}
