// 证件照规格列表
//
// 接口地址：https://api.zjzapi.com/item/list
// 请求方式：POST（application/x-www-form-urlencoded）
// 返回格式：JSON
// 接口费用：调用免费
//
// 用法（在 nodejs 目录运行）：node item/list.js
//
// 输出各规格的 item_id、名称、像素尺寸与支持底色，
// 把 item_id 和对应底色填入制作类示例（如 v3/all.js）即可。

import process from 'node:process';
import { API_KEY } from '../config.js';
import { post, isSuccess } from '../lib/zjz-api.js';

if (API_KEY === '' || API_KEY === 'YOUR_APP_KEY') {
  console.error('提示：请先在 nodejs/config.js 中将 API_KEY 替换为你的 API Key。');
}

try {
  const result = await post('item/list', { key: API_KEY });

  // HTTP 成功不代表业务成功，需再检查返回 JSON 中的 code。
  if (!isSuccess(result)) {
    console.log(`业务失败：code=${result.code}，msg=${result.msg}`);
    process.exit(1);
  }

  const list = result.data?.list ?? [];
  if (list.length === 0) {
    console.log('返回数据中没有规格列表。');
    process.exit(0);
  }

  // 打印规格清单，便于挑选 item_id 与底色。
  console.log('item_id  名称            像素          支持底色');
  for (const item of list) {
    const size = `${item.width_px ?? ''}x${item.height_px ?? ''}`;
    console.log(
      `${String(item.item_id ?? '').padEnd(8)} ${String(item.name ?? '').padEnd(14)} ${size.padEnd(12)} ${item.color ?? ''}`,
    );
  }
  console.log(`共 ${list.length} 个规格。`);
} catch (error) {
  console.log(`调用失败：${error.message}`);
  process.exit(1);
}
