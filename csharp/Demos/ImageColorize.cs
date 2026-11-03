using System;
using System.Collections.Generic;
using System.IO;
using System.Threading.Tasks;

namespace IdPhotoDemo.Demos;

/// <summary>
/// 黑白图片上色。
///
/// 接口地址：https://api.zjzapi.com/image/colorize
/// 请求方式：POST（application/x-www-form-urlencoded）
/// 返回格式：JSON
/// 接口费用：调用成功扣除0.012元，调用失败免费。
///
/// 用法（在 csharp 目录运行）：dotnet run -- image-colorize
/// 适用于黑白图片、老照片上色。
/// </summary>
public static class ImageColorize
{
    // ===== 请求参数（修改后直接运行即可）=====
    private const string ImagePath = "../images/test.jpeg"; // 待处理图片路径（相对 csharp 目录，建议使用黑白照片）

    public static async Task RunAsync()
    {
        if (string.IsNullOrEmpty(Config.ApiKey) || Config.ApiKey == "YOUR_APP_KEY")
            Console.Error.WriteLine("提示：请先在 csharp/Config.cs 中将 ApiKey 替换为你的 API Key。");

        var result = await ZjzApi.PostAsync("image/colorize", new Dictionary<string, string>
        {
            ["key"] = Config.ApiKey,
            ["image"] = await ZjzApi.ImageToBase64Async(ImagePath),
        });
        if (!ZjzApi.IsSuccess(result))
        {
            Console.WriteLine($"业务失败：code={result.Code}，msg={result.Msg}");
            return;
        }

        var imageUrl = ZjzApi.GetString(result.Data, "image");
        if (string.IsNullOrEmpty(imageUrl))
        {
            Console.WriteLine("返回数据中没有结果图地址。");
            return;
        }

        var saveDirectory = Path.Combine("output", "image-colorize");
        try
        {
            var file = await ZjzApi.DownloadImageAsync(imageUrl, saveDirectory, "colorize.jpg");
            Console.WriteLine($"上色成功：{file}");
        }
        catch (Exception exception)
        {
            Console.WriteLine($"下载失败，结果图地址：{imageUrl}（{exception.Message}）");
        }
    }
}
