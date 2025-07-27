package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func main() {
	// 初始化游戏主结构体（包含状态）
	game := NewGame()

	// 设置窗口属性
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("贾先生的2D游戏 - 开始界面 Demo")

	// 如果需要同步FPS
	ebiten.SetScreenClearedEveryFrame(true)

	// 启动游戏主循环
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
