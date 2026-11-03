#!/usr/bin/env python3
"""订单状态通知（回调接收示例）

本示例与其它示例不同：不是主动调用平台接口，而是部署在你的服务器上，
接收平台推送的订单状态通知（订单开始审核、驳回、完成、发货时触发）。

通知以 POST（application/x-www-form-urlencoded）发送到提交订单时填写的 notice_url。

本地联调方法（在 python 目录运行）：
    python receipt/notice.py
然后把 http://你的IP:8080/ 作为 receipt/submit.py 的 NOTICE_URL。

注意：本文件是独立 HTTP 服务，不依赖 zjz_api 客户端，仅用于本地联调接收平台推送。
"""

import hashlib
import json
import sys
from http.server import BaseHTTPRequestHandler, HTTPServer
from pathlib import Path
from urllib.parse import parse_qs

# 引入全局配置（python 根目录的 config.py）。
sys.path.insert(0, str(Path(__file__).resolve().parents[1]))
from config import API_KEY

PORT = 8080


def build_response(code, msg):
    """按平台约定返回处理结果：平台收到 code=0 才认为同步完成。"""
    return json.dumps({"data": {}, "code": code, "msg": msg}, ensure_ascii=False).encode("utf-8")


class NoticeHandler(BaseHTTPRequestHandler):
    def _send(self, code, msg, http_status=200):
        body = build_response(code, msg)
        self.send_response(http_status)
        self.send_header("Content-Type", "application/json; charset=utf-8")
        self.send_header("Content-Length", str(len(body)))
        self.end_headers()
        self.wfile.write(body)

    def do_POST(self):
        length = int(self.headers.get("Content-Length", 0))
        raw = self.rfile.read(length).decode("utf-8", errors="replace")
        post = {key: value[0] for key, value in parse_qs(raw).items()}

        # 验证签名：平台约定 sign = md5(API_KEY)，不一致则拒绝接收。
        expected_sign = hashlib.md5(API_KEY.encode("utf-8")).hexdigest()
        if post.get("sign", "").lower() != expected_sign:
            self._send(-1, "签名校验失败")
            return

        try:
            status = int(post.get("status", "0"))
        except ValueError:
            status = 0
        order_no = post.get("order_no", "")

        # 各状态说明：
        #   2  审核中（开始审核）
        #   3  办理完成，返回 image_receipt / image / image_print
        #   4  打印发货，返回 logistics_info
        #  -1  审核失败，返回 reject_msg（驳回理由）
        #  -2  取消完成
        if status == 2:
            # TODO: 更新订单状态为"审核中"。
            print(f"[receipt notice] 订单 {order_no} 开始审核。")
        elif status == 3:
            # TODO: 更新订单状态为"办理完成"，保存回执照地址。
            image_receipt = post.get("image_receipt", "")
            image = post.get("image", "")
            image_print = post.get("image_print", "")
            print(f"[receipt notice] 订单 {order_no} 办理完成，回执照：{image_receipt}")
        elif status == 4:
            # TODO: 更新订单状态为"已发货"，保存物流信息。
            logistics = post.get("logistics_info", "{}")
            print(f"[receipt notice] 订单 {order_no} 已发货，物流信息：{logistics}")
        elif status == -1:
            # TODO: 更新订单状态为"审核驳回"，可参考 receipt/resubmit.py 修改资料后重新提交。
            reject_msg = post.get("reject_msg", "")
            print(f"[receipt notice] 订单 {order_no} 审核驳回：{reject_msg}")
        elif status == -2:
            # TODO: 更新订单状态为"已取消"。
            print(f"[receipt notice] 订单 {order_no} 已取消。")
        else:
            self._send(-1, f"未知的订单状态：{status}")
            return

        # 处理完成后返回 code=0，告知平台通知同步完成。
        self._send(0, "同步完成")

    def do_GET(self):
        # 方便联调时做健康检查（平台只会 POST）。
        self._send(0, "ok")

    def log_message(self, fmt, *args):
        sys.stderr.write("%s - %s\n" % (self.address_string(), fmt % args))


def main():
    if API_KEY in ("", "YOUR_APP_KEY"):
        print("提示：请先在 config.py 中将 API_KEY 替换为你的 API Key。", file=sys.stderr)
    server = HTTPServer(("0.0.0.0", PORT), NoticeHandler)
    print(f"回执通知回调服务已启动：http://0.0.0.0:{PORT}/")
    print("把它作为 receipt/submit.py 的 NOTICE_URL。按 Ctrl+C 停止。")
    try:
        server.serve_forever()
    except KeyboardInterrupt:
        server.shutdown()


if __name__ == "__main__":
    main()
