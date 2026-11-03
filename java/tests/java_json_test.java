import java.util.Map;

/** 本地回归：只验证 JSON 解析，不读取配置或请求 API。 */
public class java_json_test {
    private static void assert_equal(Object expected, Object actual) {
        if (!expected.equals(actual)) {
            throw new AssertionError("Expected " + expected + ", got " + actual);
        }
    }

    public static void main(String[] args) {
        String json = "{\"code\":0,\"data\":{\"list\":{"
                + "\"white\":\"https://example.com/white.jpg\","
                + "\"blue\":\"https:\\/\\/example.com\\/blue.jpg?a=1\\u0026b=2\","
                + "\"\\u7ea2\":\"\\u0068ttps:\\/\\/example.com\\/red.jpg\","
                + "\"note\":\"not an image URL }\"}}}";
        Map<String, String> images = ZjzApiClient.jsonStringMap(json, "list");
        assert_equal(3, images.size());
        assert_equal("https://example.com/white.jpg", images.get("white"));
        assert_equal("https://example.com/blue.jpg?a=1&b=2", images.get("blue"));
        assert_equal("https://example.com/red.jpg", images.get("红"));
        assert_equal("https://example.com/a.jpg",
                ZjzApiClient.jsonString("{\"image\":\"https:\\/\\/example.com\\/a.jpg\"}", "image"));
        assert_equal("\"\\/\b\f\n\r\t中文😀",
                ZjzApiClient.jsonString("{\"value\":\"\\\"\\\\\\/\\b\\f\\n\\r\\t\\u4e2d\\u6587\\ud83d\\ude00\"}", "value"));
        assert_equal("\\u4e2d", ZjzApiClient.jsonString("{\"value\":\"\\\\u4e2d\"}", "value"));
        assert_equal("12345", ZjzApiClient.jsonString("{\"pic_id\":12345}", "pic_id"));
        assert_equal(0, ZjzApiClient.jsonStringMap("{\"data\":{}}", "list").size());
        assert_equal(0, ZjzApiClient.jsonStringMap("{\"list\":{", "list").size());
        System.out.println("Java JSON regression checks passed.");
    }
}
