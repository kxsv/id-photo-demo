<?php
/**
 * 证件照冲印排版（相纸联排）
 *
 * 将一张已制作好的证件照，按预设的相纸模板（images/print-bg/<相纸>/<尺寸>.png）
 * 复制到模板的各个槽位上，输出一张可直接拿去冲印的图（默认 300dpi）。
 *
 * 用法（在 php 目录运行）：
 *   php print.php                  # 使用默认输入（output/v3-diy/blue.jpg，需先运行 v3/diy.php 生成）
 *   php print.php 路径/照片.jpg     # 指定一张做好的证件照
 *
 * 依赖：PHP 的 Imagick 或 GD 扩展（自动二选一，无需任何第三方库 / composer）。
 * 本机当前使用 GD（Imagick 扩展未安装时使用 GD 分支，逻辑一致）。
 *
 * 前置：请先运行 v3/diy.php 生成 295x413 的证件照到 output/v3-diy/blue.jpg，
 *       或自行用任意制作接口产出对应尺寸的证件照后传入路径。
 */

// ===================== 请求参数（按需修改）=====================
$paper        = 'print6';     // 相纸规格目录名（images/print-bg/ 下的子目录，如 print6=6寸）
$use_text_tpl = false;        // 是否使用带文字/裁切线说明的模板（*-text.png）
$dpi          = 300;          // 输出分辨率
$id_photo     = __DIR__ . '/../output/v3-diy/blue.jpg'; // 默认输入：做好的证件照（1寸 295x413）

// 命令行可覆盖输入图：php print.php 路径/照片.jpg
if (isset($argv[1]) && $argv[1] !== '') {
    $id_photo = $argv[1];
}

// ===================== 排版配置（数据驱动，便于增删规格）=====================
// 键名        = 证件照像素尺寸（宽-高）
// rotate      = 是否需要旋转 90°（模板槽位为横放时置 true）
// size        = 槽位像素尺寸（输入图会被缩放/拉伸到该尺寸后再贴入，对制作尺寸的轻微误差更宽容）
// positions   = 各槽位左上角坐标 [x, y]（与模板图严格对应）
$layouts = [
    '295-413'   => ['rotate' => false, 'size' => [295, 413],  'positions' => [[120,180],[435,180],[750,180],[1065,180],[1380,180],[120,613],[435,613],[750,613],[1065,613],[1380,613]]],
    '260-378'   => ['rotate' => false, 'size' => [260, 378],  'positions' => [[207,214],[487,214],[767,214],[1047,214],[1327,214],[207,612],[487,612],[767,612],[1047,612],[1327,612]]],
    '413-579'   => ['rotate' => true,  'size' => [413, 579],  'positions' => [[308,179],[907,179],[308,612],[907,612]]],
    '390-567'   => ['rotate' => true,  'size' => [390, 567],  'positions' => [[320,202],[907,202],[320,612],[907,612]]],
    '413-531'   => ['rotate' => true,  'size' => [413, 531],  'positions' => [[81,179],[632,179],[1183,179],[81,612],[632,612],[1183,612]]],
    '626-413'   => ['rotate' => false, 'size' => [626, 413],  'positions' => [[261,179],[907,179],[261,612],[907,612]]],
    '649-991'   => ['rotate' => false, 'size' => [649, 991],  'positions' => [[238,107],[907,107]]],
    '1050-1499' => ['rotate' => true,  'size' => [1050,1499], 'positions' => [[148,78]]],
    // 以下三组坐标由"统一边距+20px间距"的均匀网格推算，并与模板裁切线对齐确认（模板本身未画占位框）
    '358-441'   => ['rotate' => false, 'size' => [358, 441],  'positions' => [[151,151],[529,151],[907,151],[1285,151],[151,612],[529,612],[907,612],[1285,612]]],
    '378-472'   => ['rotate' => false, 'size' => [378, 472],  'positions' => [[111,120],[509,120],[907,120],[1305,120],[111,612],[509,612],[907,612],[1305,612]]],
    '600-600'   => ['rotate' => false, 'size' => [600, 600],  'positions' => [[287,302],[907,302]]],
];

// ===================== 参数校验 =====================
if (!is_file($id_photo)) {
    exit("输入证件照不存在：{$id_photo}\n"
        . "请先运行 v3/diy.php 生成，或在命令行传入路径：php print.php 路径/照片.jpg\n");
}
if (!extension_loaded('imagick') && !extension_loaded('gd')) {
    exit("需要 PHP 的 imagick 或 gd 扩展，请先在 php.ini 中启用其中之一。\n");
}

list($w0, $h0) = getimagesize($id_photo);
$print_name = $w0 . '-' . $h0;

if (!isset($layouts[$print_name])) {
    $supported = implode(', ', array_keys($layouts));
    exit("不支持的证件照尺寸 {$print_name}（实际像素 {$w0}x{$h0}）。\n"
        . "支持的尺寸：{$supported}\n"
        . "请使用对应规格制作证件照（如 v3/diy.php 设 width_px/height_px）。\n");
}

$cfg = $layouts[$print_name];
$suffix   = $use_text_tpl ? '-text' : '';
$print_bg = __DIR__ . "/../images/print-bg/{$paper}/{$print_name}{$suffix}.png";
if (!is_file($print_bg)) {
    exit("相纸模板不存在：{$print_bg}\n"
        . "请确认 images/print-bg/{$paper}/ 下存在 {$print_name}{$suffix}.png\n");
}

$save_dir  = __DIR__ . '/../output/print';
if (!is_dir($save_dir) && !@mkdir($save_dir, 0777, true) && !is_dir($save_dir)) {
    exit("无法创建输出目录：{$save_dir}\n");
}
$save_path = $save_dir . "/{$paper}-{$print_name}.jpg";

// ===================== 渲染（Imagick 优先，GD 兜底）=====================
if (extension_loaded('imagick')) {
    renderImagick($print_bg, $id_photo, $cfg, $dpi, $save_path);
} else {
    renderGd($print_bg, $id_photo, $cfg, $save_path);
}

echo "排版完成：{$save_path}\n";

// ===================== Imagick 实现 =====================
function renderImagick($bgPath, $photoPath, $cfg, $dpi, $out)
{
    $tpl = new Imagick($bgPath);
    // 先铺一张白底，避免模板透明区在 JPEG 中变黑
    $bg = new Imagick();
    $bg->newImage($tpl->getImageWidth(), $tpl->getImageHeight(), new ImagickPixel('white'), 'png');
    $bg->compositeImage($tpl, Imagick::COMPOSITE_OVER, 0, 0);
    $tpl->clear();

    $photo = new Imagick($photoPath);
    list($sw, $sh) = $cfg['size'];
    $photo->resizeImage($sw, $sh, Imagick::FILTER_LANCZOS, 1); // 拉伸到槽位尺寸
    if ($cfg['rotate']) {
        $photo->rotateImage(new ImagickPixel('white'), 90);
    }

    foreach ($cfg['positions'] as $pos) {
        list($x, $y) = $pos;
        $bg->compositeImage($photo, Imagick::COMPOSITE_OVER, $x, $y);
    }

    $bg->setImageResolution($dpi, $dpi);
    $bg->setImageUnits(Imagick::RESOLUTION_PIXELSPERINCH);
    $bg->setImageFormat('jpeg');
    $bg->setImageCompressionQuality(95);
    $bg->writeImage($out);
    $bg->clear();
    $photo->clear();
}

// ===================== GD 实现 =====================
function renderGd($bgPath, $photoPath, $cfg, $out)
{
    $tpl = imagecreatefrompng($bgPath);
    if ($tpl === false) {
        exit("无法读取相纸模板：{$bgPath}\n");
    }
    $tw = imagesx($tpl);
    $th = imagesy($tpl);

    // 白色画布，避免模板透明区在 JPEG 中变黑
    $bg = imagecreatetruecolor($tw, $th);
    $white = imagecolorallocate($bg, 255, 255, 255);
    imagefill($bg, 0, 0, $white);
    imagecopy($bg, $tpl, 0, 0, 0, 0, $tw, $th);
    imagedestroy($tpl);

    $photo = openImageGd($photoPath);
    list($sw, $sh) = $cfg['size'];
    $tmp = imagecreatetruecolor($sw, $sh);
    imagecopyresampled($tmp, $photo, 0, 0, 0, 0, $sw, $sh, imagesx($photo), imagesy($photo));
    imagedestroy($photo);
    $photo = $tmp;

    if ($cfg['rotate']) {
        $rot = imagerotate($photo, -90, $white); // GD 逆时针为正，故 -90 表示顺时针 90°
        imagedestroy($photo);
        $photo = $rot;
    }

    foreach ($cfg['positions'] as $pos) {
        list($x, $y) = $pos;
        imagecopy($bg, $photo, $x, $y, 0, 0, imagesx($photo), imagesy($photo));
    }

    // GD 不写 DPI，但模板像素即 300dpi 物理尺寸，冲印店按像素/相纸即可。
    imagejpeg($bg, $out, 95);
    imagedestroy($bg);
    imagedestroy($photo);
}

// 按 MIME 打开图片（支持 jpeg/png/gif/webp）
function openImageGd($path)
{
    $info = getimagesize($path);
    switch ($info['mime']) {
        case 'image/jpeg': return imagecreatefromjpeg($path);
        case 'image/png':  return imagecreatefrompng($path);
        case 'image/gif':  return imagecreatefromgif($path);
        case 'image/webp': return function_exists('imagecreatefromwebp') ? imagecreatefromwebp($path) : null;
        default: exit("不支持的图片格式：{$info['mime']}\n");
    }
}
