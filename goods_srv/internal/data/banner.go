package data

import (
	"context"
	errors "github.com/go-kratos/kratos/v2/errors"
	"goods_srv/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)


type Banner struct {
	BaseModel
	Image string `gorm:"type:varchar(200);not null"`
	Url   string `gorm:"type:varchar(200);not null"`
	Index int32  `gorm:"type:int;default:1;not null"`
}

type bannerRepo struct {
	data *Data
	log  *log.Helper
}

// NewGreeterRepo .
func NewBannerRepo(data *Data, logger log.Logger) biz.BannerRepo {
	return &bannerRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}


func (b *bannerRepo) BannerList(context.Context) (*biz.BannerListResponse,error){
	bannerLIstResponse := biz.BannerListResponse{}

	var banners []Banner
	result := b.data.db.Find(&banners)
	bannerLIstResponse.Total = int32(result.RowsAffected)

	var bannerResponses []*biz.BannerResponse
	for _,banner := range banners{
		bannerResponses = append(bannerResponses, &biz.BannerResponse{
			Id: banner.ID,
			Image: banner.Image,
			Index: banner.Index,
			Url: banner.Url,
		})
	}

	bannerLIstResponse.Data = bannerResponses
	return &bannerLIstResponse,nil
}

func (b *bannerRepo)CreateBanner(ctx context.Context,req *biz.BannerRequest) (*biz.BannerResponse,error){
	banner := Banner{}
	banner.Image = req.Image
	banner.Index = req.Index
	banner.Url = req.Url

	b.data.db.Save(&banner)
	return &biz.BannerResponse{Id: banner.ID},nil
}

func (b *bannerRepo) DeleteBanner(ctx context.Context,req *biz.BannerRequest)error{
	if result := b.data.db.Delete(&Banner{},req.Id);result.RowsAffected == 0{
		return errors.New(404,"NotFound","轮播图不存在")
	}
	return nil
}

func (b *bannerRepo) UpdateBanner(ctx context.Context,req *biz.BannerRequest)error{
	var banner Banner

	if result := b.data.db.First(&banner,req.Id);result.RowsAffected == 0{
		return errors.New(404,"NotFound","轮播图不存在")
	}

	if req.Url != "" {
		banner.Url = req.Url
	}

	if req.Image != "" {
		banner.Image = req.Image
	}

	if req.Index != 0{
		banner.Index = req.Index
	}

	b.data.db.Save(&banner)
	return nil
}