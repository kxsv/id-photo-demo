import java.nio.file.Path;
import java.util.Map;

/**
 * V5 获取无水印证件照
 *
 * 接口地址：https://api.zjzapi.com/idcardv5/get
 * 请求方式：POST（application/x-www-form-urlencoded）
 * 返回格式：JSON
 * 接口费用：按次扣费（取无水印）
 *
 * 用法（在 java 目录编译后运行）：
 *   javac -encoding UTF-8 *.java
 *   java V5GetDemo <pic_id> [gift_hd] [clothe] [gift_hd_color]
 *
 * 参数（命令行，按位置）：
 *   pic_id        制作接口返回的 pic_id（有效期50分钟），必填
 *   gift_hd       是否获取高清精修照：0=否（默认） 1=是
 *   clothe        换装服装 key（GIFT_HD=1 时可选）
 *   gift_hd_color 高清照底色（GIFT_HD=1 时可选）
 */
public class V5GetDemo {

    public static void main(String[] args) throws Exception {
        if (Config.API_KEY.isEmpty() || "YOUR_APP_KEY".equals(Config.API_KEY)) {
            System.err.println("提示：请先在 Config.java 中将 API_KEY 替换为你的 API Key。");
        }
        if (args.length < 1 || args[0].isEmpty()) {
            System.err.println("用法：java V5GetDemo <pic_id> [gift_hd] [clothe] [gift_hd_color]");
            System.err.println("  pic_id 为必填，先运行 V5MakeDemo 获取。");
            return;
        }

        String PIC_ID = args[0];
        String GIFT_HD = args.length > 1 ? args[1] : "0";
        String CLOTHE = args.length > 2 ? args[2] : "";
        String GIFT_HD_COLOR = args.length > 3 ? args[3] : "";

        Map<String, String> form = ZjzApiClient.form(
                "key", Config.API_KEY,
                "pic_id", PIC_ID,
                "gift_hd", GIFT_HD
        );
        if (!CLOTHE.isEmpty()) form.put("clothe", CLOTHE);
        if (!GIFT_HD_COLOR.isEmpty()) form.put("gift_hd_color", GIFT_HD_COLOR);

        String body = ZjzApiClient.postForm("idcardv5/get", form);
        if (ZjzApiClient.jsonInt(body, "code") != 0) {
            System.out.println("业务失败：" + body);
            return;
        }

        String picId = ZjzApiClient.jsonString(body, "pic_id");
        System.out.println("取图成功，pic_id：" + (picId == null ? "" : picId));

        Path outputDir = ZjzApiClient.outputDir("v5-get");
        Map<String, String> list = ZjzApiClient.jsonStringMap(body, "list");
        if (list.isEmpty()) {
            System.out.println("返回数据中没有图片。");
        } else {
            for (Map.Entry<String, String> entry : list.entrySet()) {
                Path target = outputDir.resolve(entry.getKey() + ".jpg");
                try {
                    ZjzApiClient.downloadImage(entry.getValue(), target);
                    System.out.println("无水印证件照（" + entry.getKey() + "）：" + target);
                } catch (Exception e) {
                    System.out.println("下载失败：" + e.getMessage());
                }
            }
        }

        // 高清精修照（gift_hd=1 时返回）
        String giftHd = ZjzApiClient.jsonString(body, "gift_hd");
        // gift_hd 为对象（含多底色）或 null 时，jsonString 取不到有效 URL，跳过以免误下载 "null"。
        if (giftHd != null && !giftHd.isEmpty() && !"null".equals(giftHd)) {
            Path target = outputDir.resolve("gift_hd.jpg");
            try {
                ZjzApiClient.downloadImage(giftHd, target);
                System.out.println("高清精修照：" + target);
            } catch (Exception e) {
                System.out.println("高清照下载失败：" + e.getMessage());
            }
        }
    }
}
