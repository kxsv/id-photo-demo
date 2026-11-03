using System;
using System.Collections.Generic;
using System.IO;
using System.Threading.Tasks;
using SixLabors.ImageSharp;
using SixLabors.ImageSharp.PixelFormats;
using SixLabors.ImageSharp.Processing;
using SixLabors.ImageSharp.Formats.Jpeg;

namespace IdPhotoDemo.Demos;

/// <summary>
/// 证件照冲印排版（相纸联排）。
///
/// 将一张已制作好的证件照，按预设的相纸模板（images/print-bg/&lt;相纸&gt;/&lt;尺寸&gt;.png）
/// 复制到模板的各个槽位上，输出一张可直接拿去冲印的图（默认 300dpi）。
///
/// 用法（在 csharp 目录运行）：
///   dotnet run -- print                  # 使用默认输入（output/v3-diy/blue.jpg，需先运行 v3-diy 生成）
///   dotnet run -- print 路径/照片.jpg     # 指定一张做好的证件照
///
/// 依赖：SixLabors.ImageSharp（纯托管绘图库，首次运行自动还原；其余示例仍零依赖）。
/// 前置：请先运行 v3-diy 生成 295x413 的证件照到 output/v3-diy/blue.jpg，
///       或自行用任意制作接口产出对应尺寸的证件照后传入路径。
/// </summary>
public static class Print
{
    // ===== 请求参数（按需修改）=====
    private const string Paper = "print6";            // 相纸规格目录名（images/print-bg/ 下的子目录，print6=6寸）
    private const bool UseTextTpl = false;              // 是否使用带文字/裁切线说明的模板（*-text.png）
    private const int Dpi = 300;                        // 输出分辨率
    private const string DefaultIdPhoto = "output/v3-diy/blue.jpg"; // 默认输入：做好的证件照（1寸 295x413）

    // 排版配置（数据驱动，便于增删规格）
    // key        = 证件照像素尺寸（宽-高）
    // Rotate     = 是否需要旋转 90°（模板槽位为横放时置 true）
    // SlotW/SlotH = 槽位像素尺寸（输入图会被缩放到该尺寸后再贴入）
    // Positions  = 各槽位左上角坐标 [x, y]（与模板图严格对应）
    private sealed record Layout(bool Rotate, int SlotW, int SlotH, int[][] Positions);

    private static readonly Dictionary<string, Layout> Layouts = new()
    {
        ["295-413"]   = new(false, 295, 413, new int[][]
            { new[]{120,180}, new[]{435,180}, new[]{750,180}, new[]{1065,180}, new[]{1380,180},
              new[]{120,613}, new[]{435,613}, new[]{750,613}, new[]{1065,613}, new[]{1380,613} }),
        ["260-378"]   = new(false, 260, 378, new int[][]
            { new[]{207,214}, new[]{487,214}, new[]{767,214}, new[]{1047,214}, new[]{1327,214},
              new[]{207,612}, new[]{487,612}, new[]{767,612}, new[]{1047,612}, new[]{1327,612} }),
        ["413-579"]   = new(true, 413, 579, new int[][]
            { new[]{308,179}, new[]{907,179}, new[]{308,612}, new[]{907,612} }),
        ["390-567"]   = new(true, 390, 567, new int[][]
            { new[]{320,202}, new[]{907,202}, new[]{320,612}, new[]{907,612} }),
        ["413-531"]   = new(true, 413, 531, new int[][]
            { new[]{81,179}, new[]{632,179}, new[]{1183,179}, new[]{81,612}, new[]{632,612}, new[]{1183,612} }),
        ["626-413"]   = new(false, 626, 413, new int[][]
            { new[]{261,179}, new[]{907,179}, new[]{261,612}, new[]{907,612} }),
        ["649-991"]   = new(false, 649, 991, new int[][]
            { new[]{238,107}, new[]{907,107} }),
        ["1050-1499"] = new(true, 1050, 1499, new int[][]
            { new[]{148,78} }),
        // 以下三组坐标由"统一边距+20px间距"的均匀网格推算，并与模板裁切线对齐确认（模板本身未画占位框）
        ["358-441"]   = new(false, 358, 441, new int[][]
            { new[]{151,151}, new[]{529,151}, new[]{907,151}, new[]{1285,151},
              new[]{151,612}, new[]{529,612}, new[]{907,612}, new[]{1285,612} }),
        ["378-472"]   = new(false, 378, 472, new int[][]
            { new[]{111,120}, new[]{509,120}, new[]{907,120}, new[]{1305,120},
              new[]{111,612}, new[]{509,612}, new[]{907,612}, new[]{1305,612} }),
        ["600-600"]   = new(false, 600, 600, new int[][]
            { new[]{287,302}, new[]{907,302} }),
    };

    public static async Task RunAsync()
    {
        var idPhoto = Program.ExtraArgs.Length > 0 ? Program.ExtraArgs[0] : DefaultIdPhoto;

        if (!File.Exists(idPhoto))
        {
            Console.WriteLine($"输入证件照不存在：{idPhoto}");
            Console.WriteLine("请先运行 v3-diy 生成，或在命令行传入路径：dotnet run -- print 路径/照片.jpg");
            return;
        }

        using var photo = await Image.LoadAsync(idPhoto);
        var printName = $"{photo.Width}-{photo.Height}";

        if (!Layouts.TryGetValue(printName, out var layout))
        {
            Console.WriteLine($"不支持的证件照尺寸 {printName}（实际像素 {photo.Width}x{photo.Height}）。");
            Console.WriteLine("支持：" + string.Join(", ", Layouts.Keys));
            Console.WriteLine("请使用对应规格制作证件照（如 v3-diy 设 WidthPx/HeightPx）。");
            return;
        }

        var suffix = UseTextTpl ? "-text" : "";
        var bgPath = $"../images/print-bg/{Paper}/{printName}{suffix}.png";
        if (!File.Exists(bgPath))
        {
            Console.WriteLine($"相纸模板不存在：{bgPath}");
            Console.WriteLine($"请确认 images/print-bg/{Paper}/ 下存在 {printName}{suffix}.png");
            return;
        }

        // 先铺白底，避免模板透明区在 JPEG 中变黑
        using var template = await Image.LoadAsync(bgPath);
        using var canvas = new Image<Rgba32>(template.Width, template.Height, Color.White);
        canvas.Mutate(x => x.DrawImage(template, new Point(0, 0), 1f));

        // 缩放（拉伸）到槽位尺寸，需要旋转则顺时针 90°
        photo.Mutate(x => x.Resize(layout.SlotW, layout.SlotH));
        if (layout.Rotate)
            photo.Mutate(x => x.Rotate(90));

        foreach (var pos in layout.Positions)
            canvas.Mutate(x => x.DrawImage(photo, new Point(pos[0], pos[1]), 1f));

        // 设置输出 DPI（冲印店据此换算物理尺寸）
        canvas.Metadata.HorizontalResolution = Dpi;
        canvas.Metadata.VerticalResolution = Dpi;

        var saveDir = Path.Combine("output", "print");
        Directory.CreateDirectory(saveDir);
        var savePath = Path.Combine(saveDir, $"{Paper}-{printName}.jpg");
        await canvas.SaveAsync(savePath, new JpegEncoder { Quality = 95 });

        Console.WriteLine($"排版完成：{savePath}");
    }
}
