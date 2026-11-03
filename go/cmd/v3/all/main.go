// V3 证件照检测+制作示例。
//
// 接口地址：https://api.zjzapi.com/idcardv3/all
// 请求方式：POST（application/x-www-form-urlencoded）
// 接口费用：制作成功按次扣费；检测不通过扣除检测费用；beauty=1 时美颜另行计费。
//
// 用法（在 go 目录运行）：
//
//	go run ./cmd/v3/all
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
	imagePath = "../images/test.jpeg" // 待处理图片路径（相对 go 目录）
	itemID    = "1"                // 规格id，可通过「证件照规格列表」示例获取
	colors    = "white"              // 底色，如 white/blue/red；留空返回该规格支持的全部底色
	beauty    = "0"                  // 美颜：0=关闭 1=开启（开启后另行计费）
)

// v3AllData data 部分的字段（接口示例按需定义即可）。
type v3AllData struct {
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

	// HTTP 成功不代表业务成功，需再检查返回 JSON 中的 code。
	if !result.IsSuccess() {
		fmt.Printf("业务失败：code=%d，msg=%s\n", result.Code, result.Msg)
		os.Exit(1)
	}

	var data v3AllData
	if err := json.Unmarshal(result.Data, &data); err != nil {
		fmt.Println("解析 data 失败：", err)
		os.Exit(1)
	}
	fmt.Printf("制作成功，pic_id：%s\n", data.PicID)

	// 把各底色结果图下载到 output/v3-all/ 目录。
	zjzapi.SaveImages(filepath.Join("output", "v3-all"), data.List)
}
