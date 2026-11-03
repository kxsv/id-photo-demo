import java.nio.file.Path;

/**
 * 黑白图片上色
 *
 * 接口地址：https://api.zjzapi.com/image/colorize
 * 请求方式：POST（application/x-www-form-urlencoded）
 * 返回格式：JSON
 * 接口费用：按次扣费
 *
 * 用法（在 java 目录编译后运行）：
 *   javac -encoding UTF-8 *.java
 *   java ImageColorizeDemo
 */
public class ImageColorizeDemo {

    // ===== 请求参数（修改后直接运行即可）=====
    private static final String IMAGE_PATH = "../images/test.jpeg";

    public static void main(String[] args) throws Exception {
        if (Config.API_KEY.isEmpty() || "YOUR_APP_KEY".equals(Config.API_KEY)) {
            System.err.println("提示：请先在 Config.java 中将 API_KEY 替换为你的 API Key。");
        }

        java.util.Map<String, String> form = ZjzApiClient.form(
                "key", Config.API_KEY,
                "image", ZjzApiClient.imageToBase64(IMAGE_PATH)
        );

        String body = ZjzApiClient.postForm("image/colorize", form);
        if (ZjzApiClient.jsonInt(body, "code") != 0) {
            System.out.println("业务失败：" + body);
            return;
        }

        String image = ZjzApiClient.jsonString(body, "image");
        if (image != null && !image.isEmpty()) {
            Path target = ZjzApiClient.outputDir("image-colorize").resolve("colorize.jpg");
            try {
                ZjzApiClient.downloadImage(image, target);
                System.out.println("黑白上色结果：" + target);
            } catch (Exception e) {
                System.out.println("下载失败：" + e.getMessage());
            }
        } else {
            System.out.println("返回数据中没有图片。");
        }
    }
}
