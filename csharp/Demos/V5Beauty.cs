using System;
using System.Collections.Generic;
using System.IO;
using System.Threading.Tasks;

namespace IdPhotoDemo.Demos;

/// <summary>
/// V5 证件照美颜（有水印）。
///
/// 接口地址：https://api.zjzapi.com/idcardv5/beauty
/// 请求方式：POST（application/x-www-form-urlencoded）
/// 返回格式：JSON
/// 接口费用：对制作出来的证件照进行美颜；不支持换装后再美颜。
///
/// 用法（在 csharp 目录运行）：dotnet run -- v5-beauty <pic_id>
/// 注意：美颜基于制作接口返回的 pic_id（有效期50分钟），可命令行传入，或填到下方 PicId 常量。
/// </summary>
public static class V5Beauty
{
    // ===== 请求参数（可命令行传入，或在此处填写）=====
    // 用法：dotnet run -- v5-beauty <pic_id>
    // 不传时回退到下方 PicId 常量
    private const string PicId = "";   // 制作接口返回的 pic_id（有效期50分钟），先运行 v5-make 获取

    public static async Task RunAsync()
    {
        var picId = Program.ExtraArgs.Length > 0 ? Program.ExtraArgs[0] : PicId;

        if (string.IsNullOrEmpty(Config.ApiKey) || Config.ApiKey == "YOUR_APP_KEY")
            Console.Error.WriteLine("提示：请先在 csharp/Config.cs 中将 ApiKey 替换为你的 API Key。");
        if (string.IsNullOrEmpty(picId))
        {
            Console.Error.WriteLine("请传入 pic_id：dotnet run -- v5-beauty <pic_id>（运行 v5-make 获取），或填写 V5Beauty.cs 的 PicId。");
            return;
        }

        var result = await ZjzApi.PostAsync("idcardv5/beauty", new Dictionary<string, string>
        {
            ["key"] = Config.ApiKey,
            ["pic_id"] = picId,
        });
        if (!ZjzApi.IsSuccess(result))
        {
            Console.WriteLine($"业务失败：code={result.Code}，msg={result.Msg}");
            return;
        }

        Console.WriteLine("美颜成功。");
        if (result.Data.TryGetProperty("list", out var listValue))
            await ZjzApi.SaveImagesAsync(Path.Combine("output", "v5-beauty"), listValue);
    }
}
