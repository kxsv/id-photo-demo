// 图片编辑
//
// 接口地址：https://api.zjzapi.com/image/edit
// 请求方式：POST（application/x-www-form-urlencoded）
// 返回格式：JSON
// 接口费用：调用免费（账户余额最低需有10元才可调用）。
//
// 用法（在 nodejs 目录运行）：node image/edit.js
//
// 修改图片像素（dpi）与文件大小（kb）；图片大小不超过1M，分辨率不超过1024px。

import process from 'node:process';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import { API_KEY } from '../config.js';
import { post, isSuccess, imageToBase64, downloadImage } from '../lib/zjz-api.js';

// ===== 请求参数（修改后直接运行即可）=====
const imagePath = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '../../images/test.jpeg'); // 待处理图片路径
const dpi = 300;           // 图片dpi，支持72,96,150,200,300,350,500,1000；填0表示不修改
const fileSizeMin = 0;     // 文件最小(kb)，最小5kb；填0表示不限制
const fileSizeMax = 0;     // 文件最大(kb)，最大300kb；填0表示不限制

if (API_KEY === '' || API_KEY === 'YOUR_APP_KEY') {
  console.error('提示：请先在 nodejs/config.js 中将 API_KEY 替换为你的 API Key。');
}

try {
  const params = {
    key: API_KEY,
    image: await imageToBase64(imagePath),
  };
  if (dpi > 0) {
    params.dpi = dpi;
  }
  if (fileSizeMin > 0) {
    params.file_size_min = fileSizeMin;
  }
  if (fileSizeMax > 0) {
    params.file_size_max = fileSizeMax;
  }

  const result = await post('image/edit', params);

  if (!isSuccess(result)) {
    console.log(`业务失败：code=${result.code}，msg=${result.msg}`);
    process.exit(1);
  }

  const image = result.data?.image ?? '';
  if (!image) {
    console.log('返回数据中没有结果图地址。');
    process.exit(0);
  }

  const saveDir = path.join('output', 'image-edit');
  try {
    const file = await downloadImage(image, saveDir, 'edit.jpg');
    console.log(`图片编辑成功：${file}`);
  } catch (error) {
    console.log(`下载失败，结果图地址：${image}`);
  }
} catch (error) {
  console.log(`调用失败：${error.message}`);
  process.exit(1);
}
