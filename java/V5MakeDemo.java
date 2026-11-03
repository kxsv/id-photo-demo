import java.nio.file.Path;
import java.util.Map;

/**
 * V5 检测+制作（结果图带水印，附 pic_id 调 V5GetDemo 取无水印）
 *
 * 接口地址：https://api.zjzapi.com/idcardv5/make
 * 请求方式：POST（application/x-www-form-urlencoded）
 * 返回格式：JSON
 * 接口费用：制作成功按次扣费
 *
 * 用法（在 java 目录编译后运行）：
 *   javac -encoding UTF-8 *.java
 *   java V5MakeDemo
 */
public class V5MakeDemo {

    // ===== 请求参数（修改后直接运行即可）=====
    private static final String IMAGE_PATH = "../images/test.jpeg";
    private static final String ITEM_ID = "1";     // 规格 id，可通过 ItemListDemo 获取
    private static final String COLORS = "white";  // 底色，如 white/blue/red
    private static final String ENHANCE = "0";     // 画质增强：0=关闭 1=开启
    private static final String BEAUTY = "0";      // 美颜：0=关闭 1=开启
    private static final String WHITENING = "0";  // 美白：0=关闭 1=开启

    public static void main(String[] args) throws Exception {
        if (Config.API_KEY.isEmpty() || "YOUR_APP_KEY".equals(Config.API_KEY)) {
            System.err.println("提示：请先在 Config.java 中将 API_KEY 替换为你的 API Key。");
        }

        Map<String, String> form = ZjzApiClient.form(
                "key", Config.API_KEY,
                "item_id", ITEM_ID,
                "image", ZjzApiClient.imageToBase64(IMAGE_PATH),
                "enhance", ENHANCE,
                "beauty", BEAUTY,
                "whitening", WHITENING
        );
        if (!COLORS.isEmpty()) {
            form.put("colors", COLORS);
        }

        String body = ZjzApiClient.postForm("idcardv5/make", form);
        if (ZjzApiClient.jsonInt(body, "code") != 0) {
            System.out.println("业务失败：" + body);
            return;
        }

        String picId = ZjzApiClient.jsonString(body, "pic_id");
        System.out.println("制作成功，pic_id（用于 V5GetDemo 取无水印）：" + (picId == null ? "" : picId));

        Map<String, String> list = ZjzApiClient.jsonStringMap(body, "list");
        if (list.isEmpty()) {
            System.out.println("返回数据中没有图片列表。");
            return;
        }
        Path outputDir = ZjzApiClient.outputDir("v5-make");
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
