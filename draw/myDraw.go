package draw

import (
	"aiedge/post"
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

// type Rectangle struct {
// 	Bottom float64
// 	Left   float64
// 	Right  float64
// 	Top    float64
// }

//TextBrush 字体相关
type TextBrush struct {
	FontType  *truetype.Font
	FontSize  float64
	FontColor *image.Uniform
	TextWidth int
}

//Image2RGBA Image2RGBA
func Image2RGBA(img image.Image) *image.RGBA {

	baseSrcBounds := img.Bounds().Max

	newWidth := baseSrcBounds.X
	newHeight := baseSrcBounds.Y

	des := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight)) // 底板
	//首先将一个图片信息存入jpg
	draw.Draw(des, des.Bounds(), img, img.Bounds().Min, draw.Over)

	return des
}

//NewTextBrush 新生成笔刷
func NewTextBrush(FontFilePath string, FontSize float64, FontColor *image.Uniform, textWidth int) (*TextBrush, error) {
	fontFile, err := ioutil.ReadFile(FontFilePath)
	if err != nil {
		return nil, err
	}
	fontType, err := truetype.Parse(fontFile)
	if err != nil {
		return nil, err
	}
	if textWidth <= 0 {
		textWidth = 20
	}
	return &TextBrush{FontType: fontType, FontSize: FontSize, FontColor: FontColor, TextWidth: textWidth}, nil
}

func DrawRectOnImageAndSave(imageSavePath string, imageData []byte, infos []post.Rectangle) (err error) {
	//判断图片类型
	var backgroud image.Image
	filetype := http.DetectContentType(imageData)
	switch filetype {
	case "image/jpeg", "image/jpg":
		backgroud, err = jpeg.Decode(bytes.NewReader(imageData))
		if err != nil {
			fmt.Println("jpeg error")
			return err
		}

	case "image/gif":
		backgroud, err = gif.Decode(bytes.NewReader(imageData))
		if err != nil {
			return err
		}

	case "image/png":
		backgroud, err = png.Decode(bytes.NewReader(imageData))
		if err != nil {
			return err
		}
	default:
		return err
	}
	des := Image2RGBA(backgroud)
	//新建笔刷
	textBrush, _ := NewTextBrush("./fonts/arial.ttf", 15, image.Black, 15)
	for _, info := range infos {
		var c *freetype.Context
		c = freetype.NewContext()
		c.SetDPI(72)
		c.SetFont(textBrush.FontType)
		c.SetHinting(font.HintingFull)
		c.SetFontSize(textBrush.FontSize)
		c.SetClip(des.Bounds())
		c.SetDst(des)
		cGreen := image.NewUniform(color.RGBA{
			R: 0xFF,
			G: 0,
			B: 0,
			A: 255,
		})

		c.SetSrc(cGreen)
		for i := info.Left; i < info.Right; i++ { //画固定Botton和Top的两条横线
			c.DrawString("·", freetype.Pt(int(i), int(info.Bottom)))
			c.DrawString("·", freetype.Pt(int(i), int(info.Top)))
		}

		for j := info.Top; j < info.Bottom; j++ { //画固定Left和Right的两条竖线
			c.DrawString("·", freetype.Pt(int(info.Left), int(j)))
			c.DrawString("·", freetype.Pt(int(info.Right), int(j)))
		}
	}

	//保存图片
	fSave, err := os.Create(imageSavePath)
	if err != nil {
		return err
	}
	defer fSave.Close()

	err = jpeg.Encode(fSave, des, nil)

	if err != nil {
		return err
	}
	return nil
}
