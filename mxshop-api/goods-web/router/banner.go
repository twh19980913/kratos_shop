package router

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/goods-web/api/banners"
)

func InitBannerRouter(Router *gin.RouterGroup) {
	BannerRouter := Router.Group("banners")
	{
		BannerRouter.GET("", banners.List)          // 轮播图列表页
		BannerRouter.DELETE("/:id", banners.Delete) // 删除轮播图
		BannerRouter.POST("",  banners.New)       //新建轮播图
		BannerRouter.PUT("/:id", banners.Update) //修改轮播图信息
	}
}