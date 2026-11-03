using System;
using System.Collections.Generic;
using System.Threading.Tasks;

namespace IdPhotoDemo.Demos;

/// <summary>
/// 证件照规格详情。
///
/// 接口地址：https://api.zjzapi.com/item/get
/// 请求方式：POST（application/x-www-form-urlencoded）
/// 返回格式：JSON
/// 接口费用：调用免费
///
/// 用法（在 csharp 目录运行）：dotnet run -- item-detail
/// 输出单个规格的像素尺寸、毫米尺寸、dpi、文件大小要求与支持底色。
/// </summary>
public static class ItemDetail
{
    // ===== 请求参数（修改后直接运行即可）=====
    private const string ItemId = "1";   // 规格id，可通过 item-list 示例获取

    public static async Task RunAsync()
    {
        if (string.IsNullOrEmpty(Config.ApiKey) || Config.ApiKey == "YOUR_APP_KEY")
            Console.Error.WriteLine("提示：请先在 csharp/Config.cs 中将 ApiKey 替换为你的 API Key。");

        var result = await ZjzApi.PostAsync("item/get", new Dictionary<string, string>
        {
            ["key"] = Config.ApiKey,
            ["item_id"] = ItemId,
        });
        if (!ZjzApi.IsSuccess(result))
        {
            Console.WriteLine($"业务失败：code={result.Code}，msg={result.Msg}");
            return;
        }

        var data = result.Data;
        Console.WriteLine($"规格id：{ZjzApi.GetString(data, "id")}");
        Console.WriteLine($"名称：{ZjzApi.GetString(data, "name")}");
        Console.WriteLine($"像素-宽：{ZjzApi.GetString(data, "width_px")}");
        Console.WriteLine($"像素-高：{ZjzApi.GetString(data, "height_px")}");
        Console.WriteLine($"毫米-宽：{ZjzApi.GetString(data, "width_mm")}");
        Console.WriteLine($"毫米-高：{ZjzApi.GetString(data, "height_mm")}");
        Console.WriteLine($"分辨率：{ZjzApi.GetString(data, "dpi")}");
        Console.WriteLine($"文件最小(kb)：{ZjzApi.GetString(data, "file_size_min")}");
        Console.WriteLine($"文件最大(kb)：{ZjzApi.GetString(data, "file_size_max")}");
        Console.WriteLine($"文件大小说明：{ZjzApi.GetString(data, "file_size_msg")}");
        Console.WriteLine($"支持底色：{ZjzApi.GetString(data, "color")}");
        Console.WriteLine($"分类：{ZjzApi.GetString(data, "category")}");
        Console.WriteLine($"分类id：{ZjzApi.GetString(data, "category_id")}");
    }
}
