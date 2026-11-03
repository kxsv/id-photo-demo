<?php
/**
 * V5 证件照检测+自定义制作（有水印）
 *
 * 接口地址：https://api.zjzapi.com/idcardv5/diy
 * 请求方式：POST（application/x-www-form-urlencoded）
 * 返回格式：JSON
 * 接口费用：调用免费，结果图带水印；可使用返回的 pic_id 调用 v5/get 获取无水印证件照。
 *
 * 用法（在 php 目录运行）：php v5/diy.php
 *
 * 与 v5/make 的区别：不使用规格id，自定义像素尺寸、底色、文件大小与 dpi。
 */

require __DIR__ . '/../bootstrap.php';

// ===== 请求参数（修改后直接运行即可）=====
$image_path    = __DIR__ . '/../../images/test.jpeg'; // 待处理图片路径
$width_px      = 295;    // 像素尺寸-宽，不超2000px
$height_px     = 413;    // 像素尺寸-高，不超2000px
$color         = 'blue'; // 底色：颜色值（#FFFFFF）或名称（blue、white、red、tint、grey、gradient、dark_blue、dark_red）
$file_size_min = 0;      // 文件最小(kb)，最小5kb，0=不限
$file_size_max = 0;      // 文件最大(kb)，最大600kb，0=不限
$dpi           = 300;    // 分辨率，默认300
$enhance       = 1;      // 人像增强：默认已支持，不需要可传0
$beauty        = 1;      // 美颜：默认已支持，不需要可传0
$whitening     = 60;     // 美颜-美白等级：0-100，默认60（beauty=1 时生效）

try {
    $params = [
        'key' => ZJZ_DEMO_KEY,
        'image' => demo_image_base64($image_path),
        'width_px' => $width_px,
        'height_px' => $height_px,
        'color' => $color,
        'dpi' => $dpi,
        'enhance' => $enhance,
        'beauty' => $beauty,
        'whitening' => $whitening,
    ];
    if ($file_size_min > 0) {
        $params['file_size_min'] = $file_size_min;
    }
    if ($file_size_max > 0) {
        $params['file_size_max'] = $file_size_max;
    }

    $result = ZjzApi::post('idcardv5/diy', $params);
} catch (RuntimeException $e) {
    exit('调用失败：' . $e->getMessage() . "\n");
}

// HTTP 成功不代表业务成功，需再检查返回 JSON 中的 code。
if (!ZjzApi::isSuccess($result)) {
    exit('业务失败：code=' . (isset($result['code']) ? $result['code'] : 'null')
        . '，msg=' . (isset($result['msg']) ? $result['msg'] : 'null') . "\n");
}

$pic_id = isset($result['data']['pic_id']) ? (string)$result['data']['pic_id'] : '';
echo "制作成功（有水印），pic_id：{$pic_id}\n";
echo "提示：把 pic_id 填入 v5/get.php 可获取无水印证件照。\n";

// 把结果图下载到 output/v5-diy/ 目录。
$save_dir = __DIR__ . '/../../output/v5-diy';
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
