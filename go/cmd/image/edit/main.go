// 图片编辑
//
// 接口地址：https://api.zjzapi.com/image/edit
// 请求方式：POST（application/x-www-form-urlencoded）
// 接口费用：调用免费（账户余额最低需有10元才可调用）。
//
// 用法（在 go 目录运行）：
//
//	go run ./cmd/image/edit
//
// 修改图片像素（dpi）与文件大小（kb）；图片大小不超过1M，分辨率不超过1024px。
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
	dpi         = "300"                 // 图片dpi，支持72,96,150,200,300,350,500,1000；填0表示不修改
	fileSizeMin = "0"                   // 文件最小(kb)，最小5kb；填0表示不限制
	fileSizeMax = "0"                   // 文件最大(kb)，最大300kb；填0表示不限制
)

// imageEditData data 部分的字段。
type imageEditData struct {
	Image string `json:"image"`
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
		"key":   {zjzapi.APIKey},
		"image": {image},
	}
	if dpi != "0" {
		params.Set("dpi", dpi)
	}
	if fileSizeMin != "0" {
		params.Set("file_size_min", fileSizeMin)
	}
	if fileSizeMax != "0" {
		params.Set("file_size_max", fileSizeMax)
	}

	result, err := zjzapi.Post("image/edit", params)
	if err != nil {
		fmt.Println("调用失败：", err)
		os.Exit(1)
	}

	if !result.IsSuccess() {
		fmt.Printf("业务失败：code=%d，msg=%s\n", result.Code, result.Msg)
		os.Exit(1)
	}

	var data imageEditData
	if err := json.Unmarshal(result.Data, &data); err != nil {
		fmt.Println("解析 data 失败：", err)
		os.Exit(1)
	}
	if data.Image == "" {
		fmt.Println("返回数据中没有结果图地址。")
		return
	}

	saveDir := filepath.Join("output", "image-edit")
	file, err := zjzapi.DownloadImage(data.Image, saveDir, "edit.jpg")
	if err != nil {
		fmt.Printf("下载失败，结果图地址：%s\n", data.Image)
		return
	}
	fmt.Printf("图片编辑成功：%s\n", file)
}
