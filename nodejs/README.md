# 鹧应证件照 API 接口调用示例（Node.js）

基于 Node.js 内置 `fetch` 的鹧应证件照 API（`https://api.zjzapi.com`）命令行调用示例，覆盖证件照制作、检测、换底色、换装、美颜、人像增强、老照片修复、纸质回执等全部接口。除 `print.js`（冲印排版，需 `sharp` 库）外，其余接口示例无需任何第三方依赖。

## 环境要求

- Node.js >= 18（需原生 `fetch` 支持；项目已声明 `"type": "module"`，使用 ESM 语法）

## 快速开始

1. 配置 API Key：编辑 `config.js`，将 `API_KEY` 替换为你在官网「用户中心 - 应用管理」创建应用后获取的 API Key。`ACCESS_TOKEN` 仅 `account/info.js`、`account/app.js` 需要，在官网「用户中心 - 账户信息」获取。

   ```js
   // config.js
   export const API_KEY = '你的API Key';
   // 仅 account/info.js、account/app.js 需要
   export const ACCESS_TOKEN = 'YOUR_ACCESS_TOKEN';
   ```

2. 在 `nodejs` 目录下运行示例（均在命令行执行）：

   ```bash
   # 查看证件照规格列表，获取 item_id 与支持底色（免费）
   node item/list.js

   # V3 证件照检测+制作
   node v3/all.js

   # V5 证件照制作（有水印，免费；再用 pic_id 获取无水印图）
   node v5/make.js
   node v5/get.js
   ```

3. 待处理图片默认使用压缩包根目录 `images/test.jpeg`，可在各示例顶部修改 `imagePath`；结果图下载到 `output/` 对应子目录（如 `output/v3-all/`、`output/v5-make/`）。

## 目录结构

```
nodejs/
├── config.js            # 全局配置（API Key、access_token）
├── print.js             # 冲印排版（相纸联排，需 sharp）
├── lib/zjz-api.js       # 公共客户端：post()/isSuccess()/imageToBase64()/downloadImage()
├── account/             # 账户相关
├── check/               # 证件照合规性检测
├── face/                # 人像增强
├── image/               # 图片编辑、老照片上色
├── item/                # 证件照规格列表 / 详情
├── receipt/             # 纸质回执（提交、重传、制作、取消、通知）
├── v3/                  # V3 证件照制作系列
└── v5/                  # V5 证件照制作系列
```

## 示例清单

| 目录 | 文件 | 对应接口 | 说明 |
| --- | --- | --- | --- |
| `item/` | `list.js` | `item/list` | 证件照规格列表（免费） |
| `item/` | `detail.js` | `item/get` | 证件照规格详情 |
| `check/` | `check.js` | `idphoto/check` | 证件照合规性检测，可按规格或自定义参数 |
| `v3/` | `all.js` | `idcardv3/all` | V3 检测+制作 |
| `v3/` | `beauty.js` | `idcardv3/beauty` | 美颜 |
| `v3/` | `bg.js` | `idcardv3/bg` | 换底色 |
| `v3/` | `clothe.js` | `idcardv3/clothe` | 换装 |
| `v3/` | `diy.js` | `idcardv3/diy` | 自定义规格制作 |
| `v3/` | `pose.js` | `idcardv3/all` | 形象照/职业照（职业照规格 + 专属底色） |
| `v5/` | `make.js` | `idcardv5/make` | V5 检测+制作（有水印，免费） |
| `v5/` | `get.js` | `idcardv5/get` | 用 `pic_id` 获取无水印证件照 |
| `v5/` | `beauty.js` | `idcardv5/beauty` | V5 美颜 |
| `v5/` | `bg.js` | `idcardv5/bg` | V5 换底色 |
| `v5/` | `clothe.js` | `idcardv5/clothe` | V5 换装 |
| `v5/` | `diy.js` | `idcardv5/diy` | V5 自定义规格制作 |
| `v5/` | `pose.js` | `idcardv5/make` | V5 形象照/职业照 |
| `face/` | `enhance.js` | `face/enhance` | 人像增强 |
| `image/` | `edit.js` | `image/edit` | 图片编辑（修改 dpi / 文件大小） |
| `image/` | `colorize.js` | `image/colorize` | 老照片修复上色 |
| `receipt/` | `submit.js` | `receipt/submit` | 提交纸质回执 |
| `receipt/` | `resubmit.js` | `receipt/submit` | 重新提交（携带 `order_no`） |
| `receipt/` | `make.js` | `receipt/make` | 回执制作 |
| `receipt/` | `cancel.js` | `receipt/cancel` | 取消回执订单 |
| `receipt/` | `notice.js` | — | 回执状态通知回调示例（需独立运行的服务） |
| `account/` | `info.js` | `user/info` | 账户信息（需 `access_token`） |
| `account/` | `app.js` | `user/app` | 应用信息（需 `access_token`） |
| `—` | `print.js` | — | 冲印排版（相纸联排，需 `sharp`） |

## 使用说明

- 所有示例仅支持命令行运行（`node xxx.js`），不支持浏览器访问。
- 各示例文件顶部均有「请求参数」区块，修改后直接运行即可；无需命令行传参。
- 请求统一走 `lib/zjz-api.js`：
  - `post(path, params)`：POST 表单（`application/x-www-form-urlencoded`）调用接口，返回解析后的 JSON；
  - `isSuccess(result)`：判断业务是否成功（`code === 0`）；
  - `imageToBase64(imagePath)`：读取本地图片转 Base64（不带 `data:image/...` 前缀，不超过 8M）；
  - `downloadImage(imageUrl, saveDir, name)`：下载结果图到本地。
- HTTP 成功不代表业务成功，务必检查返回 JSON 中的 `code` 是否为 `0`。
- 制作/换装/美颜等基于 `pic_id` 的接口，`pic_id` 有效期 50 分钟，请先运行对应制作接口获取后填入参数区。
- 各接口的费用说明见对应示例文件头部注释及官网文档。

## 冲印排版（print.js）

把一张已制作好的证件照，按预设的相纸模板（位于仓库根 `images/print-bg/<相纸>/<尺寸>.png`）复制到模板的各个槽位上，输出一张可直接拿去冲印的图（默认 300dpi）。适用于把单张证件照批量联排到 6寸等相纸上打印。

### 依赖
- 需要 `sharp` 图像处理库（已在 `package.json` 声明，`npm install` 一并安装）；其余接口示例仍零依赖。

### 用法
```bash
# 默认输入 output/v3-diy/blue.jpg（需先运行 v3/diy.js 生成 1寸 295x413 证件照）
node print.js

# 指定任意一张做好的证件照
node print.js 路径/照片.jpg
```

### 支持的证件照尺寸与排版
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
- 输出文件：`nodejs/output/print/print6-<尺寸>.jpg`（`print6` 表示 6寸相纸，像素 1795×1205@300dpi）。
- 带文字 / 裁切线说明的模板（`*-text.png`）通过 `print.js` 顶部的 `USE_TEXT_TPL = true` 开启。
- 新增 / 调整规格：直接编辑 `print.js` 顶部的 `LAYOUTS` 配置（`size` 槽位尺寸、`rotate` 是否旋转 90°、`positions` 各槽位左上角坐标）。
- 新增相纸：在仓库根 `images/print-bg/` 下新建子目录（如 `print5/`），把 `print.js` 顶部的 `PAPER` 改为对应目录名，并补充该相纸下的模板图与对应 `positions`。
