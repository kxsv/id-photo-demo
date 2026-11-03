using System;
using System.Collections.Generic;
using System.Text.Json;
using System.Threading.Tasks;

namespace IdPhotoDemo.Demos;

/// <summary>
/// 证件照规格列表。
///
/// 接口地址：https://api.zjzapi.com/item/list
/// 请求方式：POST（application/x-www-form-urlencoded）
/// 返回格式：JSON
/// 接口费用：调用免费
///
/// 用法（在 csharp 目录运行）：dotnet run -- item-list
///
/// 输出各规格的 item_id、名称、像素尺寸与支持底色，
/// 把 item_id 和对应底色填入制作类示例（如 v3-all）即可。
/// </summary>
public static class ItemList
{
    public static async Task RunAsync()
    {
        if (string.IsNullOrEmpty(Config.ApiKey) || Config.ApiKey == "YOUR_APP_KEY")
        {
            Console.Error.WriteLine("提示：请先在 csharp/Config.cs 中将 ApiKey 替换为你的 API Key。");
        }

        var form = new Dictionary<string, string>
        {
            ["key"] = Config.ApiKey,
        };

        var result = await ZjzApi.PostAsync("item/list", form);

        // HTTP 成功不代表业务成功，需再检查返回 JSON 中的 code。
        if (!ZjzApi.IsSuccess(result))
        {
            Console.WriteLine($"业务失败：code={result.Code}，msg={result.Msg}");
            return;
        }

        if (!result.Data.TryGetProperty("list", out var listValue) || listValue.ValueKind != JsonValueKind.Array)
        {
            Console.WriteLine("返回数据中没有规格列表。");
            return;
        }

        // 打印规格清单，便于挑选 item_id 与底色。
        Console.WriteLine("item_id  名称            像素          支持底色");
        var count = 0;
        foreach (var item in listValue.EnumerateArray())
        {
            var itemId = ZjzApi.GetString(item, "item_id");
            var name = ZjzApi.GetString(item, "name");
            var size = $"{ZjzApi.GetString(item, "width_px")}x{ZjzApi.GetString(item, "height_px")}";
            var color = ZjzApi.GetString(item, "color");
            Console.WriteLine($"{itemId.PadRight(8)} {name.PadRight(14)} {size.PadRight(12)} {color}");
            count++;
        }
        Console.WriteLine($"共 {count} 个规格。");
    }
}
