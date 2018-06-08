package utils

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"log"
	"os"
	"strconv"
	"testing"
)

func TestGIF(t *testing.T) {
	file, err := os.Open("test.jpeg")
	checkError(err)
	defer file.Close()

	srcImg, err := jpeg.Decode(file)
	checkError(err)

	minX := srcImg.Bounds().Min.X
	minY := srcImg.Bounds().Min.Y
	maxX := srcImg.Bounds().Max.X
	maxY := srcImg.Bounds().Max.Y

	w := (maxX - minX) / 3
	h := (maxY - minY) / 3

	imgs := make([]image.Image, 0)
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			newFile, err := os.Create("test_" + strconv.Itoa(i*3+j) + ".jpg")
			checkError(err)

			// 裁剪
			subImg, err := Clip(srcImg, i*w, j*h, w, h)
			checkError(err)

			// 先保存成jpeg
			jpeg.Encode(newFile, subImg, &jpeg.Options{100})
			newFile.Close()

			// 读取jpeg
			fuckFile, err := os.Open("test_" + strconv.Itoa(i*3+j) + ".jpg")
			checkError(err)
			defer fuckFile.Close()
			subImg, _ = jpeg.Decode(fuckFile)

			imgs = append(imgs, subImg)
		}
	}

	anim := MergeGif(imgs, 50)

	gifFile, _ := os.Create("test.gif")
	defer gifFile.Close()
	gif.EncodeAll(gifFile, anim)

	fmt.Println("hello")
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
