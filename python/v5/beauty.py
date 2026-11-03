#!/usr/bin/env python3
"""V5 证件照美颜（有水印）

接口地址：https://api.zjzapi.com/idcardv5/beauty
请求方式：POST（application/x-www-form-urlencoded）
返回格式：JSON
接口费用：对制作出来的证件照进行美颜；不支持换装后再美颜。

用法（在 python 目录运行）：python v5/beauty.py

注意：美颜基于制作接口返回的 pic_id（有效期50分钟），
请先把 v5/make.py 运行得到的 pic_id 填到下方参数区。
"""

import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))

from config import API_KEY
from zjz_api import ZjzApi, ZjzApiError

# ===== 请求参数（修改后直接运行即可）=====
PIC_ID = ""   # 制作接口返回的 pic_id（有效期50分钟），先运行 python v5/make.py 获取


def main():
    if API_KEY in ("", "YOUR_APP_KEY"):
        print("提示：请先在 config.py 中将 API_KEY 替换为你的 API Key。", file=sys.stderr)
    if not PIC_ID:
        print("请先在参数区填写 pic_id（运行 python v5/make.py 获取）。", file=sys.stderr)
        return

    try:
        result = ZjzApi.post("idcardv5/beauty", {
            "key": API_KEY,
            "pic_id": PIC_ID,
        })
    except ZjzApiError as error:
        print(f"调用失败：{error}")
        return

    if not ZjzApi.is_success(result):
        print(f"业务失败：code={result.get('code')}，msg={result.get('msg')}")
        return

    data = result.get("data") or {}
    print("美颜成功。")

    save_dir = Path(__file__).resolve().parents[2] / "output" / "v5-beauty"
    ZjzApi.save_images(str(save_dir), data.get("list") or {})


if __name__ == "__main__":
    main()
