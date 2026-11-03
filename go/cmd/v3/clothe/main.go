// V3 证件照换装
//
// 接口地址：https://api.zjzapi.com/idcardv3/clothe
// 请求方式：POST（application/x-www-form-urlencoded）
// 接口费用：换装成功扣除0.036元。
//
// 用法（在 go 目录运行）：
//
//	go run ./cmd/v3/clothe
//
// 注意：换装基于制作接口返回的 pic_id（有效期50分钟），
// 请先把 v3/all 运行得到的 pic_id 填到下方参数区。
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
	picID  = "" // 制作接口返回的 pic_id（有效期50分钟），先运行 go run ./cmd/v3/all 获取
	clothe = "" // 服装key，见官网《服装列表》文档
)

// v3ClotheData data 部分的字段。
type v3ClotheData struct {
	List map[string]string `json:"list"`
}

func main() {
	if zjzapi.APIKey == "" || zjzapi.APIKey == "YOUR_APP_KEY" {
		fmt.Fprintln(os.Stderr, "提示：请先在 go/zjzapi/config.go 中将 APIKey 替换为你的 API Key。")
	}
	if picID == "" || clothe == "" {
		fmt.Fprintln(os.Stderr, "请先在参数区填写 pic_id（运行 go run ./cmd/v3/all 获取）与服装key clothe。")
		os.Exit(1)
	}

	result, err := zjzapi.Post("idcardv3/clothe", url.Values{
		"key":    {zjzapi.APIKey},
		"pic_id": {picID},
		"clothe": {clothe},
	})
	if err != nil {
		fmt.Println("调用失败：", err)
		os.Exit(1)
	}

	if !result.IsSuccess() {
		fmt.Printf("业务失败：code=%d，msg=%s\n", result.Code, result.Msg)
		os.Exit(1)
	}

	var data v3ClotheData
	if err := json.Unmarshal(result.Data, &data); err != nil {
		fmt.Println("解析 data 失败：", err)
		os.Exit(1)
	}
	fmt.Println("换装成功。")

	zjzapi.SaveImages(filepath.Join("output", "v3-clothe"), data.List)
}
