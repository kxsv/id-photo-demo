#!/usr/bin/env python3
"""V3 职业形象照检测+制作

接口地址：https://api.zjzapi.com/idcardv3/all（与证件照制作相同，通过职业照规格 + 专属底色区分）
请求方式：POST（application/x-www-form-urlencoded）
返回格式：JSON
接口费用：检测并制作成功扣除0.012元，检测失败扣除0.006元；若使用"美颜"功能将多扣除0.012元。

用法（在 python 目录运行）：python v3/pose.py

与 v3/all.py 的区别：item_id 使用职业形象照类规格（规格列表中的"职业形象照"分类），
底色为职业照专属：pose_black、pose_blue、pose_darkgray、pose_gray、pose_white。
"""

import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))

from config import API_KEY
from zjz_api import ZjzApi, ZjzApiError

# ===== 请求参数（修改后直接运行即可）=====
IMAGE_PATH = Path(__file__).resolve().parents[2] / "images" / "test.jpeg"  # 待处理图片路径
ITEM_ID = "1"       # 职业形象照规格id，可通过「证件照规格列表」示例获取：python item/list.py
COLORS = "pose_blue"  # 底色：pose_black、pose_blue、pose_darkgray、pose_gray、pose_white，留空返回全部
BEAUTY = 0         # 美颜：0=关闭 1=开启（开启后另行计费）


def main():
    if API_KEY in ("", "YOUR_APP_KEY"):
        print("提示：请先在 config.py 中将 API_KEY 替换为你的 API Key。", file=sys.stderr)

    try:
        params = {
            "key": API_KEY,
            "item_id": ITEM_ID,
            "image": ZjzApi.image_to_base64(str(IMAGE_PATH)),
            "beauty": BEAUTY,
        }
        # colors 为空时不传，由服务端返回全部底色。
        if COLORS:
            params["colors"] = COLORS
        result = ZjzApi.post("idcardv3/all", params)
    except ZjzApiError as error:
        print(f"调用失败：{error}")
        return

    if not ZjzApi.is_success(result):
        print(f"业务失败：code={result.get('code')}，msg={result.get('msg')}")
        return

    data = result.get("data") or {}
    print(f"制作成功，pic_id：{data.get('pic_id', '')}")

    save_dir = Path(__file__).resolve().parents[2] / "output" / "v3-pose"
    ZjzApi.save_images(str(save_dir), data.get("list") or {})


if __name__ == "__main__":
    main()
