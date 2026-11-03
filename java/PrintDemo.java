import java.awt.Color;
import java.awt.Graphics2D;
import java.awt.Image;
import java.awt.image.BufferedImage;
import java.io.File;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import javax.imageio.IIOImage;
import javax.imageio.ImageIO;
import javax.imageio.ImageWriteParam;
import javax.imageio.ImageWriter;
import javax.imageio.metadata.IIOMetadata;
import javax.imageio.metadata.IIOMetadataNode;
import javax.imageio.stream.ImageOutputStream;
import org.w3c.dom.Node;
import java.util.Iterator;
import java.util.LinkedHashMap;
import java.util.Map;

/**
 * 证件照冲印排版（相纸联排）。
 *
 * 将一张已制作好的证件照，按预设的相纸模板（../images/print-bg/&lt;相纸&gt;/&lt;尺寸&gt;.png）
 * 复制到模板的各个槽位上，输出一张可直接拿去冲印的图（默认 300dpi）。
 *
 * 用法（在 java 目录编译后运行）：
 *   javac -encoding UTF-8 *.java
 *   java PrintDemo                  # 使用默认输入（output/v3-diy/blue.jpg，需先运行 V3DiyDemo 生成）
 *   java PrintDemo 路径/照片.jpg     # 指定一张做好的证件照
 *
 * 依赖：JDK 自带 javax.imageio，无需任何第三方库。
 * 前置：请先运行 V3DiyDemo 生成 295x413 的证件照到 output/v3-diy/blue.jpg，
 *       或自行用任意制作接口产出对应尺寸的证件照后传入路径。
 */
public class PrintDemo {

    // ===== 请求参数（按需修改）=====
    private static final String PAPER = "print6";          // 相纸规格目录名（images/print-bg/ 下的子目录，print6=6寸）
    private static final boolean USE_TEXT_TPL = false;      // 是否使用带文字/裁切线说明的模板（*-text.png）
    private static final int DPI = 300;                     // 输出分辨率
    private static final String DEFAULT_ID_PHOTO = "output/v3-diy/blue.jpg"; // 默认输入：做好的证件照（1寸 295x413）

    // 排版配置（数据驱动，便于增删规格）
    // rotate = 是否需要旋转 90°；w/h = 槽位像素尺寸；positions = 各槽位左上角坐标 [x, y]
    private static final class Layout {
        final boolean rotate;
        final int w, h;
        final int[][] positions;
        Layout(boolean rotate, int w, int h, int[][] positions) {
            this.rotate = rotate; this.w = w; this.h = h; this.positions = positions;
        }
    }

    private static final Map<String, Layout> LAYOUTS = new LinkedHashMap<>();
    static {
        LAYOUTS.put("295-413",   new Layout(false, 295, 413, new int[][]{{120,180},{435,180},{750,180},{1065,180},{1380,180},{120,613},{435,613},{750,613},{1065,613},{1380,613}}));
        LAYOUTS.put("260-378",   new Layout(false, 260, 378, new int[][]{{207,214},{487,214},{767,214},{1047,214},{1327,214},{207,612},{487,612},{767,612},{1047,612},{1327,612}}));
        LAYOUTS.put("413-579",   new Layout(true,  413, 579, new int[][]{{308,179},{907,179},{308,612},{907,612}}));
        LAYOUTS.put("390-567",   new Layout(true,  390, 567, new int[][]{{320,202},{907,202},{320,612},{907,612}}));
        LAYOUTS.put("413-531",   new Layout(true,  413, 531, new int[][]{{81,179},{632,179},{1183,179},{81,612},{632,612},{1183,612}}));
        LAYOUTS.put("626-413",   new Layout(false, 626, 413, new int[][]{{261,179},{907,179},{261,612},{907,612}}));
        LAYOUTS.put("649-991",   new Layout(false, 649, 991, new int[][]{{238,107},{907,107}}));
        LAYOUTS.put("1050-1499", new Layout(true, 1050,1499, new int[][]{{148,78}}));
        // 以下三组坐标由"统一边距+20px间距"的均匀网格推算，并与模板裁切线对齐确认（模板本身未画占位框）
        LAYOUTS.put("358-441",   new Layout(false, 358, 441, new int[][]{{151,151},{529,151},{907,151},{1285,151},{151,612},{529,612},{907,612},{1285,612}}));
        LAYOUTS.put("378-472",   new Layout(false, 378, 472, new int[][]{{111,120},{509,120},{907,120},{1305,120},{111,612},{509,612},{907,612},{1305,612}}));
        LAYOUTS.put("600-600",   new Layout(false, 600, 600, new int[][]{{287,302},{907,302}}));
    }

    public static void main(String[] args) throws Exception {
        String idPhoto = args.length > 0 ? args[0] : DEFAULT_ID_PHOTO;
        File in = new File(idPhoto);
        if (!in.isFile()) {
            System.err.println("输入证件照不存在：" + idPhoto);
            System.err.println("请先运行 V3DiyDemo 生成，或在命令行传入路径：java PrintDemo 路径/照片.jpg");
            return;
        }

        BufferedImage photo = ImageIO.read(in);
        String printName = photo.getWidth() + "-" + photo.getHeight();
        Layout layout = LAYOUTS.get(printName);
        if (layout == null) {
            System.err.println("不支持的证件照尺寸 " + printName + "（实际像素 " + photo.getWidth() + "x" + photo.getHeight() + "）。");
            System.err.println("支持：" + String.join(", ", LAYOUTS.keySet()));
            System.err.println("请使用对应规格制作证件照（如 V3DiyDemo 设 WIDTH_PX/HEIGHT_PX）。");
            return;
        }

        String suffix = USE_TEXT_TPL ? "-text" : "";
        String bgPath = "../images/print-bg/" + PAPER + "/" + printName + suffix + ".png";
        File bgFile = new File(bgPath);
        if (!bgFile.isFile()) {
            System.err.println("相纸模板不存在：" + bgPath);
            System.err.println("请确认 images/print-bg/" + PAPER + "/ 下存在 " + printName + suffix + ".png");
            return;
        }

        // 先铺白底，避免模板透明区在 JPEG 中变黑
        BufferedImage tpl = ImageIO.read(bgFile);
        int tw = tpl.getWidth(), th = tpl.getHeight();
        BufferedImage canvas = new BufferedImage(tw, th, BufferedImage.TYPE_INT_RGB);
        Graphics2D g = canvas.createGraphics();
        g.setColor(Color.WHITE);
        g.fillRect(0, 0, tw, th);
        g.drawImage(tpl, 0, 0, null);

        // 缩放（拉伸）到槽位尺寸，需要旋转则顺时针 90°
        Image scaled = photo.getScaledInstance(layout.w, layout.h, Image.SCALE_SMOOTH);
        BufferedImage slot = toBufferedImage(scaled, layout.w, layout.h);
        if (layout.rotate) {
            slot = rotate90CW(slot);
        }

        for (int[] pos : layout.positions) {
            g.drawImage(slot, pos[0], pos[1], null);
        }
        g.dispose();

        // 设置输出 DPI（冲印店据此换算物理尺寸）
        Path outDir = ZjzApiClient.outputDir("print");
        Files.createDirectories(outDir);
        Path out = outDir.resolve(PAPER + "-" + printName + ".jpg");
        writeJpegWithDpi(canvas, out, DPI);

        System.out.println("排版完成：" + out);
    }

    private static BufferedImage toBufferedImage(Image img, int w, int h) {
        BufferedImage b = new BufferedImage(w, h, BufferedImage.TYPE_INT_RGB);
        Graphics2D g = b.createGraphics();
        g.drawImage(img, 0, 0, null);
        g.dispose();
        return b;
    }

    private static BufferedImage rotate90CW(BufferedImage src) {
        int w = src.getWidth(), h = src.getHeight();
        BufferedImage dst = new BufferedImage(h, w, BufferedImage.TYPE_INT_RGB);
        Graphics2D g = dst.createGraphics();
        g.translate(h / 2.0, w / 2.0);
        g.rotate(Math.toRadians(90));
        g.translate(-w / 2.0, -h / 2.0);
        g.drawImage(src, 0, 0, null);
        g.dispose();
        return dst;
    }

    private static void writeJpegWithDpi(BufferedImage img, Path out, int dpi) throws IOException {
        Iterator<ImageWriter> writers = ImageIO.getImageWritersByFormatName("jpeg");
        if (!writers.hasNext()) {
            ImageIO.write(img, "jpeg", out.toFile());
            return;
        }
        ImageWriter writer = writers.next();
        ImageWriteParam param = writer.getDefaultWriteParam();
        param.setCompressionMode(ImageWriteParam.MODE_EXPLICIT);
        param.setCompressionQuality(0.95f);
        IIOMetadata metadata = writer.getDefaultImageMetadata(
                new javax.imageio.ImageTypeSpecifier(img), param);
        setDpi(metadata, dpi);
        try (ImageOutputStream ios = ImageIO.createImageOutputStream(out.toFile())) {
            writer.setOutput(ios);
            writer.write(null, new IIOImage(img, null, metadata), param);
        } finally {
            writer.dispose();
        }
    }

    private static void setDpi(IIOMetadata metadata, int dpi) {
        try {
            String format = metadata.getNativeMetadataFormatName();
            IIOMetadataNode root = (IIOMetadataNode) metadata.getAsTree(format);
            IIOMetadataNode jfif = null;
            for (int i = 0; i < root.getLength(); i++) {
                Node n = root.item(i);
                if (n instanceof IIOMetadataNode
                        && "JPEGvariety".equals(((IIOMetadataNode) n).getNodeName())) {
                    IIOMetadataNode variety = (IIOMetadataNode) n;
                    for (int j = 0; j < variety.getLength(); j++) {
                        Node c = variety.item(j);
                        if (c instanceof IIOMetadataNode
                                && "app0JFIF".equals(((IIOMetadataNode) c).getNodeName())) {
                            jfif = (IIOMetadataNode) c;
                        }
                    }
                }
            }
            if (jfif == null) {
                jfif = new IIOMetadataNode("app0JFIF");
                IIOMetadataNode variety = new IIOMetadataNode("JPEGvariety");
                variety.appendChild(jfif);
                root.appendChild(variety);
            }
            jfif.setAttribute("Xdensity", String.valueOf(dpi));
            jfif.setAttribute("Ydensity", String.valueOf(dpi));
            jfif.setAttribute("resUnits", "1"); // 1 = 像素/英寸
            metadata.setFromTree(format, root);
        } catch (Exception e) {
            // DPI 设置失败时退化为不带 DPI 的普通 JPEG
            System.err.println("警告：设置 DPI 失败（" + e.getMessage() + "），输出将不含 DPI 信息。");
        }
    }
}
