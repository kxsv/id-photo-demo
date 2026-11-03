import java.util.Map;

/**
 * 回执驳回后重新提交
 *
 * 接口地址：https://api.zjzapi.com/receipt/submit（复用提交接口，额外带 order_no）
 * 请求方式：POST（application/x-www-form-urlencoded）
 * 返回格式：JSON
 *
 * 用法（在 java 目录编译后运行）：
 *   javac -encoding UTF-8 *.java
 *   java ReceiptResubmitDemo
 */
public class ReceiptResubmitDemo {

    // ===== 请求参数（修改后直接运行即可）=====
    private static final String PIC_ID = "";     // 回执制作返回的 pic_id
    private static final String ORDER_NO = "";   // 原订单号（驳回后重新提交）
    private static final String NOTICE_URL = "";  // 回执回调地址
    private static final String PARAM = "";       // 回执附加参数（json 字符串），留空不传

    public static void main(String[] args) throws Exception {
        if (Config.API_KEY.isEmpty() || "YOUR_APP_KEY".equals(Config.API_KEY)) {
            System.err.println("提示：请先在 Config.java 中将 API_KEY 替换为你的 API Key。");
        }
        if (PIC_ID.isEmpty() || ORDER_NO.isEmpty() || NOTICE_URL.isEmpty()) {
            System.err.println("请先在 PIC_ID / ORDER_NO / NOTICE_URL 中填入参数。");
            return;
        }

        Map<String, String> form = ZjzApiClient.form(
                "key", Config.API_KEY,
                "pic_id", PIC_ID,
                "order_no", ORDER_NO,
                "notice_url", NOTICE_URL
        );
        if (!PARAM.isEmpty()) form.put("param", PARAM);

        String body = ZjzApiClient.postForm("receipt/submit", form);
        if (ZjzApiClient.jsonInt(body, "code") != 0) {
            System.out.println("业务失败：" + body);
            return;
        }

        String orderNo = ZjzApiClient.jsonString(body, "order_no");
        System.out.println("回执重新提交成功，order_no：" + (orderNo == null ? "" : orderNo));
    }
}
