<?php
/**
 * 订单取消
 *
 * 接口地址：https://api.zjzapi.com/receipt/cancel
 * 请求方式：POST（application/x-www-form-urlencoded）
 * 返回格式：JSON
 * 接口费用：调用免费
 *
 * 用法（在 php 目录运行）：php receipt/cancel.php
 *
 * 注意：订单在"待处理"状态下才可以取消。
 */

require __DIR__ . '/../bootstrap.php';

// ===== 请求参数（修改后直接运行即可）=====
$order_no = ''; // 要取消的订单号（receipt/submit.php 的返回结果）

if ($order_no === '') {
    exit("请先在参数区填写要取消的订单号 order_no。\n");
}

try {
    $result = ZjzApi::post('receipt/cancel', [
        'key' => ZJZ_DEMO_KEY,
        'order_no' => $order_no,
    ]);
} catch (RuntimeException $e) {
    exit('调用失败：' . $e->getMessage() . "\n");
}

// HTTP 成功不代表业务成功，需再检查返回 JSON 中的 code。
if (!ZjzApi::isSuccess($result)) {
    exit('业务失败：code=' . (isset($result['code']) ? $result['code'] : 'null')
        . '，msg=' . (isset($result['msg']) ? $result['msg'] : 'null') . "\n");
}

echo "订单 {$order_no} 取消成功。\n";
