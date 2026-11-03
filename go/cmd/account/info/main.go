// 获取账号信息
//
// 接口地址：https://api.zjzapi.com/user/info
// 请求方式：POST（application/x-www-form-urlencoded）
// 接口费用：调用免费，数据有10分钟延迟。
//
// 用法（在 go 目录运行）：
//
//	go run ./cmd/account/info
//
// 注意：本接口不使用 API Key，而是使用用户中心的 access_token
// （在官网「用户中心 - 账户信息」获取），请填到 config.go 的 AccessToken。
package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"zjz-demo/zjzapi"
)

// accountInfoData data 部分的字段。
type accountInfoData struct {
	Mobile  zjzapi.FlexString `json:"mobile"`
	Money   zjzapi.FlexString `json:"money"`
	AppNum  zjzapi.FlexString `json:"app_num"`
	DayMake zjzapi.FlexString `json:"day_make"`
}

func main() {
	if zjzapi.AccessToken == "" || zjzapi.AccessToken == "YOUR_ACCESS_TOKEN" {
		fmt.Fprintln(os.Stderr, "提示：请先在 go/zjzapi/config.go 中将 AccessToken 替换为你的 access_token（官网「用户中心 - 账户信息」获取）。")
		os.Exit(1)
	}

	result, err := zjzapi.Post("user/info", url.Values{
		"access_token": {zjzapi.AccessToken},
	})
	if err != nil {
		fmt.Println("调用失败：", err)
		os.Exit(1)
	}

	if !result.IsSuccess() {
		fmt.Printf("业务失败：code=%d，msg=%s\n", result.Code, result.Msg)
		os.Exit(1)
	}

	var data accountInfoData
	if err := json.Unmarshal(result.Data, &data); err != nil {
		fmt.Println("解析 data 失败：", err)
		os.Exit(1)
	}
	fmt.Printf("登录账号：%s\n", data.Mobile)
	fmt.Printf("账户余额：%s 元\n", data.Money)
	fmt.Printf("应用数量：%s\n", data.AppNum)
	fmt.Printf("今日调用：%s 次\n", data.DayMake)
}
