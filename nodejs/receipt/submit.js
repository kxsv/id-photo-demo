// 证件照回执提交
//
// 接口地址：https://api.zjzapi.com/receipt/submit
// 请求方式：POST（application/x-www-form-urlencoded）
// 返回格式：JSON
// 接口费用：订单办理完成扣费，7元每单；若选择打印邮寄服务再增加7元。
//
// 用法（在 nodejs 目录运行）：node receipt/submit.js
//
// 注意：pic_id 来自回执制作接口的返回结果（有效期50分钟）；
// 提交后订单状态变化会通过 notice_url 回调通知，参考 receipt/notice 服务。

import process from 'node:process';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import { API_KEY } from '../config.js';
import { post, isSuccess } from '../lib/zjz-api.js';

// ===== 请求参数（修改后直接运行即可）=====
const picId = '';     // 制作接口返回的 pic_id（有效期50分钟），先运行 node receipt/make.js 获取
const noticeUrl = ''; // 回调地址：你的服务器接收订单状态通知的 URL
const param = '';     // 附加参数（json字符串），具体看每个回执规格的要求，如广东驾驶证：
                      // {"text_mobile":"13800138000","text_city":"深圳","text_name":"姓名","text_idcard":"身份证号码"}
const isPrint = 0;    // 是否打印邮寄：1=是 0=否
const printParam = ''; // 打印邮寄（is_print=1 时必填），如：
                       // {"name":"张三","phone":"13800138000","province":"广东","city":"深圳","district":"南山区","address":"xx路xx小区"}

if (API_KEY === '' || API_KEY === 'YOUR_APP_KEY') {
  console.error('提示：请先在 nodejs/config.js 中将 API_KEY 替换为你的 API Key。');
}
if (picId === '') {
  console.error('请先在参数区填写 pic_id（运行 node receipt/make.js 获取）。');
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
    notice_url: noticeUrl,
  };
  if (param) {
    params.param = param;
  }
  if (isPrint === 1) {
    params.is_print = 1;
    if (printParam) {
      params.print_param = printParam;
    }
  }

  const result = await post('receipt/submit', params);

  if (!isSuccess(result)) {
    console.log(`业务失败：code=${result.code}，msg=${result.msg}`);
    process.exit(1);
  }

  const orderNo = result.data?.order_no ?? '';
  console.log(`提交成功，订单号：${orderNo}`);
  console.log('提示：订单状态变化将通过 notice_url 回调通知；如需取消订单见 receipt/cancel.js。');
} catch (error) {
  console.log(`调用失败：${error.message}`);
  process.exit(1);
}
