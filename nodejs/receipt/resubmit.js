// 订单审核驳回后再次提交
//
// 接口地址：https://api.zjzapi.com/receipt/submit（与首次提交相同接口，额外携带 order_no）
// 请求方式：POST（application/x-www-form-urlencoded）
// 返回格式：JSON
// 接口费用：订单办理完成扣费，7元每单；若选择打印邮寄服务再增加7元。
//
// 用法（在 nodejs 目录运行）：node receipt/resubmit.js
//
// 注意：驳回通知中的 reject_msg 会说明驳回原因，按要求补充/修改 param（如身份证照片 image_idcard）后，
// 携带原订单号重新提交。

import process from 'node:process';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import { API_KEY } from '../config.js';
import { post, isSuccess } from '../lib/zjz-api.js';

// ===== 请求参数（修改后直接运行即可）=====
const picId = '';     // 制作接口返回的 pic_id（有效期50分钟），先运行 node receipt/make.js 获取
const orderNo = '';   // 上次提交返回的订单号（receipt/submit.js 的返回结果）
const noticeUrl = ''; // 回调地址：你的服务器接收订单状态通知的 URL
const param = '';     // 附加参数（json字符串），按驳回要求补充，如广东驾驶证驳回后补身份证照片：
                      // {"text_mobile":"13800138000","text_city":"深圳","image_idcard":"iVBORw0KGgo....（身份证照片base64）"}

if (API_KEY === '' || API_KEY === 'YOUR_APP_KEY') {
  console.error('提示：请先在 nodejs/config.js 中将 API_KEY 替换为你的 API Key。');
}
if (picId === '') {
  console.error('请先在参数区填写 pic_id（运行 node receipt/make.js 获取）。');
  process.exit(1);
}
if (orderNo === '') {
  console.error('请先在参数区填写 order_no（上次提交的订单号）。');
  process.exit(1);
}
if (noticeUrl === '') {
  console.error('请先在参数区填写 notice_url（你的订单状态回调地址）。');
  process.exit(1);
}

try {
  const params = {
    key: API_KEY,
    pic_id: picId,
    order_no: orderNo,
    notice_url: noticeUrl,
  };
  if (param) {
    params.param = param;
  }

  const result = await post('receipt/submit', params);

  if (!isSuccess(result)) {
    console.log(`业务失败：code=${result.code}，msg=${result.msg}`);
    process.exit(1);
  }

  const newOrderNo = result.data?.order_no ?? '';
  console.log(`再次提交成功，订单号：${newOrderNo}`);
} catch (error) {
  console.log(`调用失败：${error.message}`);
  process.exit(1);
}
