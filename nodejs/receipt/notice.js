// 订单状态通知（回调接收示例）
//
// 本示例与其它示例不同：不是主动调用平台接口，而是部署在你的服务器上，
// 接收平台推送的订单状态通知（订单开始审核、驳回、完成、发货时触发）。
//
// 通知以 POST（application/x-www-form-urlencoded）发送到提交订单时填写的 notice_url。
//
// 本地联调方法（在 nodejs 目录运行）：
//   node receipt/notice.js
// 然后把 http://你的IP:8080/ 作为 receipt/submit.js 的 noticeUrl。
//
// 注意：本文件是独立 HTTP 服务，不依赖 lib/zjz-api.js，仅用于本地联调接收平台推送。

import http from 'node:http';
import crypto from 'node:crypto';
import { API_KEY } from '../config.js';

const PORT = 8080;

/**
 * 按平台约定返回处理结果：平台收到 code=0 才认为同步完成。
 */
function buildResponse(code, msg) {
  return JSON.stringify({ data: {}, code, msg });
}

const server = http.createServer((req, res) => {
  res.setHeader('Content-Type', 'application/json; charset=utf-8');

  if (req.method !== 'POST') {
    // 方便联调时做健康检查（平台只会 POST）。
    res.end(buildResponse(0, 'ok'));
    return;
  }

  let body = '';
  req.on('data', (chunk) => {
    body += chunk;
  });
  req.on('end', () => {
    const post = new URLSearchParams(body);

    // 验证签名：平台约定 sign = md5(API_KEY)，不一致则拒绝接收。
    const expectedSign = crypto.createHash('md5').update(API_KEY).digest('hex');
    if ((post.get('sign') || '').toLowerCase() !== expectedSign) {
      res.end(buildResponse(-1, '签名校验失败'));
      return;
    }

    const status = parseInt(post.get('status') || '0', 10);
    const orderNo = post.get('order_no') || '';

    // 各状态说明：
    //   2  审核中（开始审核）
    //   3  办理完成，返回 image_receipt / image / image_print
    //   4  打印发货，返回 logistics_info
    //  -1  审核失败，返回 reject_msg（驳回理由）
    //  -2  取消完成
    switch (status) {
      case 2:
        // TODO: 更新订单状态为"审核中"。
        console.log(`[receipt notice] 订单 ${orderNo} 开始审核。`);
        break;
      case 3:
        // TODO: 更新订单状态为"办理完成"，保存回执照地址。
        console.log(`[receipt notice] 订单 ${orderNo} 办理完成，回执照：${post.get('image_receipt') || ''}`);
        break;
      case 4:
        // TODO: 更新订单状态为"已发货"，保存物流信息。
        console.log(`[receipt notice] 订单 ${orderNo} 已发货，物流信息：${post.get('logistics_info') || '{}'}`);
        break;
      case -1:
        // TODO: 更新订单状态为"审核驳回"，可参考 receipt/resubmit.js 修改资料后重新提交。
        console.log(`[receipt notice] 订单 ${orderNo} 审核驳回：${post.get('reject_msg') || ''}`);
        break;
      case -2:
        // TODO: 更新订单状态为"已取消"。
        console.log(`[receipt notice] 订单 ${orderNo} 已取消。`);
        break;
      default:
        res.end(buildResponse(-1, `未知的订单状态：${status}`));
        return;
    }

    // 处理完成后返回 code=0，告知平台通知同步完成。
    res.end(buildResponse(0, '同步完成'));
  });
});

server.listen(PORT, '0.0.0.0', () => {
  if (API_KEY === '' || API_KEY === 'YOUR_APP_KEY') {
    console.error('提示：请先在 nodejs/config.js 中将 API_KEY 替换为你的 API Key。');
  }
  console.log(`回执通知回调服务已启动：http://0.0.0.0:${PORT}/`);
  console.log('把它作为 receipt/submit.js 的 noticeUrl。按 Ctrl+C 停止。');
});
