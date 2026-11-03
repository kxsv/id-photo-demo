// 获取账号信息
//
// 接口地址：https://api.zjzapi.com/user/info
// 请求方式：POST（application/x-www-form-urlencoded）
// 返回格式：JSON
// 接口费用：调用免费，数据有10分钟延迟。
//
// 用法（在 nodejs 目录运行）：node account/info.js
//
// 注意：本接口不使用 API Key，而是使用用户中心的 access_token
// （在官网「用户中心 - 账户信息」获取），请填到 config.js 的 ACCESS_TOKEN。

import process from 'node:process';
import { API_KEY, ACCESS_TOKEN } from '../config.js';
import { post, isSuccess } from '../lib/zjz-api.js';

if (ACCESS_TOKEN === '' || ACCESS_TOKEN === 'YOUR_ACCESS_TOKEN') {
  console.error('提示：请先在 nodejs/config.js 中将 ACCESS_TOKEN 替换为你的 access_token（官网「用户中心 - 账户信息」获取）。');
  process.exit(1);
}

try {
  const result = await post('user/info', {
    access_token: ACCESS_TOKEN,
  });

  if (!isSuccess(result)) {
    console.log(`业务失败：code=${result.code}，msg=${result.msg}`);
    process.exit(1);
  }

  const data = result.data ?? {};
  console.log(`登录账号：${data.mobile ?? ''}`);
  console.log(`账户余额：${data.money ?? ''} 元`);
  console.log(`应用数量：${data.app_num ?? ''}`);
  console.log(`今日调用：${data.day_make ?? ''} 次`);
} catch (error) {
  console.log(`调用失败：${error.message}`);
  process.exit(1);
}
