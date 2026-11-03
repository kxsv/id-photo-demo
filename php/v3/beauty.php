<?php
/**
 * V3 证件照美颜
 *
 * 接口地址：https://api.zjzapi.com/idcardv3/beauty
 * 请求方式：POST（application/x-www-form-urlencoded）
 * 返回格式：JSON
 * 接口费用：美颜成功扣除0.012元；不支持换装后再美颜。
 *
 * 用法（在 php 目录运行）：php v3/beauty.php
 *
 * 注意：美颜基于制作接口返回的 pic_id（有效期50分钟），
 * 请先把 v3/all.php 运行得到的 pic_id 填到下方参数区。
 */

require __DIR__ . '/../bootstrap.php';

// ===== 请求参数（修改后直接运行即可）=====
$pic_id = ''; // 制作接口返回的 pic_id（有效期50分钟），先运行 php v3/all.php 获取

if ($pic_id === '') {
    exit("请先在参数区填写 pic_id（运行 php v3/all.php 获取）。\n");
}

try {
    $result = ZjzApi::post('idcardv3/beauty', [
        'key' => ZJZ_DEMO_KEY,
        'pic_id' => $pic_id,
    ]);
} catch (RuntimeException $e) {
    exit('调用失败：' . $e->getMessage() . "\n");
}

// HTTP 成功不代表业务成功，需再检查返回 JSON 中的 code。
if (!ZjzApi::isSuccess($result)) {
    exit('业务失败：code=' . (isset($result['code']) ? $result['code'] : 'null')
        . '，msg=' . (isset($result['msg']) ? $result['msg'] : 'null') . "\n");
}

echo "美颜成功。\n";

// 把美颜结果图下载到 output/v3-beauty/ 目录。
$save_dir = __DIR__ . '/../../output/v3-beauty';
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
