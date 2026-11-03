<?php
/**
 * 订单状态通知（回调接收示例）
 *
 * 本示例与其它示例不同：不是主动调用平台接口，而是部署在你的服务器上，
 * 接收平台推送的订单状态通知（订单开始审核、驳回、完成、发货时触发）。
 *
 * 通知以 POST（application/x-www-form-urlencoded）发送到提交订单时填写的 notice_url。
 *
 * 本地联调方法（在 php 目录运行内置服务器）：
 *   php -S 0.0.0.0:8080 receipt/notice.php
 * 然后把 http://你的IP:8080/receipt/notice.php 作为 receipt/submit.php 的 notice_url。
 *
 * 注意：本文件运行在 Web 环境，不加载 bootstrap.php（它仅允许命令行运行）。
 */

// 直接读取配置（只需 key，无需 HTTP 客户端）。
$config = require __DIR__ . '/../config.php';
$api_key = isset($config['key']) ? (string)$config['key'] : '';

header('Content-Type: application/json; charset=utf-8');

/**
 * 按平台约定返回处理结果：平台收到 code=0 才认为同步完成。
 */
function notice_response($code, $msg)
{
    echo json_encode(['data' => new stdClass(), 'code' => $code, 'msg' => $msg], JSON_UNESCAPED_UNICODE);
    exit;
}

// 只接受 POST 请求。
if (!isset($_SERVER['REQUEST_METHOD']) || $_SERVER['REQUEST_METHOD'] !== 'POST') {
    notice_response(-1, '仅支持POST请求');
}

$post = $_POST;

// 验证签名：平台约定 sign = md5(API Key)，不一致则拒绝接收。
if (!isset($post['sign']) || strcasecmp($post['sign'], md5($api_key)) !== 0) {
    notice_response(-1, '签名校验失败');
}

$status   = isset($post['status']) ? (int)$post['status'] : 0;
$order_no = isset($post['order_no']) ? (string)$post['order_no'] : '';

/**
 * 各状态说明：
 *   2  审核中（开始审核）
 *   3  办理完成，返回 image_receipt（回执照）、image（调整后的证件照，仅有调整时返回）、
 *      image_print（排版照，仅广东驾驶证、港澳通行证支持）
 *   4  打印发货，返回 logistics_info：{"express_name":"韵达","express_no":"1234567890"}
 *  -1  审核失败，返回 reject_msg（驳回理由）
 *  -2  取消完成
 */
switch ($status) {
    case 2:
        // TODO: 更新订单状态为"审核中"。
        error_log("[receipt notice] 订单 {$order_no} 开始审核。");
        break;

    case 3:
        // TODO: 更新订单状态为"办理完成"，保存回执照地址。
        $image_receipt = isset($post['image_receipt']) ? (string)$post['image_receipt'] : '';
        $image         = isset($post['image']) ? (string)$post['image'] : '';
        $image_print   = isset($post['image_print']) ? (string)$post['image_print'] : '';
        error_log("[receipt notice] 订单 {$order_no} 办理完成，回执照：{$image_receipt}");
        break;

    case 4:
        // TODO: 更新订单状态为"已发货"，保存物流信息。
        $logistics = isset($post['logistics_info']) ? (string)$post['logistics_info'] : '{}';
        error_log("[receipt notice] 订单 {$order_no} 已发货，物流信息：{$logistics}");
        break;

    case -1:
        // TODO: 更新订单状态为"审核驳回"，可参考 receipt/resubmit.php 修改资料后重新提交。
        $reject_msg = isset($post['reject_msg']) ? (string)$post['reject_msg'] : '';
        error_log("[receipt notice] 订单 {$order_no} 审核驳回：{$reject_msg}");
        break;

    case -2:
        // TODO: 更新订单状态为"已取消"。
        error_log("[receipt notice] 订单 {$order_no} 已取消。");
        break;

    default:
        notice_response(-1, '未知的订单状态：' . $status);
}

// 处理完成后返回 code=0，告知平台通知同步完成。
notice_response(0, '同步完成');
