import java.nio.file.Path;

/**
 * 人像修复增强
 *
 * 接口地址：https://api.zjzapi.com/face/enhance
 * 请求方式：POST（application/x-www-form-urlencoded）
 * 返回格式：JSON
 * 接口费用：按次扣费
 *
 * 用法（在 java 目录编译后运行）：
 *   javac -encoding UTF-8 *.java
 *   java FaceEnhanceDemo
 */
public class FaceEnhanceDemo {

    // ===== 请求参数（修改后直接运行即可）=====
    private static final String IMAGE_PATH = "../images/test.jpeg";
    private static final String SIZE = "original";  // original=原尺寸，x2=2倍超分

    public static void main(String[] args) throws Exception {
        if (Config.API_KEY.isEmpty() || "YOUR_APP_KEY".equals(Config.API_KEY)) {
            System.err.println("提示：请先在 Config.java 中将 API_KEY 替换为你的 API Key。");
        }

        java.util.Map<String, String> form = ZjzApiClient.form(
                "key", Config.API_KEY,
                "image", ZjzApiClient.imageToBase64(IMAGE_PATH),
                "size", SIZE
        );

        String body = ZjzApiClient.postForm("face/enhance", form);
        if (ZjzApiClient.jsonInt(body, "code") != 0) {
            System.out.println("业务失败：" + body);
            return;
        }

        String image = ZjzApiClient.jsonString(body, "image");
        if (image != null && !image.isEmpty()) {
            Path target = ZjzApiClient.outputDir("face-enhance").resolve("enhance.jpg");
            try {
                ZjzApiClient.downloadImage(image, target);
                System.out.println("人像修复增强结果：" + target);
            } catch (Exception e) {
                System.out.println("下载失败：" + e.getMessage());
            }
        } else {
            System.out.println("返回数据中没有图片。");
        }
    }
}
