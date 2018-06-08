package utils

import (
	"errors"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
)

// 图片裁剪
func Clip(srcImg image.Image, x, y, w, h int) (image.Image, error) {
	var subImg image.Image
	if rgbImg, ok := srcImg.(*image.YCbCr); ok {
		subImg = rgbImg.SubImage(image.Rect(x, y, x+w, y+h)).(*image.YCbCr)
	} else if rgbImg, ok := srcImg.(*image.RGBA); ok {
		subImg = rgbImg.SubImage(image.Rect(x, y, x+w, y+h)).(*image.RGBA)
	} else if rgbImg, ok := srcImg.(*image.NRGBA); ok {
		subImg = rgbImg.SubImage(image.Rect(x, y, x+w, y+h)).(*image.NRGBA)
	} else {
		return subImg, errors.New("not support format")
	}
	return subImg, nil
}

// 将多张图片合成一张gif，每一帧的间隔时间为delay，单位为1/100秒
func MergeGif(imgs []image.Image, delay int) *gif.GIF {
	anim := &gif.GIF{}
	for i := 0; i < len(imgs); i++ {
		bounds := imgs[i].Bounds()
		palettedImage := image.NewPaletted(bounds, palette.WebSafe) // or palette.Plan9
		draw.Draw(palettedImage, palettedImage.Bounds(), imgs[i], bounds.Min, draw.Over)

		anim.Image = append(anim.Image, palettedImage)
		anim.Delay = append(anim.Delay, delay)
	}

	return anim
}
