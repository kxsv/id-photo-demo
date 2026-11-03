# 鹧应证件照 API · Python Demo

基于 Python 的接口调用示例，覆盖证件照制作（V3/V5）、回执办理、合规检测、人像增强、图片编辑、规格查询与账号信息等全部接口。每个示例都是独立的 `.py` 文件，修改顶部「请求参数」后直接运行即可。

## 目录结构

```text
python/
├── README.md           本说明
├── requirements.txt   依赖（requests + Pillow；Pillow 仅 print.py 冲印排版需要）
├── config.py          配置文件（API_KEY / ACCESS_TOKEN）
├── zjz_api.py         公共客户端封装
│   ├── ZjzApi.post()          以 POST 表单方式调用接口
│   ├── ZjzApi.is_success()    判断业务是否成功（code 为 0）
│   ├── ZjzApi.image_to_base64()  读取本地图片并转 Base64
│   ├── ZjzApi.download_image()   下载单张结果图到本地
│   └── ZjzApi.save_images()      批量下载一组结果图
├── print.py           冲印排版（相纸联排，需 Pillow）
└── <分类>/<接口>.py    各接口示例，目录结构与 nodejs/php 的接口子路径一一对应
    ├── v3/            all / bg / diy / pose / clothe / beauty
    ├── v5/            make / bg / diy / pose / clothe / beauty / get
    ├── receipt/       make / submit / resubmit / cancel / notice（回执状态通知回调）
    ├── check/         check（证件照合规性检测）
    ├── face/          enhance（人像修复增强）
    ├── image/         colorize（黑白上色）/ edit（图片编辑）
    ├── item/          list（规格列表）/ detail（规格详情）
    └── account/       info（账号信息）/ app（应用信息）
```

## 环境准备

需要 Python 3.7+。安装依赖：

```bash
pip install -r requirements.txt
```

## 配置

编辑 `config.py`：

- `API_KEY`：应用的 API Key，在官网「用户中心 - 应用管理」创建应用后获取。**除账号相关接口外，其余接口只用它。**
- `ACCESS_TOKEN`：用户中心的 access_token，在官网「用户中心 - 账户信息」获取。**仅 `account/info`、`account/app` 两个接口需要。**

> 发布前请将其恢复为占位符 `YOUR_APP_KEY` / `YOUR_ACCESS_TOKEN`，切勿把真实 key 打进压缩包。

## 运行方式

在 `python` 目录下，统一用 `python <分类>/<接口>.py` 运行；示例文件顶部的「请求参数」区（图片路径、规格 id、底色、美颜等）修改后直接运行，无需命令行传参。

```bash
cd python
python item/list.py        # 证件照规格列表（先跑它拿到 item_id / 底色）
python v3/all.py           # V3 检测+制作
python v3/bg.py            # V3 更换底色
python v3/diy.py           # V3 自定义制作
python v3/pose.py          # V3 职业形象照制作
python v3/clothe.py        # V3 换装（需 pic_id）
python v3/beauty.py        # V3 美颜（需 pic_id）
python v5/make.py          # V5 检测+制作（有水印）
python v5/bg.py            # V5 更换底色
python v5/diy.py           # V5 自定义制作
python v5/pose.py          # V5 职业形象照制作
python v5/clothe.py        # V5 换装（需 pic_id）
python v5/beauty.py        # V5 美颜（需 pic_id）
python v5/get.py           # V5 获取无水印证件照（需 pic_id）
python receipt/make.py     # 回执检测+制作
python receipt/submit.py      # 回执提交（需 pic_id + notice_url）
python receipt/resubmit.py    # 回执驳回后再次提交（需 order_no）
python receipt/cancel.py      # 回执订单取消（需 order_no）
python receipt/notice.py      # 回执状态通知回调（需独立运行的服务，作 notice_url 接收端）
python check/check.py      # 证件照合规性检测
python face/enhance.py     # 人像修复增强
python image/colorize.py   # 黑白图片上色
python image/edit.py       # 图片编辑（改 dpi / 文件大小）
python item/detail.py      # 证件照规格详情
python account/info.py     # 获取账号信息（需 access_token）
python account/app.py      # 获取应用信息（需 access_token）
python print.py            # 冲印排版（相纸联排）
```

> 需 `pic_id` / `order_no` / `access_token` 的示例，先运行前置接口拿到对应参数再填入示例顶部「请求参数」区。

## 通用约定

- **请求方式**：所有接口均为 `POST`（表单：`application/x-www-form-urlencoded`），由 `ZjzApi.post()` 统一封装。
- **图片入参**：`image` 字段为本地图片的 Base64（不带 `data:image/...` 前缀），图片大小上限 8M，由 `ZjzApi.image_to_base64()` 处理；不支持远程 URL。
- **成功判定**：HTTP 200 不代表业务成功，必须再判断返回 JSON 中的 `code == 0`（`ZjzApi.is_success()`）。`code != 0` 时打印 `code` 与 `msg`。
- **结果图输出**：制作类接口的结果图默认下载到仓库根的 `output/<名称>/` 目录（如 `output/v3-all/`、`output/receipt-make/`），单个图片下载失败会打印错误并继续，不影响其余图片。
- **pic_id 有效期**：制作/回执类接口返回的 `pic_id` 有效期约 50 分钟，换装、美颜、V5 获取无水印等后续接口都基于它，请尽快使用。
- **接口费用**：各接口费用见每个示例文件顶部的 docstring；账户余额不足或被风控时接口会返回对应 `code`，请留意控制台提示。

## 公共封装（zjz_api.py）

| 方法 | 说明 |
| --- | --- |
| `ZjzApi.post(path, params)` | 调用接口，返回 JSON 转成的 dict；网络或 HTTP 错误抛 `ZjzApiError`。 |
| `ZjzApi.is_success(response)` | 业务是否成功（`code == 0`）。 |
| `ZjzApi.image_to_base64(path)` | 读取本地图片并转 Base64（不含前缀），校验文件存在与 8M 上限。 |
| `ZjzApi.download_image(url, save_dir, name)` | 下载单张结果图到本地目录。 |
| `ZjzApi.save_images(save_dir, images)` | 批量下载「名称 -> 图片URL」的结果图，并打印每个结果；单个失败仅提示后继续。 |

每个示例的标准处理流程：

```python
result = ZjzApi.post("<接口路径>", {"key": API_KEY, ...})
if not ZjzApi.is_success(result):
    print(f"业务失败：code={result.get('code')}，msg={result.get('msg')}")
    return
data = result.get("data") or {}
# 读取 data 中的字段 / 调用 ZjzApi.save_images(...) 下载结果图

## 冲印排版（print.py）

把一张已制作好的证件照，按预设的相纸模板（位于仓库根 `images/print-bg/<相纸>/<尺寸>.png`）复制到模板的各个槽位上，输出一张可直接拿去冲印的图（默认 300dpi）。适用于把单张证件照批量联排到 6寸等相纸上打印。

### 依赖
- 需要 Pillow（已在 `requirements.txt` 声明，`pip install -r requirements.txt` 一并安装）。

### 用法
```bash
# 默认输入 output/v3-diy/blue.jpg（需先运行 v3/diy.py 生成 1寸 295x413 证件照）
python print.py

# 指定任意一张做好的证件照
python print.py 路径/照片.jpg
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
- 输出文件：`output/print/print6-<尺寸>.jpg`（仓库根 `output/`，`print6` 表示 6寸相纸，像素 1795×1205@300dpi）。
- 带文字 / 裁切线说明的模板（`*-text.png`）通过 `print.py` 顶部的 `USE_TEXT_TPL = True` 开启。
- 新增 / 调整规格：直接编辑 `print.py` 顶部的 `LAYOUTS` 配置（`size` 槽位尺寸、`rotate` 是否旋转 90°、`positions` 各槽位左上角坐标）。
- 新增相纸：在仓库根 `images/print-bg/` 下新建子目录（如 `print5/`），把 `print.py` 顶部的 `PAPER` 改为对应目录名，并补充该相纸下的模板图与对应 `positions`。
```
