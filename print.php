<?php

require 'vendor/autoload.php';
use Grafika\Grafika;
use Grafika\Color;


$print_item = ['295-413','413-579','390-567','260-378','413-531','626-413'];    //支持排版的尺寸


//============== 一寸排版 ============================
$print_bg = __DIR__. '/print-bg/print-295-413.png';     //排版背景图
$id_photo = __DIR__. '/id-photo.jpeg';                  //一寸证件照

$editor = Grafika::createEditor(['Imagick']);
$editor->open($image1, $print_bg);
$editor->open($image2, $id_photo);
$pic_size = getimagesize($id_photo);

if ($pic_size[0].'-'.$pic_size[1]  == '295-413') {//一寸排版
    $editor->blend($image1, $image2, 'normal', 1, 'top-left', 120, 180);
    $editor->blend($image1, $image2, 'normal', 1, 'top-left', 435, 180);
    $editor->blend($image1, $image2, 'normal', 1, 'top-left', 750, 180);
    $editor->blend($image1, $image2, 'normal', 1, 'top-left', 1065, 180);
    $editor->blend($image1, $image2, 'normal', 1, 'top-left', 1380, 180);
    $editor->blend($image1, $image2, 'normal', 1, 'top-left', 120, 613);
    $editor->blend($image1, $image2, 'normal', 1, 'top-left', 435, 613);
    $editor->blend($image1, $image2, 'normal', 1, 'top-left', 750, 613);
    $editor->blend($image1, $image2, 'normal', 1, 'top-left', 1065, 613);
    $editor->blend($image1, $image2, 'normal', 1, 'top-left', 1380, 613);
}elseif($pic_size[0].'-'.$pic_size[1]  == '413-579'){
    $editor->rotate($image2 ,'90',new Color('#ff0000'));
    $editor->blend($image1, $image2, 'normal', 1, 'top-left', 292, 151);
    $editor->blend($image1, $image2, 'normal', 1, 'top-left', 890, 151);
    $editor->blend($image1, $image2, 'normal', 1, 'top-left', 292, 582);
    $editor->blend($image1, $image2, 'normal', 1, 'top-left', 890, 582);
}elseif($pic_size[0].'-'.$pic_size[1]  == '390-567'){
    $editor->rotate($image2 ,'90',new Color('#ff0000'));
    $editor->blend($image1, $image2, 'normal', 1, 'top-left', 304, 172);
    $editor->blend($image1, $image2, 'normal', 1, 'top-left', 884, 172);
    $editor->blend($image1, $image2, 'normal', 1, 'top-left', 304, 574);
    $editor->blend($image1, $image2, 'normal', 1, 'top-left', 884, 574);
}elseif($pic_size[0].'-'.$pic_size[1]  == '260-378'){
    $editor->blend($image1, $image2, 'normal', 1, 'top-left', 207, 214);
    $editor->blend($image1, $image2, 'normal', 1, 'top-left', 487, 214);
    $editor->blend($image1, $image2, 'normal', 1, 'top-left', 767, 214);
    $editor->blend($image1, $image2, 'normal', 1, 'top-left', 1047, 214);
    $editor->blend($image1, $image2, 'normal', 1, 'top-left', 1327, 214);
    $editor->blend($image1, $image2, 'normal', 1, 'top-left', 207, 612);
    $editor->blend($image1, $image2, 'normal', 1, 'top-left', 487, 612);
    $editor->blend($image1, $image2, 'normal', 1, 'top-left', 767, 612);
    $editor->blend($image1, $image2, 'normal', 1, 'top-left', 1047, 612);
    $editor->blend($image1, $image2, 'normal', 1, 'top-left', 1327, 612);
}elseif($pic_size[0].'-'.$pic_size[1]  == '413-531'){
    $editor->rotate($image2 ,'90',new Color('#ff0000'));
    $editor->blend($image1, $image2, 'normal', 1, 'top-left', 81, 179);
    $editor->blend($image1, $image2, 'normal', 1, 'top-left', 632, 179);
    $editor->blend($image1, $image2, 'normal', 1, 'top-left', 1183, 179);
    $editor->blend($image1, $image2, 'normal', 1, 'top-left', 81, 612);
    $editor->blend($image1, $image2, 'normal', 1, 'top-left', 632, 612);
    $editor->blend($image1, $image2, 'normal', 1, 'top-left', 1183, 612);
}elseif($pic_size[0].'-'.$pic_size[1]  == '626-413'){
    $editor->blend($image1, $image2, 'normal', 1, 'top-left', 248, 164);
    $editor->blend($image1, $image2, 'normal', 1, 'top-left', 924, 164);
    $editor->blend($image1, $image2, 'normal', 1, 'top-left', 248, 627);
    $editor->blend($image1, $image2, 'normal', 1, 'top-left', 924, 627);
}
$save_path = __DIR__. '/res-print.jpeg'; //排版后存在的目录
$editor->save($image1, $save_path);

exit('排版完成');
