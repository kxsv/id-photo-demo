// V3 职业形象照检测+制作
//
// 接口地址：https://api.zjzapi.com/idcardv3/all（与证件照制作相同，通过职业照规格 + 专属底色区分）
// 请求方式：POST（application/x-www-form-urlencoded）
// 返回格式：JSON
// 接口费用：检测并制作成功扣除0.012元，检测失败扣除0.006元；若使用"美颜"功能将多扣除0.012元。
//
// 用法（在 nodejs 目录运行）：node v3/pose.js
//
// 与 v3/all 的区别：item_id 使用职业形象照类规格（规格列表中的"职业形象照"分类），
// 底色为职业照专属：pose_black、pose_blue、pose_darkgray、pose_gray、pose_white。

import process from 'node:process';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import { API_KEY } from '../config.js';
import { post, isSuccess, imageToBase64, downloadImage } from '../lib/zjz-api.js';

// ===== 请求参数（修改后直接运行即可）=====
const imagePath = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '../../images/test.jpeg'); // 待处理图片路径
const itemId = '1';        // 职业形象照规格id，可通过「证件照规格列表」示例获取：node item/list.js
const colors = 'pose_blue'; // 底色：pose_black、pose_blue、pose_darkgray、pose_gray、pose_white，留空返回全部
const beauty = 0;          // 美颜：0=关闭 1=开启（开启后另行计费）

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

  if (!isSuccess(result)) {
    console.log(`业务失败：code=${result.code}，msg=${result.msg}`);
    process.exit(1);
  }

  const picId = result.data?.pic_id ?? '';
  console.log(`制作成功，pic_id：${picId}`);

  const saveDir = path.join('output', 'v3-pose');
  const list = result.data?.list ?? {};
  if (Object.keys(list).length === 0) {
    console.log('返回数据中没有图片列表。');
    process.exit(0);
  }
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
