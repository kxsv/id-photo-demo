#!/usr/bin/env python3
"""获取应用信息

接口地址：https://api.zjzapi.com/user/app
请求方式：POST（application/x-www-form-urlencoded）
返回格式：JSON
接口费用：调用免费，数据有10分钟延迟。

用法（在 python 目录运行）：python account/app.py

注意：本接口需要同时使用用户中心的 access_token（官网「用户中心 - 账户信息」获取，
填到 config.py 的 ACCESS_TOKEN）与你的 API Key。
"""

import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))

from config import API_KEY, ACCESS_TOKEN
from zjz_api import ZjzApi, ZjzApiError

# 应用权限：1=证件照v3，2=证件照v5，3=支持全部，4=证件照回执
TYPE_MAP = {1: "证件照v3", 2: "证件照v5", 3: "支持全部", 4: "证件照回执"}


def main():
    if ACCESS_TOKEN in ("", "YOUR_ACCESS_TOKEN"):
        print("提示：请先在 config.py 中将 ACCESS_TOKEN 替换为你的 access_token（官网「用户中心 - 账户信息」获取）。", file=sys.stderr)
        return

    try:
        result = ZjzApi.post("user/app", {
            "access_token": ACCESS_TOKEN,
            "key": API_KEY,
        })
    except ZjzApiError as error:
        print(f"调用失败：{error}")
        return

    if not ZjzApi.is_success(result):
        print(f"业务失败：code={result.get('code')}，msg={result.get('msg')}")
        return

    data = result.get("data") or {}
    app_type = TYPE_MAP.get(int(data.get("type") or 0), str(data.get("type", "")))
    status = "启用中" if int(data.get("status") or 0) == 1 else "已停用"
    print(f"应用名称：{data.get('name', '')}")
    print(f"应用权限：{app_type}")
    print(f"应用状态：{status}")
    print(f"调用次数：{data.get('make_num', '')}")


if __name__ == "__main__":
    main()
