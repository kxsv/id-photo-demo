import java.nio.file.Path;
import java.util.Map;

/**
 * V5 换装
 *
 * 接口地址：https://api.zjzapi.com/idcardv5/clothe
 * 请求方式：POST（application/x-www-form-urlencoded）
 * 返回格式：JSON
 * 接口费用：按次扣费（服装可能另计）
 *
 * 用法（在 java 目录编译后运行）：
 *   javac -encoding UTF-8 *.java
 *   java V5ClotheDemo
 */
public class V5ClotheDemo {

    // ===== 请求参数（修改后直接运行即可）=====
    private static final String PIC_ID = "";
    private static final String CLOTHE = "";

    public static void main(String[] args) throws Exception {
        if (Config.API_KEY.isEmpty() || "YOUR_APP_KEY".equals(Config.API_KEY)) {
            System.err.println("提示：请先在 Config.java 中将 API_KEY 替换为你的 API Key。");
        }
        if (PIC_ID.isEmpty() || CLOTHE.isEmpty()) {
            System.err.println("请先在 PIC_ID / CLOTHE 中填入参数（PIC_ID 来自制作接口，CLOTHE 见官网《服装列表》）。");
            return;
        }

        Map<String, String> form = ZjzApiClient.form(
                "key", Config.API_KEY,
                "pic_id", PIC_ID,
                "clothe", CLOTHE
        );

        String body = ZjzApiClient.postForm("idcardv5/clothe", form);
        if (ZjzApiClient.jsonInt(body, "code") != 0) {
            System.out.println("业务失败：" + body);
            return;
        }

        Map<String, String> list = ZjzApiClient.jsonStringMap(body, "list");
        if (list.isEmpty()) {
            System.out.println("返回数据中没有图片列表。");
            return;
        }
        Path outputDir = ZjzApiClient.outputDir("v5-clothe");
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
