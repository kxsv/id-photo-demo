#!/usr/bin/env python3
"""订单审核驳回后再次提交

接口地址：https://api.zjzapi.com/receipt/submit（与首次提交相同接口，额外携带 order_no）
请求方式：POST（application/x-www-form-urlencoded）
返回格式：JSON
接口费用：订单办理完成扣费，7元每单；若选择打印邮寄服务再增加7元。

用法（在 python 目录运行）：python receipt/resubmit.py

注意：驳回通知中的 reject_msg 会说明驳回原因，按要求补充/修改 param（如身份证照片 image_idcard）后，
携带原订单号重新提交。
"""

import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))

from config import API_KEY
from zjz_api import ZjzApi, ZjzApiError

# ===== 请求参数（修改后直接运行即可）=====
PIC_ID = ""       # 制作接口返回的 pic_id（有效期50分钟），先运行 python receipt/make.py 获取
ORDER_NO = ""     # 上次提交返回的订单号（receipt/submit.py 的返回结果）
NOTICE_URL = ""   # 回调地址：你的服务器接收订单状态通知的 URL
PARAM = ""        # 附加参数（json字符串），按驳回要求补充，如广东驾驶证驳回后补身份证照片：{"text_mobile":"13800138000","text_city":"深圳","image_idcard":"iVBORw0KGgo....（身份证照片base64）"}


def main():
    if API_KEY in ("", "YOUR_APP_KEY"):
        print("提示：请先在 config.py 中将 API_KEY 替换为你的 API Key。", file=sys.stderr)
    if not PIC_ID:
        print("请先在参数区填写 pic_id（运行 python receipt/make.py 获取）。", file=sys.stderr)
        return
    if not ORDER_NO:
        print("请先在参数区填写 order_no（上次提交的订单号）。", file=sys.stderr)
        return
    if not NOTICE_URL:
        print("请先在参数区填写 notice_url（你的订单状态回调地址）。", file=sys.stderr)
        return

    try:
        params = {
            "key": API_KEY,
            "pic_id": PIC_ID,
            "order_no": ORDER_NO,
            "notice_url": NOTICE_URL,
        }
        if PARAM:
            params["param"] = PARAM
        result = ZjzApi.post("receipt/submit", params)
    except ZjzApiError as error:
        print(f"调用失败：{error}")
        return

    if not ZjzApi.is_success(result):
        print(f"业务失败：code={result.get('code')}，msg={result.get('msg')}")
        return

    data = result.get("data") or {}
    print(f"再次提交成功，订单号：{data.get('order_no', '')}")


if __name__ == "__main__":
    main()
