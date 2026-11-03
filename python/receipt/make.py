#!/usr/bin/env python3
"""证件照回执检测+制作

接口地址：https://api.zjzapi.com/receipt/make
请求方式：POST（application/x-www-form-urlencoded）
返回格式：JSON
接口费用：调用免费（回执提交办理完成后按单扣费）。

用法（在 python 目录运行）：python receipt/make.py

图片要求：JPG/JPEG/PNG，大小不超过5M，分辨率小于6000x6000；
办理时间 8:30-20:30（节假日正常办理），预计15分钟审核完成。
"""

import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))

from config import API_KEY
from zjz_api import ZjzApi, ZjzApiError

# ===== 请求参数（修改后直接运行即可）=====
IMAGE_PATH = Path(__file__).resolve().parents[2] / "images" / "test.jpeg"  # 待处理图片路径
ITEM_ID = "1"       # 回执类规格id，可通过「证件照规格列表」示例获取：python item/list.py


def main():
    if API_KEY in ("", "YOUR_APP_KEY"):
        print("提示：请先在 config.py 中将 API_KEY 替换为你的 API Key。", file=sys.stderr)

    try:
        result = ZjzApi.post("receipt/make", {
            "key": API_KEY,
            "item_id": ITEM_ID,
            "image": ZjzApi.image_to_base64(str(IMAGE_PATH)),
        })
    except ZjzApiError as error:
        print(f"调用失败：{error}")
        return

    if not ZjzApi.is_success(result):
        print(f"业务失败：code={result.get('code')}，msg={result.get('msg')}")
        return

    data = result.get("data") or {}
    print(f"回执照制作成功，pic_id：{data.get('pic_id', '')}")
    print("提示：把 pic_id 填入 receipt/submit.py 提交办理订单。")

    save_dir = Path(__file__).resolve().parents[2] / "output" / "receipt-make"
    ZjzApi.save_images(str(save_dir), data.get("list") or {})


if __name__ == "__main__":
    main()
