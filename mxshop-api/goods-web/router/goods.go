package router

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/goods-web/api/goods"
)

func InitGoodsRouter(Router *gin.RouterGroup)  {
	GoodsRouter := Router.Group("goods")
	{
		GoodsRouter.GET("",goods.List) // 商品列表
		GoodsRouter.POST("",goods.New) // 该接口需要管理员权限
		GoodsRouter.GET("/:id",goods.Detail) // 获取商品详情
		GoodsRouter.DELETE("/:id",goods.Delete) // 删除商品
		GoodsRouter.GET("/:id/stocks",goods.Stocks) // 获取商品库存

		GoodsRouter.PUT("/:id",goods.Update)
		GoodsRouter.PATCH("/:id",goods.UpdateStatus)
	}
}