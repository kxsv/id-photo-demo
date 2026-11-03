# 鹧应证件照 API · Java 示例

可运行的 Java 接口调用示例，覆盖开放平台全部接口（证件照制作 V3/V5、回执办理、检测、人像增强、图片处理、规格与账号相关）。

不引入任何第三方依赖：使用 JDK 自带的 `java.net.http.HttpClient` 发请求，返回 JSON 用正则提取业务字段（生产环境建议替换为 Jackson/Gson 等成熟 JSON 库）。

公共客户端提取的图片 URL 和字符串字段会解码 JSON 转义，兼容 `https:\/\/...` 和 Unicode 转义形式。

## 环境要求

- JDK 11+（`java.net.http.HttpClient` 需 Java 11）

## 目录结构

```text
java/
├── README.md            本说明
├── Config.java          配置文件（API_KEY / ACCESS_TOKEN）
├── ZjzApiClient.java    请求封装（postForm / imageToBase64 / downloadImage / json* 解析）
├── ItemListDemo.java    证件照规格列表
├── ItemDetailDemo.java  证件照规格详情
├── V3AllDemo.java       V3 检测+制作
├── V3BgDemo.java        V3 更换底色
├── V3DiyDemo.java       V3 自定义制作
├── V3PoseDemo.java      V3 职业形象照制作（需职业形象照规格 id）
├── V3ClotheDemo.java    V3 换装（需 pic_id）
├── V3BeautyDemo.java    V3 美颜（需 pic_id）
├── V5MakeDemo.java      V5 检测+制作（带水印）
├── V5BgDemo.java        V5 更换底色
├── V5DiyDemo.java       V5 自定义制作
├── V5PoseDemo.java      V5 职业形象照制作（需职业形象照规格 id）
├── V5ClotheDemo.java    V5 换装（需 pic_id）
├── V5BeautyDemo.java    V5 美颜（需 pic_id）
├── V5GetDemo.java       V5 获取无水印证件照（需 pic_id）
├── ReceiptMakeDemo.java 回执检测+制作（需回执规格 id）
├── ReceiptSubmitDemo.java   回执提交（需 pic_id + notice_url）
├── ReceiptResubmitDemo.java 回执驳回后重新提交（需 order_no）
├── ReceiptCancelDemo.java   回执订单取消（需 order_no）
├── ReceiptNoticeDemo.java 回执状态通知回调（本地 HTTP 服务，需独立运行，接收平台推送的订单状态）
├── CheckDemo.java       证件照合规性检测
├── FaceEnhanceDemo.java 人像修复增强
├── ImageColorizeDemo.java   黑白图片上色
├── ImageEditDemo.java   图片编辑（改 dpi / 文件大小）
├── AccountInfoDemo.java 获取账号信息（需 access_token）
├── AccountAppDemo.java  获取应用信息（需 access_token）
├── PrintDemo.java       冲印排版（相纸联排，JDK 内置 imageio，无需第三方依赖）
└── output/              结果图输出目录
```

## 配置

编辑 `Config.java`：

- `API_KEY`：应用的 API Key，在官网「用户中心 - 应用管理」创建应用后获取。**除账号相关接口外，其余接口只用它。**
- `ACCESS_TOKEN`：用户中心的 access_token，在官网「用户中心 - 账户信息」获取。**仅 `user/info`、`user/app` 两个接口需要。**

## 编译与运行

所有示例均为默认包（无 `package` 声明），在 `java` 目录统一编译后运行：

```bash
cd java

# 编译
javac -encoding UTF-8 *.java

# 运行单个示例
java ItemListDemo
java ItemDetailDemo
java V3AllDemo
java V3BgDemo
java V3DiyDemo
java V3PoseDemo
java V3ClotheDemo
java V3BeautyDemo
java V5MakeDemo
java V5BgDemo
java V5DiyDemo
java V5PoseDemo
java V5ClotheDemo
java V5BeautyDemo
java V5GetDemo <pic_id> [gift_hd] [clothe] [gift_hd_color]
java ReceiptMakeDemo
java ReceiptSubmitDemo
java ReceiptResubmitDemo
java ReceiptCancelDemo
java ReceiptNoticeDemo
java CheckDemo
java FaceEnhanceDemo
java ImageColorizeDemo
java ImageEditDemo
java AccountInfoDemo
java AccountAppDemo
java PrintDemo                                         # 冲印排版（相纸联排）
```

> 图片路径等请求参数在各示例文件顶部的「请求参数」区修改后直接运行，无需命令行传参。
> 通用示例图片默认读取 `../images/test.jpeg`、`../images/pose.jpg`。

## 通用约定

- 所有接口均为 POST 请求，表单编码 `application/x-www-form-urlencoded`；
- `image` 参数为图片内容的 Base64 字符串（不带 `data:image/...` 前缀），原图大小不超过 8M；
- HTTP 200 不代表业务成功，还需判断返回 JSON 中的 `code` 是否为 0；
- 制作 / 处理结果图会自动下载到 `output/` 下的对应子目录（如 `output/v3-all/`）。

## 冲印排版（PrintDemo）

把一张已制作好的证件照，按预设的相纸模板（位于仓库根 `images/print-bg/<相纸>/<尺寸>.png`）复制到模板的各个槽位上，输出一张可直接拿去冲印的图（默认 300dpi）。适用于把单张证件照批量联排到 6寸等相纸上打印。

### 依赖
- JDK 自带的 `javax.imageio`（无需任何第三方库）。

### 用法
```bash
# 默认输入 output/v3-diy/blue.jpg（需先运行 V3DiyDemo 生成 1寸 295x413 证件照）
java PrintDemo

# 指定任意一张做好的证件照
java PrintDemo 路径/照片.jpg
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
- 输出文件：`java/output/print/print6-<尺寸>.jpg`（`print6` 表示 6寸相纸，像素 1795×1205@300dpi）。
- 带文字 / 裁切线说明的模板（`*-text.png`）通过 `PrintDemo.java` 顶部的 `USE_TEXT_TPL = true` 开启。
- 新增 / 调整规格：直接编辑 `PrintDemo.java` 顶部的 `LAYOUTS` 配置（`w/h` 槽位尺寸、`rotate` 是否旋转 90°、`positions` 各槽位左上角坐标）。
- 新增相纸：在仓库根 `images/print-bg/` 下新建子目录（如 `print5/`），把 `PrintDemo.java` 顶部的 `PAPER` 改为对应目录名，并补充该相纸下的模板图与对应 `positions`。
