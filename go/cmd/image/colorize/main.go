// 黑白图片上色
//
// 接口地址：https://api.zjzapi.com/image/colorize
// 请求方式：POST（application/x-www-form-urlencoded）
// 接口费用：调用成功扣除0.012元，调用失败免费。
//
// 用法（在 go 目录运行）：
//
//	go run ./cmd/image/colorize
//
// 适用于黑白图片、老照片上色；图片大小不超过4M，分辨率不超过3072px。
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
	imagePath = "../images/test.jpeg" // 待处理图片路径（建议使用黑白照片）
)

// imageColorizeData data 部分的字段。
type imageColorizeData struct {
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

	result, err := zjzapi.Post("image/colorize", url.Values{
		"key":   {zjzapi.APIKey},
		"image": {image},
	})
	if err != nil {
		fmt.Println("调用失败：", err)
		os.Exit(1)
	}

	if !result.IsSuccess() {
		fmt.Printf("业务失败：code=%d，msg=%s\n", result.Code, result.Msg)
		os.Exit(1)
	}

	var data imageColorizeData
	if err := json.Unmarshal(result.Data, &data); err != nil {
		fmt.Println("解析 data 失败：", err)
		os.Exit(1)
	}
	if data.Image == "" {
		fmt.Println("返回数据中没有结果图地址。")
		return
	}

	saveDir := filepath.Join("output", "image-colorize")
	file, err := zjzapi.DownloadImage(data.Image, saveDir, "colorize.jpg")
	if err != nil {
		fmt.Printf("下载失败，结果图地址：%s\n", data.Image)
		return
	}
	fmt.Printf("上色成功：%s\n", file)
}
