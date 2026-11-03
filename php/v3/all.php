<?php
/**
 * V3 证件照检测+制作
 *
 * 接口地址：https://api.zjzapi.com/idcardv3/all
 * 请求方式：POST（application/x-www-form-urlencoded）
 * 返回格式：JSON
 * 接口费用：制作成功按次扣费；检测不通过扣除检测费用；beauty=1 时美颜另行计费。
 *
 * 用法（在 php 目录运行）：php v3/all.php
 */

require __DIR__ . '/../bootstrap.php';

// ===== 请求参数（修改后直接运行即可）=====
$image_path = __DIR__ . '/../../images/test.jpeg'; // 待处理图片路径
$item_id    = '1';    // 规格id，可通过「证件照规格列表」示例获取：php item/list.php
$colors     = 'white'; // 底色，如 white/blue/red；留空返回该规格支持的全部底色
$beauty     = 0;       // 美颜：0=关闭 1=开启（开启后另行计费）

try {
    $result = ZjzApi::post('idcardv3/all', [
        'key' => ZJZ_DEMO_KEY,
        'item_id' => $item_id,
        'image' => demo_image_base64($image_path),
        'colors' => $colors,
        'beauty' => $beauty,
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
echo "制作成功，pic_id：{$pic_id}\n";

// 把各底色结果图下载到 output/v3-all/ 目录。
$save_dir = __DIR__ . '/../../output/v3-all';
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
