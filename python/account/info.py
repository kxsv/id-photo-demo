#!/usr/bin/env python3
"""获取账号信息

接口地址：https://api.zjzapi.com/user/info
请求方式：POST（application/x-www-form-urlencoded）
返回格式：JSON
接口费用：调用免费，数据有10分钟延迟。

用法（在 python 目录运行）：python account/info.py

注意：本接口不使用 API Key，而是使用用户中心的 access_token
（在官网「用户中心 - 账户信息」获取），请填到 config.py 的 ACCESS_TOKEN。
"""

import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))

from config import ACCESS_TOKEN
from zjz_api import ZjzApi, ZjzApiError


def main():
    if ACCESS_TOKEN in ("", "YOUR_ACCESS_TOKEN"):
        print("提示：请先在 config.py 中将 ACCESS_TOKEN 替换为你的 access_token（官网「用户中心 - 账户信息」获取）。", file=sys.stderr)
        return

    try:
        result = ZjzApi.post("user/info", {"access_token": ACCESS_TOKEN})
    except ZjzApiError as error:
        print(f"调用失败：{error}")
        return

    if not ZjzApi.is_success(result):
        print(f"业务失败：code={result.get('code')}，msg={result.get('msg')}")
        return

    data = result.get("data") or {}
    print(f"登录账号：{data.get('mobile', '')}")
    print(f"账户余额：{data.get('money', '')} 元")
    print(f"应用数量：{data.get('app_num', '')}")
    print(f"今日调用：{data.get('day_make', '')} 次")


if __name__ == "__main__":
    main()
