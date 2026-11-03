// 证件照规格详情
//
// 接口地址：https://api.zjzapi.com/item/get
// 请求方式：POST（application/x-www-form-urlencoded）
// 返回格式：JSON
// 接口费用：调用免费
//
// 用法（在 nodejs 目录运行）：node item/detail.js
//
// 输出单个规格的像素尺寸、毫米尺寸、dpi、文件大小要求与支持底色，
// 用于制作类示例（如 v3/all.js）选择 item_id 与底色。

import process from 'node:process';
import { API_KEY } from '../config.js';
import { post, isSuccess } from '../lib/zjz-api.js';

// ===== 请求参数（修改后直接运行即可）=====
const itemId = '1'; // 规格id，可通过「证件照规格列表」示例获取：node item/list.js

if (API_KEY === '' || API_KEY === 'YOUR_APP_KEY') {
  console.error('提示：请先在 nodejs/config.js 中将 API_KEY 替换为你的 API Key。');
}

try {
  const result = await post('item/get', {
    key: API_KEY,
    item_id: itemId,
  });

  if (!isSuccess(result)) {
    console.log(`业务失败：code=${result.code}，msg=${result.msg}`);
    process.exit(1);
  }

  const data = result.data ?? {};
  if (Object.keys(data).length === 0) {
    console.log('返回数据中没有规格详情。');
    process.exit(0);
  }

  const keys = [
    ['id', '规格id'],
    ['name', '名称'],
    ['width_px', '像素-宽'],
    ['height_px', '像素-高'],
    ['width_mm', '毫米-宽'],
    ['height_mm', '毫米-高'],
    ['dpi', '分辨率'],
    ['file_size_min', '文件最小(kb)'],
    ['file_size_max', '文件最大(kb)'],
    ['file_size_msg', '文件大小说明'],
    ['color', '支持底色'],
    ['category', '分类'],
    ['category_id', '分类id'],
  ];
  for (const [key, label] of keys) {
    const value = data[key];
    const text = value === undefined || value === null ? '' : (Array.isArray(value) ? JSON.stringify(value) : value);
    console.log(`${label}：${text}`);
  }
} catch (error) {
  console.log(`调用失败：${error.message}`);
  process.exit(1);
}
