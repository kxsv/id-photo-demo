// 订单取消
//
// 接口地址：https://api.zjzapi.com/receipt/cancel
// 请求方式：POST（application/x-www-form-urlencoded）
// 返回格式：JSON
// 接口费用：调用免费
//
// 用法（在 nodejs 目录运行）：node receipt/cancel.js
//
// 注意：订单在"待处理"状态下才可以取消。

import process from 'node:process';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import { API_KEY } from '../config.js';
import { post, isSuccess } from '../lib/zjz-api.js';

// ===== 请求参数（修改后直接运行即可）=====
const orderNo = ''; // 要取消的订单号（receipt/submit.js 的返回结果）

if (API_KEY === '' || API_KEY === 'YOUR_APP_KEY') {
  console.error('提示：请先在 nodejs/config.js 中将 API_KEY 替换为你的 API Key。');
}
if (orderNo === '') {
  console.error('请先在参数区填写要取消的订单号 order_no。');
  process.exit(1);
}

try {
  const result = await post('receipt/cancel', {
    key: API_KEY,
    order_no: orderNo,
  });

  if (!isSuccess(result)) {
    console.log(`业务失败：code=${result.code}，msg=${result.msg}`);
    process.exit(1);
  }

  console.log(`订单 ${orderNo} 取消成功。`);
} catch (error) {
  console.log(`调用失败：${error.message}`);
  process.exit(1);
}
