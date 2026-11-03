using System;
using System.Collections.Generic;
using System.IO;
using System.Text.Json;
using System.Threading.Tasks;

namespace IdPhotoDemo.Demos;

/// <summary>
/// V3 证件照检测+制作。
///
/// 接口地址：https://api.zjzapi.com/idcardv3/all
/// 请求方式：POST（application/x-www-form-urlencoded）
/// 返回格式：JSON
/// 接口费用：制作成功按次扣费；检测不通过扣除检测费用；beauty=1 时美颜另行计费。
///
/// 用法（在 csharp 目录运行）：dotnet run -- v3-all
/// </summary>
public static class V3All
{
    // ===== 请求参数（修改后直接运行即可）=====
    private const string ImagePath = "../images/test.jpeg"; // 待处理图片路径（相对 csharp 目录）
    private const string ItemId = "1";                    // 规格id，可通过「证件照规格列表」示例获取
    private const string Colors = "white";                  // 底色，如 white/blue/red；留空返回该规格支持的全部底色
    private const string Beauty = "0";                      // 美颜：0=关闭 1=开启（开启后另行计费）

    public static async Task RunAsync()
    {
        if (string.IsNullOrEmpty(Config.ApiKey) || Config.ApiKey == "YOUR_APP_KEY")
        {
            Console.Error.WriteLine("提示：请先在 csharp/Config.cs 中将 ApiKey 替换为你的 API Key。");
        }

        var form = new Dictionary<string, string>
        {
            ["key"] = Config.ApiKey,
            ["item_id"] = ItemId,
            ["image"] = await ZjzApi.ImageToBase64Async(ImagePath),
            ["beauty"] = Beauty,
        };
        // colors 为空时不传，由服务端返回全部底色。
        if (!string.IsNullOrEmpty(Colors))
        {
            form["colors"] = Colors;
        }

        var result = await ZjzApi.PostAsync("idcardv3/all", form);

        // HTTP 成功不代表业务成功，需再检查返回 JSON 中的 code。
        if (!ZjzApi.IsSuccess(result))
        {
            Console.WriteLine($"业务失败：code={result.Code}，msg={result.Msg}");
            return;
        }

        Console.WriteLine($"制作成功，pic_id：{ZjzApi.GetString(result.Data, "pic_id")}");

        // 把各底色结果图下载到 output/v3-all/ 目录。
        if (result.Data.TryGetProperty("list", out var listValue))
        {
            await ZjzApi.SaveImagesAsync(Path.Combine("output", "v3-all"), listValue);
        }
    }
}
