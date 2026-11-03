// 获取应用信息
//
// 接口地址：https://api.zjzapi.com/user/app
// 请求方式：POST（application/x-www-form-urlencoded）
// 返回格式：JSON
// 接口费用：调用免费，数据有10分钟延迟。
//
// 用法（在 nodejs 目录运行）：node account/app.js
//
// 注意：本接口需要同时使用用户中心的 access_token（官网「用户中心 - 账户信息」获取，
// 填到 config.js 的 ACCESS_TOKEN）与你的 API Key。

import process from 'node:process';
import { API_KEY, ACCESS_TOKEN } from '../config.js';
import { post, isSuccess } from '../lib/zjz-api.js';

if (ACCESS_TOKEN === '' || ACCESS_TOKEN === 'YOUR_ACCESS_TOKEN') {
  console.error('提示：请先在 nodejs/config.js 中将 ACCESS_TOKEN 替换为你的 access_token（官网「用户中心 - 账户信息」获取）。');
  process.exit(1);
}

try {
  const result = await post('user/app', {
    access_token: ACCESS_TOKEN,
    key: API_KEY,
  });

  if (!isSuccess(result)) {
    console.log(`业务失败：code=${result.code}，msg=${result.msg}`);
    process.exit(1);
  }

  const data = result.data ?? {};

  // 应用权限：1=证件照v3，2=证件照v5，3=支持全部，4=证件照回执
  const typeMap = { 1: '证件照v3', 2: '证件照v5', 3: '支持全部', 4: '证件照回执' };
  const type = Number(data.type ?? 0);
  const status = Number(data.status ?? 0);

  console.log(`应用名称：${data.name ?? ''}`);
  console.log(`应用权限：${typeMap[type] ?? type}`);
  console.log(`应用状态：${status === 1 ? '启用中' : '已停用'}`);
  console.log(`调用次数：${data.make_num ?? ''}`);
} catch (error) {
  console.log(`调用失败：${error.message}`);
  process.exit(1);
}
