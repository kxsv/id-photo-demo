using System;
using System.Collections.Generic;
using System.IO;
using System.Threading.Tasks;

namespace IdPhotoDemo.Demos;

/// <summary>
/// V3 职业形象照检测+制作。
///
/// 接口地址：https://api.zjzapi.com/idcardv3/all（与证件照制作相同，通过职业照规格 + 专属底色区分）
/// 请求方式：POST（application/x-www-form-urlencoded）
/// 返回格式：JSON
/// 接口费用：检测并制作成功扣除0.012元，检测失败扣除0.006元；若使用"美颜"功能将多扣除0.012元。
///
/// 用法（在 csharp 目录运行）：dotnet run -- v3-pose
/// </summary>
public static class V3Pose
{
    // ===== 请求参数（修改后直接运行即可）=====
    private const string ImagePath = "../images/pose.jpg"; // 待处理图片路径（相对 csharp 目录）
    private const string ItemId = "747";                      // 职业形象照规格id，可通过 item-list 示例获取
    private const string Colors = "pose_blue";              // 底色：pose_black/pose_blue/pose_darkgray/pose_gray/pose_white，留空返回全部
    private const string Beauty = "0";                      // 美颜：0=关闭 1=开启（开启后另行计费）

    public static async Task RunAsync()
    {
        if (string.IsNullOrEmpty(Config.ApiKey) || Config.ApiKey == "YOUR_APP_KEY")
            Console.Error.WriteLine("提示：请先在 csharp/Config.cs 中将 ApiKey 替换为你的 API Key。");

        var form = new Dictionary<string, string>
        {
            ["key"] = Config.ApiKey,
            ["item_id"] = ItemId,
            ["image"] = await ZjzApi.ImageToBase64Async(ImagePath),
            ["beauty"] = Beauty,
        };
        if (!string.IsNullOrEmpty(Colors))
            form["colors"] = Colors;

        var result = await ZjzApi.PostAsync("idcardv3/all", form);
        if (!ZjzApi.IsSuccess(result))
        {
            Console.WriteLine($"业务失败：code={result.Code}，msg={result.Msg}");
            return;
        }

        Console.WriteLine($"制作成功，pic_id：{ZjzApi.GetString(result.Data, "pic_id")}");
        if (result.Data.TryGetProperty("list", out var listValue))
            await ZjzApi.SaveImagesAsync(Path.Combine("output", "v3-pose"), listValue);
    }
}
