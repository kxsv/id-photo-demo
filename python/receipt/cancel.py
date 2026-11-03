#!/usr/bin/env python3
"""订单取消

接口地址：https://api.zjzapi.com/receipt/cancel
请求方式：POST（application/x-www-form-urlencoded）
返回格式：JSON
接口费用：调用免费

用法（在 python 目录运行）：python receipt/cancel.py

注意：订单在"待处理"状态下才可以取消。
"""

import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))

from config import API_KEY
from zjz_api import ZjzApi, ZjzApiError

# ===== 请求参数（修改后直接运行即可）=====
ORDER_NO = ""   # 要取消的订单号（receipt/submit.py 的返回结果）


def main():
    if API_KEY in ("", "YOUR_APP_KEY"):
        print("提示：请先在 config.py 中将 API_KEY 替换为你的 API Key。", file=sys.stderr)
    if not ORDER_NO:
        print("请先在参数区填写要取消的订单号 order_no。", file=sys.stderr)
        return

    try:
        result = ZjzApi.post("receipt/cancel", {
            "key": API_KEY,
            "order_no": ORDER_NO,
        })
    except ZjzApiError as error:
        print(f"调用失败：{error}")
        return

    if not ZjzApi.is_success(result):
        print(f"业务失败：code={result.get('code')}，msg={result.get('msg')}")
        return

    print(f"订单 {ORDER_NO} 取消成功。")


if __name__ == "__main__":
    main()
