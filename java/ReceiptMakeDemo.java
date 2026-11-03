import java.nio.file.Path;
import java.util.Map;

/**
 * 回执检测+制作
 *
 * 接口地址：https://api.zjzapi.com/receipt/make
 * 请求方式：POST（application/x-www-form-urlencoded）
 * 返回格式：JSON
 * 接口费用：按次扣费
 *
 * 用法（在 java 目录编译后运行）：
 *   javac -encoding UTF-8 *.java
 *   java ReceiptMakeDemo
 */
public class ReceiptMakeDemo {

    // ===== 请求参数（修改后直接运行即可）=====
    private static final String IMAGE_PATH = "../images/test.jpeg";
    private static final String ITEM_ID = "";  // 回执类规格 id（请通过 ItemListDemo 获取后填入）

    public static void main(String[] args) throws Exception {
        if (Config.API_KEY.isEmpty() || "YOUR_APP_KEY".equals(Config.API_KEY)) {
            System.err.println("提示：请先在 Config.java 中将 API_KEY 替换为你的 API Key。");
        }
        if (ITEM_ID.isEmpty()) {
            System.err.println("请先在 ITEM_ID 中填入回执类规格 id（运行 ItemListDemo 查看）。");
            return;
        }

        Map<String, String> form = ZjzApiClient.form(
                "key", Config.API_KEY,
                "item_id", ITEM_ID,
                "image", ZjzApiClient.imageToBase64(IMAGE_PATH)
        );

        String body = ZjzApiClient.postForm("receipt/make", form);
        if (ZjzApiClient.jsonInt(body, "code") != 0) {
            System.out.println("业务失败：" + body);
            return;
        }

        String picId = ZjzApiClient.jsonString(body, "pic_id");
        System.out.println("回执制作成功，pic_id（用于 ReceiptSubmitDemo）：" + (picId == null ? "" : picId));

        Map<String, String> list = ZjzApiClient.jsonStringMap(body, "list");
        if (!list.isEmpty()) {
            Path outputDir = ZjzApiClient.outputDir("receipt-make");
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
}
