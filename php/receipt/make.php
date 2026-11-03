<?php
/**
 * 证件照回执检测+制作
 *
 * 接口地址：https://api.zjzapi.com/receipt/make
 * 请求方式：POST（application/x-www-form-urlencoded）
 * 返回格式：JSON
 * 接口费用：调用免费（回执提交办理完成后按单扣费）。
 *
 * 用法（在 php 目录运行）：php receipt/make.php
 *
 * 图片要求：JPG/JPEG/PNG，大小不超过5M，分辨率小于6000x6000；
 * 办理时间 8:30-20:30（节假日正常办理），预计15分钟审核完成。
 */

require __DIR__ . '/../bootstrap.php';

// ===== 请求参数（修改后直接运行即可）=====
$image_path = __DIR__ . '/../../images/test.jpeg'; // 待处理图片路径
$item_id    = '1'; // 回执类规格id，可通过「证件照规格列表」示例获取：php item/list.php

try {
    $result = ZjzApi::post('receipt/make', [
        'key' => ZJZ_DEMO_KEY,
        'item_id' => $item_id,
        'image' => demo_image_base64($image_path),
    ]);
} catch (RuntimeException $e) {
    exit('调用失败：' . $e->getMessage() . "\n");
}

// HTTP 成功不代表业务成功，需再检查返回 JSON 中的 code。
if (!ZjzApi::isSuccess($result)) {
    exit('业务失败：code=' . (isset($result['code']) ? $result['code'] : 'null')
        . '，msg=' . (isset($result['msg']) ? $result['msg'] : 'null') . "\n");
}

$pic_id = isset($result['data']['pic_id']) ? (string)$result['data']['pic_id'] : '';
echo "回执照制作成功，pic_id：{$pic_id}\n";
echo "提示：把 pic_id 填入 receipt/submit.php 提交办理订单。\n";

// 把结果图下载到 output/receipt-make/ 目录。
$save_dir = __DIR__ . '/../../output/receipt-make';
$list = isset($result['data']['list']) && is_array($result['data']['list']) ? $result['data']['list'] : [];
if (empty($list)) {
    exit("返回数据中没有图片列表。\n");
}

foreach ($list as $color => $url) {
    try {
        $file = ZjzApi::downloadImage($url, $save_dir, $color . '.jpg');
        echo "底色 {$color}：{$file}\n";
    } catch (RuntimeException $e) {
        echo "底色 {$color} 下载失败：" . $e->getMessage() . "\n";
    }
}
