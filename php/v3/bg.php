<?php
/**
 * V3 图片检测+更换底色
 *
 * 接口地址：https://api.zjzapi.com/idcardv3/bg
 * 请求方式：POST（application/x-www-form-urlencoded）
 * 返回格式：JSON
 * 接口费用：检测并制作成功扣除0.012元，检测失败扣除0.006元；若使用"美颜"功能将多扣除0.012元。
 *
 * 用法（在 php 目录运行）：php v3/bg.php
 *
 * 与 v3/all 的区别：不按规格重排裁剪，仅对已有照片检测后更换底色；
 * 图片大小不超过4M，分辨率不超过3072px。
 */

require __DIR__ . '/../bootstrap.php';

// ===== 请求参数（修改后直接运行即可）=====
$image_path = __DIR__ . '/../../images/test.jpeg'; // 待处理图片路径
$color      = 'blue';   // 底色：颜色值（#FFFFFF）或名称（blue、white、red、tint、grey、gradient、dark_blue、dark_red），留空返回所有底色
$beauty     = 0;        // 美颜：0=关闭 1=开启（开启后另行计费）
$whitening  = 50;       // 美颜-美白等级：0-100，默认50（beauty=1 时生效）

try {
    $result = ZjzApi::post('idcardv3/bg', [
        'key' => ZJZ_DEMO_KEY,
        'image' => demo_image_base64($image_path),
        'color' => $color,
        'beauty' => $beauty,
        'whitening' => $whitening,
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
echo "更换底色成功，pic_id：{$pic_id}\n";

// 把各底色结果图下载到 output/v3-bg/ 目录。
$save_dir = __DIR__ . '/../../output/v3-bg';
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
