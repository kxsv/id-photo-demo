using System;
using System.Collections.Generic;
using System.IO;
using System.Threading.Tasks;

namespace IdPhotoDemo.Demos;

/// <summary>
/// V3 证件照检测+自定义制作。
///
/// 接口地址：https://api.zjzapi.com/idcardv3/diy
/// 请求方式：POST（application/x-www-form-urlencoded）
/// 返回格式：JSON
/// 接口费用：检测并制作成功扣除0.012元，检测失败扣除0.006元；若使用"美颜"功能将多扣除0.012元。
///
/// 用法（在 csharp 目录运行）：dotnet run -- v3-diy
/// </summary>
public static class V3Diy
{
    // ===== 请求参数（修改后直接运行即可）=====
    private const string ImagePath = "../images/test.jpeg"; // 待处理图片路径（相对 csharp 目录）
    private const string WidthPx = "295";                   // 像素尺寸-宽，不超2000px
    private const string HeightPx = "413";                  // 像素尺寸-高，不超2000px
    private const string Color = "blue";                    // 底色：颜色值或名称
    private static readonly string FileSizeMin = "0";                 // 文件最小(kb)，最小5kb，0=不限
    private static readonly string FileSizeMax = "0";                 // 文件最大(kb)，最大300kb，0=不限
    private const string Dpi = "300";                       // 分辨率，默认300
    private const string Beauty = "0";                      // 美颜：0=关闭 1=开启（开启后另行计费）
    private const string Whitening = "50";                  // 美颜-美白等级：0-100，默认50

    public static async Task RunAsync()
    {
        if (string.IsNullOrEmpty(Config.ApiKey) || Config.ApiKey == "YOUR_APP_KEY")
            Console.Error.WriteLine("提示：请先在 csharp/Config.cs 中将 ApiKey 替换为你的 API Key。");

        var form = new Dictionary<string, string>
        {
            ["key"] = Config.ApiKey,
            ["image"] = await ZjzApi.ImageToBase64Async(ImagePath),
            ["width_px"] = WidthPx,
            ["height_px"] = HeightPx,
            ["color"] = Color,
            ["dpi"] = Dpi,
            ["beauty"] = Beauty,
            ["whitening"] = Whitening,
        };
        if (FileSizeMin != "0")
            form["file_size_min"] = FileSizeMin;
        if (FileSizeMax != "0")
            form["file_size_max"] = FileSizeMax;

        var result = await ZjzApi.PostAsync("idcardv3/diy", form);
        if (!ZjzApi.IsSuccess(result))
        {
            Console.WriteLine($"业务失败：code={result.Code}，msg={result.Msg}");
            return;
        }

        Console.WriteLine($"制作成功，pic_id：{ZjzApi.GetString(result.Data, "pic_id")}");
        if (result.Data.TryGetProperty("list", out var listValue))
            await ZjzApi.SaveImagesAsync(Path.Combine("output", "v3-diy"), listValue);
    }
}
