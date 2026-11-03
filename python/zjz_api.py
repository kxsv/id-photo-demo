"""鹧应证件照 API 公共客户端。

所有接口示例共用的请求封装：
- post()          以 POST 表单方式调用接口
- is_success()    判断业务是否成功（code 为 0）
- image_to_base64()  读取本地图片并转 Base64
- download_image()   下载接口返回的结果图到本地
"""

import base64
from pathlib import Path

import requests

BASE_URL = "https://api.zjzapi.com"
CONNECT_TIMEOUT = 10
REQUEST_TIMEOUT = 60
MAX_IMAGE_BYTES = 8 * 1024 * 1024


class ZjzApiError(RuntimeError):
    """HTTP 失败或返回内容不是合法 JSON 等异常。"""


class ZjzApi:
    @staticmethod
    def post(path, params):
        """调用接口，返回 JSON 转成的 dict。"""
        try:
            response = requests.post(
                f"{BASE_URL}/{path.lstrip('/')}",
                headers={"Content-Type": "application/x-www-form-urlencoded"},
                data=params,
                timeout=(CONNECT_TIMEOUT, REQUEST_TIMEOUT),
            )
        except requests.RequestException as error:
            raise ZjzApiError(f"请求失败：{error}") from error

        if not 200 <= response.status_code < 300:
            raise ZjzApiError(f"HTTP {response.status_code}：{response.text}")

        try:
            return response.json()
        except ValueError as error:
            raise ZjzApiError(f"返回内容不是合法 JSON：{response.text}") from error

    @staticmethod
    def is_success(response):
        """业务是否成功（code 为 0）。"""
        return int(response.get("code", -1)) == 0

    @staticmethod
    def image_to_base64(path):
        """读取本地图片并转 Base64（不带 data:image/... 前缀）。"""
        image = Path(path)
        if not image.is_file():
            raise ZjzApiError(f"图片不存在：{path}")
        content = image.read_bytes()
        if len(content) > MAX_IMAGE_BYTES:
            raise ZjzApiError(f"图片不能超过 8M：{path}")
        return base64.b64encode(content).decode()

    @staticmethod
    def download_image(image_url, save_dir, name):
        """下载接口返回的结果图到本地。"""
        save_dir = Path(save_dir)
        save_dir.mkdir(parents=True, exist_ok=True)
        try:
            response = requests.get(image_url, timeout=(CONNECT_TIMEOUT, REQUEST_TIMEOUT))
        except requests.RequestException as error:
            raise ZjzApiError(f"下载失败：{error}") from error
        if not 200 <= response.status_code < 300:
            raise ZjzApiError(f"下载失败：HTTP {response.status_code}")
        file = save_dir / name
        file.write_bytes(response.content)
        return str(file)

    @staticmethod
    def save_images(save_dir, images):
        """下载一组「名称 -> 图片URL」的结果图到本地目录，并打印每个结果；
        单个图片下载失败不影响其余图片，仅打印错误后继续。
        """
        save_dir = Path(save_dir)
        save_dir.mkdir(parents=True, exist_ok=True)
        if not images:
            print("返回数据中没有图片列表。")
            return
        for name, image_url in images.items():
            try:
                file = ZjzApi.download_image(image_url, str(save_dir), f"{name}.jpg")
                print(f"{name}：{file}")
            except ZjzApiError as error:
                print(f"{name} 下载失败：{error}")
