<?php
/**
 * 获取应用信息
 *
 * 接口地址：https://api.zjzapi.com/user/app
 * 请求方式：POST（application/x-www-form-urlencoded）
 * 返回格式：JSON
 * 接口费用：调用免费，数据有10分钟延迟。
 *
 * 用法（在 php 目录运行）：php account/app.php
 *
 * 注意：本接口需要同时使用用户中心的 access_token（官网「用户中心 - 账户信息」获取，
 * 填到 config.php 的 access_token）与你的 API Key。
 */

require __DIR__ . '/../bootstrap.php';

if (ZJZ_DEMO_ACCESS_TOKEN === '' || ZJZ_DEMO_ACCESS_TOKEN === 'YOUR_ACCESS_TOKEN') {
    exit("请先在 config.php 中将 access_token 替换为你的 access_token（官网「用户中心 - 账户信息」获取）。\n");
}

try {
    $result = ZjzApi::post('user/app', [
        'access_token' => ZJZ_DEMO_ACCESS_TOKEN,
        'key' => ZJZ_DEMO_KEY,
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

// 应用权限：1=证件照v3，2=证件照v5，3=支持全部，4=证件照回执
$type_map = [1 => '证件照v3', 2 => '证件照v5', 3 => '支持全部', 4 => '证件照回执'];
$type  = isset($data['type']) ? (int)$data['type'] : 0;
$status = isset($data['status']) ? (int)$data['status'] : 0;

echo '应用名称：' . (isset($data['name']) ? $data['name'] : '') . "\n";
echo '应用权限：' . (isset($type_map[$type]) ? $type_map[$type] : $type) . "\n";
echo '应用状态：' . ($status === 1 ? '启用中' : '已停用') . "\n";
echo '调用次数：' . (isset($data['make_num']) ? $data['make_num'] : '') . "\n";
