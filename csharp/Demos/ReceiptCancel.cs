using System;
using System.Collections.Generic;
using System.Threading.Tasks;

namespace IdPhotoDemo.Demos;

/// <summary>
/// 订单取消。
///
/// 接口地址：https://api.zjzapi.com/receipt/cancel
/// 请求方式：POST（application/x-www-form-urlencoded）
/// 返回格式：JSON
/// 接口费用：调用免费
///
/// 用法（在 csharp 目录运行）：dotnet run -- receipt-cancel
/// 注意：订单在"待处理"状态下才可以取消。
/// </summary>
public static class ReceiptCancel
{
    // ===== 请求参数（修改后直接运行即可）=====
    private const string OrderNo = "";   // 要取消的订单号（receipt-submit 的返回结果）

    public static async Task RunAsync()
    {
        if (string.IsNullOrEmpty(Config.ApiKey) || Config.ApiKey == "YOUR_APP_KEY")
            Console.Error.WriteLine("提示：请先在 csharp/Config.cs 中将 ApiKey 替换为你的 API Key。");
        if (string.IsNullOrEmpty(OrderNo))
        {
            Console.Error.WriteLine("请先在参数区填写要取消的订单号 order_no。");
            return;
        }

        var result = await ZjzApi.PostAsync("receipt/cancel", new Dictionary<string, string>
        {
            ["key"] = Config.ApiKey,
            ["order_no"] = OrderNo,
        });
        if (!ZjzApi.IsSuccess(result))
        {
            Console.WriteLine($"业务失败：code={result.Code}，msg={result.Msg}");
            return;
        }

        Console.WriteLine($"订单 {OrderNo} 取消成功。");
    }
}
