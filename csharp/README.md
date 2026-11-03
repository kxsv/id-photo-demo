# 鹧应证件照 API · C# Demo

基于 .NET 8 的接口调用示例，覆盖证件照制作（V3/V5）、回执办理、合规检测、人像增强、图片编辑、规格查询与账号信息等全部接口。采用「单项目 + 子命令」结构：统一入口 `Program.cs` 按 `dotnet run -- <demo名称>` 分发到 `Demos/` 下对应的示例类。

## 目录结构

```text
csharp/
├── README.md           本说明
├── id-photo-demo.csproj 项目文件（net8.0，自带 HttpClient / System.Text.Json，其余示例无需额外依赖；仅 `print` 冲印排版示例需 SixLabors.ImageSharp，已在 csproj 中声明，首次运行自动还原）
├── Config.cs           配置文件（ApiKey / AccessToken）
├── ZjzApi.cs           公共客户端封装
│   ├── PostAsync()         以 POST 表单方式调用接口
│   ├── IsSuccess()         判断业务是否成功（code 为 0）
│   ├── ImageToBase64Async() 读取本地图片并转 Base64
│   ├── DownloadImageAsync() 下载单张结果图到本地
│   ├── GetString()         兼容数字/字符串字段，统一按字符串读取
│   └── SaveImagesAsync()   批量下载一组结果图
└── Demos/<Name>.cs     各接口示例（每个类含 RunAsync()），名称与 nodejs/php 的接口子路径一一对应
    ├── V3All/V3Bg/V3Diy/V3Pose/V3Clothe/V3Beauty
    ├── V5Make/V5Bg/V5Diy/V5Pose/V5Clothe/V5Beauty/V5Get
    ├── ReceiptMake/ReceiptSubmit/ReceiptResubmit/ReceiptCancel/ReceiptNotice
    ├── Check         证件照合规性检测
    ├── FaceEnhance    人像修复增强
    ├── ImageColorize  黑白图片上色
    ├── ImageEdit      图片编辑
    ├── ItemList/ItemDetail
    ├── AccountInfo/AccountApp
    ├── Print          冲印排版（相纸联排，需 SixLabors.ImageSharp）
```

## 环境准备

需要 .NET 8 SDK。在 `csharp` 目录下可直接运行（首次会自行还原项目）：

```bash
cd csharp
dotnet run -- item-list
```

## 编译

本项目为单文件 .NET 8 控制台程序，无需额外依赖（HttpClient / System.Text.Json 已内置）。
在 `csharp` 目录下执行以下命令即可编译与运行。

```bash
# 1. 还原依赖（首次或拉取新代码后；SDK 风格项目一般会自动还原，可省略）
dotnet restore

# 2. 编译（生成到 bin/Debug/net8.0/，不运行）
dotnet build

# 3. 编译并直接运行某个 demo
dotnet run -- <demo名称>

# 4. 发布为独立可执行文件（可选，含运行时，免去目标机装 SDK）
dotnet publish -c Release -o ./publish
```

常用 `dotnet run` 示例：

```bash
dotnet run -- item-list      # 规格列表
dotnet run -- v3-all         # V3 检测+制作
```

### 冲印排版（print）

把一张已制作好的证件照，按预设的相纸模板（位于仓库根 `images/print-bg/<相纸>/<尺寸>.png`）复制到模板的各个槽位上，输出一张可直接拿去冲印的图（默认 300dpi）。适用于把单张证件照批量联排到 6寸等相纸上打印。

> 依赖说明：该示例需 `SixLabors.ImageSharp`（纯托管绘图库，无原生依赖），已在 `id-photo-demo.csproj` 中声明，`dotnet run` 首次会自动还原；其余示例仍零依赖。

### 用法

```bash
# 默认输入 output/v3-diy/blue.jpg（需先运行 v3-diy 生成 1寸 295x413 证件照）
dotnet run -- print

# 指定任意一张做好的证件照
dotnet run -- print 路径/照片.jpg
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

- 输出文件：`csharp/output/print/print6-<尺寸>.jpg`（`print6` 表示 6寸相纸，像素 1795×1205@300dpi）。
- 带文字 / 裁切线说明的模板（`*-text.png`）通过 `Demos/Print.cs` 顶部的 `UseTextTpl = true` 开启。
- 新增 / 调整规格：直接编辑 `Demos/Print.cs` 顶部的 `Layouts` 配置（`SlotW/SlotH` 槽位尺寸、`Rotate` 是否旋转 90°、`Positions` 各槽位左上角坐标）。
- 新增相纸：在仓库根 `images/print-bg/` 下新建子目录（如 `print5/`），把 `Demos/Print.cs` 顶部的 `Paper` 改为对应目录名，并补充该相纸下的模板图与对应 `Positions`。

## 清理编译产物

`bin/`（编译输出）与 `obj/`（中间产物，含 `.cache`、`.dll`、`.exe`、`.json` 等）均由 `dotnet build` / `dotnet run` 自动生成，**可随时安全删除**，下次构建会自动重建。等价于：

```bash
dotnet clean                 # 删除 bin/ 与 obj/ 中的生成文件
# 或手动删除这两个目录
rm -rf bin obj
```

> 源码仅 `Config.cs`、`Program.cs`、`ZjzApi.cs`、`id-photo-demo.csproj`、`Demos/*.cs`、`README.md`，请勿删除。

## 配置

编辑 `Config.cs`：

- `ApiKey`：应用的 API Key，在官网「用户中心 - 应用管理」创建应用后获取。**除账号相关接口外，其余接口只用它。**
- `AccessToken`：用户中心的 access_token，在官网「用户中心 - 账户信息」获取。**仅 `account-info`、`account-app` 两个接口需要。**

> 发布前请将其恢复为占位符 `YOUR_APP_KEY` / `YOUR_ACCESS_TOKEN`，切勿把真实 key 打进压缩包。

## 运行方式

在 `csharp` 目录下，统一用 `dotnet run -- <demo名称>` 运行；示例文件顶部的「请求参数」区（图片路径、规格 id、底色、美颜等）修改后直接运行，无需命令行传参。

```bash
cd csharp
dotnet run -- item-list         # 证件照规格列表（先跑它拿到 item_id / 底色）
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
dotnet run -- print             # 冲印排版（相纸联排，把做好的证件照排到 6寸相纸）
```

> 需 `pic_id` / `order_no` / `access_token` 的示例，先运行前置接口拿到对应参数再填入示例顶部「请求参数」区。

## 通用约定

- **请求方式**：所有接口均为 `POST`（表单：`application/x-www-form-urlencoded`），由 `ZjzApi.PostAsync()` 统一封装。
- **图片入参**：`image` 字段为本地图片的 Base64（不带 `data:image/...` 前缀），图片大小上限 8M，由 `ZjzApi.ImageToBase64Async()` 处理；不支持远程 URL。示例中图片路径默认 `../images/test.jpeg`（相对 `csharp` 目录，即仓库根的 `images/`）。
- **成功判定**：HTTP 200 不代表业务成功，必须再判断返回 JSON 中的 `code == 0`（`ZjzApi.IsSuccess()`）。`code != 0` 时打印 `code` 与 `msg`。
- **结果图输出**：制作类接口的结果图默认下载到 `csharp/output/<名称>/` 目录（如 `output/v3-all/`、`output/receipt-make/`）；单个图片下载失败会打印错误并继续，不影响其余图片。
- **pic_id 有效期**：制作/回执类接口返回的 `pic_id` 有效期约 50 分钟，换装、美颜、V5 获取无水印等后续接口都基于它，请尽快使用。
- **接口费用**：各接口费用见每个示例文件顶部的 docstring；账户余额不足或被风控时接口会返回对应 `code`，请留意控制台提示。

## 公共封装（ZjzApi.cs）

| 方法 | 说明 |
| --- | --- |
| `PostAsync(path, form)` | 调用接口，返回 `Response`（含 `Code`、`Msg`、`Data`）；网络或 HTTP 错误抛 `HttpRequestException`。 |
| `IsSuccess(response)` | 业务是否成功（`Code == 0`）。 |
| `ImageToBase64Async(path)` | 读取本地图片并转 Base64（不含前缀），校验文件存在与 8M 上限。 |
| `DownloadImageAsync(url, saveDirectory, name)` | 下载单张结果图到本地目录。 |
| `GetString(element, name)` | 提取 JSON 字段，数字与字符串两种写法都兼容（数字按原始文本返回）；`pic_id`、`order_no`、账号数值等字段优先用本方法读取，避免 `GetString()` 在数值上抛异常。 |
| `SaveImagesAsync(saveDirectory, listValue)` | 批量下载「名称 -> 图片URL」的结果图，并打印每个结果；单个失败仅提示后继续。 |

每个示例的标准处理流程：

```csharp
var result = await ZjzApi.PostAsync("<接口路径>", new Dictionary<string, string>
{
    ["key"] = Config.ApiKey,
    // ... 其他参数
});
if (!ZjzApi.IsSuccess(result))
{
    Console.WriteLine($"业务失败：code={result.Code}，msg={result.Msg}");
    return;
}
var picId = ZjzApi.GetString(result.Data, "pic_id");
// 下载结果图
if (result.Data.TryGetProperty("list", out var listValue))
    await ZjzApi.SaveImagesAsync(Path.Combine("output", "<名称>"), listValue);
```
