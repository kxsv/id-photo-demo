using System;
using System.Collections.Generic;
using System.IO;
using System.Net;
using System.Text;
using System.Text.Json;
using System.Security.Cryptography;
using System.Threading.Tasks;

namespace IdPhotoDemo.Demos;

/// <summary>
/// 订单状态通知（回调接收示例）。
///
/// 本示例与其它示例不同：不是主动调用平台接口，而是部署在你的服务器上，
/// 接收平台推送的订单状态通知（订单开始审核、驳回、完成、发货时触发）。
///
/// 通知以 POST（application/x-www-form-urlencoded）发送到提交订单时填写的 notice_url。
///
/// 本地联调方法（在 csharp 目录运行）：dotnet run -- receipt-notice
/// 然后把 http://你的IP:8080/ 作为 receipt-submit 的 NoticeUrl。
///
/// 注意：本示例是独立 HTTP 服务，使用 .NET 自带的 HttpListener，接收推送并校验签名（sign = md5(ApiKey)）。
/// Windows 下 HttpListener 可能需要管理员权限或以管理员身份运行；若启动报权限错误，可执行：
///   netsh http add urlacl url=http://+:8080/ user=你的Windows用户名
/// </summary>
public static class ReceiptNotice
{
    private const int Port = 8080;
    // + 接收通过本机各 IP/域名发来的回调；也可改成实际对外提供服务的域名或 IP。
    private const string listener_host = "+";

    public static async Task RunAsync()
    {
        if (string.IsNullOrEmpty(Config.ApiKey) || Config.ApiKey == "YOUR_APP_KEY")
            Console.Error.WriteLine("提示：请先在 csharp/Config.cs 中将 ApiKey 替换为你的 API Key。");

        var prefix = $"http://{listener_host}:{Port}/";
        var listener = new HttpListener();
        listener.Prefixes.Add(prefix);
        try
        {
            listener.Start();
        }
        catch (HttpListenerException ex)
        {
            Console.Error.WriteLine($"启动回调服务失败（可能需要管理员权限）：{ex.Message}");
            Console.Error.WriteLine($"请以管理员身份运行，或执行：netsh http add urlacl url={prefix} user=你的Windows用户名");
            return;
        }

        Console.WriteLine($"回执通知回调服务已启动，监听前缀：{prefix}");
        Console.WriteLine($"请将 http://你的服务器公网IP或域名:{Port}/ 填入 receipt-submit 的 NoticeUrl。按 Ctrl+C 停止。");

        try
        {
            while (listener.IsListening)
            {
                var context = await listener.GetContextAsync();
                _ = Task.Run(() => HandleAsync(context));
            }
        }
        catch (Exception ex) when (ex is HttpListenerException || ex is OperationCanceledException)
        {
            // 服务被停止（如 Ctrl+C）。
        }
        finally
        {
            listener.Close();
        }
    }

    private static async Task HandleAsync(HttpListenerContext context)
    {
        try
        {
            if (context.Request.HttpMethod != "POST")
            {
                // 方便联调时做健康检查（平台只会 POST）。
                await ReplyAsync(context, 0, "ok");
                return;
            }

            using var reader = new StreamReader(context.Request.InputStream, context.Request.ContentEncoding);
            var body = await reader.ReadToEndAsync();
            var post = ParseForm(body);

            // 验证签名：平台约定 sign = md5(ApiKey)，不一致则拒绝接收。
            using var md5 = MD5.Create();
            var expectedBytes = md5.ComputeHash(Encoding.UTF8.GetBytes(Config.ApiKey));
            var expectedSign = Convert.ToHexString(expectedBytes).ToLowerInvariant();
            if (post.GetValueOrDefault("sign", "").ToLowerInvariant() != expectedSign)
            {
                await ReplyAsync(context, -1, "签名校验失败");
                return;
            }

            int.TryParse(post.GetValueOrDefault("status", "0"), out var status);
            var orderNo = post.GetValueOrDefault("order_no", "");

            // 各状态说明：
            //   2  审核中（开始审核）
            //   3  办理完成，返回 image_receipt / image / image_print
            //   4  打印发货，返回 logistics_info
            //  -1  审核失败，返回 reject_msg（驳回理由）
            //  -2  取消完成
            if (status == 2)
            {
                // TODO: 更新订单状态为"审核中"。
                Console.WriteLine($"[receipt notice] 订单 {orderNo} 开始审核。");
            }
            else if (status == 3)
            {
                // TODO: 更新订单状态为"办理完成"，保存回执照地址。
                Console.WriteLine($"[receipt notice] 订单 {orderNo} 办理完成，回执照：{post.GetValueOrDefault("image_receipt", "")}");
            }
            else if (status == 4)
            {
                // TODO: 更新订单状态为"已发货"，保存物流信息。
                Console.WriteLine($"[receipt notice] 订单 {orderNo} 已发货，物流信息：{post.GetValueOrDefault("logistics_info", "{}")}");
            }
            else if (status == -1)
            {
                // TODO: 更新订单状态为"审核驳回"，可参考 receipt-resubmit 修改资料后重新提交。
                Console.WriteLine($"[receipt notice] 订单 {orderNo} 审核驳回：{post.GetValueOrDefault("reject_msg", "")}");
            }
            else if (status == -2)
            {
                // TODO: 更新订单状态为"已取消"。
                Console.WriteLine($"[receipt notice] 订单 {orderNo} 已取消。");
            }
            else
            {
                await ReplyAsync(context, -1, $"未知的订单状态：{status}");
                return;
            }

            // 处理完成后返回 code=0，告知平台通知同步完成。
            await ReplyAsync(context, 0, "同步完成");
        }
        catch (Exception ex)
        {
            await ReplyAsync(context, -1, ex.Message ?? "处理异常");
        }
    }

    /// <summary>按平台约定返回处理结果：平台收到 code=0 才认为同步完成。</summary>
    private static async Task ReplyAsync(HttpListenerContext context, int code, string msg)
    {
        // 用 JsonSerializer 序列化，自动转义换行等控制字符，避免手工拼接产生非法 JSON。
        var payload = new { data = new { }, code, msg };
        var json = JsonSerializer.Serialize(payload);
        var bytes = Encoding.UTF8.GetBytes(json);
        context.Response.ContentType = "application/json; charset=utf-8";
        context.Response.StatusCode = 200;
        context.Response.ContentLength64 = bytes.Length;
        await context.Response.OutputStream.WriteAsync(bytes);
        context.Response.Close();
    }

    private static Dictionary<string, string> ParseForm(string body)
    {
        var dict = new Dictionary<string, string>();
        if (string.IsNullOrEmpty(body))
            return dict;
        foreach (var pair in body.Split('&'))
        {
            if (string.IsNullOrEmpty(pair))
                continue;
            var kv = pair.Split('=', 2);
            var key = Uri.UnescapeDataString(kv[0]);
            var value = kv.Length > 1 ? Uri.UnescapeDataString(kv[1]) : "";
            dict[key] = value;
        }
        return dict;
    }

}
