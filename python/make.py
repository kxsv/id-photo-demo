import json

import requests
import os
import base64
from PIL import Image

def make():
    api_url = 'http://api.zheyings.cn/idcardv5/make'    # 接口地址 - v5
    item_id = 1                                         # 规格id
    key = ""                                            # api key 获取路径：官网（https://www.zjzapi.com/）/用户中心/应用管理
    path = os.getcwd() + '/../images/test.jpeg'         # 图片地址

    with open(path, 'rb') as f:
        image = f.read()

    base64_image = base64.b64encode(image).decode()
    response = requests.post(url=api_url, data={"item_id": item_id, "key": key, "image": base64_image}, headers={
        "Content-Type": "application/x-www-form-urlencoded"
    })
    data = json.loads(response.text)
    print(data)
    blue_url = data['data']['list']['blue']

    r = requests.get(blue_url)
    path = os.getcwd() + "/blue.jpeg"
    with open(path, 'wb') as f:
        f.write(r.content)

    blue_path = Image.open(path)
    blue_path.show()

if __name__ == '__main__':
    make()
