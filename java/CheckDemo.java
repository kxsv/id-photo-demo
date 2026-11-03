import java.util.Map;

/**
 * 证件照合规性检测
 *
 * 接口地址：https://api.zjzapi.com/idphoto/check
 * 请求方式：POST（application/x-www-form-urlencoded）
 * 返回格式：JSON
 * 接口费用：调用免费
 *
 * 两种模式：
 *   1) 按规格检测：ITEM_ID 设为某规格 id（非 0），服务端按该规格要求检测；
 *   2) 自定义参数检测：ITEM_ID 设为 "0"，下方 WIDTH_PX/HEIGHT_PX/... 等自定义项生效。
 *
 * 用法（在 java 目录编译后运行）：
 *   javac -encoding UTF-8 *.java
 *   java CheckDemo
 */
public class CheckDemo {

    // ===== 请求参数（修改后直接运行即可）=====
    private static final String IMAGE_PATH = "../images/test.jpeg";
    private static final String ITEM_ID = "1";   // 非 0：按规格模式；设为 "0"：走自定义参数模式

    // ===== 自定义参数模式（ITEM_ID="0" 时生效）=====
    private static final String WIDTH_PX = "295";
    private static final String HEIGHT_PX = "413";
    private static final String DPI = "300";
    private static final String FILE_FORMAT = "jpg";
    private static final String BG_COLORS = "white,blue";
    private static final String SOLID_BG = "1";
    private static final String FACE_WIDTH_RATIO_MIN = "0.2";
    private static final String FACE_WIDTH_RATIO_MAX = "0.4";
    private static final String SHOULDER_HORIZONTAL = "1";
    private static final String FACE_POSE_YAW = "1";
    private static final String FACE_CLARITY = "1";
    private static final String FACE_BRIGHTNESS = "1";
    private static final String FACE_MASK = "0";
    private static final String GLASSES = "0";

    public static void main(String[] args) throws Exception {
        if (Config.API_KEY.isEmpty() || "YOUR_APP_KEY".equals(Config.API_KEY)) {
            System.err.println("提示：请先在 Config.java 中将 API_KEY 替换为你的 API Key。");
        }

        Map<String, String> form = ZjzApiClient.form(
                "key", Config.API_KEY,
                "image", ZjzApiClient.imageToBase64(IMAGE_PATH)
        );
        if (!"0".equals(ITEM_ID)) {
            form.put("item_id", ITEM_ID);
        } else {
            form.put("width_px", WIDTH_PX);
            form.put("height_px", HEIGHT_PX);
            form.put("dpi", DPI);
            form.put("file_format", FILE_FORMAT);
            form.put("bg_colors", BG_COLORS);
            form.put("solid_bg", SOLID_BG);
            form.put("face_width_ratio_min", FACE_WIDTH_RATIO_MIN);
            form.put("face_width_ratio_max", FACE_WIDTH_RATIO_MAX);
            form.put("shoulder_horizontal", SHOULDER_HORIZONTAL);
            form.put("face_pose_yaw", FACE_POSE_YAW);
            form.put("face_clarity", FACE_CLARITY);
            form.put("face_brightness", FACE_BRIGHTNESS);
            form.put("face_mask", FACE_MASK);
            form.put("glasses", GLASSES);
        }

        String body = ZjzApiClient.postForm("idphoto/check", form);
        if (ZjzApiClient.jsonInt(body, "code") != 0) {
            System.out.println("业务失败：" + body);
            return;
        }
        System.out.println("检测通过项 / 状态：" + body);
    }
}
