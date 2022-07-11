<?php

require 'vendor/autoload.php';


$api_url = 'http://api.zheyings.cn/idcardv5/make';     //接口地址 - v5
$item_id = 1;                                           //规格id
$key = '';                                              //api key 获取路径：官网（https://www.zjzapi.com/）/用户中心/应用管理
$path = __DIR__.'/pic.png';                             //图片地址
try {
    $client = new \GuzzleHttp\Client();
    $response = $client->request('POST', $api_url, [
        'form_params' => [
            'item_id' => $item_id,
            'key' => $key,
            'image' => base64_encode(file_get_contents($path)),
        ]
    ]);
    $json = $response->getBody();
    $res = json_decode($json, true);
    var_dump($res);
} catch (\Exception $e) {
    var_dump($e->getMessage());
}
