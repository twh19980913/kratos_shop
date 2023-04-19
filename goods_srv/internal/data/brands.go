package data

import (
	"context"
	errors "github.com/go-kratos/kratos/v2/errors"
	"goods_srv/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)



type Brands struct {
	BaseModel
	Name string `gorm:"type:varchar(20);not null"`
	Logo string `gorm:"type:varchar(200);default:'';not null"`
}

type brandsRepo struct {
	data *Data
	log  *log.Helper
}

// NewGreeterRepo .
func NewBrandsRepo(data *Data, logger log.Logger) biz.BrandsRepo {
	return &brandsRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (b *brandsRepo)BrandList(ctx context.Context,req *biz.BrandFilterRequest) (*biz.BrandListResponse,error){
	brandListResponse := biz.BrandListResponse{}
	var brands []Brands
	result := b.data.db.Scopes(Paginate(int(req.Pages),int(req.PagePerNums))).Find(&brands)
	if result.Error != nil{
		return nil,result.Error
	}

	b.log.Debug(brands)

	var total int64
	b.data.db.Model(&Brands{}).Count(&total)
	brandListResponse.Total = int32(total)

	var brandResponses []*biz.BrandInfoResponse
	for _,brand := range brands{
		brandResponse := biz.BrandInfoResponse{
			Id: brand.ID,
			Name: brand.Name,
			Logo: brand.Logo,
		}
		brandResponses = append(brandResponses, &brandResponse)
	}

	brandListResponse.Data = brandResponses
	return &brandListResponse,nil
}

func (b *brandsRepo)CreateBrand(ctx context.Context,req *biz.BrandRequest) (*biz.BrandInfoResponse,error){
	// 新建品牌
	if result := b.data.db.First(&Brands{},req.Id);result.RowsAffected == 1{
		return nil,errors.New(22,"InvalidArgument","InvalidArgument")
	}
	// 
	brand := Brands{}
	brand.Name = req.Name
	brand.Logo = req.Logo

	b.data.db.Save(&brand)
	return &biz.BrandInfoResponse{
		Id: brand.ID,
	},nil
}

func (b *brandsRepo)DeleteBrand(ctx context.Context,req *biz.BrandRequest) error{
	if result := b.data.db.Delete(&Brands{},req.Id);result.RowsAffected == 0{
		return errors.New(404,"NotFound","品牌不存在")
	}
	return nil
}

func (b *brandsRepo)UpdateBrand(ctx context.Context,req *biz.BrandRequest) error{
	brands := Brands{}
	if result := b.data.db.First(&brands,req.Id);result.RowsAffected == 0{
		return errors.New(404,"NotFound","品牌不存在")
	}

	if req.Name != ""{
		brands.Name = req.Name
	}

	if req.Logo != ""{
		brands.Logo = req.Logo
	}
	b.data.db.Save(&brands)
	return nil
}