<?php
/**
 * 黑白图片上色
 *
 * 接口地址：https://api.zjzapi.com/image/colorize
 * 请求方式：POST（application/x-www-form-urlencoded）
 * 返回格式：JSON
 * 接口费用：调用成功扣除0.012元，调用失败免费。
 *
 * 用法（在 php 目录运行）：php image/colorize.php
 *
 * 适用于黑白图片、老照片上色；图片大小不超过4M，分辨率不超过3072px。
 */

require __DIR__ . '/../bootstrap.php';

// ===== 请求参数（修改后直接运行即可）=====
$image_path = __DIR__ . '/../../images/test.jpeg'; // 待处理图片路径（建议使用黑白照片）

try {
    $result = ZjzApi::post('image/colorize', [
        'key' => ZJZ_DEMO_KEY,
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

$image = isset($result['data']['image']) ? (string)$result['data']['image'] : '';
if ($image === '') {
    exit("返回数据中没有结果图地址。\n");
}

// 把结果图下载到 output/image-colorize/ 目录。
$save_dir = __DIR__ . '/../../output/image-colorize';
try {
    $file = ZjzApi::downloadImage($image, $save_dir, 'colorize.jpg');
    echo "上色成功：{$file}\n";
} catch (RuntimeException $e) {
    echo "下载失败，结果图地址：{$image}\n";
}
