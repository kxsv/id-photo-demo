#!/usr/bin/env python3
"""黑白图片上色

接口地址：https://api.zjzapi.com/image/colorize
请求方式：POST（application/x-www-form-urlencoded）
返回格式：JSON
接口费用：调用成功扣除0.012元，调用失败免费。

用法（在 python 目录运行）：python image/colorize.py

适用于黑白图片、老照片上色；图片大小不超过4M，分辨率不超过3072px。
"""

import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))

from config import API_KEY
from zjz_api import ZjzApi, ZjzApiError

# ===== 请求参数（修改后直接运行即可）=====
IMAGE_PATH = Path(__file__).resolve().parents[2] / "images" / "test.jpeg"  # 待处理图片路径（建议使用黑白照片）


def main():
    if API_KEY in ("", "YOUR_APP_KEY"):
        print("提示：请先在 config.py 中将 API_KEY 替换为你的 API Key。", file=sys.stderr)

    try:
        result = ZjzApi.post("image/colorize", {
            "key": API_KEY,
            "image": ZjzApi.image_to_base64(str(IMAGE_PATH)),
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

    save_dir = Path(__file__).resolve().parents[2] / "output" / "image-colorize"
    try:
        file = ZjzApi.download_image(image_url, str(save_dir), "colorize.jpg")
        print(f"上色成功：{file}")
    except ZjzApiError as error:
        print(f"下载失败，结果图地址：{image_url}")


if __name__ == "__main__":
    main()
