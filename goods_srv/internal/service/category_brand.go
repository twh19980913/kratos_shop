package service

import (
	"context"
	pb "goods_srv/api/helloworld/v1"
	"goods_srv/internal/biz"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *GoodsService) CategoryBrandList(ctx context.Context, req *pb.CategoryBrandFilterRequest) (*pb.CategoryBrandListResponse, error) {
	rsp,_ := s.categoryBrandUsecase.CategoryBrandList(ctx,&biz.CategoryBrandFilterRequest{
		Pages: req.Pages,
		PagePerNums: req.PagePerNums,
	})
	categoryBrandListResponse := pb.CategoryBrandListResponse{}
	categoryBrandListResponse.Total = rsp.Total

	var categoryResponses []*pb.CategoryBrandResponse
	for _,categoryBrand := range rsp.Data{
		categoryResponses = append(categoryResponses, &pb.CategoryBrandResponse{
			Category: &pb.CategoryInfoResponse{
				Id: categoryBrand.Category.Id,
				Name: categoryBrand.Category.Name,
				Level: categoryBrand.Category.Level,
				IsTab: categoryBrand.Category.IsTab,
				ParentCategory: categoryBrand.Category.ParentCategory,
			},
			Brand: &pb.BrandInfoResponse{
				Id:   categoryBrand.Brand.Id,
				Name: categoryBrand.Brand.Name,
				Logo: categoryBrand.Brand.Logo,
			},
		})
	}

	categoryBrandListResponse.Data = categoryResponses
	return &categoryBrandListResponse, nil
}
func (s *GoodsService) GetCategoryBrandList(ctx context.Context, req *pb.CategoryInfoRequest) (*pb.BrandListResponse, error) {
	rsp ,_ := s.categoryBrandUsecase.GetCategoryBrandList(ctx,&biz.CategoryInfoRequest{
		Id: req.Id,
	})
	brandListResponse := pb.BrandListResponse{}
	brandListResponse.Total = rsp.Total

	var brandInfoResponse []*pb.BrandInfoResponse
	for _,categoryBrand := range rsp.Data{
		brandInfoResponse = append(brandInfoResponse, &pb.BrandInfoResponse{
			Id: categoryBrand.Id,
			Name: categoryBrand.Name,
			Logo: categoryBrand.Logo,
		})
	}
	brandListResponse.Data = brandInfoResponse
	return &brandListResponse,nil
}
func (s *GoodsService) CreateCategoryBrand(ctx context.Context, req *pb.CategoryBrandRequest) (*pb.CategoryBrandResponse, error) {
	return &pb.CategoryBrandResponse{}, nil
}
func (s *GoodsService) DeleteCategoryBrand(ctx context.Context, req *pb.CategoryBrandRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
func (s *GoodsService) UpdateCategoryBrand(ctx context.Context, req *pb.CategoryBrandRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}