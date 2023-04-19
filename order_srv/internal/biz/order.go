package biz

import (
	"github.com/go-kratos/kratos/v2/log"
  "context"
	"github.com/apache/rocketmq-client-go/v2/primitive"
  "github.com/apache/rocketmq-client-go/v2/consumer"
)

type UserInfo struct{
    Id int32
  }
  
  type OrderStatus struct{
    Id int32
    OrderSn string
    Status string
  }
  
  type CartItemRequest struct{
	Id int32
    UserId int32
    GoodsId int32
    GoodsName string
    GoodsImage string
    GoodsPrice float32
    Nums int32
    Checked bool
  }
  
  type OrderRequest struct{
    Id int32
    UserId int32
    Address string
    Name string
    Mobile string
    Post string
  }
  
  type OrderInfoResponse struct{
    Id int32
    UserId int32
    OrderSn string
    PayType string
    Status string
    Post string
	  Total float32
    Address string
    Name string
    Mobile string
    AddTime string
  }
  
  type ShopCartInfoResponse struct{
    Id int32
    UserId int32
    GoodsId int32
    Nums int32
	Checked bool
  }
  
  type OrderItemResponse struct{
    Id int32
	OrderId int32
    GoodsId int32
    GoodsName string
    GoodsImage string
    GoodsPrice float32
    Nums int32
  }
  
  type OrderInfoDetailResponse struct{
    OrderInfo *OrderInfoResponse
    Goods []*OrderItemResponse
  }
  
  type OrderFilterRequest struct{
    UserId int32
    Pages int32
    PagePerNums int32
  }
  
  type OrderListResponse struct{
    Total int32
    Data []*OrderInfoResponse
  }
  
  type CartItemListResponse struct{
    Total int32
    Data []*ShopCartInfoResponse
  }

  type OrderRepo interface {
    CartItemList(ctx context.Context,req *UserInfo) (*CartItemListResponse,error)
    CreateCartItem(ctx context.Context,req *CartItemRequest) (*ShopCartInfoResponse,error)
    UpdateCartItem(ctx context.Context,req *CartItemRequest) error
    DeleteCartItem(ctx context.Context,req *CartItemRequest) error
    OrderList(ctx context.Context,req *OrderFilterRequest) (*OrderListResponse,error)
    OrderDetail(ctx context.Context,req *OrderRequest) (*OrderInfoDetailResponse,error)
    CreateOrder(ctx context.Context,req *OrderRequest) (*OrderInfoResponse,error)
    UpdateOrderStatus(ctx context.Context,req *OrderStatus) error
    OrderTimeout (ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error)
}

// GreeterUsecase is a Greeter usecase.
type OrderUsecase struct {
	repo OrderRepo
	log  *log.Helper
}

// NewGreeterUsecase new a Greeter usecase.
func NewOrderUsecase(repo OrderRepo, logger log.Logger) *OrderUsecase {
	return &OrderUsecase{repo: repo, log: log.NewHelper(logger)}
} 


func(o *OrderUsecase)CartItemList(ctx context.Context,req *UserInfo) (*CartItemListResponse,error){
  return o.repo.CartItemList(ctx,req)
}
func(o *OrderUsecase)CreateCartItem(ctx context.Context,req *CartItemRequest) (*ShopCartInfoResponse,error){
  return o.repo.CreateCartItem(ctx,req)
}
func(o *OrderUsecase)UpdateCartItem(ctx context.Context,req *CartItemRequest) error{
  return o.repo.UpdateCartItem(ctx,req)
}
func(o *OrderUsecase)DeleteCartItem(ctx context.Context,req *CartItemRequest) error{
  return o.repo.DeleteCartItem(ctx,req)
}
func(o *OrderUsecase)OrderList(ctx context.Context,req *OrderFilterRequest) (*OrderListResponse,error){
  return o.repo.OrderList(ctx,req)
}
func(o *OrderUsecase)OrderDetail(ctx context.Context,req *OrderRequest) (*OrderInfoDetailResponse,error){
  return o.repo.OrderDetail(ctx,req)
}
func(o *OrderUsecase)CreateOrder(ctx context.Context,req *OrderRequest) (*OrderInfoResponse,error){
  return o.repo.CreateOrder(ctx,req)
}
func(o *OrderUsecase)UpdateOrderStatus(ctx context.Context,req *OrderStatus) error{
  return o.repo.UpdateOrderStatus(ctx,req)
}

func(o *OrderUsecase)OrderTimeout (ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error){
  return o.repo.OrderTimeout(ctx,msgs...)
}