<?php
/**
 * 获取账号信息
 *
 * 接口地址：https://api.zjzapi.com/user/info
 * 请求方式：POST（application/x-www-form-urlencoded）
 * 返回格式：JSON
 * 接口费用：调用免费，数据有10分钟延迟。
 *
 * 用法（在 php 目录运行）：php account/info.php
 *
 * 注意：本接口不使用 API Key，而是使用用户中心的 access_token
 * （在官网「用户中心 - 账户信息」获取），请填到 config.php 的 access_token。
 */

require __DIR__ . '/../bootstrap.php';

if (ZJZ_DEMO_ACCESS_TOKEN === '' || ZJZ_DEMO_ACCESS_TOKEN === 'YOUR_ACCESS_TOKEN') {
    exit("请先在 config.php 中将 access_token 替换为你的 access_token（官网「用户中心 - 账户信息」获取）。\n");
}

try {
    $result = ZjzApi::post('user/info', [
        'access_token' => ZJZ_DEMO_ACCESS_TOKEN,
    ]);
} catch (RuntimeException $e) {
    exit('调用失败：' . $e->getMessage() . "\n");
}

// HTTP 成功不代表业务成功，需再检查返回 JSON 中的 code。
if (!ZjzApi::isSuccess($result)) {
    exit('业务失败：code=' . (isset($result['code']) ? $result['code'] : 'null')
        . '，msg=' . (isset($result['msg']) ? $result['msg'] : 'null') . "\n");
}

$data = isset($result['data']) && is_array($result['data']) ? $result['data'] : [];

$mobile   = isset($data['mobile']) ? $data['mobile'] : '';
$money    = isset($data['money']) ? $data['money'] : '';
$app_num  = isset($data['app_num']) ? $data['app_num'] : '';
$day_make = isset($data['day_make']) ? $data['day_make'] : '';

echo "登录账号：{$mobile}\n";
echo "账户余额：{$money} 元\n";
echo "应用数量：{$app_num}\n";
echo "今日调用：{$day_make} 次\n";
