// V3 证件照检测+制作
//
// 接口地址：https://api.zjzapi.com/idcardv3/all
// 请求方式：POST（application/x-www-form-urlencoded）
// 返回格式：JSON
// 接口费用：制作成功按次扣费；检测不通过扣除检测费用；beauty=1 时美颜另行计费。
//
// 用法（在 nodejs 目录运行）：node v3/all.js

import process from 'node:process';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import { API_KEY } from '../config.js';
import { post, isSuccess, imageToBase64, downloadImage } from '../lib/zjz-api.js';

// ===== 请求参数（修改后直接运行即可）=====
const imagePath = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '../../images/test.jpeg'); // 待处理图片路径
const itemId = '1';                   // 规格id，可通过「证件照规格列表」示例获取：node item/list.js
const colors = 'white';                 // 底色，如 white/blue/red；留空返回该规格支持的全部底色
const beauty = 0;                       // 美颜：0=关闭 1=开启（开启后另行计费）

if (API_KEY === '' || API_KEY === 'YOUR_APP_KEY') {
  console.error('提示：请先在 nodejs/config.js 中将 API_KEY 替换为你的 API Key。');
}

try {
  const params = {
    key: API_KEY,
    item_id: itemId,
    image: await imageToBase64(imagePath),
    beauty,
  };
  // colors 为空时不传，由服务端返回全部底色。
  if (colors) {
    params.colors = colors;
  }

  const result = await post('idcardv3/all', params);

  // HTTP 成功不代表业务成功，需再检查返回 JSON 中的 code。
  if (!isSuccess(result)) {
    console.log(`业务失败：code=${result.code}，msg=${result.msg}`);
    process.exit(1);
  }

  const picId = result.data?.pic_id ?? '';
  console.log(`制作成功，pic_id：${picId}`);

  // 把各底色结果图下载到 output/v3-all/ 目录。
  const saveDir = path.join('output', 'v3-all');
  const list = result.data?.list ?? {};
  if (Object.keys(list).length === 0) {
    console.log('返回数据中没有图片列表。');
    process.exit(0);
  }
  for (const [color, imageUrl] of Object.entries(list)) {
    try {
      const file = await downloadImage(imageUrl, saveDir, `${color}.jpg`);
      console.log(`底色 ${color}：${file}`);
    } catch (error) {
      console.log(`底色 ${color} 下载失败：${error.message}`);
    }
  }
} catch (error) {
  console.log(`调用失败：${error.message}`);
  process.exit(1);
}
