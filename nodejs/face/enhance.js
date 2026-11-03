// 人像修复增强
//
// 接口地址：https://api.zjzapi.com/face/enhance
// 请求方式：POST（application/x-www-form-urlencoded）
// 返回格式：JSON
// 接口费用：返回成功扣除0.06元。
//
// 用法（在 nodejs 目录运行）：node face/enhance.js
//
// 适用于模糊图片清晰度修复：拍照没对焦好、图片模糊、光线差、
// 老照片不清晰、图片像素低等；图片大小不超过4M，分辨率不超过3072px。

import process from 'node:process';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import { API_KEY } from '../config.js';
import { post, isSuccess, imageToBase64, downloadImage } from '../lib/zjz-api.js';

// ===== 请求参数（修改后直接运行即可）=====
const imagePath = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '../../images/test.jpeg'); // 待处理图片路径
const size = 'original'; // 图片大小：original=原始大小、x2=高清图，默认 original

if (API_KEY === '' || API_KEY === 'YOUR_APP_KEY') {
  console.error('提示：请先在 nodejs/config.js 中将 API_KEY 替换为你的 API Key。');
}

try {
  const result = await post('face/enhance', {
    key: API_KEY,
    image: await imageToBase64(imagePath),
    size,
  });

  if (!isSuccess(result)) {
    console.log(`业务失败：code=${result.code}，msg=${result.msg}`);
    process.exit(1);
  }

  const image = result.data?.image ?? '';
  if (!image) {
    console.log('返回数据中没有结果图地址。');
    process.exit(0);
  }

  const saveDir = path.join('output', 'face-enhance');
  try {
    const file = await downloadImage(image, saveDir, 'enhance.jpg');
    console.log(`修复增强成功：${file}`);
  } catch (error) {
    console.log(`下载失败，结果图地址：${image}`);
  }
} catch (error) {
  console.log(`调用失败：${error.message}`);
  process.exit(1);
}
