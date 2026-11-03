#!/usr/bin/env python3
"""人像修复增强

接口地址：https://api.zjzapi.com/face/enhance
请求方式：POST（application/x-www-form-urlencoded）
返回格式：JSON
接口费用：返回成功扣除0.06元。

用法（在 python 目录运行）：python face/enhance.py

适用于模糊图片清晰度修复：拍照没对焦好、图片模糊、光线差、
老照片不清晰、图片像素低等；图片大小不超过4M，分辨率不超过3072px。
"""

import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))

from config import API_KEY
from zjz_api import ZjzApi, ZjzApiError

# ===== 请求参数（修改后直接运行即可）=====
IMAGE_PATH = Path(__file__).resolve().parents[2] / "images" / "test.jpeg"  # 待处理图片路径
SIZE = "original"   # 图片大小：original=原始大小、x2=高清图，默认 original


def main():
    if API_KEY in ("", "YOUR_APP_KEY"):
        print("提示：请先在 config.py 中将 API_KEY 替换为你的 API Key。", file=sys.stderr)

    try:
        result = ZjzApi.post("face/enhance", {
            "key": API_KEY,
            "image": ZjzApi.image_to_base64(str(IMAGE_PATH)),
            "size": SIZE,
        })
    except ZjzApiError as error:
        print(f"调用失败：{error}")
        return

    if not ZjzApi.is_success(result):
        print(f"业务失败：code={result.get('code')}，msg={result.get('msg')}")
        return

    data = result.get("data") or {}
    image_url = data.get("image", "")
    if not image_url:
        print("返回数据中没有结果图地址。")
        return

    save_dir = Path(__file__).resolve().parents[2] / "output" / "face-enhance"
    try:
        file = ZjzApi.download_image(image_url, str(save_dir), "enhance.jpg")
        print(f"修复增强成功：{file}")
    except ZjzApiError as error:
        print(f"下载失败，结果图地址：{image_url}")


if __name__ == "__main__":
    main()
