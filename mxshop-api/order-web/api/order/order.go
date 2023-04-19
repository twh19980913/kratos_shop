package order

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
	"go.uber.org/zap"
	"mxshop-api/order-web/api"
	"mxshop-api/order-web/forms"
	"mxshop-api/order-web/global"
	"mxshop-api/order-web/models"
	"mxshop-api/order-web/proto"
	"net/http"
	"strconv"
)

func List(ctx *gin.Context)  {
	//订单列表
	userId,_ := ctx.Get("userId")
	//如果是管理员不传userId 如果不是就传userId
	claims,_ := ctx.Get("claims")

	request := proto.OrderFilterRequest{}

	//如果是管理员用户则返回所有的订单
	model := claims.(*models.CustomClaims)
	if model.AuthorityId == 1 {
		//不是管理员
		request.UserId = int32(userId.(uint))
	}

	pages := ctx.DefaultQuery("p","0")
	pagesInt,_ := strconv.Atoi(pages)
	request.Pages = int32(pagesInt)

	perNums := ctx.DefaultQuery("pnum","0")
	perNumsInt,_ := strconv.Atoi(perNums)
	request.PagePerNums = int32(perNumsInt)

	rsp,err := global.OrderSrvClient.OrderList(context.Background(),&request)
	if err != nil {
		zap.S().Errorw("获取订单列表失败")
		api.HandleGrpcErrorToHttp(err,ctx)
		return
	}

	/*
	{
		"total":100,
		"data":[
			{
				"id":xxx,
				"status":xxx
			},
			{
			}
		]
	}
	*/
	reMap := gin.H{
		"total": rsp.Total,
	}
	orderList := make([]interface{},0)
	for _,item := range rsp.Data{
		tmpMap := map[string]interface{}{}

		tmpMap["id"] = item.Id
		tmpMap["status"] = item.Status
		tmpMap["pay_type"] = item.PayType
		tmpMap["user"] = item.UserId
		tmpMap["post"] = item.Post
		tmpMap["total"] = item.Total
		tmpMap["address"] = item.Address
		tmpMap["name"] = item.Name
		tmpMap["mobile"] = item.Mobile
		tmpMap["order_sn"] = item.OrderSn
		tmpMap["add_time"] = item.AddTime

		orderList = append(orderList, tmpMap)
	}

	reMap["data"] = orderList
	ctx.JSON(http.StatusOK,reMap)
}

func New(ctx *gin.Context) {
	orderForm := forms.CreateOrderForm{}
	if err := ctx.ShouldBindJSON(&orderForm);err != nil{
		api.HandleValidatorError(ctx,err)
		return
	}
	userId,_ := ctx.Get("userId")
	rsp ,err := global.OrderSrvClient.CreateOrder(context.Background(),&proto.OrderRequest{
		UserId: int32(userId.(uint)),
		Name: orderForm.Name,
		Mobile: orderForm.Mobile,
		Address: orderForm.Address,
		Post: orderForm.Post,
	})
	if err != nil {
		zap.S().Errorw("新建订单失败")
		api.HandleGrpcErrorToHttp(err,ctx)
		return
	}
	//TODO 返回支付宝的支付url
	//============================生成支付宝的支付URL========================
	client,err := alipay.New(global.ServerConfig.AliPayInfo.AppID,
		global.ServerConfig.AliPayInfo.PrivateKey,false)
	if err != nil {
		zap.S().Errorw("实例化支付宝失败")
		ctx.JSON(http.StatusInternalServerError,gin.H{
			"msg":err.Error(),
		})
		return
	}

	err = client.LoadAliPayPublicKey(global.ServerConfig.AliPayInfo.AliPublicKey)
	if err != nil {
		zap.S().Errorw("加载支付宝公钥失败")
		ctx.JSON(http.StatusInternalServerError,gin.H{
			"msg":err.Error(),
		})
		return
	}

	var p = alipay.TradePagePay{}
	p.NotifyURL = global.ServerConfig.AliPayInfo.NotifyURL
	p.ReturnURL = global.ServerConfig.AliPayInfo.ReturnURL // 跳转页面
	p.Subject = "慕学生鲜订单-" + rsp.OrderSn
	p.OutTradeNo = rsp.OrderSn
	p.TotalAmount = strconv.FormatFloat(float64(rsp.Total),'f',2,64)
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	url,err := client.TradePagePay(p)
	if err != nil {
		zap.S().Errorw("生成支付url失败")
		ctx.JSON(http.StatusInternalServerError,gin.H{
			"msg":err.Error(),
		})
		return
	}

	//=====================================================================
	ctx.JSON(http.StatusOK,gin.H{
		"id": rsp.Id,
		"alipay_url": url.String(),
	})
}

func Detail(ctx *gin.Context) {
	id := ctx.Param("id") // 订单ID
	i,err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound,gin.H{
			"msg":"url格式出错",
		})
		return
	}

	userId,_ := ctx.Get("userId")
	request := proto.OrderRequest{
		Id: int32(i),
	}
	//如果是管理员不传userId 如果不是就传userId
	claims,_ := ctx.Get("claims")
	//如果是管理员用户则返回所有的订单
	model := claims.(*models.CustomClaims)
	if model.AuthorityId == 1 {
		//不是管理员
		request.UserId = int32(userId.(uint))
	}

	rsp,err := global.OrderSrvClient.OrderDetail(context.Background(),&request)
	if err != nil {
		zap.S().Errorw("获取订单详情失败")
		api.HandleGrpcErrorToHttp(err,ctx)
		return
	}

	reMap := gin.H{}
	reMap["id"] = rsp.OrderInfo.Id
	reMap["status"] = rsp.OrderInfo.Status
	reMap["user"] = rsp.OrderInfo.UserId
	reMap["post"] = rsp.OrderInfo.Post
	reMap["total"] = rsp.OrderInfo.Total
	reMap["address"] = rsp.OrderInfo.Address
	reMap["name"] = rsp.OrderInfo.Name
	reMap["mobile"] = rsp.OrderInfo.Mobile
	reMap["pay_type"] = rsp.OrderInfo.PayType
	reMap["order_sn"] = rsp.OrderInfo.OrderSn

	goodsList := make([]interface{},0)
	for _,item := range rsp.Goods{
		tmpMap := gin.H{
			"id":item.GoodsId,
			"name":item.GoodsName,
			"image":item.GoodsImage,
			"price":item.GoodsPrice,
			"nums":item.Nums,
		}
		goodsList = append(goodsList, tmpMap)
	}
	reMap["goods"] = goodsList

	//==================================生成支付路径============================
	client,err := alipay.New(global.ServerConfig.AliPayInfo.AppID,
		global.ServerConfig.AliPayInfo.PrivateKey,false)
	if err != nil {
		zap.S().Errorw("实例化支付宝失败")
		ctx.JSON(http.StatusInternalServerError,gin.H{
			"msg":err.Error(),
		})
		return
	}

	err = client.LoadAliPayPublicKey(global.ServerConfig.AliPayInfo.AliPublicKey)
	if err != nil {
		zap.S().Errorw("加载支付宝公钥失败")
		ctx.JSON(http.StatusInternalServerError,gin.H{
			"msg":err.Error(),
		})
		return
	}

	var p = alipay.TradePagePay{}
	p.NotifyURL = global.ServerConfig.AliPayInfo.NotifyURL
	p.ReturnURL = global.ServerConfig.AliPayInfo.ReturnURL // 跳转页面
	p.Subject = "慕学生鲜订单-" + rsp.OrderInfo.OrderSn
	p.OutTradeNo = rsp.OrderInfo.OrderSn
	p.TotalAmount = strconv.FormatFloat(float64(rsp.OrderInfo.Total),'f',2,64)
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	url,err := client.TradePagePay(p)
	if err != nil {
		zap.S().Errorw("生成支付url失败")
		ctx.JSON(http.StatusInternalServerError,gin.H{
			"msg":err.Error(),
		})
		return
	}
	//================================================================
	reMap["alipay_url"] = url.String()

	ctx.JSON(http.StatusOK,reMap)
}