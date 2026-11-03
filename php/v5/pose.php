<?php
/**
 * V5 职业形象照检测+制作（有水印）
 *
 * 接口地址：https://api.zjzapi.com/idcardv5/make（与证件照制作相同，通过职业照规格 + 专属底色区分）
 * 请求方式：POST（application/x-www-form-urlencoded）
 * 返回格式：JSON
 * 接口费用：调用免费，结果图带水印；可使用返回的 pic_id 调用 v5/get 获取无水印证件照。
 *
 * 用法（在 php 目录运行）：php v5/pose.php
 *
 * 与 v5/make 的区别：item_id 使用职业形象照类规格（规格列表中的"职业形象照"分类），
 * 底色为职业照专属：pose_black、pose_blue、pose_darkgray、pose_gray、pose_white。
 */

require __DIR__ . '/../bootstrap.php';

// ===== 请求参数（修改后直接运行即可）=====
$image_path = __DIR__ . '/../../images/test.jpeg'; // 待处理图片路径
$item_id    = '1';          // 职业形象照规格id，可通过「证件照规格列表」示例获取：php item/list.php
$colors     = 'pose_blue';  // 底色：pose_black、pose_blue、pose_darkgray、pose_gray、pose_white，留空返回全部
$enhance    = 1;            // 人像增强：默认已支持，不需要可传0
$beauty     = 1;            // 美颜：默认已支持，不需要可传0
$whitening  = 50;           // 美颜-美白等级：0-100，默认50（beauty=1 时生效）

try {
    $result = ZjzApi::post('idcardv5/make', [
        'key' => ZJZ_DEMO_KEY,
        'item_id' => $item_id,
        'image' => demo_image_base64($image_path),
        'colors' => $colors,
        'enhance' => $enhance,
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
echo "制作成功（有水印），pic_id：{$pic_id}\n";
echo "提示：把 pic_id 填入 v5/get.php 可获取无水印证件照。\n";

// 把各底色结果图下载到 output/v5-pose/ 目录。
$save_dir = __DIR__ . '/../../output/v5-pose';
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
