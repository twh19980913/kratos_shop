package service

import (
	"context"
	"fmt"
	pb "goods_srv/api/helloworld/v1"
	"goods_srv/internal/biz"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *GoodsService) BrandList(ctx context.Context, req *pb.BrandFilterRequest) (*pb.BrandListResponse, error) {
	rsp,err := s.bu.BrandList(ctx,&biz.BrandFilterRequest{
		Pages: req.Pages,
		PagePerNums: req.PagePerNums,
	})
	if err != nil {
		return nil,err
	}
	brandListResponse := pb.BrandListResponse{}
	brandListResponse.Total = rsp.Total

	var brandResponses []*pb.BrandInfoResponse
	for _,data := range rsp.Data{
		fmt.Println(data.Name)
		brandResponse := pb.BrandInfoResponse{
			Id: data.Id,
			Name: data.Name,
			Logo: data.Logo,
		}
		brandResponses = append(brandResponses, &brandResponse)
	}
	brandListResponse.Data = brandResponses
	return &brandListResponse,nil
}
func (s *GoodsService) CreateBrand(ctx context.Context, req *pb.BrandRequest) (*pb.BrandInfoResponse, error) {
	rsp,_ := s.bu.CreateBrand(ctx,&biz.BrandRequest{
		Name: req.Name,
		Logo: req.Logo,
	})
	return &pb.BrandInfoResponse{
		Id: rsp.Id,
	},nil
}
func (s *GoodsService) DeleteBrand(ctx context.Context, req *pb.BrandRequest) (*emptypb.Empty, error) {
	s.bu.DeleteBrand(ctx,&biz.BrandRequest{
		Id: req.Id,
	})
	return &emptypb.Empty{},nil
}
func (s *GoodsService) UpdateBrand(ctx context.Context, req *pb.BrandRequest) (*emptypb.Empty, error) {
	s.bu.UpdateBrand(ctx,&biz.BrandRequest{
		Id: req.Id,
		Name: req.Name,
		Logo: req.Logo,
	})
	return &emptypb.Empty{},nil
}