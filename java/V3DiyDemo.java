import java.nio.file.Path;
import java.util.Map;

/**
 * V3 自定义制作
 *
 * 接口地址：https://api.zjzapi.com/idcardv3/diy
 * 请求方式：POST（application/x-www-form-urlencoded）
 * 返回格式：JSON
 * 接口费用：按次扣费
 *
 * 用法（在 java 目录编译后运行）：
 *   javac -encoding UTF-8 *.java
 *   java V3DiyDemo
 */
public class V3DiyDemo {

    // ===== 请求参数（修改后直接运行即可）=====
    private static final String IMAGE_PATH = "../images/test.jpeg";
    private static final String WIDTH_PX = "295";    // 目标像素宽
    private static final String HEIGHT_PX = "413";   // 目标像素高
    private static final String COLOR = "blue";      // 底色（与 PrintDemo 默认输入 blue.jpg 保持一致）
    private static final String DPI = "300";         // 输出 DPI
    private static final String BEAUTY = "0";        // 美颜：0=关闭 1=开启
    private static final String WHITENING = "0";     // 美白：0=关闭 1=开启
    private static final String FILE_SIZE_MIN = "";  // 最小文件大小(KB)，留空不限制
    private static final String FILE_SIZE_MAX = "";  // 最大文件大小(KB)，留空不限制

    public static void main(String[] args) throws Exception {
        if (Config.API_KEY.isEmpty() || "YOUR_APP_KEY".equals(Config.API_KEY)) {
            System.err.println("提示：请先在 Config.java 中将 API_KEY 替换为你的 API Key。");
        }

        Map<String, String> form = ZjzApiClient.form(
                "key", Config.API_KEY,
                "image", ZjzApiClient.imageToBase64(IMAGE_PATH),
                "width_px", WIDTH_PX,
                "height_px", HEIGHT_PX,
                "color", COLOR,
                "dpi", DPI
        );
        if (!BEAUTY.isEmpty()) form.put("beauty", BEAUTY);
        if (!WHITENING.isEmpty()) form.put("whitening", WHITENING);
        if (!FILE_SIZE_MIN.isEmpty()) form.put("file_size_min", FILE_SIZE_MIN);
        if (!FILE_SIZE_MAX.isEmpty()) form.put("file_size_max", FILE_SIZE_MAX);

        String body = ZjzApiClient.postForm("idcardv3/diy", form);
        if (ZjzApiClient.jsonInt(body, "code") != 0) {
            System.out.println("业务失败：" + body);
            return;
        }

        Map<String, String> list = ZjzApiClient.jsonStringMap(body, "list");
        if (list.isEmpty()) {
            System.out.println("返回数据中没有图片列表。");
            return;
        }
        Path outputDir = ZjzApiClient.outputDir("v3-diy");
        for (Map.Entry<String, String> entry : list.entrySet()) {
            Path target = outputDir.resolve(entry.getKey() + ".jpg");
            try {
                ZjzApiClient.downloadImage(entry.getValue(), target);
                System.out.println("底色 " + entry.getKey() + "：" + target);
            } catch (Exception e) {
                System.out.println("底色 " + entry.getKey() + " 下载失败：" + e.getMessage());
            }
        }
    }
}
