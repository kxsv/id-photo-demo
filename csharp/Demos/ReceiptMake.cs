using System;
using System.Collections.Generic;
using System.IO;
using System.Threading.Tasks;

namespace IdPhotoDemo.Demos;

/// <summary>
/// 证件照回执检测+制作。
///
/// 接口地址：https://api.zjzapi.com/receipt/make
/// 请求方式：POST（application/x-www-form-urlencoded）
/// 返回格式：JSON
/// 接口费用：调用免费（回执提交办理完成后按单扣费）。
///
/// 用法（在 csharp 目录运行）：dotnet run -- receipt-make
/// </summary>
public static class ReceiptMake
{
    // ===== 请求参数（修改后直接运行即可）=====
    private const string ImagePath = "../images/test.jpeg"; // 待处理图片路径（相对 csharp 目录）
    private const string ItemId = "1";                      // 回执类规格id，可通过 item-list 示例获取

    public static async Task RunAsync()
    {
        if (string.IsNullOrEmpty(Config.ApiKey) || Config.ApiKey == "YOUR_APP_KEY")
            Console.Error.WriteLine("提示：请先在 csharp/Config.cs 中将 ApiKey 替换为你的 API Key。");

        var result = await ZjzApi.PostAsync("receipt/make", new Dictionary<string, string>
        {
            ["key"] = Config.ApiKey,
            ["item_id"] = ItemId,
            ["image"] = await ZjzApi.ImageToBase64Async(ImagePath),
        });
        if (!ZjzApi.IsSuccess(result))
        {
            Console.WriteLine($"业务失败：code={result.Code}，msg={result.Msg}");
            return;
        }

        Console.WriteLine($"回执照制作成功，pic_id：{ZjzApi.GetString(result.Data, "pic_id")}");
        Console.WriteLine("提示：把 pic_id 填入 receipt-submit 提交办理订单。");
        if (result.Data.TryGetProperty("list", out var listValue))
            await ZjzApi.SaveImagesAsync(Path.Combine("output", "receipt-make"), listValue);
    }
}
