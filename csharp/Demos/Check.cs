using System;
using System.Collections.Generic;
using System.IO;
using System.Text.Json;
using System.Threading.Tasks;

namespace IdPhotoDemo.Demos;

/// <summary>
/// 证件照合规性检测。
///
/// 接口地址：https://api.zjzapi.com/idphoto/check
/// 请求方式：POST（application/x-www-form-urlencoded）
/// 返回格式：JSON
/// 接口费用：检测成功扣除0.012元。
///
/// 用法（在 csharp 目录运行）：dotnet run -- check
/// 提供规格id（item_id）即可按该规格检测；也可令 ItemId 为 "0"，改用下方自定义规格参数。
/// </summary>
public static class Check
{
    // ===== 请求参数（修改后直接运行即可）=====
    private const string ImagePath = "../images/test.jpeg"; // 待检测图片路径（相对 csharp 目录）
    private static readonly string ItemId = "1";            // 规格id，提供后可省略下方规格参数；改为 "0" 表示使用自定义规格参数

    // 方式二（ItemId 为 "0" 时生效）的自定义检测规格，按需增删。
    private const string CustomWidthPx = "295";
    private const string CustomHeightPx = "413";
    private const string CustomDpi = "300";
    private const string CustomFileFormat = "jpg";
    private const string CustomBgColors = "white,blue";
    private const string CustomSolidBg = "1";
    private const string FaceWidthRatioMin = "0.4";
    private const string FaceWidthRatioMax = "0.6";
    private const string ShoulderHorizontal = "5";
    private const string FacePoseYaw = "15";
    private const string FaceClarity = "0.5";
    private const string FaceBrightness = "0.3";
    private const string FaceMask = "1";
    private const string Glasses = "1";

    public static async Task RunAsync()
    {
        if (string.IsNullOrEmpty(Config.ApiKey) || Config.ApiKey == "YOUR_APP_KEY")
            Console.Error.WriteLine("提示：请先在 csharp/Config.cs 中将 ApiKey 替换为你的 API Key。");

        var form = new Dictionary<string, string>
        {
            ["key"] = Config.ApiKey,
            ["image"] = await ZjzApi.ImageToBase64Async(ImagePath),
        };
        if (ItemId != "0")
        {
            // 方式一：按规格id检测（推荐，规格要求自动带上）。
            form["item_id"] = ItemId;
        }
        else
        {
            // 方式二：自定义检测规格（不传 item_id 时生效），按需增删。
            form["width_px"] = CustomWidthPx;
            form["height_px"] = CustomHeightPx;
            form["dpi"] = CustomDpi;
            form["file_format"] = CustomFileFormat;
            form["bg_colors"] = CustomBgColors;
            form["solid_bg"] = CustomSolidBg;
            form["face_width_ratio_min"] = FaceWidthRatioMin;
            form["face_width_ratio_max"] = FaceWidthRatioMax;
            form["shoulder_horizontal"] = ShoulderHorizontal;
            form["face_pose_yaw"] = FacePoseYaw;
            form["face_clarity"] = FaceClarity;
            form["face_brightness"] = FaceBrightness;
            form["face_mask"] = FaceMask;
            form["glasses"] = Glasses;
        }

        var result = await ZjzApi.PostAsync("idphoto/check", form);
        if (!ZjzApi.IsSuccess(result))
        {
            Console.WriteLine($"业务失败：code={result.Code}，msg={result.Msg}");
            return;
        }

        // passed_status 可能为数字也可能为字符串（接口部分字段两种写法都有），统一按字符串读取后解析。
        var passedStatusRaw = ZjzApi.GetString(result.Data, "passed_status");
        var passedStatus = passedStatusRaw.Length > 0 && int.TryParse(passedStatusRaw, out var passedStatusVal) && passedStatusVal == 1;
        Console.WriteLine(passedStatus ? "总检测结果：通过" : "总检测结果：不通过");

        if (!result.Data.TryGetProperty("passed_list", out var list) || list.ValueKind != JsonValueKind.Array)
        {
            Console.WriteLine("返回数据中没有检测项列表。");
            return;
        }

        Console.WriteLine($"{"检测项",-16}{"结果",-6}{"提示",-10}检测值");
        foreach (var item in list.EnumerateArray())
        {
            // passed 同样可能为数字或字符串，统一按字符串读取后解析。
            var passedRaw = ZjzApi.GetString(item, "passed");
            var passed = passedRaw.Length > 0 && int.TryParse(passedRaw, out var passedVal) && passedVal == 1;
            var key = ZjzApi.GetString(item, "key");
            var msg = ZjzApi.GetString(item, "msg");
            var actual = "";
            var unit = "";
            if (item.TryGetProperty("detail", out var detail) && detail.ValueKind == JsonValueKind.Object)
            {
                actual = ZjzApi.GetString(detail, "actual");
                unit = ZjzApi.GetString(detail, "unit");
            }
            var actualText = actual.Length > 0 ? actual + unit : "";
            Console.WriteLine($"{key,-16}{(passed ? "通过" : "不通过"),-6}{msg,-10}{actualText}");
        }
    }
}
