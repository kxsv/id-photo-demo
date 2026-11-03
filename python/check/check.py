#!/usr/bin/env python3
"""证件照合规性检测

接口地址：https://api.zjzapi.com/idphoto/check
请求方式：POST（application/x-www-form-urlencoded）
返回格式：JSON
接口费用：检测成功扣除0.012元。

用法（在 python 目录运行）：python check/check.py

提供规格id（item_id）即可按该规格检测；也可不传 item_id，
自定义像素、dpi、文件大小、底色、格式以及各项人像构图/姿态/画质/遮挡/着装检测项，
完整参数列表见官网《证件照合规性检测》文档。
"""

import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))

from config import API_KEY
from zjz_api import ZjzApi, ZjzApiError

# ===== 请求参数（修改后直接运行即可）=====
IMAGE_PATH = Path(__file__).resolve().parents[2] / "images" / "test.jpeg"  # 待检测图片路径（不超过4M，分辨率不超过3072px）
ITEM_ID = "1"    # 规格id，提供后可省略 width_px/height_px/dpi 等规格参数；填 0 表示使用下方自定义规格参数

# 方式二（item_id 为 0 时生效）的自定义检测规格，按需增删。
CUSTOM_WIDTH_PX = "295"          # 证件照宽（px）
CUSTOM_HEIGHT_PX = "413"         # 证件照高（px）
CUSTOM_DPI = "300"               # 图片分辨率，常见300、350
CUSTOM_FILE_FORMAT = "jpg"       # 照片格式：jpg、jpeg、png
CUSTOM_BG_COLORS = "white,blue"  # 允许的底色，逗号分隔
CUSTOM_SOLID_BG = "1"            # 纯色背景检测：0/1
FACE_WIDTH_RATIO_MIN = "0.4"     # 人像构图-脸宽占比下限（0.01-1.0）
FACE_WIDTH_RATIO_MAX = "0.6"     # 人像构图-脸宽占比上限（0.01-1.0）
SHOULDER_HORIZONTAL = "5"        # 姿态-双肩水平角度（推荐值见文档）
FACE_POSE_YAW = "15"             # 姿态-左右转头角度（推荐值见文档）
FACE_CLARITY = "0.5"             # 画质-清晰度（推荐值见文档）
FACE_BRIGHTNESS = "0.3"          # 画质-亮度（推荐值见文档）
FACE_MASK = "1"                  # 遮挡-口罩检测：0=不检测 1=检测
GLASSES = "1"                    # 遮挡-眼镜检测：0=不检测 1=检测


def main():
    if API_KEY in ("", "YOUR_APP_KEY"):
        print("提示：请先在 config.py 中将 API_KEY 替换为你的 API Key。", file=sys.stderr)

    try:
        params = {
            "key": API_KEY,
            "image": ZjzApi.image_to_base64(str(IMAGE_PATH)),
        }
        if ITEM_ID != "0":
            # 方式一：按规格id检测（推荐，规格要求自动带上）。
            params["item_id"] = ITEM_ID
        else:
            # 方式二：自定义检测规格（不传 item_id 时生效），按需增删。
            params.update({
                "width_px": CUSTOM_WIDTH_PX,
                "height_px": CUSTOM_HEIGHT_PX,
                "dpi": CUSTOM_DPI,
                "file_format": CUSTOM_FILE_FORMAT,
                "bg_colors": CUSTOM_BG_COLORS,
                "solid_bg": CUSTOM_SOLID_BG,
                "face_width_ratio_min": FACE_WIDTH_RATIO_MIN,
                "face_width_ratio_max": FACE_WIDTH_RATIO_MAX,
                "shoulder_horizontal": SHOULDER_HORIZONTAL,
                "face_pose_yaw": FACE_POSE_YAW,
                "face_clarity": FACE_CLARITY,
                "face_brightness": FACE_BRIGHTNESS,
                "face_mask": FACE_MASK,
                "glasses": GLASSES,
            })
        result = ZjzApi.post("idphoto/check", params)
    except ZjzApiError as error:
        print(f"调用失败：{error}")
        return

    if not ZjzApi.is_success(result):
        print(f"业务失败：code={result.get('code')}，msg={result.get('msg')}")
        return

    data = result.get("data") or {}
    if int(data.get("passed_status") or -1) == 1:
        print("总检测结果：通过")
    else:
        print("总检测结果：不通过")

    passed_list = data.get("passed_list") or []
    if not passed_list:
        print("返回数据中没有检测项列表。")
        return

    print(f"{'检测项':<16}{'结果':<6}{'提示':<10}检测值")
    for item in passed_list:
        detail = item.get("detail") or {}
        actual = detail.get("actual", "")
        unit = detail.get("unit", "")
        if actual != "" and actual is not None:
            actual = f"{actual}{unit}"
        passed = "通过" if int(item.get("passed") or 0) == 1 else "不通过"
        print(f"{str(item.get('key', '')):<16}{passed:<6}{str(item.get('msg', '')):<10}{actual}")


if __name__ == "__main__":
    main()
