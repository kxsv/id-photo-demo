<?php
/**
 * 图片编辑
 *
 * 接口地址：https://api.zjzapi.com/image/edit
 * 请求方式：POST（application/x-www-form-urlencoded）
 * 返回格式：JSON
 * 接口费用：调用免费（账户余额最低需有10元才可调用）。
 *
 * 用法（在 php 目录运行）：php image/edit.php
 *
 * 修改图片像素（dpi）与文件大小（kb）；图片大小不超过1M，分辨率不超过1024px。
 */

require __DIR__ . '/../bootstrap.php';

// ===== 请求参数（修改后直接运行即可）=====
$image_path    = __DIR__ . '/../../images/test.jpeg'; // 待处理图片路径
$dpi           = 300; // 图片dpi，支持72,96,150,200,300,350,500,1000；填0表示不修改
$file_size_min = 0;   // 文件最小(kb)，最小5kb；填0表示不限制
$file_size_max = 0;   // 文件最大(kb)，最大300kb；填0表示不限制

try {
    $params = [
        'key' => ZJZ_DEMO_KEY,
        'image' => demo_image_base64($image_path),
    ];
    if ($dpi > 0) {
        $params['dpi'] = $dpi;
    }
    if ($file_size_min > 0) {
        $params['file_size_min'] = $file_size_min;
    }
    if ($file_size_max > 0) {
        $params['file_size_max'] = $file_size_max;
    }

    $result = ZjzApi::post('image/edit', $params);
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

// 把结果图下载到 output/image-edit/ 目录。
$save_dir = __DIR__ . '/../../output/image-edit';
try {
    $file = ZjzApi::downloadImage($image, $save_dir, 'edit.jpg');
    echo "图片编辑成功：{$file}\n";
} catch (RuntimeException $e) {
    echo "下载失败，结果图地址：{$image}\n";
}
