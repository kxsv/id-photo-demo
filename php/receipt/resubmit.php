<?php
/**
 * 订单审核驳回后再次提交
 *
 * 接口地址：https://api.zjzapi.com/receipt/submit（与首次提交相同接口，额外携带 order_no）
 * 请求方式：POST（application/x-www-form-urlencoded）
 * 返回格式：JSON
 * 接口费用：订单办理完成扣费，7元每单；若选择打印邮寄服务再增加7元。
 *
 * 用法（在 php 目录运行）：php receipt/resubmit.php
 *
 * 注意：驳回通知中的 reject_msg 会说明驳回原因（见 receipt/notice.php 的 status=-1 分支），
 * 按要求补充/修改 param（如身份证照片 image_idcard）后，携带原订单号重新提交。
 */

require __DIR__ . '/../bootstrap.php';

// ===== 请求参数（修改后直接运行即可）=====
$pic_id     = ''; // 制作接口返回的 pic_id（有效期50分钟），先运行 php receipt/make.php 获取
$order_no   = ''; // 上次提交返回的订单号（receipt/submit.php 的返回结果）
$notice_url = ''; // 回调地址：你的服务器接收订单状态通知的 URL，参考 receipt/notice.php
$param      = ''; // 附加参数（json字符串），按驳回要求补充，如广东驾驶证驳回后补身份证照片：
                  // {"text_mobile":"13800138000","text_city":"深圳","image_idcard":"iVBORw0KGgo....（身份证照片base64）"}

if ($pic_id === '') {
    exit("请先在参数区填写 pic_id（运行 php receipt/make.php 获取）。\n");
}
if ($order_no === '') {
    exit("请先在参数区填写 order_no（上次提交的订单号）。\n");
}
if ($notice_url === '') {
    exit("请先在参数区填写 notice_url（你的订单状态回调地址，参考 receipt/notice.php）。\n");
}

$params = [
    'key' => ZJZ_DEMO_KEY,
    'pic_id' => $pic_id,
    'order_no' => $order_no,
    'notice_url' => $notice_url,
];
if ($param !== '') {
    $params['param'] = $param;
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

$new_order_no = isset($result['data']['order_no']) ? (string)$result['data']['order_no'] : '';
echo "再次提交成功，订单号：{$new_order_no}\n";
