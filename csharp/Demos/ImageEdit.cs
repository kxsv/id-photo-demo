using System;
using System.Collections.Generic;
using System.IO;
using System.Threading.Tasks;

namespace IdPhotoDemo.Demos;

/// <summary>
/// 图片编辑。
///
/// 接口地址：https://api.zjzapi.com/image/edit
/// 请求方式：POST（application/x-www-form-urlencoded）
/// 返回格式：JSON
/// 接口费用：调用免费（账户余额最低需有10元才可调用）。
///
/// 用法（在 csharp 目录运行）：dotnet run -- image-edit
/// 修改图片像素（dpi）与文件大小（kb）。
/// </summary>
public static class ImageEdit
{
    // ===== 请求参数（修改后直接运行即可）=====
    private const string ImagePath = "../images/test.jpeg"; // 待处理图片路径（相对 csharp 目录）
    private const string Dpi = "300";           // 图片dpi，支持72,96,150,200,300,350,500,1000；填 "0" 表示不修改
    private static readonly string FileSizeMin = "0";     // 文件最小(kb)，最小5kb；填 "0" 表示不限制
    private static readonly string FileSizeMax = "0";     // 文件最大(kb)，最大300kb；填 "0" 表示不限制

    public static async Task RunAsync()
    {
        if (string.IsNullOrEmpty(Config.ApiKey) || Config.ApiKey == "YOUR_APP_KEY")
            Console.Error.WriteLine("提示：请先在 csharp/Config.cs 中将 ApiKey 替换为你的 API Key。");

        var form = new Dictionary<string, string>
        {
            ["key"] = Config.ApiKey,
            ["image"] = await ZjzApi.ImageToBase64Async(ImagePath),
        };
        if (Dpi != "0")
            form["dpi"] = Dpi;
        if (FileSizeMin != "0")
            form["file_size_min"] = FileSizeMin;
        if (FileSizeMax != "0")
            form["file_size_max"] = FileSizeMax;

        var result = await ZjzApi.PostAsync("image/edit", form);
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

        var saveDirectory = Path.Combine("output", "image-edit");
        try
        {
            var file = await ZjzApi.DownloadImageAsync(imageUrl, saveDirectory, "edit.jpg");
            Console.WriteLine($"图片编辑成功：{file}");
        }
        catch (Exception exception)
        {
            Console.WriteLine($"下载失败，结果图地址：{imageUrl}（{exception.Message}）");
        }
    }
}
