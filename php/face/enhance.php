<?php
/**
 * 人像修复增强
 *
 * 接口地址：https://api.zjzapi.com/face/enhance
 * 请求方式：POST（application/x-www-form-urlencoded）
 * 返回格式：JSON
 * 接口费用：返回成功扣除0.06元。
 *
 * 用法（在 php 目录运行）：php face/enhance.php
 *
 * 适用于模糊图片清晰度修复：拍照没对焦好、图片模糊、光线差、
 * 老照片不清晰、图片像素低等；图片大小不超过4M，分辨率不超过3072px。
 */

require __DIR__ . '/../bootstrap.php';

// ===== 请求参数（修改后直接运行即可）=====
$image_path = __DIR__ . '/../../images/test.jpeg'; // 待处理图片路径
$size       = 'original'; // 图片大小：original=原始大小、x2=高清图，默认 original

try {
    $result = ZjzApi::post('face/enhance', [
        'key' => ZJZ_DEMO_KEY,
        'image' => demo_image_base64($image_path),
        'size' => $size,
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

// 把结果图下载到 output/face-enhance/ 目录。
$save_dir = __DIR__ . '/../../output/face-enhance';
try {
    $file = ZjzApi::downloadImage($image, $save_dir, 'enhance.jpg');
    echo "修复增强成功：{$file}\n";
} catch (RuntimeException $e) {
    echo "下载失败，结果图地址：{$image}\n";
}
