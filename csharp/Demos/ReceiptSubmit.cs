using System;
using System.Collections.Generic;
using System.Threading.Tasks;

namespace IdPhotoDemo.Demos;

/// <summary>
/// 证件照回执提交。
///
/// 接口地址：https://api.zjzapi.com/receipt/submit
/// 请求方式：POST（application/x-www-form-urlencoded）
/// 返回格式：JSON
/// 接口费用：订单办理完成扣费，7元每单；若选择打印邮寄服务再增加7元。
///
/// 用法（在 csharp 目录运行）：dotnet run -- receipt-submit
/// 注意：pic_id 来自回执制作接口的返回结果（有效期50分钟）；
/// 提交后订单状态变化会通过 notice_url 回调通知。
/// </summary>
public static class ReceiptSubmit
{
    // ===== 请求参数（修改后直接运行即可）=====
    private const string PicId = "";       // 制作接口返回的 pic_id（有效期50分钟），先运行 receipt-make 获取
    private const string NoticeUrl = "";   // 回调地址：你的服务器接收订单状态通知的 URL
    private const string Param = "";       // 附加参数（json字符串），如广东驾驶证：{"text_mobile":"13800138000","text_city":"深圳","text_name":"姓名","text_idcard":"身份证号码"}
    private const string IsPrint = "0";    // 是否打印邮寄：1=是 0=否
    private const string PrintParam = "";  // 打印邮寄（is_print=1 时必填），如：{"name":"张三","phone":"13800138000","province":"广东","city":"深圳","district":"南山区","address":"xx路xx小区"}

    public static async Task RunAsync()
    {
        if (string.IsNullOrEmpty(Config.ApiKey) || Config.ApiKey == "YOUR_APP_KEY")
            Console.Error.WriteLine("提示：请先在 csharp/Config.cs 中将 ApiKey 替换为你的 API Key。");
        if (string.IsNullOrEmpty(PicId))
        {
            Console.Error.WriteLine("请先在参数区填写 pic_id（运行 receipt-make 获取）。");
            return;
        }
        if (string.IsNullOrEmpty(NoticeUrl))
        {
            Console.Error.WriteLine("请先在参数区填写 notice_url（你的订单状态回调地址）。");
            return;
        }

        var form = new Dictionary<string, string>
        {
            ["key"] = Config.ApiKey,
            ["pic_id"] = PicId,
            ["notice_url"] = NoticeUrl,
        };
        if (!string.IsNullOrEmpty(Param))
            form["param"] = Param;
        form["is_print"] = IsPrint;
        if (!string.IsNullOrEmpty(PrintParam))
            form["print_param"] = PrintParam;

        var result = await ZjzApi.PostAsync("receipt/submit", form);
        if (!ZjzApi.IsSuccess(result))
        {
            Console.WriteLine($"业务失败：code={result.Code}，msg={result.Msg}");
            return;
        }

        Console.WriteLine($"提交成功，订单号：{ZjzApi.GetString(result.Data, "order_no")}");
        Console.WriteLine("提示：订单状态变化将通过 notice_url 回调通知；如需取消订单见 receipt-cancel。");
    }
}
