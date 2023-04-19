package data

import (
	"goods_srv/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
	errors "github.com/go-kratos/kratos/v2/errors"
	"context"
)


type GoodsCategoryBrand struct {
	BaseModel
	CategoryID int32 `gorm:"type:int;index:idx_category_brand,unique"`
	Category   Category

	BrandsID int32 `gorm:"type:int;index:idx_category_brand,unique"`
	Brands   Brands
}

func (GoodsCategoryBrand) TableName() string {
	return "goodscategorybrand"
}

type goodsCategoryBrandRepo struct {
	data *Data
	log  *log.Helper
}

// NewGreeterRepo .
func NewGoodsCategoryBrandRepoRepo(data *Data, logger log.Logger) biz.GoodsCategoryBrandRepo {
	return &goodsCategoryBrandRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (g *goodsCategoryBrandRepo)CategoryBrandList(ctx context.Context,req *biz.CategoryBrandFilterRequest)(*biz.CategoryBrandListResponse,error){
	var categoryBrands []GoodsCategoryBrand
	categoryBrandListResponse := biz.CategoryBrandListResponse{}

	var total int64
	g.data.db.Model(&GoodsCategoryBrand{}).Count(&total)
	categoryBrandListResponse.Total = int32(total)

	g.data.db.Preload("Category").Preload("Brands").Scopes(Paginate(int(req.Pages),int(req.PagePerNums))).Find(&categoryBrands)

	var categoryResponses []*biz.CategoryBrandResponse
	for _,categoryBrand := range categoryBrands{
		categoryResponses = append(categoryResponses, &biz.CategoryBrandResponse{
			Category: &biz.CategoryInfoResponse{
				Id: categoryBrand.Category.ID,
				Name: categoryBrand.Category.Name,
				Level: categoryBrand.Category.Level,
				IsTab: categoryBrand.Category.IsTab,
				ParentCategory: categoryBrand.Category.ParentCategoryID,
			},
			Brand: &biz.BrandInfoResponse{
				Id:   categoryBrand.Brands.ID,
				Name: categoryBrand.Brands.Name,
				Logo: categoryBrand.Brands.Logo,
			},
		})
	}

	categoryBrandListResponse.Data = categoryResponses
	return &categoryBrandListResponse,nil
}
func (g *goodsCategoryBrandRepo)GetCategoryBrandList(ctx context.Context,req *biz.CategoryInfoRequest) (*biz.BrandListResponse,error){
	brandListResponse := biz.BrandListResponse{}

	var category Category
	if result := g.data.db.Find(&category,req.Id).First(&category);result.RowsAffected == 0{
		return nil, errors.New(22,"InvalidArgument","商品分类不存在")
	}

	var categoryBrands []GoodsCategoryBrand
	if result := g.data.db.Preload("Brands").Where(&GoodsCategoryBrand{CategoryID: req.Id}).Find(&categoryBrands);result.RowsAffected > 0{
		brandListResponse.Total = int32(result.RowsAffected)
	}

	var brandInfoResponse []*biz.BrandInfoResponse
	for _,categoryBrand := range categoryBrands{
		brandInfoResponse = append(brandInfoResponse, &biz.BrandInfoResponse{
			Id: categoryBrand.Brands.ID,
			Name: categoryBrand.Brands.Name,
			Logo: categoryBrand.Brands.Logo,
		})
	}
	brandListResponse.Data = brandInfoResponse
	return &brandListResponse,nil
}
func (g *goodsCategoryBrandRepo)CreateCategoryBrand(ctx context.Context,req *biz.CategoryBrandRequest) (*biz.CategoryBrandResponse,error){
	var category Category
	if result := g.data.db.First(&category,req.CategoryId);result.RowsAffected == 0{
		return nil, errors.New(22,"InvalidArgument","商品分类不存在")
	}

	var brand Brands
	if result := g.data.db.First(&brand,req.BrandId);result.RowsAffected == 0{
		return nil, errors.New(22,"InvalidArgument","品牌不存在")
	}

	categoryBrand := GoodsCategoryBrand{
		CategoryID: req.CategoryId,
		BrandsID: req.BrandId,
	}

	g.data.db.Save(&categoryBrand)
	return &biz.CategoryBrandResponse{Id: categoryBrand.ID},nil
}
func (g *goodsCategoryBrandRepo)DeleteCategoryBrand(ctx context.Context,req *biz.CategoryBrandRequest) error{
	if result := g.data.db.Delete(&GoodsCategoryBrand{}, req.Id); result.RowsAffected == 0 {
		return  errors.New(22,"InvalidArgument","品牌不存在")
	}
	return nil
}
func (g *goodsCategoryBrandRepo)UpdateCategoryBrand(ctx context.Context,req *biz.CategoryBrandRequest) error{
	var categoryBrand GoodsCategoryBrand

	if result := g.data.db.First(&categoryBrand, req.Id); result.RowsAffected == 0 {
		return errors.New(22,"InvalidArgument","品牌不存在")
	}

	var category Category
	if result := g.data.db.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return errors.New(22,"InvalidArgument","品牌不存在")
	}

	var brand Brands
	if result := g.data.db.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return errors.New(22,"InvalidArgument","品牌不存在")
	}

	categoryBrand.CategoryID = req.CategoryId
	categoryBrand.BrandsID = req.BrandId

	g.data.db.Save(&categoryBrand)

	return  nil
}