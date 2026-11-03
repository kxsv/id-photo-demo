// 证件照回执检测+制作
//
// 接口地址：https://api.zjzapi.com/receipt/make
// 请求方式：POST（application/x-www-form-urlencoded）
// 接口费用：调用免费（回执提交办理完成后按单扣费）。
//
// 用法（在 go 目录运行）：
//
//	go run ./cmd/receipt/make
//
// 图片要求：JPG/JPEG/PNG，大小不超过5M，分辨率小于6000x6000；
// 办理时间 8:30-20:30（节假日正常办理），预计15分钟审核完成。
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
	itemID    = "1"                   // 回执类规格id，可通过「证件照规格列表」示例获取（go run ./cmd/item/list）
)

// receiptMakeData data 部分的字段。
type receiptMakeData struct {
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

	result, err := zjzapi.Post("receipt/make", url.Values{
		"key":     {zjzapi.APIKey},
		"item_id": {itemID},
		"image":   {image},
	})
	if err != nil {
		fmt.Println("调用失败：", err)
		os.Exit(1)
	}

	if !result.IsSuccess() {
		fmt.Printf("业务失败：code=%d，msg=%s\n", result.Code, result.Msg)
		os.Exit(1)
	}

	var data receiptMakeData
	if err := json.Unmarshal(result.Data, &data); err != nil {
		fmt.Println("解析 data 失败：", err)
		os.Exit(1)
	}
	fmt.Printf("回执照制作成功，pic_id：%s\n", data.PicID)
	fmt.Println("提示：把 pic_id 填入 receipt/submit 提交办理订单。")

	zjzapi.SaveImages(filepath.Join("output", "receipt-make"), data.List)
}
