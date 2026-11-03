<?php
/**
 * 鹧应证件照 API 公共客户端（基于 GuzzleHttp）
 *
 * 所有接口示例共用的请求封装：
 * - post()          以 POST 表单方式调用接口
 * - isSuccess()      判断业务是否成功（code === 0）
 * - imageToBase64()  读取本地图片并转 Base64
 * - downloadImage()  下载接口返回的结果图到本地
 */

declare(strict_types=1);

use GuzzleHttp\Client;
use GuzzleHttp\Exception\RequestException;

class ZjzApi
{
    /** API 网关地址 */
    const BASE_URL = 'https://api.zjzapi.com';

    /** 连接超时（秒） */
    const CONNECT_TIMEOUT = 10;

    /** 请求超时（秒） */
    const TIMEOUT = 60;

    /** @var Client|null 复用的 HTTP 客户端（保持连接复用） */
    private static $client;

    /**
     * 获取 HTTP 客户端实例
     *
     * @return Client
     */
    private static function client()
    {
        if (self::$client === null) {
            self::$client = new Client([
                'connect_timeout' => self::CONNECT_TIMEOUT,
                'timeout' => self::TIMEOUT,
                // 4xx/5xx 时抛出 RequestException（http_errors 默认即 true）。
                'http_errors' => true,
            ]);
        }

        return self::$client;
    }

    /**
     * 调用接口
     *
     * @param string $path   接口路径，如 idcardv3/all
     * @param array  $params 请求参数（已含 key 等）
     * @return array 接口返回的 JSON（已转为数组）
     * @throws RuntimeException HTTP 失败或返回内容不是合法 JSON 时抛出
     */
    public static function post($path, array $params)
    {
        try {
            $response = self::client()->post(self::BASE_URL . '/' . ltrim($path, '/'), [
                'form_params' => $params,
            ]);
        } catch (RequestException $e) {
            if ($e->hasResponse()) {
                $http_response = $e->getResponse();
                throw new RuntimeException(
                    'HTTP ' . $http_response->getStatusCode() . '：' . (string)$http_response->getBody()
                );
            }
            // 连接失败、超时等无响应的异常（含 ConnectException）。
            throw new RuntimeException('请求失败：' . $e->getMessage());
        }

        $body = (string)$response->getBody();
        $data = json_decode($body, true);
        if (!is_array($data)) {
            throw new RuntimeException('返回内容不是合法 JSON：' . $body);
        }

        return $data;
    }

    /**
     * 业务是否成功
     *
     * @param array $response post() 的返回结果
     * @return bool
     */
    public static function isSuccess(array $response)
    {
        return (int)(isset($response['code']) ? $response['code'] : -1) === 0;
    }

    /**
     * 读取本地图片并转 Base64（不带 data:image/... 前缀）
     *
     * @param string $path 图片路径
     * @return string
     * @throws RuntimeException 文件不存在或读取失败时抛出
     */
    public static function imageToBase64($path)
    {
        if ($path === '' || !is_file($path)) {
            throw new RuntimeException('图片不存在：' . ($path === '' ? '(未指定)' : $path));
        }
        $content = file_get_contents($path);
        if ($content === false) {
            throw new RuntimeException('图片读取失败：' . $path);
        }
        if (strlen($content) > 8 * 1024 * 1024) {
            throw new RuntimeException('图片不能超过 8M：' . $path);
        }
        return base64_encode($content);
    }

    /**
     * 下载接口返回的结果图到本地
     *
     * @param string $url      图片地址
     * @param string $save_dir 保存目录（不存在时自动创建）
     * @param string $name     保存文件名，如 white.jpg
     * @return string 保存后的文件路径
     * @throws RuntimeException 目录创建或下载失败时抛出
     */
    public static function downloadImage($url, $save_dir, $name)
    {
        if (!is_dir($save_dir) && !mkdir($save_dir, 0755, true) && !is_dir($save_dir)) {
            throw new RuntimeException('目录创建失败：' . $save_dir);
        }

        $file = rtrim($save_dir, '/\\') . DIRECTORY_SEPARATOR . $name;
        try {
            // sink 直接把响应体写入文件，不占用内存。
            self::client()->get($url, ['sink' => $file]);
        } catch (RequestException $e) {
            throw new RuntimeException('下载失败：' . $e->getMessage());
        }

        return $file;
    }
}
