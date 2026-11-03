import java.util.Map;

/**
 * 回执订单取消
 *
 * 接口地址：https://api.zjzapi.com/receipt/cancel
 * 请求方式：POST（application/x-www-form-urlencoded）
 * 返回格式：JSON
 * 说明：仅“待处理”状态的订单可取消
 *
 * 用法（在 java 目录编译后运行）：
 *   javac -encoding UTF-8 *.java
 *   java ReceiptCancelDemo
 */
public class ReceiptCancelDemo {

    // ===== 请求参数（修改后直接运行即可）=====
    private static final String ORDER_NO = "";  // 待取消的订单号（来自 ReceiptSubmitDemo）

    public static void main(String[] args) throws Exception {
        if (Config.API_KEY.isEmpty() || "YOUR_APP_KEY".equals(Config.API_KEY)) {
            System.err.println("提示：请先在 Config.java 中将 API_KEY 替换为你的 API Key。");
        }
        if (ORDER_NO.isEmpty()) {
            System.err.println("请先在 ORDER_NO 中填入待取消的订单号（来自 ReceiptSubmitDemo）。");
            return;
        }

        Map<String, String> form = ZjzApiClient.form(
                "key", Config.API_KEY,
                "order_no", ORDER_NO
        );

        String body = ZjzApiClient.postForm("receipt/cancel", form);
        if (ZjzApiClient.jsonInt(body, "code") != 0) {
            System.out.println("业务失败：" + body);
            return;
        }
        System.out.println("回执订单已取消：" + body);
    }
}
