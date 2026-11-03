using System;
using System.Collections.Generic;
using System.IO;
using System.Threading.Tasks;

namespace IdPhotoDemo.Demos;

/// <summary>
/// V5 职业形象照检测+制作（有水印）。
///
/// 接口地址：https://api.zjzapi.com/idcardv5/make（与证件照制作相同，通过职业照规格 + 专属底色区分）
/// 请求方式：POST（application/x-www-form-urlencoded）
/// 返回格式：JSON
/// 接口费用：调用免费，结果图带水印；可使用返回的 pic_id 调用 v5-get 获取无水印证件照。
///
/// 用法（在 csharp 目录运行）：dotnet run -- v5-pose
/// </summary>
public static class V5Pose
{
    // ===== 请求参数（修改后直接运行即可）=====
    private const string ImagePath = "../images/pose.jpg"; // 待处理图片路径（相对 csharp 目录）
    private const string ItemId = "747";                      // 职业形象照规格id，可通过 item-list 示例获取
    private const string Colors = "pose_blue";              // 底色：pose_black/pose_blue/pose_darkgray/pose_gray/pose_white，留空返回全部
    private const string Enhance = "1";                     // 人像增强：默认已支持，不需要可传0
    private const string Beauty = "1";                      // 美颜：默认已支持，不需要可传0
    private const string Whitening = "50";                  // 美颜-美白等级：0-100，默认50

    public static async Task RunAsync()
    {
        if (string.IsNullOrEmpty(Config.ApiKey) || Config.ApiKey == "YOUR_APP_KEY")
            Console.Error.WriteLine("提示：请先在 csharp/Config.cs 中将 ApiKey 替换为你的 API Key。");

        var form = new Dictionary<string, string>
        {
            ["key"] = Config.ApiKey,
            ["item_id"] = ItemId,
            ["image"] = await ZjzApi.ImageToBase64Async(ImagePath),
            ["enhance"] = Enhance,
            ["beauty"] = Beauty,
            ["whitening"] = Whitening,
        };
        if (!string.IsNullOrEmpty(Colors))
            form["colors"] = Colors;

        var result = await ZjzApi.PostAsync("idcardv5/make", form);
        if (!ZjzApi.IsSuccess(result))
        {
            Console.WriteLine($"业务失败：code={result.Code}，msg={result.Msg}");
            return;
        }

        Console.WriteLine($"制作成功（有水印），pic_id：{ZjzApi.GetString(result.Data, "pic_id")}");
        Console.WriteLine("提示：把 pic_id 填入 v5-get 可获取无水印证件照。");
        if (result.Data.TryGetProperty("list", out var listValue))
            await ZjzApi.SaveImagesAsync(Path.Combine("output", "v5-pose"), listValue);
    }
}
