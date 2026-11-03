#!/usr/bin/env python3
"""证件照冲印排版（相纸联排）。

将一张已制作好的证件照，按预设的相纸模板（images/print-bg/<相纸>/<尺寸>.png）
复制到模板的各个槽位上，输出一张可直接拿去冲印的图（默认 300dpi）。

用法（在 python 目录运行）：
    python print.py                  # 使用默认输入（output/v3-diy/blue.jpg，需先运行 v3/diy.py 生成）
    python print.py 路径/照片.jpg     # 指定一张做好的证件照

依赖：Pillow（已在 requirements.txt 中声明，pip install -r requirements.txt 即可）。
前置：请先运行 v3/diy.py 生成 295x413 的证件照到 output/v3-diy/blue.jpg，
      或自行用任意制作接口产出对应尺寸的证件照后传入路径。
"""

import sys
from pathlib import Path

# ===== 请求参数（按需修改）=====
PAPER = "print6"            # 相纸规格目录名（images/print-bg/ 下的子目录，print6=6寸）
USE_TEXT_TPL = False        # 是否使用带文字/裁切线说明的模板（*-text.png）
DPI = 300                   # 输出分辨率

# 默认输入：仓库根（id-photo-demo） output/v3-diy/blue.jpg（与 v3/diy.py 输出位置一致）
# 注意：本文件位于 python/ 下（比 v3/diy.py 浅一层），故用 parents[1] 取到 id-photo-demo 目录
DEFAULT_ID_PHOTO = Path(__file__).resolve().parents[1] / "output" / "v3-diy" / "blue.jpg"
TEMPLATE_DIR = Path(__file__).resolve().parents[1] / "images" / "print-bg"
OUTPUT_DIR = Path(__file__).resolve().parents[1] / "output" / "print"

# 排版配置（数据驱动，便于增删规格）
# key        = 证件照像素尺寸（宽-高）
# rotate     = 是否需要旋转 90°（模板槽位为横放时置 True）
# size       = 槽位像素尺寸（输入图会被缩放/拉伸到该尺寸后再贴入）
# positions  = 各槽位左上角坐标 [x, y]（与模板图严格对应）
LAYOUTS = {
    "295-413":   {"rotate": False, "size": (295, 413),  "positions": [[120,180],[435,180],[750,180],[1065,180],[1380,180],[120,613],[435,613],[750,613],[1065,613],[1380,613]]},
    "260-378":   {"rotate": False, "size": (260, 378),  "positions": [[207,214],[487,214],[767,214],[1047,214],[1327,214],[207,612],[487,612],[767,612],[1047,612],[1327,612]]},
    "413-579":   {"rotate": True,  "size": (413, 579),  "positions": [[308,179],[907,179],[308,612],[907,612]]},
    "390-567":   {"rotate": True,  "size": (390, 567),  "positions": [[320,202],[907,202],[320,612],[907,612]]},
    "413-531":   {"rotate": True,  "size": (413, 531),  "positions": [[81,179],[632,179],[1183,179],[81,612],[632,612],[1183,612]]},
    "626-413":   {"rotate": False, "size": (626, 413),  "positions": [[261,179],[907,179],[261,612],[907,612]]},
    "649-991":   {"rotate": False, "size": (649, 991),  "positions": [[238,107],[907,107]]},
    "1050-1499": {"rotate": True,  "size": (1050,1499), "positions": [[148,78]]},
    # 以下三组坐标由"统一边距+20px间距"的均匀网格推算，并与模板裁切线对齐确认（模板本身未画占位框）
    "358-441":   {"rotate": False, "size": (358, 441),  "positions": [[151,151],[529,151],[907,151],[1285,151],[151,612],[529,612],[907,612],[1285,612]]},
    "378-472":   {"rotate": False, "size": (378, 472),  "positions": [[111,120],[509,120],[907,120],[1305,120],[111,612],[509,612],[907,612],[1305,612]]},
    "600-600":   {"rotate": False, "size": (600, 600),  "positions": [[287,302],[907,302]]},
}


def main():
    id_photo = Path(sys.argv[1]) if len(sys.argv) > 1 else DEFAULT_ID_PHOTO

    if not id_photo.is_file():
        print(f"输入证件照不存在：{id_photo}")
        print("请先运行 v3/diy.py 生成，或在命令行传入路径：python print.py 路径/照片.jpg")
        return

    from PIL import Image

    photo = Image.open(id_photo).convert("RGB")
    name = f"{photo.width}-{photo.height}"
    if name not in LAYOUTS:
        print(f"不支持的证件照尺寸 {name}（实际像素 {photo.width}x{photo.height}）。")
        print("支持：" + ", ".join(LAYOUTS.keys()))
        print("请使用对应规格制作证件照（如 v3/diy.py 设 WIDTH_PX/HEIGHT_PX）。")
        return

    cfg = LAYOUTS[name]
    suffix = "-text" if USE_TEXT_TPL else ""
    bg_path = TEMPLATE_DIR / PAPER / f"{name}{suffix}.png"
    if not bg_path.is_file():
        print(f"相纸模板不存在：{bg_path}")
        print(f"请确认 {TEMPLATE_DIR / PAPER} 下存在 {name}{suffix}.png")
        return

    # 先铺白底，避免模板透明区在 JPEG 中变黑
    tpl = Image.open(bg_path).convert("RGB")
    canvas = Image.new("RGB", tpl.size, (255, 255, 255))
    canvas.paste(tpl, (0, 0))

    # 缩放（拉伸）到槽位尺寸，需要旋转则顺时针 90°（Pillow 正角度为逆时针，故用 -90）
    slot = photo.resize(cfg["size"])
    if cfg["rotate"]:
        slot = slot.rotate(-90, expand=True)

    for x, y in cfg["positions"]:
        canvas.paste(slot, (x, y))

    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)
    out = OUTPUT_DIR / f"{PAPER}-{name}.jpg"
    canvas.save(out, "JPEG", quality=95, dpi=(DPI, DPI))

    print(f"排版完成：{out}")


if __name__ == "__main__":
    main()
