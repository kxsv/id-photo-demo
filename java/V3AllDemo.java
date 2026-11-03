import java.util.Map;

/**
 * V3 证件照检测+制作
 *
 * 接口地址：https://api.zjzapi.com/idcardv3/all
 * 请求方式：POST（application/x-www-form-urlencoded）
 * 返回格式：JSON
 * 接口费用：制作成功按次扣费；检测不通过扣除检测费用；beauty=1 时美颜另行计费。
 *
 * 用法（在 java 目录编译后运行）：
 *   javac -encoding UTF-8 *.java
 *   java V3AllDemo
 */
public class V3AllDemo {

    // ===== 请求参数（修改后直接运行即可）=====
    private static final String IMAGE_PATH = "../images/test.jpeg"; // 待处理图片路径（相对 java 目录）
    private static final String ITEM_ID = "1";     // 规格id，可通过「证件照规格列表」示例获取
    private static final String COLORS = "white";    // 底色，如 white/blue/red；留空返回该规格支持的全部底色
    private static final String BEAUTY = "0";       // 美颜：0=关闭 1=开启（开启后另行计费）

    public static void main(String[] args) throws Exception {
        if (Config.API_KEY.isEmpty() || "YOUR_APP_KEY".equals(Config.API_KEY)) {
            System.err.println("提示：请先在 Config.java 中将 API_KEY 替换为你的 API Key。");
        }

        Map<String, String> form = ZjzApiClient.form(
                "key", Config.API_KEY,
                "item_id", ITEM_ID,
                "image", ZjzApiClient.imageToBase64(IMAGE_PATH),
                "beauty", BEAUTY
        );
        // colors 为空时不传，由服务端返回全部底色。
        if (!COLORS.isEmpty()) {
            form.put("colors", COLORS);
        }

        String body = ZjzApiClient.postForm("idcardv3/all", form);

        // HTTP 成功不代表业务成功，需再检查返回 JSON 中的 code。
        if (ZjzApiClient.jsonInt(body, "code") != 0) {
            System.out.println("业务失败：" + body);
            return;
        }

        String picId = ZjzApiClient.jsonString(body, "pic_id");
        System.out.println("制作成功，pic_id：" + (picId == null ? "" : picId));

        // 把各底色结果图下载到 output/v3-all/ 目录。
        Map<String, String> list = ZjzApiClient.jsonStringMap(body, "list");
        if (list.isEmpty()) {
            System.out.println("返回数据中没有图片列表。");
            return;
        }
        java.nio.file.Path outputDir = ZjzApiClient.outputDir("v3-all");
        for (Map.Entry<String, String> entry : list.entrySet()) {
            java.nio.file.Path target = outputDir.resolve(entry.getKey() + ".jpg");
            try {
                ZjzApiClient.downloadImage(entry.getValue(), target);
                System.out.println("底色 " + entry.getKey() + "：" + target);
            } catch (Exception e) {
                System.out.println("底色 " + entry.getKey() + " 下载失败：" + e.getMessage());
            }
        }
    }
}
