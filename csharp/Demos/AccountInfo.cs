using System;
using System.Collections.Generic;
using System.Threading.Tasks;

namespace IdPhotoDemo.Demos;

/// <summary>
/// 获取账号信息。
///
/// 接口地址：https://api.zjzapi.com/user/info
/// 请求方式：POST（application/x-www-form-urlencoded）
/// 返回格式：JSON
/// 接口费用：调用免费，数据有10分钟延迟。
///
/// 用法（在 csharp 目录运行）：dotnet run -- account-info
/// 注意：本接口不使用 API Key，而是使用用户中心的 access_token（官网「用户中心 - 账户信息」获取），
/// 请填到 Config.cs 的 AccessToken。
/// </summary>
public static class AccountInfo
{
    public static async Task RunAsync()
    {
        if (string.IsNullOrEmpty(Config.AccessToken) || Config.AccessToken == "YOUR_ACCESS_TOKEN")
        {
            Console.Error.WriteLine("提示：请先在 csharp/Config.cs 中将 AccessToken 替换为你的 access_token（官网「用户中心 - 账户信息」获取）。");
            return;
        }

        var result = await ZjzApi.PostAsync("user/info", new Dictionary<string, string>
        {
            ["access_token"] = Config.AccessToken,
        });
        if (!ZjzApi.IsSuccess(result))
        {
            Console.WriteLine($"业务失败：code={result.Code}，msg={result.Msg}");
            return;
        }

        var data = result.Data;
        Console.WriteLine($"登录账号：{ZjzApi.GetString(data, "mobile")}");
        Console.WriteLine($"账户余额：{ZjzApi.GetString(data, "money")} 元");
        Console.WriteLine($"应用数量：{ZjzApi.GetString(data, "app_num")}");
        Console.WriteLine($"今日调用：{ZjzApi.GetString(data, "day_make")} 次");
    }
}
