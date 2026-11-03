#!/usr/bin/env python3
"""证件照规格详情

接口地址：https://api.zjzapi.com/item/get
请求方式：POST（application/x-www-form-urlencoded）
返回格式：JSON
接口费用：调用免费

用法（在 python 目录运行）：python item/detail.py

输出单个规格的像素尺寸、毫米尺寸、dpi、文件大小要求与支持底色，
用于制作类示例（如 v3/all.py）选择 item_id 与底色。
"""

import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))

from config import API_KEY
from zjz_api import ZjzApi, ZjzApiError

# ===== 请求参数（修改后直接运行即可）=====
ITEM_ID = "1"   # 规格id，可通过「证件照规格列表」示例获取：python item/list.py


def main():
    if API_KEY in ("", "YOUR_APP_KEY"):
        print("提示：请先在 config.py 中将 API_KEY 替换为你的 API Key。", file=sys.stderr)

    try:
        result = ZjzApi.post("item/get", {"key": API_KEY, "item_id": ITEM_ID})
    except ZjzApiError as error:
        print(f"调用失败：{error}")
        return

    if not ZjzApi.is_success(result):
        print(f"业务失败：code={result.get('code')}，msg={result.get('msg')}")
        return

    data = result.get("data") or {}
    print(f"规格id：{data.get('id', '')}")
    print(f"名称：{data.get('name', '')}")
    print(f"像素-宽：{data.get('width_px', '')}")
    print(f"像素-高：{data.get('height_px', '')}")
    print(f"毫米-宽：{data.get('width_mm', '')}")
    print(f"毫米-高：{data.get('height_mm', '')}")
    print(f"分辨率：{data.get('dpi', '')}")
    print(f"文件最小(kb)：{data.get('file_size_min', '')}")
    print(f"文件最大(kb)：{data.get('file_size_max', '')}")
    print(f"文件大小说明：{data.get('file_size_msg', '')}")
    print(f"支持底色：{data.get('color', '')}")
    print(f"分类：{data.get('category', '')}")
    print(f"分类id：{data.get('category_id', '')}")


if __name__ == "__main__":
    main()
