import java.util.Map;

/**
 * 获取应用信息（账号相关接口，需 access_token）
 *
 * 接口地址：https://api.zjzapi.com/user/app
 * 请求方式：POST（application/x-www-form-urlencoded）
 * 返回格式：JSON
 * 说明：本接口使用 access_token + key
 *
 * 用法（在 java 目录编译后运行）：
 *   javac -encoding UTF-8 *.java
 *   java AccountAppDemo
 */
public class AccountAppDemo {

    public static void main(String[] args) throws Exception {
        if (Config.ACCESS_TOKEN.isEmpty() || "YOUR_ACCESS_TOKEN".equals(Config.ACCESS_TOKEN)) {
            System.err.println("提示：请先在 Config.java 中将 ACCESS_TOKEN 替换为你的 access_token。");
            return;
        }

        Map<String, String> form = ZjzApiClient.form(
                "access_token", Config.ACCESS_TOKEN,
                "key", Config.API_KEY
        );

        String body = ZjzApiClient.postForm("user/app", form);
        if (ZjzApiClient.jsonInt(body, "code") != 0) {
            System.out.println("业务失败：" + body);
            return;
        }
        System.out.println("应用信息：" + body);
    }
}
