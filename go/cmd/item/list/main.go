// 证件照规格列表示例。
//
// 接口地址：https://api.zjzapi.com/item/list
// 请求方式：POST（application/x-www-form-urlencoded）
// 接口费用：调用免费
//
// 用法（在 go 目录运行）：
//
//	go run ./cmd/item/list
//
// 输出各规格的 item_id、名称、像素尺寸与支持底色，
// 把 item_id 和对应底色填入制作类示例（如 cmd/v3/all）即可。
package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"zjz-demo/zjzapi"
)

// item 单个规格的字段。
type item struct {
	ItemID   string `json:"item_id"`
	Name     string `json:"name"`
	Color    string `json:"color"`
	WidthPx  string `json:"width_px"`
	HeightPx string `json:"height_px"`
}

// itemListData data 部分的字段（接口示例按需定义即可）。
type itemListData struct {
	List []item `json:"list"`
}

func main() {
	if zjzapi.APIKey == "" || zjzapi.APIKey == "YOUR_APP_KEY" {
		fmt.Fprintln(os.Stderr, "提示：请先在 go/zjzapi/config.go 中将 APIKey 替换为你的 API Key。")
	}

	result, err := zjzapi.Post("item/list", url.Values{
		"key": {zjzapi.APIKey},
	})
	if err != nil {
		fmt.Println("调用失败：", err)
		os.Exit(1)
	}

	// HTTP 成功不代表业务成功，需再检查返回 JSON 中的 code。
	if !result.IsSuccess() {
		fmt.Printf("业务失败：code=%d，msg=%s\n", result.Code, result.Msg)
		os.Exit(1)
	}

	var data itemListData
	if err := json.Unmarshal(result.Data, &data); err != nil {
		fmt.Println("解析 data 失败：", err)
		os.Exit(1)
	}
	if len(data.List) == 0 {
		fmt.Println("返回数据中没有规格列表。")
		return
	}

	// 打印规格清单，便于挑选 item_id 与底色。
	fmt.Println("item_id  名称            像素          支持底色")
	for _, it := range data.List {
		fmt.Printf("%-8s %-14s %-12s %s\n", it.ItemID, it.Name, it.WidthPx+"x"+it.HeightPx, it.Color)
	}
	fmt.Printf("共 %d 个规格。\n", len(data.List))
}
