using System;
using System.Threading.Tasks;

namespace IdPhotoDemo;

/// <summary>
/// 示例入口：按 demo 名称分发到对应接口示例。
/// 用法（在 csharp 目录运行）：
///   dotnet run -- &lt;demo名称&gt;
/// 示例：
///   dotnet run -- v3-all
///   dotnet run -- item-list
/// </summary>
public static class Program
{
    /// <summary>
    /// 命令行中 demo 名称之后的剩余参数，供需要参数的 demo 读取。
    /// 例如 `dotnet run -- v5-beauty ABC123` 时，ExtraArgs = ["ABC123"]。
    /// </summary>
    public static string[] ExtraArgs { get; private set; } = Array.Empty<string>();

    public static async Task<int> Main(string[] args)
    {
        try
        {
            if (args.Length == 0)
            {
                PrintUsage();
                return 1;
            }

            ExtraArgs = args.Length > 1 ? args[1..] : Array.Empty<string>();

            switch (args[0])
            {
                case "v3-all":
                    await Demos.V3All.RunAsync();
                    return 0;
                case "v3-bg":
                    await Demos.V3Bg.RunAsync();
                    return 0;
                case "v3-diy":
                    await Demos.V3Diy.RunAsync();
                    return 0;
                case "v3-pose":
                    await Demos.V3Pose.RunAsync();
                    return 0;
                case "v3-clothe":
                    await Demos.V3Clothe.RunAsync();
                    return 0;
                case "v3-beauty":
                    await Demos.V3Beauty.RunAsync();
                    return 0;
                case "v5-make":
                    await Demos.V5Make.RunAsync();
                    return 0;
                case "v5-bg":
                    await Demos.V5Bg.RunAsync();
                    return 0;
                case "v5-diy":
                    await Demos.V5Diy.RunAsync();
                    return 0;
                case "v5-pose":
                    await Demos.V5Pose.RunAsync();
                    return 0;
                case "v5-clothe":
                    await Demos.V5Clothe.RunAsync();
                    return 0;
                case "v5-beauty":
                    await Demos.V5Beauty.RunAsync();
                    return 0;
                case "v5-get":
                    await Demos.V5Get.RunAsync();
                    return 0;
                case "receipt-make":
                    await Demos.ReceiptMake.RunAsync();
                    return 0;
                case "receipt-submit":
                    await Demos.ReceiptSubmit.RunAsync();
                    return 0;
                case "receipt-resubmit":
                    await Demos.ReceiptResubmit.RunAsync();
                    return 0;
                case "receipt-cancel":
                    await Demos.ReceiptCancel.RunAsync();
                    return 0;
                case "receipt-notice":
                    await Demos.ReceiptNotice.RunAsync();
                    return 0;
                case "check":
                    await Demos.Check.RunAsync();
                    return 0;
                case "face-enhance":
                    await Demos.FaceEnhance.RunAsync();
                    return 0;
                case "image-colorize":
                    await Demos.ImageColorize.RunAsync();
                    return 0;
                case "image-edit":
                    await Demos.ImageEdit.RunAsync();
                    return 0;
                case "item-list":
                    await Demos.ItemList.RunAsync();
                    return 0;
                case "item-detail":
                    await Demos.ItemDetail.RunAsync();
                    return 0;
                case "account-info":
                    await Demos.AccountInfo.RunAsync();
                    return 0;
                case "account-app":
                    await Demos.AccountApp.RunAsync();
                    return 0;
                case "print":
                    await Demos.Print.RunAsync();
                    return 0;
                default:
                    Console.WriteLine($"未知 demo：{args[0]}");
                    PrintUsage();
                    return 1;
            }
        }
        catch (Exception exception)
        {
            Console.Error.WriteLine("执行失败：" + exception.Message);
            return 1;
        }
    }

    private static void PrintUsage()
    {
        Console.WriteLine(
            """
            用法：dotnet run -- <demo名称>
            当前可用的 demo：
              v3-all          V3 证件照检测+制作
              v3-bg           V3 更换底色
              v3-diy          V3 自定义制作
              v3-pose         V3 职业形象照制作
              v3-clothe       V3 换装（需 pic_id）
              v3-beauty       V3 美颜（需 pic_id）
              v5-make         V5 检测+制作（有水印）
              v5-bg           V5 更换底色
              v5-diy          V5 自定义制作
              v5-pose         V5 职业形象照制作
              v5-clothe       V5 换装（需 pic_id）
              v5-beauty       V5 美颜（需 pic_id）
              v5-get          V5 获取无水印证件照（需 pic_id）
              receipt-make    回执检测+制作
              receipt-submit  回执提交（需 pic_id + notice_url）
              receipt-resubmit 回执驳回后再次提交（需 order_no）
              receipt-cancel  回执订单取消（需 order_no）
              receipt-notice  回执状态通知回调（需独立运行的服务）
              check           证件照合规性检测
              face-enhance    人像修复增强
              image-colorize  黑白图片上色
              image-edit      图片编辑（改 dpi / 文件大小）
              item-list       证件照规格列表
              item-detail     证件照规格详情
              account-info    获取账号信息（需 access_token）
              account-app     获取应用信息（需 access_token）
              print           冲印排版（相纸联排，需 SixLabors.ImageSharp）
            """);
    }
}
