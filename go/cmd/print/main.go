// 证件照冲印排版（相纸联排）。
//
// 将一张已制作好的证件照，按预设的相纸模板（../images/print-bg/<相纸>/<尺寸>.png）
// 复制到模板的各个槽位上，输出一张可直接拿去冲印的图（默认 300dpi）。
//
// 用法（在 go 目录运行）：
//   go run ./cmd/print                 # 使用默认输入（output/v3-diy/blue.jpg，需先运行 v3/diy 生成）
//   go run ./cmd/print 路径/照片.jpg     # 指定一张做好的证件照
//
// 依赖：仅使用 Go 标准库（image / image/draw / image/jpeg / image/png），无需任何第三方库。
// 前置：请先运行 v3/diy 生成 295x413 的证件照到 output/v3-diy/blue.jpg，
//       或自行用任意制作接口产出对应尺寸的证件照后传入路径。
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"
	"path/filepath"

	_ "image/png" // 注册 PNG 解码器
)

const (
	paper       = "print6"
	useTextTpl  = false
	dpi         = 300
	defaultIn   = "output/v3-diy/blue.jpg"
	outputSub   = "output/print"
	templateDir = "../images/print-bg/print6"
)

// layout 排版配置（数据驱动，便于增删规格）
type layout struct {
	rotate    bool
	w, h      int
	positions [][2]int
}

var layouts = map[string]layout{
	"295-413":   {false, 295, 413, [][2]int{{120, 180}, {435, 180}, {750, 180}, {1065, 180}, {1380, 180}, {120, 613}, {435, 613}, {750, 613}, {1065, 613}, {1380, 613}}},
	"260-378":   {false, 260, 378, [][2]int{{207, 214}, {487, 214}, {767, 214}, {1047, 214}, {1327, 214}, {207, 612}, {487, 612}, {767, 612}, {1047, 612}, {1327, 612}}},
	"413-579":   {true, 413, 579, [][2]int{{308, 179}, {907, 179}, {308, 612}, {907, 612}}},
	"390-567":   {true, 390, 567, [][2]int{{320, 202}, {907, 202}, {320, 612}, {907, 612}}},
	"413-531":   {true, 413, 531, [][2]int{{81, 179}, {632, 179}, {1183, 179}, {81, 612}, {632, 612}, {1183, 612}}},
	"626-413":   {false, 626, 413, [][2]int{{261, 179}, {907, 179}, {261, 612}, {907, 612}}},
	"649-991":   {false, 649, 991, [][2]int{{238, 107}, {907, 107}}},
	"1050-1499": {true, 1050, 1499, [][2]int{{148, 78}}},
	// 以下三组坐标由"统一边距+20px间距"的均匀网格推算，并与模板裁切线对齐确认（模板本身未画占位框）
	"358-441": {false, 358, 441, [][2]int{{151, 151}, {529, 151}, {907, 151}, {1285, 151}, {151, 612}, {529, 612}, {907, 612}, {1285, 612}}},
	"378-472": {false, 378, 472, [][2]int{{111, 120}, {509, 120}, {907, 120}, {1305, 120}, {111, 612}, {509, 612}, {907, 612}, {1305, 612}}},
	"600-600": {false, 600, 600, [][2]int{{287, 302}, {907, 302}}},
}

func main() {
	idPhoto := defaultIn
	if len(os.Args) > 1 {
		idPhoto = os.Args[1]
	}

	f, err := os.Open(idPhoto)
	if err != nil {
		fmt.Printf("输入证件照不存在：%s\n请先运行 v3/diy 生成，或在命令行传入路径：go run ./cmd/print 路径/照片.jpg\n", idPhoto)
		return
	}
	srcImg, _, err := image.Decode(f)
	f.Close()
	if err != nil {
		fmt.Println("读取证件照失败：", err)
		return
	}

	b := srcImg.Bounds()
	name := fmt.Sprintf("%d-%d", b.Dx(), b.Dy())
	lay, ok := layouts[name]
	if !ok {
		fmt.Printf("不支持的证件照尺寸 %s（实际像素 %dx%d）。\n支持：", name, b.Dx(), b.Dy())
		for k := range layouts {
			fmt.Print(k, " ")
		}
		fmt.Println()
		return
	}

	suffix := ""
	if useTextTpl {
		suffix = "-text"
	}
	bgPath := filepath.Join(templateDir, name+suffix+".png")
	bgFile, err := os.Open(bgPath)
	if err != nil {
		fmt.Println("相纸模板不存在：", bgPath)
		return
	}
	tplImg, _, err := image.Decode(bgFile)
	bgFile.Close()
	if err != nil {
		fmt.Println("读取模板失败：", err)
		return
	}
	tb := tplImg.Bounds()
	tw, th := tb.Dx(), tb.Dy()

	// 先铺白底，避免模板透明区在 JPEG 中变黑
	canvas := image.NewRGBA(image.Rect(0, 0, tw, th))
	draw.Draw(canvas, canvas.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
	draw.Draw(canvas, canvas.Bounds(), tplImg, tb.Min, draw.Over)

	// 缩放（拉伸）到槽位尺寸，需要旋转则顺时针 90°
	slot := resize(srcImg, lay.w, lay.h)
	if lay.rotate {
		slot = rotate90CW(slot)
	}
	for _, p := range lay.positions {
		r := image.Rect(p[0], p[1], p[0]+slot.Bounds().Dx(), p[1]+slot.Bounds().Dy())
		draw.Draw(canvas, r, slot, slot.Bounds().Min, draw.Over)
	}

	_ = os.MkdirAll(outputSub, 0o755)
	outPath := filepath.Join(outputSub, fmt.Sprintf("%s-%s.jpg", paper, name))
	if err := writeJpegWithDpi(canvas, outPath, dpi); err != nil {
		fmt.Println("保存失败：", err)
		return
	}
	fmt.Println("排版完成：", outPath)
}

// resize 用最近邻将图片缩放到指定尺寸（有效输入尺寸本就一致，此处兼容轻微误差）
func resize(src image.Image, w, h int) image.Image {
	sb := src.Bounds()
	sw, sh := sb.Dx(), sb.Dy()
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		sy := y * sh / h
		for x := 0; x < w; x++ {
			sx := x * sw / w
			dst.Set(x, y, src.At(sb.Min.X+sx, sb.Min.Y+sy))
		}
	}
	return dst
}

// rotate90CW 顺时针旋转 90°
func rotate90CW(src image.Image) image.Image {
	sb := src.Bounds()
	w, h := sb.Dx(), sb.Dy()
	dst := image.NewRGBA(image.Rect(0, 0, h, w))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dx := h - 1 - y
			dy := x
			dst.Set(dx, dy, src.At(sb.Min.X+x, sb.Min.Y+y))
		}
	}
	return dst
}

// writeJpegWithDpi 编码 JPEG 并写入 300dpi 的 JFIF 密度信息
func writeJpegWithDpi(img image.Image, path string, dpi int) error {
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 95}); err != nil {
		return err
	}
	out := patchJfifDpi(buf.Bytes(), dpi)
	return os.WriteFile(path, out, 0o644)
}

// patchJfifDpi 设置 JFIF 分辨率；标准库 JPEG 编码结果没有 APP0 时补入该段。
func patchJfifDpi(data []byte, dpi int) []byte {
	if len(data) < 4 || data[0] != 0xFF || data[1] != 0xD8 {
		return data
	}
	i := 2
	for i+4 <= len(data) {
		if data[i] != 0xFF {
			break
		}
		marker := data[i+1]
		if marker == 0xDA || marker == 0xD9 { // 扫描数据或文件结束
			break
		}
		segment_length := int(binary.BigEndian.Uint16(data[i+2 : i+4]))
		if segment_length < 2 || segment_length > len(data)-i-2 {
			return data // 截断或非法段，不修改原始内容
		}
		if marker == 0xE0 && segment_length >= 16 && string(data[i+4:i+9]) == "JFIF\x00" {
			out := append([]byte(nil), data...)
			out[i+11] = 1 // resUnits = 1 (dots per inch)
			binary.BigEndian.PutUint16(out[i+12:i+14], uint16(dpi))
			binary.BigEndian.PutUint16(out[i+14:i+16], uint16(dpi))
			return out
		}
		i += 2 + segment_length
	}
	// APP0: 长度16，JFIF 1.01，像素/英寸，无缩略图。
	app0 := []byte{0xFF, 0xE0, 0, 16, 'J', 'F', 'I', 'F', 0, 1, 1, 1, 0, 0, 0, 0, 0, 0}
	binary.BigEndian.PutUint16(app0[12:14], uint16(dpi))
	binary.BigEndian.PutUint16(app0[14:16], uint16(dpi))
	out := make([]byte, 0, len(data)+len(app0))
	out = append(out, data[:2]...)
	out = append(out, app0...)
	return append(out, data[2:]...)
}
