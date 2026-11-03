#!/usr/bin/env python3
"""V5 职业形象照检测+制作（有水印）

接口地址：https://api.zjzapi.com/idcardv5/make（与证件照制作相同，通过职业照规格 + 专属底色区分）
请求方式：POST（application/x-www-form-urlencoded）
返回格式：JSON
接口费用：调用免费，结果图带水印；可使用返回的 pic_id 调用 v5/get 获取无水印证件照。

用法（在 python 目录运行）：python v5/pose.py

与 v5/make.py 的区别：item_id 使用职业形象照类规格（规格列表中的"职业形象照"分类），
底色为职业照专属：pose_black、pose_blue、pose_darkgray、pose_gray、pose_white。
"""

import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))

from config import API_KEY
from zjz_api import ZjzApi, ZjzApiError

# ===== 请求参数（修改后直接运行即可）=====
IMAGE_PATH = Path(__file__).resolve().parents[2] / "images" / "test.jpeg"  # 待处理图片路径
ITEM_ID = "1"          # 职业形象照规格id，可通过「证件照规格列表」示例获取：python item/list.py
COLORS = "pose_blue"   # 底色：pose_black、pose_blue、pose_darkgray、pose_gray、pose_white，留空返回全部
ENHANCE = 1            # 人像增强：默认已支持，不需要可传0
BEAUTY = 1             # 美颜：默认已支持，不需要可传0
WHITENING = 50         # 美颜-美白等级：0-100，默认50（beauty=1 时生效）


def main():
    if API_KEY in ("", "YOUR_APP_KEY"):
        print("提示：请先在 config.py 中将 API_KEY 替换为你的 API Key。", file=sys.stderr)

    try:
        params = {
            "key": API_KEY,
            "item_id": ITEM_ID,
            "image": ZjzApi.image_to_base64(str(IMAGE_PATH)),
            "enhance": ENHANCE,
            "beauty": BEAUTY,
            "whitening": WHITENING,
        }
        if COLORS:
            params["colors"] = COLORS
        result = ZjzApi.post("idcardv5/make", params)
    except ZjzApiError as error:
        print(f"调用失败：{error}")
        return

    if not ZjzApi.is_success(result):
        print(f"业务失败：code={result.get('code')}，msg={result.get('msg')}")
        return

    data = result.get("data") or {}
    print(f"制作成功（有水印），pic_id：{data.get('pic_id', '')}")
    print("提示：把 pic_id 填入 v5/get.py 可获取无水印证件照。")

    save_dir = Path(__file__).resolve().parents[2] / "output" / "v5-pose"
    ZjzApi.save_images(str(save_dir), data.get("list") or {})


if __name__ == "__main__":
    main()
