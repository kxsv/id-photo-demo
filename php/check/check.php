<?php
/**
 * 证件照合规性检测
 *
 * 接口地址：https://api.zjzapi.com/idphoto/check
 * 请求方式：POST（application/x-www-form-urlencoded）
 * 返回格式：JSON
 * 接口费用：检测成功扣除0.012元。
 *
 * 用法（在 php 目录运行）：php check/check.php
 *
 * 提供规格id（item_id）即可按该规格检测；也可不传 item_id，
 * 自定义像素、dpi、文件大小、底色、格式以及各项人像构图/姿态/画质/遮挡/着装检测项，
 * 完整参数列表见官网《证件照合规性检测》文档。
 */

require __DIR__ . '/../bootstrap.php';

// ===== 请求参数（修改后直接运行即可）=====
$image_path = __DIR__ . '/../../images/test.jpeg'; // 待检测图片路径（不超过4M，分辨率不超过3072px）
$item_id    = '1'; // 规格id，提供后可省略 width_px/height_px/dpi 等规格参数；填 0 表示使用下方自定义规格参数

try {
    $params = [
        'key' => ZJZ_DEMO_KEY,
        'image' => demo_image_base64($image_path),
    ];

    if ((int)$item_id > 0) {
        // 方式一：按规格id检测（推荐，规格要求自动带上）。
        $params['item_id'] = $item_id;
    } else {
        // 方式二：自定义检测规格（不传 item_id 时生效），按需增删。
        $params['width_px']  = 295;   // 证件照宽（px）
        $params['height_px'] = 413;   // 证件照高（px）
        $params['dpi']       = 300;   // 图片分辨率，常见300、350
        $params['file_format'] = 'jpg'; // 照片格式：jpg、jpeg、png
        $params['bg_colors'] = 'white,blue'; // 允许的底色，逗号分隔
        $params['solid_bg']  = 1;     // 纯色背景检测：0/1
        // 人像构图示例（0.01-1.0 的比例值）：
        $params['face_width_ratio_min'] = 0.4;
        $params['face_width_ratio_max'] = 0.6;
        // 姿态与完整性示例（单位：倾斜角度，推荐值见文档）：
        $params['shoulder_horizontal'] = 5;
        $params['face_pose_yaw'] = 15;
        // 画质与光照示例（推荐值见文档）：
        $params['face_clarity']   = 0.5;
        $params['face_brightness'] = 0.3;
        // 遮挡与着装示例（0=不检测 1=检测）：
        $params['face_mask'] = 1;
        $params['glasses']   = 1;
    }

    $result = ZjzApi::post('idphoto/check', $params);
} catch (RuntimeException $e) {
    exit('调用失败：' . $e->getMessage() . "\n");
}

// HTTP 成功不代表业务成功，需再检查返回 JSON 中的 code。
if (!ZjzApi::isSuccess($result)) {
    exit('业务失败：code=' . (isset($result['code']) ? $result['code'] : 'null')
        . '，msg=' . (isset($result['msg']) ? $result['msg'] : 'null') . "\n");
}

$passed_status = isset($result['data']['passed_status']) ? (int)$result['data']['passed_status'] : -1;
echo $passed_status === 1 ? "总检测结果：通过\n" : "总检测结果：不通过\n";

// 打印各检测项结果。
$passed_list = isset($result['data']['passed_list']) && is_array($result['data']['passed_list'])
    ? $result['data']['passed_list'] : [];
if (empty($passed_list)) {
    exit("返回数据中没有检测项列表。\n");
}

echo "检测项        结果    提示        检测值\n";
foreach ($passed_list as $item) {
    $actual = '';
    if (isset($item['detail']['actual'])) {
        $actual = $item['detail']['actual'];
        if (isset($item['detail']['unit'])) {
            $actual .= $item['detail']['unit'];
        }
    }
    printf(
        "%-14s %-6s %-10s %s\n",
        isset($item['key']) ? $item['key'] : '',
        (isset($item['passed']) && (int)$item['passed'] === 1) ? '通过' : '不通过',
        isset($item['msg']) ? $item['msg'] : '',
        $actual
    );
}
