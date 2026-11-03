// V3 职业形象照检测+制作
//
// 接口地址：https://api.zjzapi.com/idcardv3/all（与证件照制作相同，通过职业照规格 + 专属底色区分）
// 请求方式：POST（application/x-www-form-urlencoded）
// 接口费用：检测并制作成功扣除0.012元，检测失败扣除0.006元；若使用"美颜"功能将多扣除0.012元。
//
// 用法（在 go 目录运行）：
//
//	go run ./cmd/v3/pose
//
// 与 v3/all 的区别：item_id 使用职业形象照类规格（规格列表中的"职业形象照"分类），
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
	beauty    = "0"                   // 美颜：0=关闭 1=开启（开启后另行计费）
)

// v3PoseData data 部分的字段。
type v3PoseData struct {
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
		"key":     {zjzapi.APIKey},
		"item_id": {itemID},
		"image":   {image},
		"beauty":  {beauty},
	}
	// colors 为空时不传，由服务端返回全部底色。
	if colors != "" {
		params.Set("colors", colors)
	}

	result, err := zjzapi.Post("idcardv3/all", params)
	if err != nil {
		fmt.Println("调用失败：", err)
		os.Exit(1)
	}

	if !result.IsSuccess() {
		fmt.Printf("业务失败：code=%d，msg=%s\n", result.Code, result.Msg)
		os.Exit(1)
	}

	var data v3PoseData
	if err := json.Unmarshal(result.Data, &data); err != nil {
		fmt.Println("解析 data 失败：", err)
		os.Exit(1)
	}
	fmt.Printf("制作成功，pic_id：%s\n", data.PicID)

	zjzapi.SaveImages(filepath.Join("output", "v3-pose"), data.List)
}
