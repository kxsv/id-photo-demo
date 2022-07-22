from PIL import Image
import os


def print_main():
    # 支持排版的尺寸
    print_item = ['260-378', '295-413', '390-567', '413-531', '413-579', '626-413', '649-991', '1050-1499']

    id_photo = Image.open(os.getcwd() + '/../images/test.jpeg')  # 图片地址
    width_px, height_px = id_photo.size
    print_name = str(width_px) + '-' + str(height_px)

    if print_name not in print_item:
        print('不支持排版')
        exit(0)

    print_bg = Image.open(os.getcwd() + '/../images/print-bg/' + print_name + '.png')  # 获取排版背景图

    if print_name == '295-413':
        print_bg.paste(id_photo, (120, 180))
        print_bg.paste(id_photo, (435, 180))
        print_bg.paste(id_photo, (750, 180))
        print_bg.paste(id_photo, (1065, 180))
        print_bg.paste(id_photo, (1380, 180))
        print_bg.paste(id_photo, (120, 613))
        print_bg.paste(id_photo, (435, 613))
        print_bg.paste(id_photo, (750, 613))
        print_bg.paste(id_photo, (1065, 613))
        print_bg.paste(id_photo, (1380, 613))
    elif print_name == '413-579':
        id_photo_rotate = id_photo.transpose(Image.ROTATE_90)
        print_bg.paste(id_photo_rotate, (308, 179))
        print_bg.paste(id_photo_rotate, (907, 179))
        print_bg.paste(id_photo_rotate, (308, 612))
        print_bg.paste(id_photo_rotate, (907, 612))
    elif print_name == '390-567':
        id_photo_rotate = id_photo.transpose(Image.ROTATE_90)
        print_bg.paste(id_photo_rotate, (320, 202))
        print_bg.paste(id_photo_rotate, (907, 202))
        print_bg.paste(id_photo_rotate, (320, 612))
        print_bg.paste(id_photo_rotate, (907, 612))
    elif print_name == '260-378':
        print_bg.paste(id_photo, (207, 214))
        print_bg.paste(id_photo, (487, 214))
        print_bg.paste(id_photo, (767, 214))
        print_bg.paste(id_photo, (1047, 214))
        print_bg.paste(id_photo, (1327, 214))
        print_bg.paste(id_photo, (207, 612))
        print_bg.paste(id_photo, (487, 612))
        print_bg.paste(id_photo, (767, 612))
        print_bg.paste(id_photo, (1047, 612))
        print_bg.paste(id_photo, (1327, 612))
    elif print_name == '413-531':
        id_photo_rotate = id_photo.transpose(Image.ROTATE_90)
        print_bg.paste(id_photo_rotate, (81, 179))
        print_bg.paste(id_photo_rotate, (632, 179))
        print_bg.paste(id_photo_rotate, (1183, 179))
        print_bg.paste(id_photo_rotate, (81, 612))
        print_bg.paste(id_photo_rotate, (632, 612))
        print_bg.paste(id_photo_rotate, (1183, 612))
    elif print_name == '626-413':
        print_bg.paste(id_photo, (261, 179))
        print_bg.paste(id_photo, (907, 179))
        print_bg.paste(id_photo, (261, 612))
        print_bg.paste(id_photo, (907, 612))
    elif print_name == '649-991':
        print_bg.paste(id_photo, (238, 107))
        print_bg.paste(id_photo, (907, 107))
    elif print_name == '1050-1499':
        id_photo_rotate = id_photo.transpose(Image.ROTATE_90)
        print_bg.paste(id_photo_rotate, (148, 78))

    path = os.getcwd() + "/res-print.jpeg"
    print_bg.save(path)
    print_bg.show()


if __name__ == '__main__':
    print_main()
