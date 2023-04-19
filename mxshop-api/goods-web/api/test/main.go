package main

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"log"
)

func main() {
	client,err := dysmsapi.NewClientWithAccessKey("cn-beijing","LTAI4G93fa8kbGK4Nva8tSrh","eAwRFpEbY28AiKGILO75YuxTVTIvO1")
	if err != nil {
		log.Panic(err)
	}

	smsCode := "123456"
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https"
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = "cn-beijing"
	request.QueryParams["PhoneNumbers"] = "13223403830"                            //手机号
	request.QueryParams["SignName"] = "我的米粒在线教育网站"                                       //阿里云验证过的项目名 自己设置
	request.QueryParams["TemplateCode"] = "SMS_204111019"                          //阿里云的短信模板号 自己设置
	request.QueryParams["TemplateParam"] = "{\"code\":" + smsCode + "}" //短信模板中的验证码内容 自己生成   之前试过直接返回，但是失败，加上code成功。
	response, err := client.ProcessCommonRequest(request)
	fmt.Print(client.DoAction(request, response))

	if err != nil {
		fmt.Print(err.Error())
	}
}
