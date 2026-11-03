import java.util.Map;

/**
 * 证件照规格详情
 *
 * 接口地址：https://api.zjzapi.com/item/get
 * 请求方式：POST（application/x-www-form-urlencoded）
 * 返回格式：JSON
 * 接口费用：调用免费
 *
 * 用法（在 java 目录编译后运行）：
 *   javac -encoding UTF-8 *.java
 *   java ItemDetailDemo
 */
public class ItemDetailDemo {

    // ===== 请求参数（修改后直接运行即可）=====
    private static final String ITEM_ID = "1";  // 规格 id，可通过 ItemListDemo 获取

    public static void main(String[] args) throws Exception {
        if (Config.API_KEY.isEmpty() || "YOUR_APP_KEY".equals(Config.API_KEY)) {
            System.err.println("提示：请先在 Config.java 中将 API_KEY 替换为你的 API Key。");
        }

        Map<String, String> form = ZjzApiClient.form(
                "key", Config.API_KEY,
                "item_id", ITEM_ID
        );

        String body = ZjzApiClient.postForm("item/get", form);
        if (ZjzApiClient.jsonInt(body, "code") != 0) {
            System.out.println("业务失败：" + body);
            return;
        }

        String name = ZjzApiClient.jsonString(body, "name");
        String itemId = ZjzApiClient.jsonString(body, "item_id");
        String color = ZjzApiClient.jsonString(body, "color");
        System.out.println("规格详情：item_id=" + itemId + "，名称=" + name + "，支持底色=" + color);
        System.out.println("完整返回：" + body);
    }
}
