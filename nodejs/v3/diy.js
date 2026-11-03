// V3 证件照检测+自定义制作
//
// 接口地址：https://api.zjzapi.com/idcardv3/diy
// 请求方式：POST（application/x-www-form-urlencoded）
// 返回格式：JSON
// 接口费用：检测并制作成功扣除0.012元，检测失败扣除0.006元；若使用"美颜"功能将多扣除0.012元。
//
// 用法（在 nodejs 目录运行）：node v3/diy.js
//
// 与 v3/all 的区别：不使用规格id，自定义像素尺寸、底色、文件大小与 dpi。

import process from 'node:process';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import { API_KEY } from '../config.js';
import { post, isSuccess, imageToBase64, downloadImage } from '../lib/zjz-api.js';

// ===== 请求参数（修改后直接运行即可）=====
const imagePath = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '../../images/test.jpeg'); // 待处理图片路径
const widthPx = 295;       // 像素尺寸-宽，不超2000px
const heightPx = 413;      // 像素尺寸-高，不超2000px
const color = 'blue';      // 底色：颜色值（#FFFFFF）或名称（blue、white、red、tint、grey、gradient、dark_blue、dark_red）
const fileSizeMin = 0;     // 文件最小(kb)，最小5kb，0=不限
const fileSizeMax = 0;     // 文件最大(kb)，最大300kb，0=不限
const dpi = 300;           // 分辨率，默认300
const beauty = 0;          // 美颜：0=关闭 1=开启（开启后另行计费）
const whitening = 50;      // 美颜-美白等级：0-100，默认50（beauty=1 时生效）

if (API_KEY === '' || API_KEY === 'YOUR_APP_KEY') {
  console.error('提示：请先在 nodejs/config.js 中将 API_KEY 替换为你的 API Key。');
}

try {
  const params = {
    key: API_KEY,
    image: await imageToBase64(imagePath),
    width_px: widthPx,
    height_px: heightPx,
    color,
    dpi,
    beauty,
    whitening,
  };
  if (fileSizeMin > 0) params.file_size_min = fileSizeMin;
  if (fileSizeMax > 0) params.file_size_max = fileSizeMax;

  const result = await post('idcardv3/diy', params);

  if (!isSuccess(result)) {
    console.log(`业务失败：code=${result.code}，msg=${result.msg}`);
    process.exit(1);
  }

  const picId = result.data?.pic_id ?? '';
  console.log(`制作成功，pic_id：${picId}`);

  const saveDir = path.join('output', 'v3-diy');
  const list = result.data?.list ?? {};
  if (Object.keys(list).length === 0) {
    console.log('返回数据中没有图片列表。');
    process.exit(0);
  }
  for (const [c, imageUrl] of Object.entries(list)) {
    try {
      const file = await downloadImage(imageUrl, saveDir, `${c}.jpg`);
      console.log(`底色 ${c}：${file}`);
    } catch (error) {
      console.log(`底色 ${c} 下载失败：${error.message}`);
    }
  }
} catch (error) {
  console.log(`调用失败：${error.message}`);
  process.exit(1);
}
