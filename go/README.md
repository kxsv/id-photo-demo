# 鹧应证件照 API · Go 示例

可运行的 Go 接口调用示例，覆盖开放平台全部接口（证件照制作 V3/V5、回执办理、检测、人像增强、图片处理、规格与账号相关），与 PHP / Python / Java / Node.js / C# 示例能力一致。

## 环境要求

- Go 1.18+

## 目录结构

```text
go/
├── README.md           本说明
├── go.mod              模块定义（module: zjz-demo）
├── zjzapi/             公共客户端封装
│   ├── config.go       配置文件（APIKey / AccessToken）
│   └── client.go       请求封装（Post / IsSuccess / ImageToBase64 / DownloadImage / SaveImages）
└── cmd/                各接口示例，每个「分类/子路径」子目录一个 main 包
    ├── v3/
    │   ├── all/        V3 检测+制作
    │   ├── bg/         更换底色
    │   ├── diy/        自定义制作
    │   ├── pose/       职业形象照制作
    │   ├── clothe/     换装（需 pic_id）
    │   └── beauty/     美颜（需 pic_id）
    ├── v5/             make/bg/diy/pose/clothe/beauty/get（同上，结果带/去水印）
    ├── receipt/        make/submit/resubmit/cancel/notice（回执办理 / 状态通知回调）
    ├── check/check/    证件照合规性检测
    ├── face/enhance/   人像修复增强
    ├── image/          colorize/上色、edit/编辑
    ├── item/           list/规格列表、detail/规格详情
    ├── account/        info/账号信息、app/应用信息（需 access_token）
    └── print/          冲印排版（相纸联排，标准库 image，无需第三方依赖）
```

## 配置

编辑 `zjzapi/config.go`：

- `APIKey`：应用的 API Key，在官网「用户中心 - 应用管理」创建应用后获取。**除账号相关接口外，其余接口只用它。**
- `AccessToken`：用户中心的 access_token，在官网「用户中心 - 账户信息」获取。**仅 `account/info`、`account/app` 两个接口需要。**

> 发布前请将其恢复为占位符 `YOUR_APP_KEY` / `YOUR_ACCESS_TOKEN`，切勿把真实 key 打进压缩包。

## 运行方式

每个示例都是独立的可执行 `main` 包，目录结构与 nodejs/php 的接口子路径一一对应；在 `go` 目录下统一用 `go run ./cmd/<分类>/<子路径>` 运行。示例文件顶部的「请求参数」区（图片路径、规格 id、底色、美颜等）修改后直接运行，无需命令行传参。

```bash
cd go
go run ./cmd/item/list      # 证件照规格列表（先跑它拿到 item_id / 底色）
go run ./cmd/v3/all         # V3 检测+制作
go run ./cmd/v3/bg          # V3 更换底色
go run ./cmd/v3/diy         # V3 自定义制作
go run ./cmd/v3/pose        # V3 职业形象照制作
go run ./cmd/v3/clothe      # V3 换装（需 pic_id）
go run ./cmd/v3/beauty      # V3 美颜（需 pic_id）
go run ./cmd/v5/make        # V5 检测+制作（有水印）
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
go run ./cmd/item/detail    # 证件照规格详情
go run ./cmd/account/info   # 获取账号信息（需 access_token）
go run ./cmd/account/app    # 获取应用信息（需 access_token）
go run ./cmd/print          # 冲印排版（相纸联排）
```

## 通用约定

- 所有接口均为 POST 请求，表单编码 `application/x-www-form-urlencoded`；
- `image` 参数为图片内容的 Base64 字符串（不带 `data:image/...` 前缀），原图大小不超过 8M；
- HTTP 200 不代表业务成功，还需判断返回 JSON 中的 `code` 是否为 0（`Response.IsSuccess()`）；
- 制作 / 处理结果图会自动下载到 `go/output/` 下的对应子目录（如 `output/v3-all/`）；
- 带「需 pic_id / order_no / access_token」标注的示例，先运行前置接口拿到对应参数，再填入示例顶部「请求参数」区。`pic_id` 有效期 50 分钟。

## 公共封装（zjzapi 包）

| 方法 | 说明 |
| --- | --- |
| `Post(path, params)` | 以 POST 表单方式调用接口，返回解析后的 `*Response` |
| `(*Response).IsSuccess()` | 业务是否成功（`code == 0`） |
| `ImageToBase64(path)` | 读取本地图片并转 Base64（自动校验 8M 上限） |
| `DownloadImage(url, dir, name)` | 下载单张结果图到本地 |
| `SaveImages(dir, images)` | 批量下载一组「名称 -> URL」的结果图并打印结果 |

各示例的请求参数均通过 `url.Values` 组装，与官网接口文档字段一一对应（如 `item_id`、`colors`、`beauty`、`whitening`、`width_px` 等）。

## 冲印排版（cmd/print）

把一张已制作好的证件照，按预设的相纸模板（位于仓库根 `images/print-bg/<相纸>/<尺寸>.png`）复制到模板的各个槽位上，输出一张可直接拿去冲印的图（默认 300dpi）。适用于把单张证件照批量联排到 6寸等相纸上打印。

### 依赖
- 仅使用 Go 标准库（`image` / `image/draw` / `image/jpeg` / `image/png`），无需任何第三方库。

### 用法
```bash
# 默认输入 output/v3-diy/blue.jpg（需先运行 v3/diy 生成 1寸 295x413 证件照）
go run ./cmd/print

# 指定任意一张做好的证件照
go run ./cmd/print 路径/照片.jpg
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
- 输出文件：`go/output/print/print6-<尺寸>.jpg`（`print6` 表示 6寸相纸，像素 1795×1205@300dpi）。
- 带文字 / 裁切线说明的模板（`*-text.png`）通过 `cmd/print/main.go` 顶部的 `useTextTpl = true` 开启。
- 新增 / 调整规格：直接编辑 `cmd/print/main.go` 顶部的 `layouts` 配置（`w/h` 槽位尺寸、`rotate` 是否旋转 90°、`positions` 各槽位左上角坐标）。
- 新增相纸：在仓库根 `images/print-bg/` 下新建子目录（如 `print5/`），把 `cmd/print/main.go` 顶部的 `paper` 改为对应目录名，并补充该相纸下的模板图与对应 `positions`。
