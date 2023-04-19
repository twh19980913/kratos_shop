package service

import (
	"context"
	pb "goods_srv/api/helloworld/v1"
	"goods_srv/internal/biz"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *GoodsService) GetAllCategorysList(ctx context.Context, req *emptypb.Empty) (*pb.CategoryListResponse, error) {
	rsp,_ := s.categoryUsecase.GetAllCategorysList(ctx)
	return &pb.CategoryListResponse{
		JsonData: rsp.JsonData,
	},nil
}
func (s *GoodsService) GetSubCategory(ctx context.Context, req *pb.CategoryListRequest) (*pb.SubCategoryListResponse, error) {
	rsp,_ := s.categoryUsecase.GetSubCategory(ctx,&biz.CategoryListRequest{
		Id: req.Id,
		Level: req.Level,
	})

	categoryListResponse := pb.SubCategoryListResponse{}
	categoryListResponse.Info = &pb.CategoryInfoResponse{
		Id: rsp.Info.Id,
		Name: rsp.Info.Name,
		Level: rsp.Info.Level,
		IsTab: rsp.Info.IsTab,
		ParentCategory: rsp.Info.ParentCategory,
	}

	var subCategoryResponse []*pb.CategoryInfoResponse
	for _,subCategory := range rsp.SubCategorys{
		subCategoryResponse = append(subCategoryResponse, &pb.CategoryInfoResponse{
			Id: subCategory.Id,
			Name: subCategory.Name,
			Level: subCategory.Level,
			IsTab: subCategory.IsTab,
			ParentCategory: subCategory.ParentCategory,
		})
	}
	categoryListResponse.SubCategorys = subCategoryResponse
	return &categoryListResponse,nil
}
func (s *GoodsService) CreateCategory(ctx context.Context, req *pb.CategoryInfoRequest) (*pb.CategoryInfoResponse, error) {
	rsp ,_ := s.categoryUsecase.CreateCategory(ctx,&biz.CategoryInfoRequest{
		Name: req.Name,
		ParentCategory: req.ParentCategory,
		Level: req.Level,
		IsTab: req.IsTab,
	})

	return &pb.CategoryInfoResponse{Id: rsp.Id},nil
}
func (s *GoodsService) DeleteCategory(ctx context.Context, req *pb.DeleteCategoryRequest) (*emptypb.Empty, error) {
	s.categoryUsecase.DeleteCategory(ctx,&biz.DeleteCategoryRequest{
		Id: req.Id,
	})
	return &emptypb.Empty{},nil
}
func (s *GoodsService) UpdateCategory(ctx context.Context, req *pb.CategoryInfoRequest) (*emptypb.Empty, error) {
	s.categoryUsecase.UpdateCategory(ctx,&biz.CategoryInfoRequest{
		Id: req.Id,
		Name: req.Name,
		ParentCategory: req.ParentCategory,
		Level: req.Level,
		IsTab: req.IsTab,
	})
	return &emptypb.Empty{},nil
}