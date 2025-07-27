package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"log"
)

// 定义所有的道具图片
var (
	itemImages = map[string]*ebiten.Image{
		"1001": loadImage("photos/type/jinBi.png"), // 使用相对路径
	}
)

// loadImage 加载图片文件
func loadImage(filename string) *ebiten.Image {
	// 使用相对路径
	imgPath := filename
	img, _, err := ebitenutil.NewImageFromFile(imgPath)
	if err != nil {
		log.Fatalf("加载图片失败 %s: %v", imgPath, err)
	}
	return img
}
