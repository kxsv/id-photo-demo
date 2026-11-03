// 证件照回执提交
//
// 接口地址：https://api.zjzapi.com/receipt/submit
// 请求方式：POST（application/x-www-form-urlencoded）
// 接口费用：订单办理完成扣费，7元每单；若选择打印邮寄服务再增加7元。
//
// 用法（在 go 目录运行）：
//
//	go run ./cmd/receipt/submit
//
// 注意：pic_id 来自回执制作接口的返回结果（有效期50分钟）；
// 提交后订单状态变化会通过 notice_url 回调通知。
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
	picID       = "" // 制作接口返回的 pic_id（有效期50分钟），先运行 go run ./cmd/receipt/make 获取
	noticeURL   = "" // 回调地址：你的服务器接收订单状态通知的 URL
	param       = "" // 附加参数（json字符串），具体看每个回执规格的要求，如广东驾驶证：{"text_mobile":"13800138000","text_city":"深圳","text_name":"姓名","text_idcard":"身份证号码"}
	isPrint     = "0" // 是否打印邮寄：1=是 0=否
	printParam  = ""  // 打印邮寄（is_print=1 时必填），如：{"name":"张三","phone":"13800138000","province":"广东","city":"深圳","district":"南山区","address":"xx路xx小区"}
)

// receiptSubmitData data 部分的字段。
type receiptSubmitData struct {
	OrderNo zjzapi.FlexString `json:"order_no"`
}

func main() {
	if zjzapi.APIKey == "" || zjzapi.APIKey == "YOUR_APP_KEY" {
		fmt.Fprintln(os.Stderr, "提示：请先在 go/zjzapi/config.go 中将 APIKey 替换为你的 API Key。")
	}
	if picID == "" {
		fmt.Fprintln(os.Stderr, "请先在参数区填写 pic_id（运行 go run ./cmd/receipt/make 获取）。")
		os.Exit(1)
	}
	if noticeURL == "" {
		fmt.Fprintln(os.Stderr, "请先在参数区填写 notice_url（你的订单状态回调地址）。")
		os.Exit(1)
	}

	params := url.Values{
		"key":        {zjzapi.APIKey},
		"pic_id":     {picID},
		"notice_url": {noticeURL},
	}
	if param != "" {
		params.Set("param", param)
	}
	if isPrint == "1" {
		params.Set("is_print", "1")
		if printParam != "" {
			params.Set("print_param", printParam)
		}
	}

	result, err := zjzapi.Post("receipt/submit", params)
	if err != nil {
		fmt.Println("调用失败：", err)
		os.Exit(1)
	}

	if !result.IsSuccess() {
		fmt.Printf("业务失败：code=%d，msg=%s\n", result.Code, result.Msg)
		os.Exit(1)
	}

	var data receiptSubmitData
	if err := json.Unmarshal(result.Data, &data); err != nil {
		fmt.Println("解析 data 失败：", err)
		os.Exit(1)
	}
	fmt.Printf("提交成功，订单号：%s\n", data.OrderNo)
	fmt.Println("提示：订单状态变化将通过 notice_url 回调通知；如需取消订单见 receipt/cancel。")
}
