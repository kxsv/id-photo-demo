using System;
using System.Collections.Generic;
using System.IO;
using System.Threading.Tasks;

namespace IdPhotoDemo.Demos;

/// <summary>
/// V5 获取证件照（无水印）。
///
/// 接口地址：https://api.zjzapi.com/idcardv5/get
/// 请求方式：POST（application/x-www-form-urlencoded）
/// 返回格式：JSON
/// 接口费用：获取成功扣除0.6元，重复获取不扣费。
///
/// 用法（在 csharp 目录运行）：dotnet run -- v5-get
/// 注意：pic_id 来自制作类接口（v5-make、v5-bg、v5-diy、v5-pose 等）的返回结果，有效期50分钟；
/// 若换装过，请同时填写换装时使用的 clothe。
/// </summary>
public static class V5Get
{
    // ===== 请求参数（修改后直接运行即可）=====
    private const string PicId = "";         // 制作接口返回的 pic_id（有效期50分钟），先运行 v5-make 获取
    private const string Clothe = "";        // 换装时使用的服装key，未换装留空
    private const string GiftHd = "0";       // 1=赠送一张三寸高清照（仅规格 width_px<500 支持；更换底色、自定义制作暂不支持）
    private const string GiftHdColor = "";   // 高清照底色，留空默认返回制作结果第一个底色

    public static async Task RunAsync()
    {
        if (string.IsNullOrEmpty(Config.ApiKey) || Config.ApiKey == "YOUR_APP_KEY")
            Console.Error.WriteLine("提示：请先在 csharp/Config.cs 中将 ApiKey 替换为你的 API Key。");
        if (string.IsNullOrEmpty(PicId))
        {
            Console.Error.WriteLine("请先在参数区填写 pic_id（运行 v5-make 获取）。");
            return;
        }

        var form = new Dictionary<string, string>
        {
            ["key"] = Config.ApiKey,
            ["pic_id"] = PicId,
            ["gift_hd"] = GiftHd,
        };
        if (!string.IsNullOrEmpty(Clothe))
            form["clothe"] = Clothe;
        if (!string.IsNullOrEmpty(GiftHdColor))
            form["gift_hd_color"] = GiftHdColor;

        var result = await ZjzApi.PostAsync("idcardv5/get", form);
        if (!ZjzApi.IsSuccess(result))
        {
            Console.WriteLine($"业务失败：code={result.Code}，msg={result.Msg}");
            return;
        }

        Console.WriteLine("获取无水印证件照成功。");
        if (result.Data.TryGetProperty("list", out var listValue))
            await ZjzApi.SaveImagesAsync(Path.Combine("output", "v5-get"), listValue);

        // 赠送的三寸高清照（开启 gift_hd 时才有）。
        var giftHd = ZjzApi.GetString(result.Data, "gift_hd");
        if (!string.IsNullOrEmpty(giftHd))
        {
            try
            {
                var file = await ZjzApi.DownloadImageAsync(giftHd, Path.Combine("output", "v5-get"), "gift_hd.jpg");
                Console.WriteLine($"高清照：{file}");
            }
            catch (Exception exception)
            {
                Console.WriteLine($"高清照下载失败：{exception.Message}（结果图地址：{giftHd}）");
            }
        }
    }
}
