// V5 证件照检测+制作（有水印）
//
// 接口地址：https://api.zjzapi.com/idcardv5/make
// 请求方式：POST（application/x-www-form-urlencoded）
// 接口费用：调用免费，结果图带水印；可使用返回的 pic_id 调用 v5/get 获取无水印证件照。
//
// 用法（在 go 目录运行）：
//
//	go run ./cmd/v5/make
//
// 输出 pic_id 后，可把它填入 v5/get 获取无水印图（获取成功扣除0.6元，重复获取不扣费）。
package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"zjz-demo/zjzapi"
)

// ===== 请求参数（修改后直接运行即可）=====
var (
	imagePath = "../images/test.jpeg" // 待处理图片路径
	itemID    = "1"                   // 规格id，可通过「证件照规格列表」示例获取（go run ./cmd/item/list）
	colors    = "white"               // 底色，如 white/blue/red；留空返回该规格支持的全部底色
	enhance   = "1"                   // 人像增强：默认已支持，不需要可传0
	beauty    = "0"                   // 美颜：默认已支持，不需要可传0
	whitening = "50"                  // 美颜-美白等级：0-100，默认50（beauty=1 时生效）
)

// v5MakeData data 部分的字段。
type v5MakeData struct {
	PicID zjzapi.FlexString `json:"pic_id"`
	List  map[string]string `json:"list"`
}

func main() {
	if zjzapi.APIKey == "" || zjzapi.APIKey == "YOUR_APP_KEY" {
		fmt.Fprintln(os.Stderr, "提示：请先在 go/zjzapi/config.go 中将 APIKey 替换为你的 API Key。")
	}

	image, err := zjzapi.ImageToBase64(imagePath)
	if err != nil {
		fmt.Println("读取图片失败：", err)
		os.Exit(1)
	}

	params := url.Values{
		"key":       {zjzapi.APIKey},
		"item_id":   {itemID},
		"image":     {image},
		"enhance":   {enhance},
		"beauty":    {beauty},
		"whitening": {whitening},
	}
	if colors != "" {
		params.Set("colors", colors)
	}

	result, err := zjzapi.Post("idcardv5/make", params)
	if err != nil {
		fmt.Println("调用失败：", err)
		os.Exit(1)
	}

	if !result.IsSuccess() {
		fmt.Printf("业务失败：code=%d，msg=%s\n", result.Code, result.Msg)
		os.Exit(1)
	}

	var data v5MakeData
	if err := json.Unmarshal(result.Data, &data); err != nil {
		fmt.Println("解析 data 失败：", err)
		os.Exit(1)
	}
	fmt.Printf("制作成功（有水印），pic_id：%s\n", data.PicID)
	fmt.Println("提示：把 pic_id 填入 v5/get 可获取无水印证件照。")

	zjzapi.SaveImages(filepath.Join("output", "v5-make"), data.List)
}
