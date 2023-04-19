package service

import (
	"context"
	pb "goods_srv/api/helloworld/v1"
	"goods_srv/internal/biz"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *GoodsService) BannerList(ctx context.Context, req *emptypb.Empty) (*pb.BannerListResponse, error) {
	bannerListReponse := pb.BannerListResponse{}
	rsp,_ := s.bannerUsecase.BannerList(ctx)
	bannerListReponse.Total = rsp.Total

	var bannerResponses []*pb.BannerResponse

	for _,banner := range rsp.Data{
		bannerResponses = append(bannerResponses, &pb.BannerResponse{
			Id: banner.Id,
			Image: banner.Image,
			Index: banner.Index,
			Url: banner.Url,
		})
	}

	bannerListReponse.Data = bannerResponses
	return &bannerListReponse,nil
}
func (s *GoodsService) CreateBanner(ctx context.Context, req *pb.BannerRequest) (*pb.BannerResponse, error) {
	rsp,_ := s.bannerUsecase.CreateBanner(ctx,&biz.BannerRequest{
		Image: req.Image,
		Index: req.Index,
		Url: req.Url,
	})
	return &pb.BannerResponse{Id: rsp.Id},nil
}
func (s *GoodsService) DeleteBanner(ctx context.Context, req *pb.BannerRequest) (*emptypb.Empty, error) {
	s.bannerUsecase.DeleteBanner(ctx,&biz.BannerRequest{
		Id: req.Id,
	})
	return &emptypb.Empty{},nil
}
func (s *GoodsService) UpdateBanner(ctx context.Context, req *pb.BannerRequest) (*emptypb.Empty, error) {
	s.bannerUsecase.UpdateBanner(ctx,&biz.BannerRequest{
		Id: req.Id,
		Image: req.Image,
		Url: req.Url,
		Index: req.Index,
	})
	return &emptypb.Empty{},nil
}