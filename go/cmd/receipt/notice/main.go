// 订单状态通知（回调接收示例）
//
// 本示例与其它示例不同：不是主动调用平台接口，而是部署在你的服务器上，
// 接收平台推送的订单状态通知（订单开始审核、驳回、完成、发货时触发）。
//
// 通知以 POST（application/x-www-form-urlencoded）发送到提交订单时填写的 notice_url。
//
// 本地联调方法（在 go 目录运行）：
//
//	go run ./cmd/receipt/notice
//
// 然后把 http://你的IP:8080/ 作为 go run ./cmd/receipt/submit 的 noticeURL。
//
// 注意：本示例是独立 HTTP 服务，使用 zjzapi.APIKey 配置，接收推送并校验签名（sign = md5(APIKey)）。
package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"zjz-demo/zjzapi"
)

// noticePort 回调服务监听端口。
const noticePort = 8080

// noticeResp 与平台约定的返回结构（data 始终为空对象）。
type noticeResp struct {
	Data interface{} `json:"data"`
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
}

func writeJSON(w http.ResponseWriter, code int, msg string) {
	resp := noticeResp{Data: struct{}{}, Code: code, Msg: msg}
	body, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(body)
}

func noticeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// 方便联调时做健康检查（平台只会 POST）。
		writeJSON(w, 0, "ok")
		return
	}

	raw, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSON(w, -1, "读取请求体失败")
		return
	}
	post, _ := url.ParseQuery(string(raw))

	// 验证签名：平台约定 sign = md5(APIKey)，不一致则拒绝接收。
	sum := md5.Sum([]byte(zjzapi.APIKey))
	expected := hex.EncodeToString(sum[:])
	if strings.ToLower(post.Get("sign")) != expected {
		writeJSON(w, -1, "签名校验失败")
		return
	}

	status, _ := strconv.Atoi(post.Get("status"))
	orderNo := post.Get("order_no")

	// 各状态说明：
	//   2  审核中（开始审核）
	//   3  办理完成，返回 image_receipt / image / image_print
	//   4  打印发货，返回 logistics_info
	//  -1  审核失败，返回 reject_msg（驳回理由）
	//  -2  取消完成
	switch status {
	case 2:
		// TODO: 更新订单状态为"审核中"。
		fmt.Printf("[receipt notice] 订单 %s 开始审核。\n", orderNo)
	case 3:
		// TODO: 更新订单状态为"办理完成"，保存回执照地址。
		fmt.Printf("[receipt notice] 订单 %s 办理完成，回执照：%s\n", orderNo, post.Get("image_receipt"))
	case 4:
		// TODO: 更新订单状态为"已发货"，保存物流信息。
		fmt.Printf("[receipt notice] 订单 %s 已发货，物流信息：%s\n", orderNo, post.Get("logistics_info"))
	case -1:
		// TODO: 更新订单状态为"审核驳回"，可参考 receipt/resubmit 修改资料后重新提交。
		fmt.Printf("[receipt notice] 订单 %s 审核驳回：%s\n", orderNo, post.Get("reject_msg"))
	case -2:
		// TODO: 更新订单状态为"已取消"。
		fmt.Printf("[receipt notice] 订单 %s 已取消。\n", orderNo)
	default:
		writeJSON(w, -1, fmt.Sprintf("未知的订单状态：%d", status))
		return
	}

	// 处理完成后返回 code=0，告知平台通知同步完成。
	writeJSON(w, 0, "同步完成")
}

func main() {
	if zjzapi.APIKey == "" || zjzapi.APIKey == "YOUR_APP_KEY" {
		fmt.Fprintln(os.Stderr, "提示：请先在 go/zjzapi/config.go 中将 APIKey 替换为你的 API Key。")
	}

	addr := fmt.Sprintf("0.0.0.0:%d", noticePort)
	http.HandleFunc("/", noticeHandler)
	fmt.Printf("回执通知回调服务已启动：http://%s/\n", addr)
	fmt.Println("把它作为 go run ./cmd/receipt/submit 的 noticeURL。按 Ctrl+C 停止。")
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Println("服务启动失败：", err)
		os.Exit(1)
	}
}
