// 订单取消
//
// 接口地址：https://api.zjzapi.com/receipt/cancel
// 请求方式：POST（application/x-www-form-urlencoded）
// 接口费用：调用免费
//
// 用法（在 go 目录运行）：
//
//	go run ./cmd/receipt/cancel
//
// 注意：订单在"待处理"状态下才可以取消。
package main

import (
	"fmt"
	"net/url"
	"os"

	"zjz-demo/zjzapi"
)

// ===== 请求参数（修改后直接运行即可）=====
var (
	orderNo = "" // 要取消的订单号（receipt/submit 的返回结果）
)

func main() {
	if zjzapi.APIKey == "" || zjzapi.APIKey == "YOUR_APP_KEY" {
		fmt.Fprintln(os.Stderr, "提示：请先在 go/zjzapi/config.go 中将 APIKey 替换为你的 API Key。")
	}
	if orderNo == "" {
		fmt.Fprintln(os.Stderr, "请先在参数区填写要取消的订单号 order_no。")
		os.Exit(1)
	}

	result, err := zjzapi.Post("receipt/cancel", url.Values{
		"key":      {zjzapi.APIKey},
		"order_no": {orderNo},
	})
	if err != nil {
		fmt.Println("调用失败：", err)
		os.Exit(1)
	}

	if !result.IsSuccess() {
		fmt.Printf("业务失败：code=%d，msg=%s\n", result.Code, result.Msg)
		os.Exit(1)
	}

	fmt.Printf("订单 %s 取消成功。\n", orderNo)
}
