// V3 证件照换装
//
// 接口地址：https://api.zjzapi.com/idcardv3/clothe
// 请求方式：POST（application/x-www-form-urlencoded）
// 返回格式：JSON
// 接口费用：换装成功扣除0.036元。
//
// 用法（在 nodejs 目录运行）：node v3/clothe.js
//
// 注意：换装基于制作接口返回的 pic_id（有效期50分钟），
// 请先把 v3/all.js 运行得到的 pic_id 填到下方参数区。

import process from 'node:process';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import { API_KEY } from '../config.js';
import { post, isSuccess, downloadImage } from '../lib/zjz-api.js';

// ===== 请求参数（修改后直接运行即可）=====
const picId = ''; // 制作接口返回的 pic_id（有效期50分钟），先运行 node v3/all.js 获取
const clothe = ''; // 服装key，见官网《服装列表》文档

if (API_KEY === '' || API_KEY === 'YOUR_APP_KEY') {
  console.error('提示：请先在 nodejs/config.js 中将 API_KEY 替换为你的 API Key。');
}
if (picId === '' || clothe === '') {
  console.error('请先在参数区填写 pic_id（运行 node v3/all.js 获取）与服装key clothe。');
  process.exit(1);
}

try {
  const result = await post('idcardv3/clothe', {
    key: API_KEY,
    pic_id: picId,
    clothe,
  });

  if (!isSuccess(result)) {
    console.log(`业务失败：code=${result.code}，msg=${result.msg}`);
    process.exit(1);
  }

  console.log('换装成功。');

  const saveDir = path.join('output', 'v3-clothe');
  const list = result.data?.list ?? {};
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
