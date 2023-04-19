package data

import (
	"context"
	"encoding/json"
	errors "github.com/go-kratos/kratos/v2/errors"
	"goods_srv/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type Category struct {
	BaseModel
	Name             string      `gorm:"type:varchar(20);not null" json:"name"`
	ParentCategoryID int32       `json:"parent"`
	ParentCategory   *Category   `json:"-"`
	SubCategory      []*Category `gorm:"foreignKey:ParentCategoryID;references:ID" json:"sub_category"`
	Level            int32       `gorm:"type:int;not null;default:1" json:"level"`
	IsTab            bool        `gorm:"default:false;not null" json:"is_tab"`
}

type categoryRepo struct {
	data *Data
	log  *log.Helper
}

// NewGreeterRepo .
func NewCategoryRepo(data *Data, logger log.Logger) biz.CategoryRepo {
	return &categoryRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

 func(c *categoryRepo) GetAllCategorysList(ctx context.Context) (*biz.CategoryListResponse,error){
	// 获取所有的分类数据
	/*
	[
		{
			"id":xxx,
			"name":"",
			"level":1,
			"parent":xxxId
			"sub_category":[
				"id":xxx,
				"name":"",
				"level":2,
				"parent":xxxId
				"sub_category":[]
			]
		}
	]
	*/
	var categorys []Category
	c.data.db.Where(&Category{Level: 1}).Preload("SubCategory.SubCategory").Find(&categorys)
	b,_ := json.Marshal(&categorys)
	return &biz.CategoryListResponse{
		JsonData: string(b),
	},nil
 }
 func(c *categoryRepo) GetSubCategory(ctx context.Context,req *biz.CategoryListRequest) (*biz.SubCategoryListResponse,error){
	categoryListResponse := biz.SubCategoryListResponse{}
	var category Category
	if result := c.data.db.First(&category,req.Id);result.RowsAffected == 0{
		return nil,errors.New(404,"NotFound","商品分类不存在")
	}

	categoryListResponse.Info = &biz.CategoryInfoResponse{
		Id: category.ID,
		Name: category.Name,
		Level: category.Level,
		IsTab: category.IsTab,
		ParentCategory: category.ParentCategoryID,
	}
	var subCategorys []Category
	var subCategoryResponse []*biz.CategoryInfoResponse
	preloads := "SubCategory"
	//构造子分类
	if category.Level == 1 {
		preloads = "SubCategory.SubCategory"
	}

	c.data.db.Where(&Category{ParentCategoryID: req.Id}).Preload(preloads).Find(&subCategorys)

	for _,subCategory := range subCategorys{
		subCategoryResponse = append(subCategoryResponse, &biz.CategoryInfoResponse{
			Id: subCategory.ID,
			Name: subCategory.Name,
			Level: category.Level,
			IsTab: subCategory.IsTab,
			ParentCategory: subCategory.ParentCategoryID,
		})
	}

	categoryListResponse.SubCategorys = subCategoryResponse
	return &categoryListResponse,nil
 }
 func(c *categoryRepo) CreateCategory(ctx context.Context,req *biz.CategoryInfoRequest) (*biz.CategoryInfoResponse,error){
	category := Category{}
	cMap := map[string]interface{}{}
	cMap["name"] = req.Name
	cMap["level"] = req.Level
	cMap["is_tab"] = req.IsTab
	if req.Level != 1 {
		cMap["parent_category_id"] = req.ParentCategory
	}
	c.data.db.Create(cMap)
	return &biz.CategoryInfoResponse{
		Id: category.ID,
	},nil
 }
 func(c *categoryRepo) DeleteCategory(ctx context.Context,req *biz.DeleteCategoryRequest) error{
	if result := c.data.db.Delete(&Category{},req.Id);result.RowsAffected == 0{
		return errors.New(404,"NotFound","商品分类不存在")
	}
	return nil
 }
 func(c *categoryRepo) UpdateCategory(ctx context.Context,req *biz.CategoryInfoRequest) error{
	var category Category

	if result := c.data.db.First(&category,req.Id);result.RowsAffected == 0{
		return errors.New(404,"NotFound","商品分类不存在")
	}

	if req.Name != ""{
		category.Name = req.Name
	}

	if req.ParentCategory != 0 {
		category.ParentCategoryID = req.ParentCategory
	}

	if req.Level != 0 {
		category.Level = req.Level
	}

	if req.IsTab {
		category.IsTab = req.IsTab
	}

	c.data.db.Save(&category)
	return nil
 }