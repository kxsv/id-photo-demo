#!/usr/bin/env python3
"""V5 证件照换装（有水印）

接口地址：https://api.zjzapi.com/idcardv5/clothe
请求方式：POST（application/x-www-form-urlencoded）
返回格式：JSON
接口费用：调用免费，结果图带水印；可使用返回的 pic_id 调用 v5/get 获取无水印证件照。

用法（在 python 目录运行）：python v5/clothe.py

注意：换装基于制作接口返回的 pic_id（有效期50分钟），
请先把 v5/make.py 运行得到的 pic_id 填到下方参数区。
"""

import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))

from config import API_KEY
from zjz_api import ZjzApi, ZjzApiError

# ===== 请求参数（修改后直接运行即可）=====
PIC_ID = ""   # 制作接口返回的 pic_id（有效期50分钟），先运行 python v5/make.py 获取
CLOTHE = ""   # 服装key，见官网《服装列表》文档


def main():
    if API_KEY in ("", "YOUR_APP_KEY"):
        print("提示：请先在 config.py 中将 API_KEY 替换为你的 API Key。", file=sys.stderr)
    if not PIC_ID or not CLOTHE:
        print("请先在参数区填写 pic_id（运行 python v5/make.py 获取）与服装key clothe。", file=sys.stderr)
        return

    try:
        result = ZjzApi.post("idcardv5/clothe", {
            "key": API_KEY,
            "pic_id": PIC_ID,
            "clothe": CLOTHE,
        })
    except ZjzApiError as error:
        print(f"调用失败：{error}")
        return

    if not ZjzApi.is_success(result):
        print(f"业务失败：code={result.get('code')}，msg={result.get('msg')}")
        return

    data = result.get("data") or {}
    print("换装成功（有水印）。")

    save_dir = Path(__file__).resolve().parents[2] / "output" / "v5-clothe"
    ZjzApi.save_images(str(save_dir), data.get("list") or {})


if __name__ == "__main__":
    main()
