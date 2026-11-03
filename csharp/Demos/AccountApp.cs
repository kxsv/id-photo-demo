using System;
using System.Collections.Generic;
using System.Threading.Tasks;

namespace IdPhotoDemo.Demos;

/// <summary>
/// 获取应用信息。
///
/// 接口地址：https://api.zjzapi.com/user/app
/// 请求方式：POST（application/x-www-form-urlencoded）
/// 返回格式：JSON
/// 接口费用：调用免费，数据有10分钟延迟。
///
/// 用法（在 csharp 目录运行）：dotnet run -- account-app
/// 注意：本接口需要同时使用用户中心的 access_token（官网「用户中心 - 账户信息」获取，填到 Config.cs 的 AccessToken）
/// 与你的 API Key。
/// </summary>
public static class AccountApp
{
    // 应用权限：1=证件照v3，2=证件照v5，3=支持全部，4=证件照回执
    private static string TypeText(string raw) => raw switch
    {
        "1" => "证件照v3",
        "2" => "证件照v5",
        "3" => "支持全部",
        "4" => "证件照回执",
        _ => raw,
    };

    public static async Task RunAsync()
    {
        if (string.IsNullOrEmpty(Config.AccessToken) || Config.AccessToken == "YOUR_ACCESS_TOKEN")
        {
            Console.Error.WriteLine("提示：请先在 csharp/Config.cs 中将 AccessToken 替换为你的 access_token（官网「用户中心 - 账户信息」获取）。");
            return;
        }

        var result = await ZjzApi.PostAsync("user/app", new Dictionary<string, string>
        {
            ["access_token"] = Config.AccessToken,
            ["key"] = Config.ApiKey,
        });
        if (!ZjzApi.IsSuccess(result))
        {
            Console.WriteLine($"业务失败：code={result.Code}，msg={result.Msg}");
            return;
        }

        var data = result.Data;
        var statusRaw = ZjzApi.GetString(data, "status");
        var statusText = statusRaw == "1" ? "启用中" : "已停用";
        Console.WriteLine($"应用名称：{ZjzApi.GetString(data, "name")}");
        Console.WriteLine($"应用权限：{TypeText(ZjzApi.GetString(data, "type"))}");
        Console.WriteLine($"应用状态：{statusText}");
        Console.WriteLine($"调用次数：{ZjzApi.GetString(data, "make_num")}");
    }
}
