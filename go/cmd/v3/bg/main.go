// V3 图片检测+更换底色
//
// 接口地址：https://api.zjzapi.com/idcardv3/bg
// 请求方式：POST（application/x-www-form-urlencoded）
// 接口费用：检测并制作成功扣除0.012元，检测失败扣除0.006元；若使用"美颜"功能将多扣除0.012元。
//
// 用法（在 go 目录运行）：
//
//	go run ./cmd/v3/bg
//
// 与 v3/all 的区别：不按规格重排裁剪，仅对已有照片检测后更换底色；
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
	beauty    = "0"                   // 美颜：0=关闭 1=开启（开启后另行计费）
	whitening = "50"                  // 美颜-美白等级：0-100，默认50（beauty=1 时生效）
)

// v3BgData data 部分的字段。
type v3BgData struct {
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
		"beauty":    {beauty},
		"whitening": {whitening},
	}
	// color 为空时不传，由服务端返回全部底色。
	if color != "" {
		params.Set("color", color)
	}

	result, err := zjzapi.Post("idcardv3/bg", params)
	if err != nil {
		fmt.Println("调用失败：", err)
		os.Exit(1)
	}

	if !result.IsSuccess() {
		fmt.Printf("业务失败：code=%d，msg=%s\n", result.Code, result.Msg)
		os.Exit(1)
	}

	var data v3BgData
	if err := json.Unmarshal(result.Data, &data); err != nil {
		fmt.Println("解析 data 失败：", err)
		os.Exit(1)
	}
	fmt.Printf("更换底色成功，pic_id：%s\n", data.PicID)

	// 把各底色结果图下载到 output/v3-bg/ 目录。
	zjzapi.SaveImages(filepath.Join("output", "v3-bg"), data.List)
}
