// 获取应用信息
//
// 接口地址：https://api.zjzapi.com/user/app
// 请求方式：POST（application/x-www-form-urlencoded）
// 接口费用：调用免费，数据有10分钟延迟。
//
// 用法（在 go 目录运行）：
//
//	go run ./cmd/account/app
//
// 注意：本接口需要同时使用用户中心的 access_token（官网「用户中心 - 账户信息」获取，
// 填到 config.go 的 AccessToken）与你的 API Key。
package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"zjz-demo/zjzapi"
)

// accountAppData data 部分的字段。
type accountAppData struct {
	Name    string            `json:"name"`
	Type    zjzapi.FlexString `json:"type"`
	Status  zjzapi.FlexString `json:"status"`
	MakeNum zjzapi.FlexString `json:"make_num"`
}

func main() {
	if zjzapi.AccessToken == "" || zjzapi.AccessToken == "YOUR_ACCESS_TOKEN" {
		fmt.Fprintln(os.Stderr, "提示：请先在 go/zjzapi/config.go 中将 AccessToken 替换为你的 access_token（官网「用户中心 - 账户信息」获取）。")
		os.Exit(1)
	}

	result, err := zjzapi.Post("user/app", url.Values{
		"access_token": {zjzapi.AccessToken},
		"key":          {zjzapi.APIKey},
	})
	if err != nil {
		fmt.Println("调用失败：", err)
		os.Exit(1)
	}

	if !result.IsSuccess() {
		fmt.Printf("业务失败：code=%d，msg=%s\n", result.Code, result.Msg)
		os.Exit(1)
	}

	var data accountAppData
	if err := json.Unmarshal(result.Data, &data); err != nil {
		fmt.Println("解析 data 失败：", err)
		os.Exit(1)
	}

	// 应用权限：1=证件照v3，2=证件照v5，3=支持全部，4=证件照回执
	typeMap := map[zjzapi.FlexString]string{"1": "证件照v3", "2": "证件照v5", "3": "支持全部", "4": "证件照回执"}
	appType := typeMap[data.Type]
	if appType == "" {
		appType = string(data.Type)
	}
	status := "已停用"
	if data.Status == "1" {
		status = "启用中"
	}

	fmt.Printf("应用名称：%s\n", data.Name)
	fmt.Printf("应用权限：%s\n", appType)
	fmt.Printf("应用状态：%s\n", status)
	fmt.Printf("调用次数：%s\n", data.MakeNum)
}
