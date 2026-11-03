// V5 获取证件照（无水印）
//
// 接口地址：https://api.zjzapi.com/idcardv5/get
// 请求方式：POST（application/x-www-form-urlencoded）
// 接口费用：获取成功扣除0.6元，重复获取不扣费。
//
// 用法（在 go 目录运行）：
//
//	go run ./cmd/v5/get
//
// 注意：pic_id 来自制作类接口（v5/make、v5/bg、v5/diy、v5/pose 等）的返回结果，
// 有效期50分钟；若换装过，请同时填写换装时使用的 clothe。
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
	picID       = ""    // 制作接口返回的 pic_id（有效期50分钟），先运行 go run ./cmd/v5/make 获取
	clothe      = ""    // 换装时使用的服装key，未换装留空
	giftHd      = "0"   // 1=赠送一张三寸高清照（仅规格 width_px<500 支持；更换底色、自定义制作暂不支持）
	giftHdColor = ""    // 高清照底色，留空默认返回制作结果第一个底色
)

// v5GetData data 部分的字段。
type v5GetData struct {
	List   map[string]string `json:"list"`
	GiftHD string            `json:"gift_hd"`
}

func main() {
	if zjzapi.APIKey == "" || zjzapi.APIKey == "YOUR_APP_KEY" {
		fmt.Fprintln(os.Stderr, "提示：请先在 go/zjzapi/config.go 中将 APIKey 替换为你的 API Key。")
	}
	if picID == "" {
		fmt.Fprintln(os.Stderr, "请先在参数区填写 pic_id（运行 go run ./cmd/v5/make 获取）。")
		os.Exit(1)
	}

	params := url.Values{
		"key":      {zjzapi.APIKey},
		"pic_id":   {picID},
		"gift_hd":  {giftHd},
	}
	if clothe != "" {
		params.Set("clothe", clothe)
	}
	if giftHdColor != "" {
		params.Set("gift_hd_color", giftHdColor)
	}

	result, err := zjzapi.Post("idcardv5/get", params)
	if err != nil {
		fmt.Println("调用失败：", err)
		os.Exit(1)
	}

	if !result.IsSuccess() {
		fmt.Printf("业务失败：code=%d，msg=%s\n", result.Code, result.Msg)
		os.Exit(1)
	}

	var data v5GetData
	if err := json.Unmarshal(result.Data, &data); err != nil {
		fmt.Println("解析 data 失败：", err)
		os.Exit(1)
	}
	fmt.Println("获取无水印证件照成功。")

	saveDir := filepath.Join("output", "v5-get")
	zjzapi.SaveImages(saveDir, data.List)

	// 赠送的三寸高清照（开启 gift_hd 时才有）。
	if data.GiftHD != "" {
		file, err := zjzapi.DownloadImage(data.GiftHD, saveDir, "gift_hd.jpg")
		if err != nil {
			fmt.Printf("高清照下载失败：%v（结果图地址：%s）\n", err, data.GiftHD)
		} else {
			fmt.Printf("高清照：%s\n", file)
		}
	}
}
