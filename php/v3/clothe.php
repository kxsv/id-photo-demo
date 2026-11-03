<?php
/**
 * V3 证件照换装
 *
 * 接口地址：https://api.zjzapi.com/idcardv3/clothe
 * 请求方式：POST（application/x-www-form-urlencoded）
 * 返回格式：JSON
 * 接口费用：换装成功扣除0.036元。
 *
 * 用法（在 php 目录运行）：php v3/clothe.php
 *
 * 注意：换装基于制作接口返回的 pic_id（有效期50分钟），
 * 请先把 v3/all.php 运行得到的 pic_id 填到下方参数区。
 */

require __DIR__ . '/../bootstrap.php';

// ===== 请求参数（修改后直接运行即可）=====
$pic_id = ''; // 制作接口返回的 pic_id（有效期50分钟），先运行 php v3/all.php 获取
$clothe = ''; // 服装key，见官网《服装列表》文档

if ($pic_id === '' || $clothe === '') {
    exit("请先在参数区填写 pic_id（运行 php v3/all.php 获取）与服装key clothe。\n");
}

try {
    $result = ZjzApi::post('idcardv3/clothe', [
        'key' => ZJZ_DEMO_KEY,
        'pic_id' => $pic_id,
        'clothe' => $clothe,
    ]);
} catch (RuntimeException $e) {
    exit('调用失败：' . $e->getMessage() . "\n");
}

// HTTP 成功不代表业务成功，需再检查返回 JSON 中的 code。
if (!ZjzApi::isSuccess($result)) {
    exit('业务失败：code=' . (isset($result['code']) ? $result['code'] : 'null')
        . '，msg=' . (isset($result['msg']) ? $result['msg'] : 'null') . "\n");
}

echo "换装成功。\n";

// 把换装结果图下载到 output/v3-clothe/ 目录。
$save_dir = __DIR__ . '/../../output/v3-clothe';
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
