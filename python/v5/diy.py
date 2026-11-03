#!/usr/bin/env python3
"""V5 证件照检测+自定义制作（有水印）

接口地址：https://api.zjzapi.com/idcardv5/diy
请求方式：POST（application/x-www-form-urlencoded）
返回格式：JSON
接口费用：调用免费，结果图带水印；可使用返回的 pic_id 调用 v5/get 获取无水印证件照。

用法（在 python 目录运行）：python v5/diy.py

与 v5/make.py 的区别：不使用规格id，自定义像素尺寸、底色、文件大小与 dpi。
"""

import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))

from config import API_KEY
from zjz_api import ZjzApi, ZjzApiError

# ===== 请求参数（修改后直接运行即可）=====
IMAGE_PATH = Path(__file__).resolve().parents[2] / "images" / "test.jpeg"  # 待处理图片路径
WIDTH_PX = "295"      # 像素尺寸-宽，不超2000px
HEIGHT_PX = "413"     # 像素尺寸-高，不超2000px
COLOR = "blue"        # 底色：颜色值（#FFFFFF）或名称（blue、white、red、tint、grey、gradient、dark_blue、dark_red）
FILE_SIZE_MIN = "0"   # 文件最小(kb)，最小5kb，0=不限
FILE_SIZE_MAX = "0"   # 文件最大(kb)，最大600kb，0=不限
DPI = "300"           # 分辨率，默认300
ENHANCE = 1           # 人像增强：默认已支持，不需要可传0
BEAUTY = 1            # 美颜：默认已支持，不需要可传0
WHITENING = 60        # 美颜-美白等级：0-100，默认60（beauty=1 时生效）


def main():
    if API_KEY in ("", "YOUR_APP_KEY"):
        print("提示：请先在 config.py 中将 API_KEY 替换为你的 API Key。", file=sys.stderr)

    try:
        params = {
            "key": API_KEY,
            "image": ZjzApi.image_to_base64(str(IMAGE_PATH)),
            "width_px": WIDTH_PX,
            "height_px": HEIGHT_PX,
            "color": COLOR,
            "dpi": DPI,
            "enhance": ENHANCE,
            "beauty": BEAUTY,
            "whitening": WHITENING,
        }
        if FILE_SIZE_MIN != "0":
            params["file_size_min"] = FILE_SIZE_MIN
        if FILE_SIZE_MAX != "0":
            params["file_size_max"] = FILE_SIZE_MAX

        result = ZjzApi.post("idcardv5/diy", params)
    except ZjzApiError as error:
        print(f"调用失败：{error}")
        return

    if not ZjzApi.is_success(result):
        print(f"业务失败：code={result.get('code')}，msg={result.get('msg')}")
        return

    data = result.get("data") or {}
    print(f"制作成功（有水印），pic_id：{data.get('pic_id', '')}")
    print("提示：把 pic_id 填入 v5/get.py 可获取无水印证件照。")

    save_dir = Path(__file__).resolve().parents[2] / "output" / "v5-diy"
    ZjzApi.save_images(str(save_dir), data.get("list") or {})


if __name__ == "__main__":
    main()
