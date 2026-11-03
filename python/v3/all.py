#!/usr/bin/env python3
"""V3 证件照检测+制作

接口地址：https://api.zjzapi.com/idcardv3/all
请求方式：POST（application/x-www-form-urlencoded）
返回格式：JSON
接口费用：制作成功按次扣费；检测不通过扣除检测费用；beauty=1 时美颜另行计费。

用法（在 python 目录运行）：python v3/all.py
"""

import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))

from config import API_KEY
from zjz_api import ZjzApi, ZjzApiError

# ===== 请求参数（修改后直接运行即可）=====
IMAGE_PATH = Path(__file__).resolve().parents[2] / "images" / "test.jpeg"  # 待处理图片路径
ITEM_ID = "1"    # 规格id，可通过「证件照规格列表」示例获取：python item/list.py
COLORS = "white"   # 底色，如 white/blue/red；留空返回该规格支持的全部底色
BEAUTY = 0         # 美颜：0=关闭 1=开启（开启后另行计费）


def main():
    if API_KEY in ("", "YOUR_APP_KEY"):
        print("提示：请先在 config.py 中将 API_KEY 替换为你的 API Key。", file=sys.stderr)

    try:
        result = ZjzApi.post("idcardv3/all", {
            "key": API_KEY,
            "item_id": ITEM_ID,
            "image": ZjzApi.image_to_base64(str(IMAGE_PATH)),
            "colors": COLORS,
            "beauty": BEAUTY,
        })
    except ZjzApiError as error:
        print(f"调用失败：{error}")
        return

    # HTTP 成功不代表业务成功，需再检查返回 JSON 中的 code。
    if not ZjzApi.is_success(result):
        print(f"业务失败：code={result.get('code')}，msg={result.get('msg')}")
        return

    data = result.get("data") or {}
    print(f"制作成功，pic_id：{data.get('pic_id', '')}")

    # 把各底色结果图下载到 output/v3-all/ 目录。
    save_dir = Path(__file__).resolve().parents[2] / "output" / "v3-all"
    ZjzApi.save_images(str(save_dir), data.get("list") or {})


if __name__ == "__main__":
    main()
