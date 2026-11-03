// V5 证件照检测+制作（有水印）
//
// 接口地址：https://api.zjzapi.com/idcardv5/make
// 请求方式：POST（application/x-www-form-urlencoded）
// 返回格式：JSON
// 接口费用：调用免费，结果图带水印；可使用返回的 pic_id 调用 v5/get 获取无水印证件照。
//
// 用法（在 nodejs 目录运行）：node v5/make.js
//
// 输出 pic_id 后，可把它填入 v5/get.js 获取无水印图（获取成功扣除0.6元，重复获取不扣费）。

import process from 'node:process';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import { API_KEY } from '../config.js';
import { post, isSuccess, imageToBase64, downloadImage } from '../lib/zjz-api.js';

// ===== 请求参数（修改后直接运行即可）=====
const imagePath = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '../../images/test.jpeg'); // 待处理图片路径
const itemId = '1';        // 规格id，可通过「证件照规格列表」示例获取：node item/list.js
const colors = 'white';    // 底色，如 white/blue/red；留空返回该规格支持的全部底色
const enhance = 1;         // 人像增强：默认已支持，不需要可传0
const beauty = 1;          // 美颜：默认已支持，不需要可传0
const whitening = 50;      // 美颜-美白等级：0-100，默认50（beauty=1 时生效）

if (API_KEY === '' || API_KEY === 'YOUR_APP_KEY') {
  console.error('提示：请先在 nodejs/config.js 中将 API_KEY 替换为你的 API Key。');
}

try {
  const params = {
    key: API_KEY,
    item_id: itemId,
    image: await imageToBase64(imagePath),
    enhance,
    beauty,
    whitening,
  };
  if (colors) {
    params.colors = colors;
  }

  const result = await post('idcardv5/make', params);

  if (!isSuccess(result)) {
    console.log(`业务失败：code=${result.code}，msg=${result.msg}`);
    process.exit(1);
  }

  const picId = result.data?.pic_id ?? '';
  console.log(`制作成功（有水印），pic_id：${picId}`);
  console.log('提示：把 pic_id 填入 v5/get.js 可获取无水印证件照。');

  const saveDir = path.join('output', 'v5-make');
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
