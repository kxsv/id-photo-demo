// V3 证件照检测+自定义制作
//
// 接口地址：https://api.zjzapi.com/idcardv3/diy
// 请求方式：POST（application/x-www-form-urlencoded）
// 接口费用：检测并制作成功扣除0.012元，检测失败扣除0.006元；若使用"美颜"功能将多扣除0.012元。
//
// 用法（在 go 目录运行）：
//
//	go run ./cmd/v3/diy
//
// 与 v3/all 的区别：不使用规格id，自定义像素尺寸、底色、文件大小与 dpi。
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
	imagePath   = "../images/test.jpeg" // 待处理图片路径
	widthPx     = "295"                 // 像素尺寸-宽，不超2000px
	heightPx    = "413"                 // 像素尺寸-高，不超2000px
	color       = "blue"                // 底色：颜色值（#FFFFFF）或名称（blue、white、red、tint、grey、gradient、dark_blue、dark_red）
	fileSizeMin = "0"                   // 文件最小(kb)，最小5kb，0=不限
	fileSizeMax = "0"                   // 文件最大(kb)，最大300kb，0=不限
	dpi         = "300"                 // 分辨率，默认300
	beauty      = "0"                   // 美颜：0=关闭 1=开启（开启后另行计费）
	whitening   = "50"                  // 美颜-美白等级：0-100，默认50（beauty=1 时生效）
)

// v3DiyData data 部分的字段。
type v3DiyData struct {
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
		"width_px":  {widthPx},
		"height_px": {heightPx},
		"color":     {color},
		"dpi":       {dpi},
		"beauty":    {beauty},
		"whitening": {whitening},
	}
	if fileSizeMin != "0" {
		params.Set("file_size_min", fileSizeMin)
	}
	if fileSizeMax != "0" {
		params.Set("file_size_max", fileSizeMax)
	}

	result, err := zjzapi.Post("idcardv3/diy", params)
	if err != nil {
		fmt.Println("调用失败：", err)
		os.Exit(1)
	}

	if !result.IsSuccess() {
		fmt.Printf("业务失败：code=%d，msg=%s\n", result.Code, result.Msg)
		os.Exit(1)
	}

	var data v3DiyData
	if err := json.Unmarshal(result.Data, &data); err != nil {
		fmt.Println("解析 data 失败：", err)
		os.Exit(1)
	}
	fmt.Printf("制作成功，pic_id：%s\n", data.PicID)

	zjzapi.SaveImages(filepath.Join("output", "v3-diy"), data.List)
}
