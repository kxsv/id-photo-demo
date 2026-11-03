// 证件照回执检测+制作
//
// 接口地址：https://api.zjzapi.com/receipt/make
// 请求方式：POST（application/x-www-form-urlencoded）
// 返回格式：JSON
// 接口费用：调用免费（回执提交办理完成后按单扣费）。
//
// 用法（在 nodejs 目录运行）：node receipt/make.js
//
// 图片要求：JPG/JPEG/PNG，大小不超过5M，分辨率小于6000x6000；
// 办理时间 8:30-20:30（节假日正常办理），预计15分钟审核完成。

import process from 'node:process';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import { API_KEY } from '../config.js';
import { post, isSuccess, imageToBase64, downloadImage } from '../lib/zjz-api.js';

// ===== 请求参数（修改后直接运行即可）=====
const imagePath = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '../../images/test.jpeg'); // 待处理图片路径
const itemId = '1'; // 回执类规格id，可通过「证件照规格列表」示例获取：node item/list.js

if (API_KEY === '' || API_KEY === 'YOUR_APP_KEY') {
  console.error('提示：请先在 nodejs/config.js 中将 API_KEY 替换为你的 API Key。');
}

try {
  const result = await post('receipt/make', {
    key: API_KEY,
    item_id: itemId,
    image: await imageToBase64(imagePath),
  });

  if (!isSuccess(result)) {
    console.log(`业务失败：code=${result.code}，msg=${result.msg}`);
    process.exit(1);
  }

  const picId = result.data?.pic_id ?? '';
  console.log(`回执照制作成功，pic_id：${picId}`);
  console.log('提示：把 pic_id 填入 receipt/submit.js 提交办理订单。');

  const saveDir = path.join('output', 'receipt-make');
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
