package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"log"
)

type item struct {
	id        int64
	count     int64 // 数量
	*ItemData       // 道具数据
}

// ItemData 初始化游戏道具数据的结构体
type ItemData struct {
	id    int64
	Name  string        // 物品名称
	Image *ebiten.Image // 物品图片
}

// ItemImages 定义所有的游戏道具
var (
	ItemImages = map[int64]*ItemData{
		1001: {
			id:    1001,
			Name:  "Gold", // 金币
			Image: loadImage("photos/type/jinBi.png"),
		},
		1002: {
			id:    1002,
			Name:  "SwordXinShou", // 新手剑
			Image: loadImage("photos/type/SwordXinShou.png"),
		},
		1003: {
			id:    1003,
			Name:  "Sword1", // 一级剑
			Image: loadImage("photos/type/Sword1.png"),
		},
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
