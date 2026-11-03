<?php
/**
 * V5 获取证件照（无水印）
 *
 * 接口地址：https://api.zjzapi.com/idcardv5/get
 * 请求方式：POST（application/x-www-form-urlencoded）
 * 返回格式：JSON
 * 接口费用：获取成功扣除0.6元，重复获取不扣费。
 *
 * 用法（在 php 目录运行）：php v5/get.php
 *
 * 注意：pic_id 来自制作类接口（v5/make、v5/bg、v5/diy、v5/pose 等）的返回结果，
 * 有效期50分钟；若换装过，请同时填写换装时使用的 clothe。
 */

require __DIR__ . '/../bootstrap.php';

// ===== 请求参数（修改后直接运行即可）=====
$pic_id         = '';  // 制作接口返回的 pic_id（有效期50分钟），先运行 php v5/make.php 获取
$clothe         = '';  // 换装时使用的服装key，未换装留空
$gift_hd        = 0;   // 1=赠送一张三寸高清照（仅规格 width_px<500 支持；更换底色、自定义制作暂不支持）
$gift_hd_color  = '';  // 高清照底色，留空默认返回制作结果第一个底色

if ($pic_id === '') {
    exit("请先在参数区填写 pic_id（运行 php v5/make.php 获取）。\n");
}

$params = [
    'key' => ZJZ_DEMO_KEY,
    'pic_id' => $pic_id,
    'gift_hd' => $gift_hd,
];
if ($clothe !== '') {
    $params['clothe'] = $clothe;
}
if ($gift_hd_color !== '') {
    $params['gift_hd_color'] = $gift_hd_color;
}

try {
    $result = ZjzApi::post('idcardv5/get', $params);
} catch (RuntimeException $e) {
    exit('调用失败：' . $e->getMessage() . "\n");
}

// HTTP 成功不代表业务成功，需再检查返回 JSON 中的 code。
if (!ZjzApi::isSuccess($result)) {
    exit('业务失败：code=' . (isset($result['code']) ? $result['code'] : 'null')
        . '，msg=' . (isset($result['msg']) ? $result['msg'] : 'null') . "\n");
}

echo "获取无水印证件照成功。\n";

// 把结果图下载到 output/v5-get/ 目录。
$save_dir = __DIR__ . '/../../output/v5-get';
$list = isset($result['data']['list']) && is_array($result['data']['list']) ? $result['data']['list'] : [];
foreach ($list as $color => $url) {
    try {
        $file = ZjzApi::downloadImage($url, $save_dir, $color . '.jpg');
        echo "底色 {$color}：{$file}\n";
    } catch (RuntimeException $e) {
        echo "底色 {$color} 下载失败：" . $e->getMessage() . "\n";
    }
}

// 赠送的三寸高清照（开启 gift_hd 时才有）。
if (isset($result['data']['gift_hd']) && is_string($result['data']['gift_hd']) && $result['data']['gift_hd'] !== '') {
    try {
        $file = ZjzApi::downloadImage($result['data']['gift_hd'], $save_dir, 'gift_hd.jpg');
        echo "高清照：{$file}\n";
    } catch (RuntimeException $e) {
        echo "高清照下载失败：" . $e->getMessage() . "\n";
    }
}
