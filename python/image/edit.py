#!/usr/bin/env python3
"""图片编辑

接口地址：https://api.zjzapi.com/image/edit
请求方式：POST（application/x-www-form-urlencoded）
返回格式：JSON
接口费用：调用免费（账户余额最低需有10元才可调用）。

用法（在 python 目录运行）：python image/edit.py

修改图片像素（dpi）与文件大小（kb）；图片大小不超过1M，分辨率不超过1024px。
"""

import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))

from config import API_KEY
from zjz_api import ZjzApi, ZjzApiError

# ===== 请求参数（修改后直接运行即可）=====
IMAGE_PATH = Path(__file__).resolve().parents[2] / "images" / "test.jpeg"  # 待处理图片路径
DPI = "300"           # 图片dpi，支持72,96,150,200,300,350,500,1000；填0表示不修改
FILE_SIZE_MIN = "0"   # 文件最小(kb)，最小5kb；填0表示不限制
FILE_SIZE_MAX = "0"   # 文件最大(kb)，最大300kb；填0表示不限制


def main():
    if API_KEY in ("", "YOUR_APP_KEY"):
        print("提示：请先在 config.py 中将 API_KEY 替换为你的 API Key。", file=sys.stderr)

    try:
        params = {
            "key": API_KEY,
            "image": ZjzApi.image_to_base64(str(IMAGE_PATH)),
        }
        if DPI != "0":
            params["dpi"] = DPI
        if FILE_SIZE_MIN != "0":
            params["file_size_min"] = FILE_SIZE_MIN
        if FILE_SIZE_MAX != "0":
            params["file_size_max"] = FILE_SIZE_MAX
        result = ZjzApi.post("image/edit", params)
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

    save_dir = Path(__file__).resolve().parents[2] / "output" / "image-edit"
    try:
        file = ZjzApi.download_image(image_url, str(save_dir), "edit.jpg")
        print(f"图片编辑成功：{file}")
    except ZjzApiError as error:
        print(f"下载失败，结果图地址：{image_url}")


if __name__ == "__main__":
    main()
