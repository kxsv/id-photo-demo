// 证件照合规性检测
//
// 接口地址：https://api.zjzapi.com/idphoto/check
// 请求方式：POST（application/x-www-form-urlencoded）
// 接口费用：检测成功扣除0.012元。
//
// 用法（在 go 目录运行）：
//
//	go run ./cmd/check/check
//
// 提供规格id（item_id）即可按该规格检测；也可不传 item_id，
// 自定义像素、dpi、文件大小、底色、格式以及各项人像构图/姿态/画质/遮挡/着装检测项，
// 完整参数列表见官网《证件照合规性检测》文档。
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
	imagePath = "../images/test.jpeg" // 待检测图片路径（不超过4M，分辨率不超过3072px）
	itemID    = "1"                   // 规格id，提供后可省略 width_px/height_px/dpi 等规格参数；填 0 表示使用下方自定义规格参数

	// 方式二（item_id 为 0 时生效）的自定义检测规格，按需增删。
	customWidthPx     = "295"             // 证件照宽（px）
	customHeightPx    = "413"             // 证件照高（px）
	customDpi         = "300"             // 图片分辨率，常见300、350
	customFileFormat  = "jpg"             // 照片格式：jpg、jpeg、png
	customBgColors    = "white,blue"      // 允许的底色，逗号分隔
	customSolidBg     = "1"               // 纯色背景检测：0/1
	faceWidthRatioMin = "0.4"             // 人像构图-脸宽占比下限（0.01-1.0）
	faceWidthRatioMax = "0.6"             // 人像构图-脸宽占比上限（0.01-1.0）
	shoulderHorizontal = "5"              // 姿态-双肩水平角度（推荐值见文档）
	facePoseYaw        = "15"             // 姿态-左右转头角度（推荐值见文档）
	faceClarity        = "0.5"            // 画质-清晰度（推荐值见文档）
	faceBrightness     = "0.3"            // 画质-亮度（推荐值见文档）
	faceMask           = "1"              // 遮挡-口罩检测：0=不检测 1=检测
	glasses            = "1"              // 遮挡-眼镜检测：0=不检测 1=检测
)

// checkDetail 单个检测项的 detail 字段。
type checkDetail struct {
	Actual interface{} `json:"actual"`
	Unit   string      `json:"unit"`
}

// checkItem 单个检测项。
type checkItem struct {
	Key    string                 `json:"key"`
	Passed zjzapi.FlexString      `json:"passed"`
	Msg    string                 `json:"msg"`
	Detail checkDetail            `json:"detail"`
}

// checkData data 部分的字段。
type checkData struct {
	// passed_status 可能为数字也可能为字符串，用 FlexString 兼容，避免字符串返回时反序列化失败。
	PassedStatus zjzapi.FlexString `json:"passed_status"`
	PassedList   []checkItem       `json:"passed_list"`
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

	if itemID != "0" {
		// 方式一：按规格id检测（推荐，规格要求自动带上）。
		params.Set("item_id", itemID)
	} else {
		// 方式二：自定义检测规格（不传 item_id 时生效），按需增删。
		params.Set("width_px", customWidthPx)
		params.Set("height_px", customHeightPx)
		params.Set("dpi", customDpi)
		params.Set("file_format", customFileFormat)
		params.Set("bg_colors", customBgColors)
		params.Set("solid_bg", customSolidBg)
		params.Set("face_width_ratio_min", faceWidthRatioMin)
		params.Set("face_width_ratio_max", faceWidthRatioMax)
		params.Set("shoulder_horizontal", shoulderHorizontal)
		params.Set("face_pose_yaw", facePoseYaw)
		params.Set("face_clarity", faceClarity)
		params.Set("face_brightness", faceBrightness)
		params.Set("face_mask", faceMask)
		params.Set("glasses", glasses)
	}

	result, err := zjzapi.Post("idphoto/check", params)
	if err != nil {
		fmt.Println("调用失败：", err)
		os.Exit(1)
	}

	if !result.IsSuccess() {
		fmt.Printf("业务失败：code=%d，msg=%s\n", result.Code, result.Msg)
		os.Exit(1)
	}

	var data checkData
	if err := json.Unmarshal(result.Data, &data); err != nil {
		fmt.Println("解析 data 失败：", err)
		os.Exit(1)
	}
	if data.PassedStatus == "1" {
		fmt.Println("总检测结果：通过")
	} else {
		fmt.Println("总检测结果：不通过")
	}

	if len(data.PassedList) == 0 {
		fmt.Println("返回数据中没有检测项列表。")
		return
	}

	fmt.Println("检测项        结果    提示        检测值")
	for _, it := range data.PassedList {
		actual := ""
		if it.Detail.Actual != nil {
			actual = fmt.Sprintf("%v", it.Detail.Actual)
			if it.Detail.Unit != "" {
				actual += it.Detail.Unit
			}
		}
		passed := "不通过"
		if it.Passed == "1" {
			passed = "通过"
		}
		key := it.Key
		if len(key) > 14 {
			key = key[:14]
		}
		fmt.Printf("%-14s %-6s %-10s %s\n", key, passed, it.Msg, actual)
	}
}
