# 鹧应证件照 API 接口调用示例（PHP）

基于 [GuzzleHttp](https://github.com/guzzle/guzzle) 的鹧应证件照 API（`https://api.zjzapi.com`）命令行调用示例，覆盖证件照制作、检测、换底色、换装、美颜、人像增强、老照片修复、纸质回执等常用接口。

## 环境要求

- PHP >= 7.2（需启用 `curl` 扩展）
- Composer（可选：压缩包已附带 `vendor/`，可直接运行）
- 冲印排版示例 `print.php` 需启用 `imagick` 或 `gd` 扩展之一（自动二选一，**无需第三方库**）

## 快速开始

1. 配置 API Key：编辑 `config.php`，将 `key` 替换为你在官网「用户中心 - 应用管理」创建应用后获取的 API Key。

   ```php
   return [
       'key' => '你的API Key',
       // 仅 account/info.php、account/app.php 需要
       'access_token' => 'YOUR_ACCESS_TOKEN',
   ];
   ```

2. 在 `php` 目录下运行示例（均在命令行执行）：

   ```bash
   # 查看证件照规格列表，获取 item_id 与支持底色（免费）
   php item/list.php

   # V3 证件照检测+制作
   php v3/all.php

   # V5 证件照制作（有水印，再用 pic_id 获取无水印图）
   php v5/make.php
   php v5/get.php

   # 冲印排版：先生成证件照，再联排到 6寸相纸
   php v3/diy.php     # 生成 output/v3-diy/blue.jpg（1寸 295x413）
   php print.php      # 输出 output/print/print6-295-413.jpg
   ```

3. 待处理图片默认使用压缩包根目录 `images/test.jpeg`，可在各示例顶部修改 `$image_path`；结果图下载到 `../output/` 对应子目录。

## 目录结构

```
php/
├── bootstrap.php      # 公共引导文件：加载配置、公共客户端与工具函数
├── config.php         # 全局配置（API Key、access_token）
├── lib/ZjzApi.php     # 公共客户端：post()/isSuccess()/imageToBase64()/downloadImage()
├── account/           # 账户相关
├── check/             # 证件照合规性检测
├── face/              # 人像增强
├── image/             # 图片编辑、老照片上色
├── item/              # 证件照规格列表 / 详情
├── receipt/           # 纸质回执（提交、重传、制作、取消、通知）
├── v3/                # V3 证件照制作系列
├── v5/                # V5 证件照制作系列
└── print.php          # 冲印排版：把做好的证件照联排到相纸模板上
```
> 相纸模板位于压缩包根目录 `images/print-bg/<相纸>/<尺寸>.png`（`print6` = 6寸）。

## 示例清单

| 目录 | 文件 | 对应接口 | 说明 |
| --- | --- | --- | --- |
| `item/` | `list.php` | `item/list` | 证件照规格列表（免费） |
| `item/` | `detail.php` | `item/get` | 证件照规格详情 |
| `check/` | `check.php` | `idphoto/check` | 证件照合规性检测，可按规格或自定义参数 |
| `v3/` | `all.php` | `idcardv3/all` | V3 检测+制作 |
| `v3/` | `beauty.php` | `idcardv3/beauty` | 美颜 |
| `v3/` | `bg.php` | `idcardv3/bg` | 换底色 |
| `v3/` | `clothe.php` | `idcardv3/clothe` | 换装 |
| `v3/` | `diy.php` | `idcardv3/diy` | 自定义规格制作 |
| `v3/` | `pose.php` | `idcardv3/all` | 形象照/职业照（职业照规格 + 专属底色） |
| `v5/` | `make.php` | `idcardv5/make` | V5 检测+制作（有水印，免费） |
| `v5/` | `get.php` | `idcardv5/get` | 用 `pic_id` 获取无水印证件照 |
| `v5/` | `beauty.php` | `idcardv5/beauty` | V5 美颜 |
| `v5/` | `bg.php` | `idcardv5/bg` | V5 换底色 |
| `v5/` | `clothe.php` | `idcardv5/clothe` | V5 换装 |
| `v5/` | `diy.php` | `idcardv5/diy` | V5 自定义规格制作 |
| `v5/` | `pose.php` | `idcardv5/make` | V5 形象照/职业照 |
| `face/` | `enhance.php` | `face/enhance` | 人像增强 |
| `image/` | `edit.php` | `image/edit` | 图片编辑（裁剪、旋转等） |
| `image/` | `colorize.php` | `image/colorize` | 老照片修复上色 |
| `receipt/` | `submit.php` | `receipt/submit` | 提交纸质回执 |
| `receipt/` | `resubmit.php` | `receipt/submit` | 重新提交（携带 `order_no`） |
| `receipt/` | `make.php` | `receipt/make` | 回执制作 |
| `receipt/` | `cancel.php` | `receipt/cancel` | 取消回执订单 |
| `receipt/` | `notice.php` | — | 回执状态通知回调示例 |
| `account/` | `info.php` | `user/info` | 账户信息（需 `access_token`） |
| `account/` | `app.php` | `user/app` | 应用信息（需 `access_token`） |
| `—` | `print.php` | — | 冲印排版（相纸联排）：把做好的证件照按模板排版到 6寸相纸（需 `imagick` 或 `gd`） |

## 使用说明

- 所有示例仅支持命令行运行（`php xxx.php`），不支持浏览器访问。
- 各示例文件顶部均有「请求参数」区块，修改后直接运行即可。
- 请求统一走 `lib/ZjzApi.php`：
  - `ZjzApi::post($path, $params)`：POST 表单调用接口，返回 JSON 数组；
  - `ZjzApi::isSuccess($result)`：判断业务是否成功（`code === 0`）；
  - `ZjzApi::imageToBase64($path)`：读取本地图片转 Base64（不超过 8M）；
  - `ZjzApi::downloadImage($url, $dir, $name)`：下载结果图到本地。
- HTTP 成功不代表业务成功，务必检查返回 JSON 中的 `code` 是否为 `0`。
- 各接口的费用说明见对应示例文件头部注释及官网文档。

## 冲印排版（print.php）

把一张已制作好的证件照，按预设的相纸模板（位于 `../images/print-bg/<相纸>/<尺寸>.png`）复制到模板的各个槽位上，输出一张可直接拿去冲印的图（默认 300dpi）。适用于把单张证件照批量联排到 6寸等相纸上打印。

### 依赖

- PHP 的 `imagick` 或 `gd` 扩展（**二选一，自动选择，无需任何第三方库 / composer**）。
- 本机无 `imagick` 时自动回退到 `gd`。

### 用法

```bash
# 默认输入 output/v3-diy/blue.jpg（需先运行 v3/diy.php 生成 1寸 295x413 证件照）
php print.php

# 指定任意一张做好的证件照
php print.php 路径/照片.jpg
```

### 支持的证件照尺寸与排版

输入证件照的像素尺寸必须匹配下表之一（即模板文件名 `<宽>-<高>.png`），脚本据此选择对应相纸槽位：

| 尺寸（宽-高 px） | 常见称呼 | 单张相纸联排数量 | 备注 |
| --- | --- | --- | --- |
| 295-413 | 1寸 | 10（5×2） | — |
| 260-378 | 小1寸 | 10（5×2） | — |
| 413-579 | 2寸 | 4（2×2） | 横放 |
| 390-567 | 大2寸 | 4（2×2） | 横放 |
| 413-531 | 小2寸 | 6（3×2） | 横放 |
| 626-413 | 大1寸 | 4（2×2） | — |
| 649-991 | 驾照/特殊 | 2（2×1） | — |
| 1050-1499 | 退伍/特殊 | 1（1×1） | 横放 |
| 358-441 | — | 8（4×2） | — |
| 378-472 | — | 8（4×2） | — |
| 600-600 | 方版 | 2（2×1） | — |

### 说明

- 输出文件：`../output/print/print6-<尺寸>.jpg`（`print6` 表示 6寸相纸，像素 1795×1205@300dpi）。
- 带文字 / 裁切线说明的模板（`*-text.png`）通过 `print.php` 顶部的 `$use_text_tpl = true` 开启。
- 新增 / 调整规格：直接编辑 `print.php` 顶部的 `$layouts` 配置（`size` 槽位尺寸、`rotate` 是否旋转 90°、`positions` 各槽位左上角坐标）。
- 新增相纸：在 `../images/print-bg/` 下新建子目录（如 `print5/`），把 `print.php` 顶部的 `$paper` 改为对应目录名，并补充该相纸下的模板图与对应 `positions`。

## 重新安装依赖

压缩包已附带 `vendor/`，如需自行安装：

```bash
composer install
```
