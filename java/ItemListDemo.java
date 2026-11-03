import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * 证件照规格列表
 *
 * 接口地址：https://api.zjzapi.com/item/list
 * 请求方式：POST（application/x-www-form-urlencoded）
 * 返回格式：JSON
 * 接口费用：调用免费
 *
 * 用法（在 java 目录编译后运行）：
 *   javac -encoding UTF-8 *.java
 *   java ItemListDemo
 *
 * 输出各规格的 item_id、名称、像素尺寸与支持底色，
 * 把 item_id 和对应底色填入制作类示例（如 V3AllDemo）即可。
 */
public class ItemListDemo {

    public static void main(String[] args) throws Exception {
        if (Config.API_KEY.isEmpty() || "YOUR_APP_KEY".equals(Config.API_KEY)) {
            System.err.println("提示：请先在 Config.java 中将 API_KEY 替换为你的 API Key。");
        }

        String body = ZjzApiClient.postForm("item/list", ZjzApiClient.form("key", Config.API_KEY));

        // HTTP 成功不代表业务成功，需再检查返回 JSON 中的 code。
        if (ZjzApiClient.jsonInt(body, "code") != 0) {
            System.out.println("业务失败：" + body);
            return;
        }

        List<Map<String, String>> items = parseItems(body);
        if (items.isEmpty()) {
            System.out.println("返回数据中没有规格列表。");
            return;
        }

        // 打印规格清单，便于挑选 item_id 与底色。
        System.out.println("item_id  名称            像素          支持底色");
        for (Map<String, String> item : items) {
            System.out.println(String.format(
                    "%-8s %-14s %-12s %s",
                    item.getOrDefault("item_id", ""),
                    item.getOrDefault("name", ""),
                    item.getOrDefault("width_px", "") + "x" + item.getOrDefault("height_px", ""),
                    item.getOrDefault("color", "")
            ));
        }
        System.out.println("共 " + items.size() + " 个规格。");
    }

    /**
     * 从返回 JSON 中解析 data.list 内的规格对象。
     *
     * 使用转义感知的字符扫描（值里可能含有 " 与 [] 等字符，不能简单按括号截取）；
     * Demo 不引入第三方 JSON 库，生产环境建议使用 Jackson/Gson 等成熟库后再做解析。
     */
    private static List<Map<String, String>> parseItems(String json) {
        List<Map<String, String>> items = new ArrayList<>();
        Matcher start = Pattern.compile("\"list\"\\s*:\\s*\\[").matcher(json);
        if (!start.find()) {
            return items;
        }

        int i = start.end();
        while (i < json.length()) {
            // 找到下一个对象开头的 '{'。
            while (i < json.length() && json.charAt(i) != '{') {
                i++;
            }
            if (i >= json.length()) {
                break;
            }
            int objectEnd = findObjectEnd(json, i);
            if (objectEnd < 0) {
                break;
            }
            String itemJson = json.substring(i + 1, objectEnd);
            Map<String, String> item = new LinkedHashMap<>();
            for (String field : new String[]{"item_id", "name", "color", "width_px", "height_px"}) {
                item.put(field, fieldValue(itemJson, field));
            }
            items.add(item);
            i = objectEnd + 1;
        }
        return items;
    }

    /** 返回 position 处 '{' 匹配的 '}' 下标（正确跳过字符串内容），未找到返回 -1 */
    private static int findObjectEnd(String json, int begin) {
        boolean inString = false;
        int depth = 0;
        for (int i = begin; i < json.length(); i++) {
            char c = json.charAt(i);
            if (inString) {
                if (c == '\\') {
                    i++;
                } else if (c == '"') {
                    inString = false;
                }
                continue;
            }
            if (c == '"') {
                inString = true;
            } else if (c == '{') {
                depth++;
            } else if (c == '}') {
                depth--;
                if (depth == 0) {
                    return i;
                }
            }
        }
        return -1;
    }

    /** 提取字段值，兼容字符串与数字两种写法，并还原转义字符 */
    private static String fieldValue(String json, String field) {
        Matcher value = Pattern.compile("\"" + field + "\"\\s*:\\s*\"((?:[^\"\\\\]|\\\\.)*)\"").matcher(json);
        if (value.find()) {
            return value.group(1).replace("\\\"", "\"").replace("\\\\", "\\");
        }
        value = Pattern.compile("\"" + field + "\"\\s*:\\s*(-?\\d+)").matcher(json);
        return value.find() ? value.group(1) : "";
    }
}
