#!/usr/bin/env python3
"""证件照规格列表

接口地址：https://api.zjzapi.com/item/list
请求方式：POST（application/x-www-form-urlencoded）
返回格式：JSON
接口费用：调用免费

用法（在 python 目录运行）：python item/list.py

输出各规格的 item_id、名称、像素尺寸与支持底色，
把 item_id 和对应底色填入制作类示例（如 v3/all.py）即可。
"""

import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))

from config import API_KEY
from zjz_api import ZjzApi, ZjzApiError


def main():
    if API_KEY in ("", "YOUR_APP_KEY"):
        print("提示：请先在 config.py 中将 API_KEY 替换为你的 API Key。", file=sys.stderr)

    try:
        result = ZjzApi.post("item/list", {"key": API_KEY})
    except ZjzApiError as error:
        print(f"调用失败：{error}")
        return

    # HTTP 成功不代表业务成功，需再检查返回 JSON 中的 code。
    if not ZjzApi.is_success(result):
        print(f"业务失败：code={result.get('code')}，msg={result.get('msg')}")
        return

    item_list = (result.get("data") or {}).get("list") or []
    if not item_list:
        print("返回数据中没有规格列表。")
        return

    # 打印规格清单，便于挑选 item_id 与底色。
    print("item_id  名称            像素          支持底色")
    for item in item_list:
        size = f"{item.get('width_px', '')}x{item.get('height_px', '')}"
        print(
            f"{str(item.get('item_id', '')):<8} {str(item.get('name', '')):<14} {size:<12} {item.get('color', '')}"
        )
    print(f"共 {len(item_list)} 个规格。")


if __name__ == "__main__":
    main()
