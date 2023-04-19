package service

import (
	"context"
	"fmt"

	pb "goods_srv/api/helloworld/v1"
	"goods_srv/internal/biz"

	"google.golang.org/protobuf/types/known/emptypb"
)

type GoodsService struct {
	bu *biz.BrandsUsecase
	bannerUsecase *biz.BannerUsecase
	categoryUsecase *biz.CategoryUsecase
	categoryBrandUsecase *biz.GoodsCategoryBrandUsecase
	goodsUsecase *biz.GoodsUsecase
	pb.UnimplementedGoodsServer
}

func NewGoodsService(bu *biz.BrandsUsecase,bannerUsecase *biz.BannerUsecase,
	categoryUsecase *biz.CategoryUsecase,categoryBrandUsecase *biz.GoodsCategoryBrandUsecase,
	goodsUsecase *biz.GoodsUsecase) *GoodsService {
	return &GoodsService{bu: bu,bannerUsecase: bannerUsecase,
		categoryUsecase: categoryUsecase,
		categoryBrandUsecase: categoryBrandUsecase,
		goodsUsecase: goodsUsecase,
	}
}

func ModelToResponse(goods *biz.GoodsInfoResponse) pb.GoodsInfoResponse {
	return pb.GoodsInfoResponse{
		Id:       goods.Id,
		CategoryId: goods.CategoryId,
		Name: goods.Name,
		GoodsSn: goods.GoodsSn,
		ClickNum: goods.ClickNum,
		SoldNum: goods.SoldNum,
		FavNum: goods.FavNum,
		MarketPrice: goods.MarketPrice,
		ShopPrice: goods.ShopPrice,
		GoodsBrief: goods.GoodsBrief,
		ShipFree: goods.ShipFree,
		GoodsFrontImage: goods.GoodsFrontImage,
		IsNew: goods.IsNew,
		IsHot: goods.IsHot,
		OnSale: goods.OnSale,
		DescImages: goods.DescImages,
		Images: goods.Images,
		Category: &pb.CategoryBriefInfoResponse{
			Id:   goods.Category.Id,
			Name: goods.Category.Name,
		},
		Brand: &pb.BrandInfoResponse{
			Id:   goods.Brand.Id,
			Name: goods.Brand.Name,
			Logo: goods.Brand.Logo,
		},
	}
}

func (s *GoodsService) GoodsList(ctx context.Context, req *pb.GoodsFilterRequest) (*pb.GoodsListResponse, error) {
	rsp,_ := s.goodsUsecase.GoodsList(ctx,&biz.GoodsFilterRequest{
		PriceMin: req.PriceMin,
		PriceMax: req.PriceMax,
		IsHot: req.IsHot,
		IsNew: req.IsNew,
		IsTab: req.IsTab,
		TopCategory: req.TopCategory,
		Pages: req.Pages,
		PagePerNums: req.PagePerNums,
		KeyWords: req.KeyWords,
		Brand: req.Brand,
	})

	goodsListResponse := &pb.GoodsListResponse{}
	goodsListResponse.Total = rsp.Total

	for _,good := range rsp.Data{
		goodsInfoResponse := ModelToResponse(good)
		goodsListResponse.Data = append(goodsListResponse.Data, &goodsInfoResponse)
	}
	return goodsListResponse,nil
}
func (s *GoodsService) BatchGetGoods(ctx context.Context, req *pb.BatchGoodsIdInfo) (*pb.GoodsListResponse, error) {
	fmt.Println("进来了")
	goodsListResponse := &pb.GoodsListResponse{}
	rsp,_ := s.goodsUsecase.BatchGetGoods(ctx,&biz.BatchGoodsIdInfo{
		Id: req.Id,
	})
	fmt.Println("出去了")
	goodsListResponse.Total = rsp.Total
	for _,good := range rsp.Data{
		goodsInfoResponse := ModelToResponse(good)
		goodsListResponse.Data = append(goodsListResponse.Data, &goodsInfoResponse)
	}
	return goodsListResponse,nil
}
func (s *GoodsService) CreateGoods(ctx context.Context, req *pb.CreateGoodsInfo) (*pb.GoodsInfoResponse, error) {
	rsp,_ := s.goodsUsecase.CreateGoods(ctx,&biz.CreateGoodsInfo{
		Name: req.Name,
		GoodsSn: req.GoodsSn,
		Stocks: req.Stocks,
		MarketPrice: req.MarketPrice,
		ShopPrice: req.ShopPrice,
		GoodsBrief: req.GoodsBrief,
		GoodsDesc: req.GoodsDesc,
		ShipFree: req.ShipFree,
		GoodsFrontImage: req.GoodsFrontImage,
		Images: req.Images,
		DescImages: req.DescImages,
		IsNew: req.IsNew,
		IsHot: req.IsHot,
		OnSale: req.OnSale,
		CategoryId: req.CategoryId,
		BrandId: req.BrandId,
	})
	return &pb.GoodsInfoResponse{Id: rsp.Id},nil
}
func (s *GoodsService) DeleteGoods(ctx context.Context, req *pb.DeleteGoodsInfo) (*emptypb.Empty, error) {
	s.goodsUsecase.DeleteGoods(ctx,&biz.DeleteGoodsInfo{
		Id: req.Id,
	})
	return &emptypb.Empty{},nil
}
func (s *GoodsService) UpdateGoods(ctx context.Context, req *pb.CreateGoodsInfo) (*emptypb.Empty, error) {
	s.goodsUsecase.UpdateGoods(ctx,&biz.CreateGoodsInfo{
		Id: req.Id,
		Name: req.Name,
		GoodsSn: req.GoodsSn,
		Stocks: req.Stocks,
		MarketPrice: req.MarketPrice,
		ShopPrice: req.ShopPrice,
		GoodsBrief: req.GoodsBrief,
		GoodsDesc: req.GoodsDesc,
		ShipFree: req.ShipFree,
		GoodsFrontImage: req.GoodsFrontImage,
		Images: req.Images,
		DescImages: req.DescImages,
		IsNew: req.IsNew,
		IsHot: req.IsHot,
		OnSale: req.OnSale,
		CategoryId: req.CategoryId,
		BrandId: req.BrandId,
	})
	return &emptypb.Empty{},nil
}
func (s *GoodsService) GetGoodsDetail(ctx context.Context, req *pb.GoodInfoRequest) (*pb.GoodsInfoResponse, error) {
	rsp,_ := s.goodsUsecase.GetGoodsDetail(ctx,&biz.GoodInfoRequest{
		Id: req.Id,
	})
	goodsInfoResponse := ModelToResponse(rsp)
	return &goodsInfoResponse,nil
}




