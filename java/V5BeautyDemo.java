import java.nio.file.Path;
import java.util.Map;

/**
 * V5 美颜
 *
 * 接口地址：https://api.zjzapi.com/idcardv5/beauty
 * 请求方式：POST（application/x-www-form-urlencoded）
 * 返回格式：JSON
 * 接口费用：按次扣费
 *
 * 用法（在 java 目录编译后运行）：
 *   javac -encoding UTF-8 *.java
 *   java V5BeautyDemo
 */
public class V5BeautyDemo {

    // ===== 请求参数（修改后直接运行即可）=====
    private static final String PIC_ID = "";

    public static void main(String[] args) throws Exception {
        if (Config.API_KEY.isEmpty() || "YOUR_APP_KEY".equals(Config.API_KEY)) {
            System.err.println("提示：请先在 Config.java 中将 API_KEY 替换为你的 API Key。");
        }
        if (PIC_ID.isEmpty()) {
            System.err.println("请先在 PIC_ID 中填入制作接口返回的 pic_id（先运行 V5MakeDemo 获取）。");
            return;
        }

        Map<String, String> form = ZjzApiClient.form(
                "key", Config.API_KEY,
                "pic_id", PIC_ID
        );

        String body = ZjzApiClient.postForm("idcardv5/beauty", form);
        if (ZjzApiClient.jsonInt(body, "code") != 0) {
            System.out.println("业务失败：" + body);
            return;
        }

        Map<String, String> list = ZjzApiClient.jsonStringMap(body, "list");
        if (list.isEmpty()) {
            System.out.println("返回数据中没有图片列表。");
            return;
        }
        Path outputDir = ZjzApiClient.outputDir("v5-beauty");
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
