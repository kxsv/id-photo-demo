#!/usr/bin/env python3
"""V5 获取证件照（无水印）

接口地址：https://api.zjzapi.com/idcardv5/get
请求方式：POST（application/x-www-form-urlencoded）
返回格式：JSON
接口费用：获取成功扣除0.6元，重复获取不扣费。

用法（在 python 目录运行）：python v5/get.py

注意：pic_id 来自制作类接口（v5/make、v5/bg、v5/diy、v5/pose 等）的返回结果，
有效期50分钟；若换装过，请同时填写换装时使用的 clothe。
"""

import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))

from config import API_KEY
from zjz_api import ZjzApi, ZjzApiError

# ===== 请求参数（修改后直接运行即可）=====
PIC_ID = ""        # 制作接口返回的 pic_id（有效期50分钟），先运行 python v5/make.py 获取
CLOTHE = ""        # 换装时使用的服装key，未换装留空
GIFT_HD = 0        # 1=赠送一张三寸高清照（仅规格 width_px<500 支持；更换底色、自定义制作暂不支持）
GIFT_HD_COLOR = "" # 高清照底色，留空默认返回制作结果第一个底色


def main():
    if API_KEY in ("", "YOUR_APP_KEY"):
        print("提示：请先在 config.py 中将 API_KEY 替换为你的 API Key。", file=sys.stderr)
    if not PIC_ID:
        print("请先在参数区填写 pic_id（运行 python v5/make.py 获取）。", file=sys.stderr)
        return

    try:
        params = {
            "key": API_KEY,
            "pic_id": PIC_ID,
            "gift_hd": GIFT_HD,
        }
        if CLOTHE:
            params["clothe"] = CLOTHE
        if GIFT_HD_COLOR:
            params["gift_hd_color"] = GIFT_HD_COLOR
        result = ZjzApi.post("idcardv5/get", params)
    except ZjzApiError as error:
        print(f"调用失败：{error}")
        return

    if not ZjzApi.is_success(result):
        print(f"业务失败：code={result.get('code')}，msg={result.get('msg')}")
        return

    data = result.get("data") or {}
    print("获取无水印证件照成功。")

    save_dir = Path(__file__).resolve().parents[2] / "output" / "v5-get"
    ZjzApi.save_images(str(save_dir), data.get("list") or {})

    # 赠送的三寸高清照（开启 gift_hd 时才有）。
    gift_hd = data.get("gift_hd", "")
    if gift_hd:
        try:
            file = ZjzApi.download_image(gift_hd, str(save_dir), "gift_hd.jpg")
            print(f"高清照：{file}")
        except ZjzApiError as error:
            print(f"高清照下载失败：{error}（结果图地址：{gift_hd}）")


if __name__ == "__main__":
    main()
