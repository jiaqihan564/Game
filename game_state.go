package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// GameState 定义所有界面状态类型
type GameState int

const (
	StateMenu GameState = iota // 开始菜单界面
	StatePlay                  // 游戏主界面
)

// Game 结构体：主程序运行载体
type Game struct {
	currentState GameState // 当前界面状态

	menuScreen *MenuScreen // 菜单界面
	playScreen *PlayScreen // 游戏界面
}

// NewGame 初始化游戏
func NewGame() *Game {
	return &Game{
		currentState: StateMenu,
		menuScreen:   NewMenuScreen(),
		playScreen:   NewPlayScreen(),
	}
}

// Update 每帧更新逻辑
func (g *Game) Update() error {
	switch g.currentState {
	case StateMenu:
		// 菜单界面更新逻辑
		if g.menuScreen.Update() {
			g.currentState = StatePlay // 切换到游戏界面
		}
	case StatePlay:
		g.playScreen.Update()
	}
	return nil
}

// Draw 渲染逻辑
func (g *Game) Draw(screen *ebiten.Image) {
	switch g.currentState {
	case StateMenu:
		g.menuScreen.Draw(screen)
	case StatePlay:
		g.playScreen.Draw(screen)
	}

	// 显示游戏帧率
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FPS: %.2f", ebiten.ActualFPS()), 10, 10)
}

// Layout 添加Layout方法实现ebiten.Game接口
func (g *Game) Layout(int, int) (screenWidth int, screenHeight int) {
	// 返回游戏逻辑屏幕尺寸，可以与窗口尺寸不同
	return 800, 600
}
