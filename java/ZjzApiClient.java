import java.io.File;
import java.io.FileNotFoundException;
import java.io.IOException;
import java.net.URI;
import java.net.URLEncoder;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.time.Duration;
import java.util.Base64;
import java.util.LinkedHashMap;
import java.util.Map;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
import java.util.stream.Collectors;

/**
 * 鹧应证件照 API 公共客户端。
 *
 * 所有接口示例共用的请求封装：
 * - postForm()       以 POST 表单方式调用接口
 * - imageToBase64()  读取本地图片并转 Base64
 * - downloadImage()  下载接口返回的结果图到本地
 * - jsonInt() / jsonStringMap()  从返回 JSON 中提取业务字段
 *
 * Demo 不引入第三方依赖，业务字段使用正则提取；
 * 生产环境建议引入 Jackson/Gson 等成熟 JSON 库后再做解析。
 */
public final class ZjzApiClient {

    /** API 网关地址 */
    public static final String BASE_URL = "https://api.zjzapi.com";

    private static final HttpClient HTTP = HttpClient.newBuilder()
            .connectTimeout(Duration.ofSeconds(10))
            .followRedirects(HttpClient.Redirect.NORMAL)
            .build();

    private ZjzApiClient() {
    }

    /**
     * 调用接口，返回 JSON 字符串
     *
     * @param path 接口路径，如 idcardv3/all
     * @param form 请求参数（已含 key 等）
     */
    public static String postForm(String path, Map<String, String> form) throws IOException, InterruptedException {
        String body = form.entrySet().stream()
                .map(entry -> URLEncoder.encode(entry.getKey(), StandardCharsets.UTF_8)
                        + "=" + URLEncoder.encode(entry.getValue(), StandardCharsets.UTF_8))
                .collect(Collectors.joining("&"));
        HttpRequest request = HttpRequest.newBuilder(URI.create(BASE_URL + "/" + path))
                .timeout(Duration.ofSeconds(60))
                .header("Content-Type", "application/x-www-form-urlencoded")
                .POST(HttpRequest.BodyPublishers.ofString(body))
                .build();
        HttpResponse<String> response = HTTP.send(request, HttpResponse.BodyHandlers.ofString(StandardCharsets.UTF_8));
        if (response.statusCode() < 200 || response.statusCode() >= 300) {
            throw new IOException("HTTP " + response.statusCode() + "：" + response.body());
        }
        return response.body();
    }

    /**
     * 读取本地图片并转 Base64（不带 data:image/... 前缀）
     */
    public static String imageToBase64(String path) throws IOException {
        File file = new File(path);
        if (!file.isFile()) {
            throw new FileNotFoundException("图片不存在：" + path);
        }
        byte[] content = Files.readAllBytes(file.toPath());
        if (content.length > 8 * 1024 * 1024) {
            throw new IOException("图片不能超过 8M：" + path);
        }
        return Base64.getEncoder().encodeToString(content);
    }

    /**
     * 下载接口返回的结果图到本地
     *
     * @param url    图片地址
     * @param target 保存路径（父目录自动创建）
     */
    public static void downloadImage(String url, Path target) throws IOException, InterruptedException {
        if (target.getParent() != null) {
            Files.createDirectories(target.getParent());
        }
        HttpRequest request = HttpRequest.newBuilder(URI.create(url))
                .timeout(Duration.ofSeconds(60))
                .GET()
                .build();
        HttpResponse<Path> response = HTTP.send(request, HttpResponse.BodyHandlers.ofFile(target));
        if (response.statusCode() < 200 || response.statusCode() >= 300) {
            throw new IOException("下载失败：HTTP " + response.statusCode());
        }
    }

    /**
     * 从返回 JSON 中提取整型字段（找不到时返回 Integer.MIN_VALUE）
     *
     * 用于业务判断，如 code == 0 表示业务成功。
     */
    public static int jsonInt(String json, String field) {
        Matcher matcher = Pattern.compile("\"" + field + "\"\\s*:\\s*(-?\\d+)").matcher(json);
        return matcher.find() ? Integer.parseInt(matcher.group(1)) : Integer.MIN_VALUE;
    }

    /**
     * 从返回 JSON 中提取 "list" 形式的对象字段，返回有序键值对
     *
     * 形如 {"list": {"white": "http://...", "blue": "http://..."}} 的结构，
     * 键值为图片地址时才会被提取。
     */
    public static Map<String, String> jsonStringMap(String json, String field) {
        Map<String, String> values = new LinkedHashMap<>();
        Matcher start = Pattern.compile("\"" + Pattern.quote(field) + "\"\\s*:\\s*\\{").matcher(json);
        if (!start.find()) {
            return values;
        }
        // 用括号匹配而非 indexOf('}')，避免对象内含嵌套结构或字面量 '}' 时截断。
        int close = matchBrace(json, start.end() - 1);
        if (close < 0) {
            return values;
        }
        // 先提取完整 JSON 字符串，再解码；合法 URL 也可能写成 https:\/\/...。
        Matcher pair = Pattern.compile("\"((?:\\\\.|[^\"\\\\])*)\"\\s*:\\s*\"((?:\\\\.|[^\"\\\\])*)\"")
                .matcher(json.substring(start.end(), close));
        while (pair.find()) {
            String image_url = decode_json_string(pair.group(2));
            if (image_url.startsWith("http://") || image_url.startsWith("https://")) {
                values.put(decode_json_string(pair.group(1)), image_url);
            }
        }
        return values;
    }

    /** 返回 openIndex 处 '{' 对应的 '}' 下标（正确跳过字符串内容），未找到返回 -1 */
    private static int matchBrace(String json, int openIndex) {
        boolean inString = false;
        int depth = 0;
        for (int i = openIndex; i < json.length(); i++) {
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

    /**
     * 从返回 JSON 中提取字符串字段（找不到时返回 null）
     */
    public static String jsonString(String json, String field) {
        if (json == null) return null;
        // 兼容字符串 / 数字 / 布尔 / null 取值（如 pic_id、order_no 为数字）
        Pattern p = Pattern.compile("(?s)\"" + Pattern.quote(field)
                + "\"\\s*:\\s*(\"(?:\\\\.|[^\"\\\\])*\"|[-0-9.]+|true|false|null)");
        Matcher matcher = p.matcher(json);
        if (!matcher.find()) return null;
        String token = matcher.group(1);
        if (token.startsWith("\"")) {
            return decode_json_string(token.substring(1, token.length() - 1));
        }
        return token; // 数字 / 布尔 / null 原样返回
    }

    /** 解码 JSON 字符串内容，支持斜杠、控制字符及 Unicode 转义。 */
    private static String decode_json_string(String value) {
        StringBuilder decoded = new StringBuilder();
        for (int index = 0; index < value.length(); index++) {
            char current = value.charAt(index);
            if (current != '\\') {
                decoded.append(current);
                continue;
            }
            if (++index >= value.length()) {
                throw new IllegalArgumentException("不完整的 JSON 字符串转义");
            }
            switch (value.charAt(index)) {
                case '"': decoded.append('"'); break;
                case '\\': decoded.append('\\'); break;
                case '/': decoded.append('/'); break;
                case 'b': decoded.append('\b'); break;
                case 'f': decoded.append('\f'); break;
                case 'n': decoded.append('\n'); break;
                case 'r': decoded.append('\r'); break;
                case 't': decoded.append('\t'); break;
                case 'u':
                    if (index + 4 >= value.length()) {
                        throw new IllegalArgumentException("不完整的 JSON Unicode 转义");
                    }
                    decoded.append((char) Integer.parseInt(value.substring(index + 1, index + 5), 16));
                    index += 4;
                    break;
                default:
                    throw new IllegalArgumentException("无效的 JSON 字符串转义");
            }
        }
        return decoded.toString();
    }

    /** 便捷方法：拼接表单参数 */
    public static Map<String, String> form(Object... pairs) {
        Map<String, String> form = new LinkedHashMap<>();
        for (int i = 0; i + 1 < pairs.length; i += 2) {
            form.put(String.valueOf(pairs[i]), String.valueOf(pairs[i + 1]));
        }
        return form;
    }

    /** 常用保存目录：output/xxx（相对运行目录） */
    public static Path outputDir(String name) {
        return Paths.get("output", name);
    }
}
