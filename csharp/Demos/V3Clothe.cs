using System;
using System.Collections.Generic;
using System.IO;
using System.Threading.Tasks;

namespace IdPhotoDemo.Demos;

/// <summary>
/// V3 证件照换装。
///
/// 接口地址：https://api.zjzapi.com/idcardv3/clothe
/// 请求方式：POST（application/x-www-form-urlencoded）
/// 返回格式：JSON
/// 接口费用：换装成功扣除0.036元。
///
/// 用法（在 csharp 目录运行）：dotnet run -- v3-clothe
/// 注意：换装基于制作接口返回的 pic_id（有效期50分钟），请先把 v3-all 运行得到的 pic_id 填到参数区。
/// </summary>
public static class V3Clothe
{
    // ===== 请求参数（修改后直接运行即可）=====
    private const string PicId = "";   // 制作接口返回的 pic_id（有效期50分钟），先运行 v3-all 获取
    private const string Clothe = "";  // 服装key，见官网《服装列表》文档

    public static async Task RunAsync()
    {
        if (string.IsNullOrEmpty(Config.ApiKey) || Config.ApiKey == "YOUR_APP_KEY")
            Console.Error.WriteLine("提示：请先在 csharp/Config.cs 中将 ApiKey 替换为你的 API Key。");
        if (string.IsNullOrEmpty(PicId) || string.IsNullOrEmpty(Clothe))
        {
            Console.Error.WriteLine("请先在参数区填写 pic_id（运行 v3-all 获取）与服装key clothe。");
            return;
        }

        var result = await ZjzApi.PostAsync("idcardv3/clothe", new Dictionary<string, string>
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

        Console.WriteLine("换装成功。");
        if (result.Data.TryGetProperty("list", out var listValue))
            await ZjzApi.SaveImagesAsync(Path.Combine("output", "v3-clothe"), listValue);
    }
}
