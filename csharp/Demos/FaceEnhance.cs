using System;
using System.Collections.Generic;
using System.IO;
using System.Threading.Tasks;

namespace IdPhotoDemo.Demos;

/// <summary>
/// 人像修复增强。
///
/// 接口地址：https://api.zjzapi.com/face/enhance
/// 请求方式：POST（application/x-www-form-urlencoded）
/// 返回格式：JSON
/// 接口费用：返回成功扣除0.06元。
///
/// 用法（在 csharp 目录运行）：dotnet run -- face-enhance
/// 适用于模糊图片清晰度修复：拍照没对焦好、图片模糊、光线差、老照片不清晰、图片像素低等。
/// </summary>
public static class FaceEnhance
{
    // ===== 请求参数（修改后直接运行即可）=====
    private const string ImagePath = "../images/test.jpeg"; // 待处理图片路径（相对 csharp 目录）
    private const string Size = "original";                 // 图片大小：original=原始大小、x2=高清图，默认 original

    public static async Task RunAsync()
    {
        if (string.IsNullOrEmpty(Config.ApiKey) || Config.ApiKey == "YOUR_APP_KEY")
            Console.Error.WriteLine("提示：请先在 csharp/Config.cs 中将 ApiKey 替换为你的 API Key。");

        var result = await ZjzApi.PostAsync("face/enhance", new Dictionary<string, string>
        {
            ["key"] = Config.ApiKey,
            ["image"] = await ZjzApi.ImageToBase64Async(ImagePath),
            ["size"] = Size,
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

        var saveDirectory = Path.Combine("output", "face-enhance");
        try
        {
            var file = await ZjzApi.DownloadImageAsync(imageUrl, saveDirectory, "enhance.jpg");
            Console.WriteLine($"修复增强成功：{file}");
        }
        catch (Exception exception)
        {
            Console.WriteLine($"下载失败，结果图地址：{imageUrl}（{exception.Message}）");
        }
    }
}
