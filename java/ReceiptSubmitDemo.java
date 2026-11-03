import java.util.Map;

/**
 * 回执提交
 *
 * 接口地址：https://api.zjzapi.com/receipt/submit
 * 请求方式：POST（application/x-www-form-urlencoded）
 * 返回格式：JSON
 * 接口费用：按次扣费
 *
 * 用法（在 java 目录编译后运行）：
 *   javac -encoding UTF-8 *.java
 *   java ReceiptSubmitDemo
 */
public class ReceiptSubmitDemo {

    // ===== 请求参数（修改后直接运行即可）=====
    private static final String PIC_ID = "";     // 回执制作返回的 pic_id（有效期50分钟）
    private static final String NOTICE_URL = "";  // 回执回调地址（服务端回执办理完成后通知）
    private static final String PARAM = "";      // 回执附加参数（json 字符串），留空不传
    private static final String IS_PRINT = "0";  // 是否打印：0=否 1=是
    private static final String PRINT_PARAM = ""; // 打印参数（IS_PRINT=1 时填写）

    public static void main(String[] args) throws Exception {
        if (Config.API_KEY.isEmpty() || "YOUR_APP_KEY".equals(Config.API_KEY)) {
            System.err.println("提示：请先在 Config.java 中将 API_KEY 替换为你的 API Key。");
        }
        if (PIC_ID.isEmpty() || NOTICE_URL.isEmpty()) {
            System.err.println("请先在 PIC_ID / NOTICE_URL 中填入参数（PIC_ID 来自 ReceiptMakeDemo）。");
            return;
        }

        Map<String, String> form = ZjzApiClient.form(
                "key", Config.API_KEY,
                "pic_id", PIC_ID,
                "notice_url", NOTICE_URL
        );
        if (!PARAM.isEmpty()) form.put("param", PARAM);
        if ("1".equals(IS_PRINT)) {
            form.put("is_print", "1");
            if (!PRINT_PARAM.isEmpty()) form.put("print_param", PRINT_PARAM);
        }

        String body = ZjzApiClient.postForm("receipt/submit", form);
        if (ZjzApiClient.jsonInt(body, "code") != 0) {
            System.out.println("业务失败：" + body);
            return;
        }

        String orderNo = ZjzApiClient.jsonString(body, "order_no");
        System.out.println("回执提交成功，order_no（用于 ReceiptResubmitDemo / ReceiptCancelDemo）：" + (orderNo == null ? "" : orderNo));
    }
}
