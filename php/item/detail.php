<?php
/**
 * 证件照规格详情
 *
 * 接口地址：https://api.zjzapi.com/item/get
 * 请求方式：POST（application/x-www-form-urlencoded）
 * 返回格式：JSON
 * 接口费用：调用免费
 *
 * 用法（在 php 目录运行）：php item/detail.php
 *
 * 输出单个规格的像素尺寸、毫米尺寸、dpi、文件大小要求与支持底色，
 * 用于制作类示例（如 v3/all.php）选择 item_id 与底色。
 */

require __DIR__ . '/../bootstrap.php';

// ===== 请求参数（修改后直接运行即可）=====
$item_id = '1'; // 规格id，可通过「证件照规格列表」示例获取：php item/list.php

try {
    $result = ZjzApi::post('item/get', [
        'key' => ZJZ_DEMO_KEY,
        'item_id' => $item_id,
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
if (empty($data)) {
    exit("返回数据中没有规格详情。\n");
}

// 打印规格详情。
$keys = [
    'id' => '规格id',
    'name' => '名称',
    'width_px' => '像素-宽',
    'height_px' => '像素-高',
    'width_mm' => '毫米-宽',
    'height_mm' => '毫米-高',
    'dpi' => '分辨率',
    'file_size_min' => '文件最小(kb)',
    'file_size_max' => '文件最大(kb)',
    'file_size_msg' => '文件大小说明',
    'color' => '支持底色',
    'category' => '分类',
    'category_id' => '分类id',
];
foreach ($keys as $key => $label) {
    $value = isset($data[$key]) ? (is_array($data[$key]) ? json_encode($data[$key], JSON_UNESCAPED_UNICODE) : $data[$key]) : '';
    echo "{$label}：{$value}\n";
}
