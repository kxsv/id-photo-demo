<?php
/**
 * Demo 公共引导文件：加载配置与公共客户端，所有接口示例都先 require 本文件。
 */

if (PHP_SAPI !== 'cli') {
    exit("Demo 仅支持在命令行运行。\n");
}

// Guzzle 依赖（vendor 已随压缩包提供；也可自行 composer install）。
require __DIR__ . '/vendor/autoload.php';

require __DIR__ . '/lib/ZjzApi.php';

$config = require __DIR__ . '/config.php';

// 全局 API Key，所有示例通过 ZJZ_DEMO_KEY 引用。
define('ZJZ_DEMO_KEY', isset($config['key']) ? (string)$config['key'] : '');

// 用户中心 access_token，仅 account/info、account/app 示例使用。
define('ZJZ_DEMO_ACCESS_TOKEN', isset($config['access_token']) ? (string)$config['access_token'] : '');

if (ZJZ_DEMO_KEY === '' || ZJZ_DEMO_KEY === 'YOUR_APP_KEY') {
    fwrite(STDERR, "提示：请先在 config.php 中将 key 替换为你的 API Key，否则接口会返回鉴权失败。\n");
}

/**
 * 读取示例图片并转为 Base64
 *
 * @param string $image_path 图片路径，留空时使用压缩包根目录 images/test.jpeg
 * @return string
 */
function demo_image_base64($image_path = '')
{
    if ($image_path === '') {
        $image_path = __DIR__ . '/../images/test.jpeg';
    }
    return ZjzApi::imageToBase64($image_path);
}
