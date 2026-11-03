// V5 获取证件照（无水印）
//
// 接口地址：https://api.zjzapi.com/idcardv5/get
// 请求方式：POST（application/x-www-form-urlencoded）
// 返回格式：JSON
// 接口费用：获取成功扣除0.6元，重复获取不扣费。
//
// 用法（在 nodejs 目录运行）：node v5/get.js
//
// 注意：pic_id 来自制作类接口（v5/make、v5/bg、v5/diy、v5/pose 等）的返回结果，
// 有效期50分钟；若换装过，请同时填写换装时使用的 clothe。

import process from 'node:process';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import { API_KEY } from '../config.js';
import { post, isSuccess, downloadImage } from '../lib/zjz-api.js';

// ===== 请求参数（修改后直接运行即可）=====
const picId = '';          // 制作接口返回的 pic_id（有效期50分钟），先运行 node v5/make.js 获取
const clothe = '';         // 换装时使用的服装key，未换装留空
const giftHd = 0;          // 1=赠送一张三寸高清照（仅规格 width_px<500 支持；更换底色、自定义制作暂不支持）
const giftHdColor = '';    // 高清照底色，留空默认返回制作结果第一个底色

if (API_KEY === '' || API_KEY === 'YOUR_APP_KEY') {
  console.error('提示：请先在 nodejs/config.js 中将 API_KEY 替换为你的 API Key。');
}
if (picId === '') {
  console.error('请先在参数区填写 pic_id（运行 node v5/make.js 获取）。');
  process.exit(1);
}

try {
  const params = {
    key: API_KEY,
    pic_id: picId,
    gift_hd: giftHd,
  };
  if (clothe) {
    params.clothe = clothe;
  }
  if (giftHdColor) {
    params.gift_hd_color = giftHdColor;
  }

  const result = await post('idcardv5/get', params);

  if (!isSuccess(result)) {
    console.log(`业务失败：code=${result.code}，msg=${result.msg}`);
    process.exit(1);
  }

  console.log('获取无水印证件照成功。');

  const saveDir = path.join('output', 'v5-get');
  const list = result.data?.list ?? {};
  for (const [c, imageUrl] of Object.entries(list)) {
    try {
      const file = await downloadImage(imageUrl, saveDir, `${c}.jpg`);
      console.log(`底色 ${c}：${file}`);
    } catch (error) {
      console.log(`底色 ${c} 下载失败：${error.message}`);
    }
  }

  // 赠送的三寸高清照（开启 gift_hd 时才有）。
  const giftHdUrl = result.data?.gift_hd ?? '';
  if (giftHdUrl) {
    try {
      const file = await downloadImage(giftHdUrl, saveDir, 'gift_hd.jpg');
      console.log(`高清照：${file}`);
    } catch (error) {
      console.log(`高清照下载失败：${error.message}`);
    }
  }
} catch (error) {
  console.log(`调用失败：${error.message}`);
  process.exit(1);
}
