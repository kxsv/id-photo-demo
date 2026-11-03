using System;
using System.Collections.Generic;
using System.Threading.Tasks;

namespace IdPhotoDemo.Demos;

/// <summary>
/// 订单审核驳回后再次提交。
///
/// 接口地址：https://api.zjzapi.com/receipt/submit（与首次提交相同接口，额外携带 order_no）
/// 请求方式：POST（application/x-www-form-urlencoded）
/// 返回格式：JSON
/// 接口费用：订单办理完成扣费，7元每单；若选择打印邮寄服务再增加7元。
///
/// 用法（在 csharp 目录运行）：dotnet run -- receipt-resubmit
/// 注意：驳回通知中的 reject_msg 会说明驳回原因，按要求补充/修改 param（如身份证照片 image_idcard）后，
/// 携带原订单号重新提交。
/// </summary>
public static class ReceiptResubmit
{
    // ===== 请求参数（修改后直接运行即可）=====
    private const string PicId = "";       // 制作接口返回的 pic_id（有效期50分钟），先运行 receipt-make 获取
    private const string OrderNo = "";     // 上次提交返回的订单号（receipt-submit 的返回结果）
    private const string NoticeUrl = "";   // 回调地址：你的服务器接收订单状态通知的 URL
    private const string Param = "";       // 附加参数（json字符串），按驳回要求补充，如广东驾驶证驳回后补身份证照片：{"text_mobile":"13800138000","text_city":"深圳","image_idcard":"iVBORw0KGgo....（身份证照片base64）"}

    public static async Task RunAsync()
    {
        if (string.IsNullOrEmpty(Config.ApiKey) || Config.ApiKey == "YOUR_APP_KEY")
            Console.Error.WriteLine("提示：请先在 csharp/Config.cs 中将 ApiKey 替换为你的 API Key。");
        if (string.IsNullOrEmpty(PicId))
        {
            Console.Error.WriteLine("请先在参数区填写 pic_id（运行 receipt-make 获取）。");
            return;
        }
        if (string.IsNullOrEmpty(OrderNo))
        {
            Console.Error.WriteLine("请先在参数区填写 order_no（上次提交的订单号）。");
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
            ["order_no"] = OrderNo,
            ["notice_url"] = NoticeUrl,
        };
        if (!string.IsNullOrEmpty(Param))
            form["param"] = Param;

        var result = await ZjzApi.PostAsync("receipt/submit", form);
        if (!ZjzApi.IsSuccess(result))
        {
            Console.WriteLine($"业务失败：code={result.Code}，msg={result.Msg}");
            return;
        }

        Console.WriteLine($"再次提交成功，订单号：{ZjzApi.GetString(result.Data, "order_no")}");
    }
}
