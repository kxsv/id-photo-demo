<?php
/**
 * 证件照规格列表
 *
 * 接口地址：https://api.zjzapi.com/item/list
 * 请求方式：POST（application/x-www-form-urlencoded）
 * 返回格式：JSON
 * 接口费用：调用免费
 *
 * 用法（在 php 目录运行）：php item/list.php
 *
 * 输出各规格的 item_id、名称、像素尺寸与支持底色，
 * 把 item_id 和对应底色填入制作类示例（如 v3/all.php）即可。
 */

require __DIR__ . '/../bootstrap.php';

try {
    $result = ZjzApi::post('item/list', [
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

$list = isset($result['data']['list']) && is_array($result['data']['list']) ? $result['data']['list'] : [];
if (empty($list)) {
    exit("返回数据中没有规格列表。\n");
}

// 打印规格清单，便于挑选 item_id 与底色。
echo "item_id  名称            像素          支持底色\n";
foreach ($list as $item) {
    printf(
        "%-8s %-14s %-12s %s\n",
        isset($item['item_id']) ? $item['item_id'] : '',
        isset($item['name']) ? $item['name'] : '',
        (isset($item['width_px']) ? $item['width_px'] : '') . 'x' . (isset($item['height_px']) ? $item['height_px'] : ''),
        isset($item['color']) ? $item['color'] : ''
    );
}
echo '共 ' . count($list) . " 个规格。\n";
