import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpHandler;
import com.sun.net.httpserver.HttpServer;

import java.io.IOException;
import java.io.OutputStream;
import java.net.InetSocketAddress;
import java.nio.charset.StandardCharsets;
import java.util.HashMap;
import java.util.Map;
import java.security.MessageDigest;

/**
 * 订单状态通知（回调接收示例）
 *
 * 本示例与其它示例不同：不是主动调用平台接口，而是部署在你的服务器上，
 * 接收平台推送的订单状态通知（订单开始审核、驳回、完成、发货时触发）。
 *
 * 通知以 POST（application/x-www-form-urlencoded）发送到提交订单时填写的 notice_url。
 *
 * 本地联调方法（在 java 目录编译后运行）：
 *   javac -encoding UTF-8 *.java
 *   java ReceiptNoticeDemo
 * 然后把 http://你的IP:8080/ 作为 ReceiptSubmitDemo 的 NOTICE_URL。
 *
 * 注意：本示例不引入第三方依赖，使用 JDK 自带的 HttpServer，接收平台推送并校验签名（sign = md5(API_KEY)）。
 */
public class ReceiptNoticeDemo {

    private static final int PORT = 8080;

    public static void main(String[] args) throws IOException {
        HttpServer server = HttpServer.create(new InetSocketAddress("0.0.0.0", PORT), 0);
        server.createContext("/", new NoticeHandler());
        server.start();
        System.out.println("回执通知回调服务已启动：http://0.0.0.0:" + PORT + "/");
        System.out.println("把它作为 ReceiptSubmitDemo 的 NOTICE_URL。按 Ctrl+C 停止。");
    }

    /** 按平台约定返回处理结果：平台收到 code=0 才认为同步完成。 */
    private static void sendJson(HttpExchange exchange, int code, String msg) throws IOException {
        String json = "{\"data\":{},\"code\":" + code + ",\"msg\":\"" + escape(msg) + "\"}";
        byte[] body = json.getBytes(StandardCharsets.UTF_8);
        exchange.getResponseHeaders().add("Content-Type", "application/json; charset=utf-8");
        exchange.sendResponseHeaders(200, body.length);
        try (OutputStream os = exchange.getResponseBody()) {
            os.write(body);
        }
    }

    /**
     * 标准 JSON 字符串转义。
     *
     * 仅转义 \" 与 \\ 会漏掉换行、制表符等控制字符（如服务端 msg 含换行会让 JSON 非法），
     * 这里按 JSON 规范完整处理：双引号、反斜杠、以及所有 ASCII 控制字符（小于 0x20 的用 Unicode 转义表示）。
     */
    private static String escape(String s) {
        if (s == null) {
            return "";
        }
        StringBuilder sb = new StringBuilder(s.length() + 16);
        for (int i = 0; i < s.length(); i++) {
            char c = s.charAt(i);
            switch (c) {
                case '"':  sb.append("\\\""); break;
                case '\\': sb.append("\\\\"); break;
                case '\n': sb.append("\\n");  break;
                case '\r': sb.append("\\r");  break;
                case '\t': sb.append("\\t");  break;
                case '\b': sb.append("\\b");  break;
                case '\f': sb.append("\\f");  break;
                default:
                    if (c < 0x20) {
                        sb.append(String.format("\\u%04x", (int) c));
                    } else {
                        sb.append(c);
                    }
            }
        }
        return sb.toString();
    }

    /**
     * 计算字符串的 MD5（小写十六进制），用于回执通知的签名校验。
     */
    private static String md5(String s) {
        try {
            MessageDigest md = MessageDigest.getInstance("MD5");
            byte[] bytes = md.digest(s.getBytes(StandardCharsets.UTF_8));
            StringBuilder sb = new StringBuilder(bytes.length * 2);
            for (byte b : bytes) {
                sb.append(String.format("%02x", b));
            }
            return sb.toString();
        } catch (Exception e) {
            return "";
        }
    }

    static class NoticeHandler implements HttpHandler {
        @Override
        public void handle(HttpExchange exchange) throws IOException {
            try {
                if (!"POST".equalsIgnoreCase(exchange.getRequestMethod())) {
                    // 方便联调时做健康检查（平台只会 POST）。
                    sendJson(exchange, 0, "ok");
                    return;
                }

                String raw = new String(exchange.getRequestBody().readAllBytes(), StandardCharsets.UTF_8);
                Map<String, String> post = parseForm(raw);

                // 验证签名：平台约定 sign = md5(API_KEY)，不一致则拒绝接收。
                String expectedSign = md5(Config.API_KEY);
                if (!expectedSign.equalsIgnoreCase(post.get("sign"))) {
                    sendJson(exchange, -1, "签名校验失败");
                    return;
                }

                int status;
                try {
                    status = Integer.parseInt(post.getOrDefault("status", "0"));
                } catch (NumberFormatException e) {
                    status = 0;
                }
                String orderNo = post.getOrDefault("order_no", "");

                // 各状态说明：
                //   2  审核中（开始审核）
                //   3  办理完成，返回 image_receipt / image / image_print
                //   4  打印发货，返回 logistics_info
                //  -1  审核失败，返回 reject_msg（驳回理由）
                //  -2  取消完成
                if (status == 2) {
                    // TODO: 更新订单状态为"审核中"。
                    System.out.println("[receipt notice] 订单 " + orderNo + " 开始审核。");
                } else if (status == 3) {
                    // TODO: 更新订单状态为"办理完成"，保存回执照地址。
                    System.out.println("[receipt notice] 订单 " + orderNo
                            + " 办理完成，回执照：" + post.getOrDefault("image_receipt", ""));
                } else if (status == 4) {
                    // TODO: 更新订单状态为"已发货"，保存物流信息。
                    System.out.println("[receipt notice] 订单 " + orderNo
                            + " 已发货，物流信息：" + post.getOrDefault("logistics_info", "{}"));
                } else if (status == -1) {
                    // TODO: 更新订单状态为"审核驳回"，可参考 ReceiptResubmitDemo 修改资料后重新提交。
                    System.out.println("[receipt notice] 订单 " + orderNo
                            + " 审核驳回：" + post.getOrDefault("reject_msg", ""));
                } else if (status == -2) {
                    // TODO: 更新订单状态为"已取消"。
                    System.out.println("[receipt notice] 订单 " + orderNo + " 已取消。");
                } else {
                    sendJson(exchange, -1, "未知的订单状态：" + status);
                    return;
                }

                // 处理完成后返回 code=0，告知平台通知同步完成。
                sendJson(exchange, 0, "同步完成");
            } catch (Exception e) {
                sendJson(exchange, -1, e.getMessage() == null ? "处理异常" : e.getMessage());
            } finally {
                exchange.close();
            }
        }
    }

    private static Map<String, String> parseForm(String raw) {
        Map<String, String> map = new HashMap<>();
        if (raw == null || raw.isEmpty()) {
            return map;
        }
        for (String pair : raw.split("&")) {
            if (pair.isEmpty()) {
                continue;
            }
            int idx = pair.indexOf('=');
            if (idx < 0) {
                map.put(decode(pair), "");
            } else {
                map.put(decode(pair.substring(0, idx)), decode(pair.substring(idx + 1)));
            }
        }
        return map;
    }

    private static String decode(String s) {
        try {
            return java.net.URLDecoder.decode(s, "UTF-8");
        } catch (Exception e) {
            return s;
        }
    }

}
