using System;
using System.Collections.Generic;
using System.IO;
using System.Threading.Tasks;

namespace IdPhotoDemo.Demos;

/// <summary>
/// V5 证件照换装（有水印）。
///
/// 接口地址：https://api.zjzapi.com/idcardv5/clothe
/// 请求方式：POST（application/x-www-form-urlencoded）
/// 返回格式：JSON
/// 接口费用：调用免费，结果图带水印；可使用返回的 pic_id 调用 v5-get 获取无水印证件照。
///
/// 用法（在 csharp 目录运行）：dotnet run -- v5-clothe
/// 注意：换装基于制作接口返回的 pic_id（有效期50分钟），请先把 v5-make 运行得到的 pic_id 填到参数区。
/// </summary>
public static class V5Clothe
{
    // ===== 请求参数（修改后直接运行即可）=====
    private const string PicId = "";   // 制作接口返回的 pic_id（有效期50分钟），先运行 v5-make 获取
    private const string Clothe = "";  // 服装key，见官网《服装列表》文档

    public static async Task RunAsync()
    {
        if (string.IsNullOrEmpty(Config.ApiKey) || Config.ApiKey == "YOUR_APP_KEY")
            Console.Error.WriteLine("提示：请先在 csharp/Config.cs 中将 ApiKey 替换为你的 API Key。");
        if (string.IsNullOrEmpty(PicId) || string.IsNullOrEmpty(Clothe))
        {
            Console.Error.WriteLine("请先在参数区填写 pic_id（运行 v5-make 获取）与服装key clothe。");
            return;
        }

        var result = await ZjzApi.PostAsync("idcardv5/clothe", new Dictionary<string, string>
        {
            ["key"] = Config.ApiKey,
            ["pic_id"] = PicId,
            ["clothe"] = Clothe,
        });
        if (!ZjzApi.IsSuccess(result))
        {
            Console.WriteLine($"业务失败：code={result.Code}，msg={result.Msg}");
            return;
        }

        Console.WriteLine("换装成功（有水印）。");
        if (result.Data.TryGetProperty("list", out var listValue))
            await ZjzApi.SaveImagesAsync(Path.Combine("output", "v5-clothe"), listValue);
    }
}
