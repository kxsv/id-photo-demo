// 订单审核驳回后再次提交
//
// 接口地址：https://api.zjzapi.com/receipt/submit（与首次提交相同接口，额外携带 order_no）
// 请求方式：POST（application/x-www-form-urlencoded）
// 接口费用：订单办理完成扣费，7元每单；若选择打印邮寄服务再增加7元。
//
// 用法（在 go 目录运行）：
//
//	go run ./cmd/receipt/resubmit
//
// 注意：驳回通知中的 reject_msg 会说明驳回原因，按要求补充/修改 param（如身份证照片 image_idcard）后，
// 携带原订单号重新提交。
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
	picID     = "" // 制作接口返回的 pic_id（有效期50分钟），先运行 go run ./cmd/receipt/make 获取
	orderNo   = "" // 上次提交返回的订单号（receipt/submit 的返回结果）
	noticeURL = "" // 回调地址：你的服务器接收订单状态通知的 URL
	param     = "" // 附加参数（json字符串），按驳回要求补充，如广东驾驶证驳回后补身份证照片：{"text_mobile":"13800138000","text_city":"深圳","image_idcard":"iVBORw0KGgo....（身份证照片base64）"}
)

// receiptResubmitData data 部分的字段。
type receiptResubmitData struct {
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
	if orderNo == "" {
		fmt.Fprintln(os.Stderr, "请先在参数区填写 order_no（上次提交的订单号）。")
		os.Exit(1)
	}
	if noticeURL == "" {
		fmt.Fprintln(os.Stderr, "请先在参数区填写 notice_url（你的订单状态回调地址）。")
		os.Exit(1)
	}

	params := url.Values{
		"key":        {zjzapi.APIKey},
		"pic_id":     {picID},
		"order_no":   {orderNo},
		"notice_url": {noticeURL},
	}
	if param != "" {
		params.Set("param", param)
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

	var data receiptResubmitData
	if err := json.Unmarshal(result.Data, &data); err != nil {
		fmt.Println("解析 data 失败：", err)
		os.Exit(1)
	}
	fmt.Printf("再次提交成功，订单号：%s\n", data.OrderNo)
}
