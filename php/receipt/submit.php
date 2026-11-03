<?php
/**
 * 证件照回执提交
 *
 * 接口地址：https://api.zjzapi.com/receipt/submit
 * 请求方式：POST（application/x-www-form-urlencoded）
 * 返回格式：JSON
 * 接口费用：订单办理完成扣费，7元每单；若选择打印邮寄服务再增加7元。
 *
 * 用法（在 php 目录运行）：php receipt/submit.php
 *
 * 注意：pic_id 来自回执制作接口的返回结果（有效期50分钟）；
 * 提交后订单状态变化会通过 notice_url 回调通知，参考 receipt/notice.php。
 */

require __DIR__ . '/../bootstrap.php';

// ===== 请求参数（修改后直接运行即可）=====
$pic_id     = ''; // 制作接口返回的 pic_id（有效期50分钟），先运行 php receipt/make.php 获取
$notice_url = ''; // 回调地址：你的服务器接收订单状态通知的 URL，参考 receipt/notice.php
$param      = ''; // 附加参数（json字符串），具体看每个回执规格的要求，如广东驾驶证：
                  // {"text_mobile":"13800138000","text_city":"深圳","text_name":"姓名","text_idcard":"身份证号码"}
$is_print   = 0;  // 是否打印邮寄：1=是 0=否
$print_param = ''; // 打印邮寄（is_print=1 时必填），如：
                   // {"name":"张三","phone":"13800138000","province":"广东","city":"深圳","district":"南山区","address":"xx路xx小区"}

if ($pic_id === '') {
    exit("请先在参数区填写 pic_id（运行 php receipt/make.php 获取）。\n");
}
if ($notice_url === '') {
    exit("请先在参数区填写 notice_url（你的订单状态回调地址，参考 receipt/notice.php）。\n");
}

$params = [
    'key' => ZJZ_DEMO_KEY,
    'pic_id' => $pic_id,
    'notice_url' => $notice_url,
];
if ($param !== '') {
    $params['param'] = $param;
}
if ($is_print === 1) {
    $params['is_print'] = 1;
    if ($print_param !== '') {
        $params['print_param'] = $print_param;
    }
}

try {
    $result = ZjzApi::post('receipt/submit', $params);
} catch (RuntimeException $e) {
    exit('调用失败：' . $e->getMessage() . "\n");
}

// HTTP 成功不代表业务成功，需再检查返回 JSON 中的 code。
if (!ZjzApi::isSuccess($result)) {
    exit('业务失败：code=' . (isset($result['code']) ? $result['code'] : 'null')
        . '，msg=' . (isset($result['msg']) ? $result['msg'] : 'null') . "\n");
}

$order_no = isset($result['data']['order_no']) ? (string)$result['data']['order_no'] : '';
echo "提交成功，订单号：{$order_no}\n";
echo "提示：订单状态变化将通过 notice_url 回调通知；如需取消订单见 receipt/cancel.php。\n";
