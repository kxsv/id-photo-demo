# 鹧应证件照 API Demo

官网地址：https://zjzapi.com/

可运行的接口调用示例，覆盖平台全部接口和证件照排版示例，提供 **PHP、Java、Python、Go、Node.js、C#** 六种语言实现。

## 目录结构

```text
id-photo-demo/
├── README.md          说明文档（接口清单与运行示例）
├── images/            待处理示例图片（test.jpeg）
├── output/            运行结果输出目录（自动创建）
├── php/               PHP 示例（PHP 7.2+，基于 Guzzle，vendor 已内置）
├── python/            Python 示例（Python 3.8+，需 requests）
├── java/              Java 示例（Java 11+）
├── go/                Go 示例（Go 1.18+）
├── nodejs/            Node.js 示例（Node.js 18+）
└── csharp/            C# 示例（.NET 8+）
```


## 快速开始

1. 选择熟悉的编程语言，配置文件中的 API Key （在官网「用户中心 - 应用管理」创建应用后获取）。
2. 默认使用 `images/test.jpeg` 测试，也可换成自己的照片（改示例文件顶部的图片路径即可）。
3. 修改示例文件顶部「请求参数」区的图片路径、规格 id 等，然后直接运行：


### Python

```bash
# 配置文件：python/config.py
pip install -r requirements.txt
cd python

python v3/all.py          # V3 检测+制作
python v3/bg.py           # V3 更换底色
python v3/diy.py          # V3 自定义制作
python v3/pose.py         # V3 职业形象照制作
python v3/clothe.py       # V3 换装（需 pic_id）
python v3/beauty.py       # V3 美颜（需 pic_id）
python v5/make.py         # V5 检测+制作（有水印）
python v5/bg.py           # V5 更换底色
python v5/diy.py          # V5 自定义制作
python v5/pose.py         # V5 职业形象照制作
python v5/clothe.py       # V5 换装（需 pic_id）
python v5/beauty.py       # V5 美颜（需 pic_id）
python v5/get.py          # V5 获取无水印证件照（需 pic_id）
python receipt/make.py    # 回执检测+制作
python receipt/submit.py     # 回执提交（需 pic_id + notice_url）
python receipt/resubmit.py   # 回执驳回后再次提交（需 order_no）
python receipt/cancel.py     # 回执订单取消（需 order_no）
python receipt/notice.py     # 回执状态通知回调（需独立运行的服务，作 notice_url 接收端）
python check/check.py     # 证件照合规性检测
python face/enhance.py    # 人像修复增强
python image/colorize.py  # 黑白图片上色
python image/edit.py      # 图片编辑（改 dpi / 文件大小）
python item/list.py       # 证件照规格列表
python item/detail.py     # 证件照规格详情
python account/info.py    # 获取账号信息（需 access_token）
python account/app.py     # 获取应用信息（需 access_token）
```

### Java

```bash
# 配置文件：java/Config.java
cd java
javac -encoding UTF-8 *.java
java V3AllDemo
java ItemListDemo
java ReceiptNoticeDemo   # 回执状态通知回调（需独立运行的服务）
```

> JDK 22+ 支持免编译直接运行源码：`java V3AllDemo.java`（JEP 458）；JDK 11~21 需先执行上面的 `javac` 编译。
> Windows 控制台如中文显示乱码，先执行 `chcp 65001` 切换到 UTF-8 代码页。

### Go

```bash
# 配置文件：go/zjzapi/config.go
cd go
go run ./cmd/v3/all         # 证件照制作 V3（检测+制作）
go run ./cmd/v3/bg          # V3 更换底色
go run ./cmd/v3/diy         # V3 自定义制作
go run ./cmd/v3/pose        # V3 职业形象照制作
go run ./cmd/v3/clothe      # V3 换装（需 pic_id）
go run ./cmd/v3/beauty      # V3 美颜（需 pic_id）
go run ./cmd/v5/make        # 证件照制作 V5（有水印）
go run ./cmd/v5/bg          # V5 更换底色
go run ./cmd/v5/diy         # V5 自定义制作
go run ./cmd/v5/pose        # V5 职业形象照制作
go run ./cmd/v5/clothe      # V5 换装（需 pic_id）
go run ./cmd/v5/beauty      # V5 美颜（需 pic_id）
go run ./cmd/v5/get         # V5 获取无水印证件照（需 pic_id）
go run ./cmd/receipt/make   # 回执检测+制作
go run ./cmd/receipt/submit     # 回执提交（需 pic_id + notice_url）
go run ./cmd/receipt/resubmit   # 回执驳回后再次提交（需 order_no）
go run ./cmd/receipt/cancel     # 回执订单取消（需 order_no）
go run ./cmd/receipt/notice     # 回执状态通知回调（需独立运行的服务）
go run ./cmd/check/check    # 证件照合规性检测
go run ./cmd/face/enhance   # 人像修复增强
go run ./cmd/image/colorize # 黑白图片上色
go run ./cmd/image/edit     # 图片编辑（改 dpi / 文件大小）
go run ./cmd/item/list      # 证件照规格列表
go run ./cmd/item/detail    # 证件照规格详情
go run ./cmd/account/info   # 获取账号信息（需 access_token）
go run ./cmd/account/app    # 获取应用信息（需 access_token）
```

> Go 下每个示例为 `go/cmd/<分类>/<子路径>/main.go`，统一用 `go run ./cmd/<分类>/<子路径>` 运行；
> 目录结构与 nodejs/php 的接口子路径一一对应（如 `cmd/v3/all` ↔ `v3/all.js` ↔ `v3/all.php`）。
> 需 pic_id / order_no / access_token 的示例，先运行前置接口拿到对应参数再填入示例顶部「请求参数」区。

### PHP

```bash
# 配置文件：php/config.php
# 依赖 Guzzle（vendor 目录已随压缩包提供，无需 composer install）
cd php
php v5/make.php # V5 检测+制作（有水印）
```



### Node.js

```bash
# 配置文件：nodejs/config.js
cd nodejs
node v3/all.js
node receipt/notice.js   # 回执状态通知回调（需独立运行的服务）
```

### C#

```bash
# 配置文件：csharp/Config.cs
cd csharp
dotnet run -- item-list         # 证件照规格列表
dotnet run -- v3-all            # V3 检测+制作
dotnet run -- v3-bg             # V3 更换底色
dotnet run -- v3-diy            # V3 自定义制作
dotnet run -- v3-pose           # V3 职业形象照制作
dotnet run -- v3-clothe         # V3 换装（需 pic_id）
dotnet run -- v3-beauty         # V3 美颜（需 pic_id）
dotnet run -- v5-make           # V5 检测+制作（有水印）
dotnet run -- v5-bg             # V5 更换底色
dotnet run -- v5-diy            # V5 自定义制作
dotnet run -- v5-pose           # V5 职业形象照制作
dotnet run -- v5-clothe         # V5 换装（需 pic_id）
dotnet run -- v5-beauty         # V5 美颜（需 pic_id）
dotnet run -- v5-get            # V5 获取无水印证件照（需 pic_id）
dotnet run -- receipt-make      # 回执检测+制作
dotnet run -- receipt-submit    # 回执提交（需 pic_id + notice_url）
dotnet run -- receipt-resubmit  # 回执驳回后再次提交（需 order_no）
dotnet run -- receipt-cancel    # 回执订单取消（需 order_no）
dotnet run -- receipt-notice     # 回执状态通知回调（需独立运行的服务）
dotnet run -- check             # 证件照合规性检测
dotnet run -- face-enhance      # 人像修复增强
dotnet run -- image-colorize    # 黑白图片上色
dotnet run -- image-edit        # 图片编辑（改 dpi / 文件大小）
dotnet run -- item-detail       # 证件照规格详情
dotnet run -- account-info      # 获取账号信息（需 access_token）
dotnet run -- account-app       # 获取应用信息（需 access_token）
```

> 每个示例文件顶部都有「请求参数」区（图片路径、规格 id、底色、美颜等），修改后直接运行，无需命令行传参。

## 通用约定

- 所有接口均为 POST 请求，表单编码 `application/x-www-form-urlencoded`；
- `image` 参数为图片内容的 Base64 字符串（不带 `data:image/...` 前缀），原图大小不超过 8M；
- HTTP 200 不代表业务成功，还需判断返回 JSON 中的 `code` 是否为 0；
- 请求参数（图片路径、规格 id、底色等）统一写在每个示例文件顶部的「请求参数」区，修改后直接运行；
- 制作/处理结果图会自动下载到对应语言的 `output/` 子目录；
- 各语言示例文件与接口子路径同名：如 `v3/all.php`、`v3/all.py`、`v3/all.js`、`go run ./cmd/v3/all`、C# 子命令 `v3-all`；Java 为对应的驼峰类名（如 `V3AllDemo`）；
- Java 示例不引入第三方依赖，业务字段使用正则提取，生产环境建议使用 Jackson/Gson 等成熟 JSON 库。

## 接口清单

### 证件照制作 V3

| 接口 | 接口地址 | 示例子路径 |
| --- | --- | --- |
| 证件照检测+制作 | https://api.zjzapi.com/idcardv3/all | `v3/all` |
| 图片检测+更换底色 | https://api.zjzapi.com/idcardv3/bg | `v3/bg` |
| 证件照检测+自定义制作 | https://api.zjzapi.com/idcardv3/diy | `v3/diy` |
| 职业形象照检测+制作 | https://api.zjzapi.com/idcardv3/all | `v3/pose` |
| 证件照换装 | https://api.zjzapi.com/idcardv3/clothe | `v3/clothe` |
| 证件照美颜 | https://api.zjzapi.com/idcardv3/beauty | `v3/beauty` |

### 证件照制作 V5

| 接口 | 接口地址 | 示例子路径 |
| --- | --- | --- |
| 证件照检测+制作（有水印） | https://api.zjzapi.com/idcardv5/make | `v5/make` |
| 图片检测+更换底色（有水印） | https://api.zjzapi.com/idcardv5/bg | `v5/bg` |
| 证件照检测+自定义制作（有水印） | https://api.zjzapi.com/idcardv5/diy | `v5/diy` |
| 职业形象照检测+制作（有水印） | https://api.zjzapi.com/idcardv5/make | `v5/pose` |
| 证件照美颜（有水印） | https://api.zjzapi.com/idcardv5/beauty | `v5/beauty` |
| 证件照换装（有水印） | https://api.zjzapi.com/idcardv5/clothe | `v5/clothe` |
| 获取证件照（无水印） | https://api.zjzapi.com/idcardv5/get | `v5/get` |

### 证件照回执办理

| 接口 | 接口地址 | 示例子路径 |
| --- | --- | --- |
| 证件照回执检测+制作 | https://api.zjzapi.com/receipt/make | `receipt/make` |
| 证件照回执提交 | https://api.zjzapi.com/receipt/submit | `receipt/submit` |
| 订单状态通知 | 回调通知（服务端接收） | `receipt/notice` |
| 审核驳回后再次提交 | https://api.zjzapi.com/receipt/submit | `receipt/resubmit` |
| 订单取消 | https://api.zjzapi.com/receipt/cancel | `receipt/cancel` |

> `receipt/notice` 是**服务端接收回调**的示例（不是主动调用），六种语言均已提供：`php/receipt/notice.php`、`python/receipt/notice.py`、`nodejs/receipt/notice.js`、`java/ReceiptNoticeDemo.java`、`go run ./cmd/receipt/notice`、`dotnet run -- receipt-notice`。它需作为独立服务运行，把其地址填入 `receipt/submit` 的 `notice_url`，平台推送订单状态变化时验证 `sign = md5(API Key)` 后返回 `{"data":{},"code":0}`。

### 证件照检测

| 接口 | 接口地址 | 示例子路径 |
| --- | --- | --- |
| 证件照合规性检测 | https://api.zjzapi.com/idphoto/check | `check/check` |

### 人像修复增强

| 接口 | 接口地址 | 示例子路径 |
| --- | --- | --- |
| 人像修复增强 | https://api.zjzapi.com/face/enhance | `face/enhance` |

### 图片处理

| 接口 | 接口地址 | 示例子路径 |
| --- | --- | --- |
| 黑白图片上色 | https://api.zjzapi.com/image/colorize | `image/colorize` |
| 图片编辑 | https://api.zjzapi.com/image/edit | `image/edit` |

### 规格相关

| 接口 | 接口地址 | 示例子路径 |
| --- | --- | --- |
| 证件照规格列表 | https://api.zjzapi.com/item/list | `item/list` |
| 证件照规格详情 | https://api.zjzapi.com/item/get | `item/detail` |

### 账号相关

| 接口 | 接口地址 | 示例子路径 |
| --- | --- | --- |
| 获取账号信息 | https://api.zjzapi.com/user/info | `account/info` |
| 获取应用信息 | https://api.zjzapi.com/user/app | `account/app` |

