package router

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/order-web/api/shop_cart"
	"mxshop-api/order-web/middlewares"
)

func InitShopCartRouter(Router *gin.RouterGroup) {
	ShopCartRouter := Router.Group("shopcarts").Use(middlewares.JWTAuth())
	{
		ShopCartRouter.GET("",shop_cart.List) // 购物车列表
		ShopCartRouter.DELETE("/:id",shop_cart.Delete) // 删除条目
		ShopCartRouter.POST("",shop_cart.New) // 添加商品到购物车
		ShopCartRouter.PATCH("/:id",shop_cart.Update) // 修改条目
	}
}