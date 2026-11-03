// V5 图片检测+更换底色（有水印）
//
// 接口地址：https://api.zjzapi.com/idcardv5/bg
// 请求方式：POST（application/x-www-form-urlencoded）
// 接口费用：调用免费，结果图带水印；可使用返回的 pic_id 调用 v5/get 获取无水印证件照。
//
// 用法（在 go 目录运行）：
//
//	go run ./cmd/v5/bg
//
// 与 v5/make 的区别：不按规格重排裁剪，仅对已有照片检测后更换底色；
// 图片大小不超过4M，分辨率不超过3072px。
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
	color     = "blue"                // 底色：颜色值（#FFFFFF）或名称（blue、white、red、tint、grey、gradient、dark_blue、dark_red），留空返回所有底色
	enhance   = "1"                   // 人像增强：默认已支持，不需要可传0
	beauty    = "1"                   // 美颜：默认已支持，不需要可传0
	whitening = "60"                  // 美颜-美白等级：0-100，默认60（beauty=1 时生效）
)

// v5BgData data 部分的字段。
type v5BgData struct {
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
		"image":     {image},
		"enhance":   {enhance},
		"beauty":    {beauty},
		"whitening": {whitening},
	}
	if color != "" {
		params.Set("color", color)
	}

	result, err := zjzapi.Post("idcardv5/bg", params)
	if err != nil {
		fmt.Println("调用失败：", err)
		os.Exit(1)
	}

	if !result.IsSuccess() {
		fmt.Printf("业务失败：code=%d，msg=%s\n", result.Code, result.Msg)
		os.Exit(1)
	}

	var data v5BgData
	if err := json.Unmarshal(result.Data, &data); err != nil {
		fmt.Println("解析 data 失败：", err)
		os.Exit(1)
	}
	fmt.Printf("更换底色成功（有水印），pic_id：%s\n", data.PicID)
	fmt.Println("提示：把 pic_id 填入 v5/get 可获取无水印证件照。")

	zjzapi.SaveImages(filepath.Join("output", "v5-bg"), data.List)
}
