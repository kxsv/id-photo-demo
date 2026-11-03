// 证件照规格详情
//
// 接口地址：https://api.zjzapi.com/item/get
// 请求方式：POST（application/x-www-form-urlencoded）
// 接口费用：调用免费
//
// 用法（在 go 目录运行）：
//
//	go run ./cmd/item/detail
//
// 输出单个规格的像素尺寸、毫米尺寸、dpi、文件大小要求与支持底色，
// 用于制作类示例（如 v3/all）选择 item_id 与底色。
package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"zjz-demo/zjzapi"
)

// ===== 请求参数（修改后直接运行即可）=====
var (
	itemID = "1" // 规格id，可通过「证件照规格列表」示例获取（go run ./cmd/item/list）
)

// itemDetail 规格详情字段。
type itemDetail struct {
	ID           zjzapi.FlexString `json:"id"`
	Name         string            `json:"name"`
	WidthPx      zjzapi.FlexString `json:"width_px"`
	HeightPx     zjzapi.FlexString `json:"height_px"`
	WidthMm      zjzapi.FlexString `json:"width_mm"`
	HeightMm     zjzapi.FlexString `json:"height_mm"`
	Dpi          zjzapi.FlexString `json:"dpi"`
	FileSizeMin  interface{}       `json:"file_size_min"`
	FileSizeMax  interface{}       `json:"file_size_max"`
	FileSizeMsg  string            `json:"file_size_msg"`
	Color        interface{}       `json:"color"`
	Category     string            `json:"category"`
	CategoryID   zjzapi.FlexString `json:"category_id"`
}

func main() {
	if zjzapi.APIKey == "" || zjzapi.APIKey == "YOUR_APP_KEY" {
		fmt.Fprintln(os.Stderr, "提示：请先在 go/zjzapi/config.go 中将 APIKey 替换为你的 API Key。")
	}

	result, err := zjzapi.Post("item/get", url.Values{
		"key":     {zjzapi.APIKey},
		"item_id": {itemID},
	})
	if err != nil {
		fmt.Println("调用失败：", err)
		os.Exit(1)
	}

	if !result.IsSuccess() {
		fmt.Printf("业务失败：code=%d，msg=%s\n", result.Code, result.Msg)
		os.Exit(1)
	}

	var data itemDetail
	if err := json.Unmarshal(result.Data, &data); err != nil {
		fmt.Println("解析 data 失败：", err)
		os.Exit(1)
	}

	fmt.Printf("规格id：%s\n", data.ID)
	fmt.Printf("名称：%s\n", data.Name)
	fmt.Printf("像素-宽：%s\n", data.WidthPx)
	fmt.Printf("像素-高：%s\n", data.HeightPx)
	fmt.Printf("毫米-宽：%s\n", data.WidthMm)
	fmt.Printf("毫米-高：%s\n", data.HeightMm)
	fmt.Printf("分辨率：%s\n", data.Dpi)
	fmt.Printf("文件最小(kb)：%v\n", data.FileSizeMin)
	fmt.Printf("文件最大(kb)：%v\n", data.FileSizeMax)
	fmt.Printf("文件大小说明：%s\n", data.FileSizeMsg)
	fmt.Printf("支持底色：%v\n", data.Color)
	fmt.Printf("分类：%s\n", data.Category)
	fmt.Printf("分类id：%s\n", data.CategoryID)
}
