import java.util.Map;

/**
 * 获取账号信息（账号相关接口，需 access_token）
 *
 * 接口地址：https://api.zjzapi.com/user/info
 * 请求方式：POST（application/x-www-form-urlencoded）
 * 返回格式：JSON
 * 说明：本接口不使用 key，仅用 access_token
 *
 * 用法（在 java 目录编译后运行）：
 *   javac -encoding UTF-8 *.java
 *   java AccountInfoDemo
 */
public class AccountInfoDemo {

    public static void main(String[] args) throws Exception {
        if (Config.ACCESS_TOKEN.isEmpty() || "YOUR_ACCESS_TOKEN".equals(Config.ACCESS_TOKEN)) {
            System.err.println("提示：请先在 Config.java 中将 ACCESS_TOKEN 替换为你的 access_token。");
            return;
        }

        Map<String, String> form = ZjzApiClient.form(
                "access_token", Config.ACCESS_TOKEN
        );

        String body = ZjzApiClient.postForm("user/info", form);
        if (ZjzApiClient.jsonInt(body, "code") != 0) {
            System.out.println("业务失败：" + body);
            return;
        }
        System.out.println("账号信息：" + body);
    }
}
