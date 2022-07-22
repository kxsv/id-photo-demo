<?php

require 'vendor/autoload.php';
use Grafika\Grafika;
use Grafika\Color;

//支持排版的尺寸
$print_item = ['260-378','295-413','390-567','413-531','413-579','626-413','649-991','1050-1499'];


$id_photo = __DIR__.'/../images/test.jpeg';                                     //一寸证件照
list($width_px,$height_px) = getimagesize($id_photo);
$print_name = $width_px.'-'.$height_px;

if(!in_array($print_name,$print_item)){
    exit('不支持排版');
}


$print_bg = __DIR__. "/../images/print-bg/{$print_name}.png";                 //获取排版背景图


try {
    $editor = Grafika::createEditor(['Imagick']);
    $editor->open($image1, $print_bg);
    $editor->open($image2, $id_photo);

    if ($print_name  == '295-413') {//一寸排版
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
    }elseif($print_name  == '413-579'){
        $editor->rotate($image2, '90', new Color('#ff0000'));
        $editor->blend($image1, $image2, 'normal', 1, 'top-left', 308, 179);
        $editor->blend($image1, $image2, 'normal', 1, 'top-left', 907, 179);
        $editor->blend($image1, $image2, 'normal', 1, 'top-left', 308, 612);
        $editor->blend($image1, $image2, 'normal', 1, 'top-left', 907, 612);
    }elseif($print_name  == '390-567'){
        $editor->rotate($image2, '90', new Color('#ff0000'));
        $editor->blend($image1, $image2, 'normal', 1, 'top-left', 320, 202);
        $editor->blend($image1, $image2, 'normal', 1, 'top-left', 907, 202);
        $editor->blend($image1, $image2, 'normal', 1, 'top-left', 320, 612);
        $editor->blend($image1, $image2, 'normal', 1, 'top-left', 907, 612);
    }elseif($print_name  == '260-378'){
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
    }elseif($print_name  == '413-531'){
        $editor->rotate($image2 ,'90',new Color('#ff0000'));
        $editor->blend($image1, $image2, 'normal', 1, 'top-left', 81, 179);
        $editor->blend($image1, $image2, 'normal', 1, 'top-left', 632, 179);
        $editor->blend($image1, $image2, 'normal', 1, 'top-left', 1183, 179);
        $editor->blend($image1, $image2, 'normal', 1, 'top-left', 81, 612);
        $editor->blend($image1, $image2, 'normal', 1, 'top-left', 632, 612);
        $editor->blend($image1, $image2, 'normal', 1, 'top-left', 1183, 612);
    }elseif($print_name  == '626-413'){
        $editor->blend($image1, $image2, 'normal', 1, 'top-left', 261, 179);
        $editor->blend($image1, $image2, 'normal', 1, 'top-left', 907, 179);
        $editor->blend($image1, $image2, 'normal', 1, 'top-left', 261, 612);
        $editor->blend($image1, $image2, 'normal', 1, 'top-left', 907, 612);
    } elseif ($print_name == '649-991') {
        $editor->blend($image1, $image2, 'normal', 1, 'top-left', 238, 107);
        $editor->blend($image1, $image2, 'normal', 1, 'top-left', 907, 107);
    } elseif ($print_name == '1050-1499') {
        $editor->rotate($image2, '90', new Color('#ff0000'));
        $editor->blend($image1, $image2, 'normal', 1, 'top-left', 148, 78);
    }
    $save_path = __DIR__. '/res-print.jpeg'; //排版后保存的目录
    $editor->save($image1, $save_path);

    exit('排版完成');
} catch (Exception $e) {
    var_dump($e->getMessage());
}

