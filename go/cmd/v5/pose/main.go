// V5 职业形象照检测+制作（有水印）
//
// 接口地址：https://api.zjzapi.com/idcardv5/make（与证件照制作相同，通过职业照规格 + 专属底色区分）
// 请求方式：POST（application/x-www-form-urlencoded）
// 接口费用：调用免费，结果图带水印；可使用返回的 pic_id 调用 v5/get 获取无水印证件照。
//
// 用法（在 go 目录运行）：
//
//	go run ./cmd/v5/pose
//
// 与 v5/make 的区别：item_id 使用职业形象照类规格（规格列表中的"职业形象照"分类），
// 底色为职业照专属：pose_black、pose_blue、pose_darkgray、pose_gray、pose_white。
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
	imagePath = "../images/pose.jpg" // 待处理图片路径
	itemID    = "1"                   // 职业形象照规格id，可通过「证件照规格列表」示例获取（go run ./cmd/item/list）
	colors    = "pose_blue"           // 底色：pose_black、pose_blue、pose_darkgray、pose_gray、pose_white，留空返回全部
	enhance   = "1"                   // 人像增强：默认已支持，不需要可传0
	beauty    = "1"                   // 美颜：默认已支持，不需要可传0
	whitening = "50"                  // 美颜-美白等级：0-100，默认50（beauty=1 时生效）
)

// v5PoseData data 部分的字段。
type v5PoseData struct {
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

	var data v5PoseData
	if err := json.Unmarshal(result.Data, &data); err != nil {
		fmt.Println("解析 data 失败：", err)
		os.Exit(1)
	}
	fmt.Printf("制作成功（有水印），pic_id：%s\n", data.PicID)
	fmt.Println("提示：把 pic_id 填入 v5/get 可获取无水印证件照。")

	zjzapi.SaveImages(filepath.Join("output", "v5-pose"), data.List)
}
