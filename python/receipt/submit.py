#!/usr/bin/env python3
"""证件照回执提交

接口地址：https://api.zjzapi.com/receipt/submit
请求方式：POST（application/x-www-form-urlencoded）
返回格式：JSON
接口费用：订单办理完成扣费，7元每单；若选择打印邮寄服务再增加7元。

用法（在 python 目录运行）：python receipt/submit.py

注意：pic_id 来自回执制作接口的返回结果（有效期50分钟）；
提交后订单状态变化会通过 notice_url 回调通知。
"""

import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))

from config import API_KEY
from zjz_api import ZjzApi, ZjzApiError

# ===== 请求参数（修改后直接运行即可）=====
PIC_ID = ""       # 制作接口返回的 pic_id（有效期50分钟），先运行 python receipt/make.py 获取
NOTICE_URL = ""   # 回调地址：你的服务器接收订单状态通知的 URL
PARAM = ""        # 附加参数（json字符串），具体看每个回执规格的要求，如广东驾驶证：{"text_mobile":"13800138000","text_city":"深圳","text_name":"姓名","text_idcard":"身份证号码"}
IS_PRINT = 0      # 是否打印邮寄：1=是 0=否
PRINT_PARAM = ""  # 打印邮寄（is_print=1 时必填），如：{"name":"张三","phone":"13800138000","province":"广东","city":"深圳","district":"南山区","address":"xx路xx小区"}


def main():
    if API_KEY in ("", "YOUR_APP_KEY"):
        print("提示：请先在 config.py 中将 API_KEY 替换为你的 API Key。", file=sys.stderr)
    if not PIC_ID:
        print("请先在参数区填写 pic_id（运行 python receipt/make.py 获取）。", file=sys.stderr)
        return
    if not NOTICE_URL:
        print("请先在参数区填写 notice_url（你的订单状态回调地址）。", file=sys.stderr)
        return

    try:
        params = {
            "key": API_KEY,
            "pic_id": PIC_ID,
            "notice_url": NOTICE_URL,
        }
        if PARAM:
            params["param"] = PARAM
        if IS_PRINT == 1:
            params["is_print"] = 1
            if PRINT_PARAM:
                params["print_param"] = PRINT_PARAM
        result = ZjzApi.post("receipt/submit", params)
    except ZjzApiError as error:
        print(f"调用失败：{error}")
        return

    if not ZjzApi.is_success(result):
        print(f"业务失败：code={result.get('code')}，msg={result.get('msg')}")
        return

    data = result.get("data") or {}
    print(f"提交成功，订单号：{data.get('order_no', '')}")
    print("提示：订单状态变化将通过 notice_url 回调通知；如需取消订单见 receipt/cancel.py。")


if __name__ == "__main__":
    main()
