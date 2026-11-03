// 证件照冲印排版（相纸联排）。
//
// 将一张已制作好的证件照，按预设的相纸模板（../images/print-bg/<相纸>/<尺寸>.png）
// 复制到模板的各个槽位上，输出一张可直接拿去冲印的图（默认 300dpi）。
//
// 用法（在 nodejs 目录运行）：
//   node print.js                  # 使用默认输入（output/v3-diy/blue.jpg，需先运行 v3/diy.js 生成）
//   node print.js 路径/照片.jpg     # 指定一张做好的证件照
//
// 依赖：sharp（图片合成与 DPI 设置；已在 package.json 中声明，npm install 即可）。
//       其余接口示例仍零依赖。
// 前置：请先运行 v3/diy.js 生成 295x413 的证件照到 output/v3-diy/blue.jpg，
//       或自行用任意制作接口产出对应尺寸的证件照后传入路径。

import path from 'node:path';
import fs from 'node:fs';
import { fileURLToPath } from 'node:url';
import sharp from 'sharp';

// ===== 请求参数（按需修改）=====
const PAPER = 'print6';        // 相纸规格目录名（images/print-bg/ 下的子目录，print6=6寸）
const USE_TEXT_TPL = false;    // 是否使用带文字/裁切线说明的模板（*-text.png）
const DPI = 300;               // 输出分辨率

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const DEFAULT_ID_PHOTO = path.resolve(__dirname, 'output/v3-diy/blue.jpg');

// 排版配置（数据驱动，便于增删规格）
// key        = 证件照像素尺寸（宽-高）
// rotate     = 是否需要旋转 90°（模板槽位为横放时置 true）
// size       = 槽位像素尺寸（输入图会被缩放/拉伸到该尺寸后再贴入）
// positions  = 各槽位左上角坐标 [x, y]（与模板图严格对应）
const LAYOUTS = {
  '295-413':   { rotate: false, size: [295, 413],  positions: [[120,180],[435,180],[750,180],[1065,180],[1380,180],[120,613],[435,613],[750,613],[1065,613],[1380,613]] },
  '260-378':   { rotate: false, size: [260, 378],  positions: [[207,214],[487,214],[767,214],[1047,214],[1327,214],[207,612],[487,612],[767,612],[1047,612],[1327,612]] },
  '413-579':   { rotate: true,  size: [413, 579],  positions: [[308,179],[907,179],[308,612],[907,612]] },
  '390-567':   { rotate: true,  size: [390, 567],  positions: [[320,202],[907,202],[320,612],[907,612]] },
  '413-531':   { rotate: true,  size: [413, 531],  positions: [[81,179],[632,179],[1183,179],[81,612],[632,612],[1183,612]] },
  '626-413':   { rotate: false, size: [626, 413],  positions: [[261,179],[907,179],[261,612],[907,612]] },
  '649-991':   { rotate: false, size: [649, 991],  positions: [[238,107],[907,107]] },
  '1050-1499': { rotate: true,  size: [1050,1499], positions: [[148,78]] },
  // 以下三组坐标由"统一边距+20px间距"的均匀网格推算，并与模板裁切线对齐确认（模板本身未画占位框）
  '358-441':   { rotate: false, size: [358, 441],  positions: [[151,151],[529,151],[907,151],[1285,151],[151,612],[529,612],[907,612],[1285,612]] },
  '378-472':   { rotate: false, size: [378, 472],  positions: [[111,120],[509,120],[907,120],[1305,120],[111,612],[509,612],[907,612],[1305,612]] },
  '600-600':   { rotate: false, size: [600, 600],  positions: [[287,302],[907,302]] },
};

async function main() {
  const idPhoto = process.argv[2] || DEFAULT_ID_PHOTO;

  if (!fs.existsSync(idPhoto)) {
    console.error(`输入证件照不存在：${idPhoto}`);
    console.error('请先运行 v3/diy.js 生成，或在命令行传入路径：node print.js 路径/照片.jpg');
    process.exit(1);
  }

  const photo = sharp(idPhoto);
  const meta = await photo.metadata();
  const name = `${meta.width}-${meta.height}`;
  const cfg = LAYOUTS[name];
  if (!cfg) {
    console.error(`不支持的证件照尺寸 ${name}（实际像素 ${meta.width}x${meta.height}）。`);
    console.error('支持：' + Object.keys(LAYOUTS).join(', '));
    process.exit(1);
  }

  const suffix = USE_TEXT_TPL ? '-text' : '';
  const bgPath = path.resolve(__dirname, '../images/print-bg', PAPER, `${name}${suffix}.png`);
  if (!fs.existsSync(bgPath)) {
    console.error(`相纸模板不存在：${bgPath}`);
    process.exit(1);
  }

  // 缩放（拉伸）到槽位尺寸，需要旋转则顺时针 90°；转 PNG 缓冲以便无损拼合
  let slot = photo.resize(cfg.size[0], cfg.size[1]);
  if (cfg.rotate) slot = slot.rotate(90);
  const slotBuf = await slot.png().toBuffer();

  const composites = cfg.positions.map(([x, y]) => ({ input: slotBuf, left: x, top: y }));

  const outDir = path.join('output', 'print');
  fs.mkdirSync(outDir, { recursive: true });
  const out = path.join(outDir, `${PAPER}-${name}.jpg`);

  // flatten 先铺白底，避免模板透明区在 JPEG 中变黑
  await sharp(bgPath)
    .flatten({ background: { r: 255, g: 255, b: 255 } })
    .composite(composites)
    .withMetadata({ density: DPI })
    .jpeg({ quality: 95 })
    .toFile(out);

  console.log(`排版完成：${out}`);
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
