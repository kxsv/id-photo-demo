#!/usr/bin/env python3
"""V3 图片检测+更换底色

接口地址：https://api.zjzapi.com/idcardv3/bg
请求方式：POST（application/x-www-form-urlencoded）
返回格式：JSON
接口费用：检测并制作成功扣除0.012元，检测失败扣除0.006元；若使用"美颜"功能将多扣除0.012元。

用法（在 python 目录运行）：python v3/bg.py

与 v3/all.py 的区别：不按规格重排裁剪，仅对已有照片检测后更换底色；
图片大小不超过4M，分辨率不超过3072px。
"""

import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))

from config import API_KEY
from zjz_api import ZjzApi, ZjzApiError

# ===== 请求参数（修改后直接运行即可）=====
IMAGE_PATH = Path(__file__).resolve().parents[2] / "images" / "test.jpeg"  # 待处理图片路径
COLOR = "blue"    # 底色：颜色值（#FFFFFF）或名称（blue、white、red、tint、grey、gradient、dark_blue、dark_red），留空返回所有底色
BEAUTY = 0         # 美颜：0=关闭 1=开启（开启后另行计费）
WHITENING = 50     # 美颜-美白等级：0-100，默认50（beauty=1 时生效）


def main():
    if API_KEY in ("", "YOUR_APP_KEY"):
        print("提示：请先在 config.py 中将 API_KEY 替换为你的 API Key。", file=sys.stderr)

    try:
        params = {
            "key": API_KEY,
            "image": ZjzApi.image_to_base64(str(IMAGE_PATH)),
            "beauty": BEAUTY,
            "whitening": WHITENING,
        }
        # color 为空时不传，由服务端返回全部底色。
        if COLOR:
            params["color"] = COLOR
        result = ZjzApi.post("idcardv3/bg", params)
    except ZjzApiError as error:
        print(f"调用失败：{error}")
        return

    if not ZjzApi.is_success(result):
        print(f"业务失败：code={result.get('code')}，msg={result.get('msg')}")
        return

    data = result.get("data") or {}
    print(f"更换底色成功，pic_id：{data.get('pic_id', '')}")

    # 把各底色结果图下载到 output/v3-bg/ 目录。
    save_dir = Path(__file__).resolve().parents[2] / "output" / "v3-bg"
    ZjzApi.save_images(str(save_dir), data.get("list") or {})


if __name__ == "__main__":
    main()
