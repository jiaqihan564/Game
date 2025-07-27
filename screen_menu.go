package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
)

// MenuScreen 定义菜单界面结构
type MenuScreen struct {
	startButtonRect [4]int        // [x, y, width, height]
	clicked         bool          // 是否被点击
	backgroundImage *ebiten.Image // 背景图片（新增字段）
}

// NewMenuScreen 构造函数
func NewMenuScreen() *MenuScreen {
	m := &MenuScreen{
		startButtonRect: [4]int{220, 200, 200, 50},
	}

	// 预加载背景图片（新增代码）
	var err error
	m.backgroundImage, _, err = ebitenutil.NewImageFromFile("photos/beijing.png")
	if err != nil {
		panic(err)
	}

	return m
}

// Update 处理鼠标点击逻辑，点击后返回 true 切换场景
func (m *MenuScreen) Update() bool {
	x, y := ebiten.CursorPosition()
	btnX, btnY, btnW, btnH := m.startButtonRect[0], m.startButtonRect[1], m.startButtonRect[2], m.startButtonRect[3]

	// 检测鼠标是否在按钮范围内且点击
	if x >= btnX && x <= btnX+btnW && y >= btnY && y <= btnY+btnH {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			m.clicked = true
		}
	}

	return m.clicked
}

// Draw 渲染菜单界面
func (m *MenuScreen) Draw(screen *ebiten.Image) {
	m.DrawBackground(screen)
	m.DrawStartButton(screen)
}

// DrawBackground 绘制背景图片并使其铺满整个屏幕
func (m *MenuScreen) DrawBackground(screen *ebiten.Image) {
	bgImage, _, err := ebitenutil.NewImageFromFile("photos/beijing.png")
	if err != nil {
		panic(err)
	}

	// 获取背景图片和屏幕的尺寸
	bgWidth, bgHeight := bgImage.Bounds().Dx(), bgImage.Bounds().Dy()
	screenWidth, screenHeight := screen.Bounds().Dx(), screen.Bounds().Dy()

	// 计算缩放比例，使背景图片铺满整个屏幕
	scaleX := float64(screenWidth) / float64(bgWidth)
	scaleY := float64(screenHeight) / float64(bgHeight)

	// 设置缩放和平移操作，使背景图片完全覆盖屏幕
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scaleX, scaleY)

	// 绘制背景图片
	screen.DrawImage(bgImage, op)
}

// DrawStartButton 绘制开始按钮及其交互效果
func (m *MenuScreen) DrawStartButton(screen *ebiten.Image) {
	btnX, btnY, btnW, btnH := m.startButtonRect[0], m.startButtonRect[1], m.startButtonRect[2], m.startButtonRect[3]

	// 检查鼠标是否在按钮上
	x, y := ebiten.CursorPosition()
	mouseHover := x >= btnX && x <= btnX+btnW && y >= btnY && y <= btnY+btnH

	// 根据鼠标状态设置按钮颜色和大小
	btnColor := color.RGBA{R: 100, G: 200, B: 100, A: 255}
	scale := 1.0
	if mouseHover {
		btnColor = color.RGBA{R: 150, G: 255, B: 150, A: 255} // 鼠标悬停时改变颜色
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			scale = 0.95 // 点击时缩小
		} else {
			scale = 1.05 // 悬停时放大
		}
	}

	// 绘制按钮背景
	btnImage := ebiten.NewImage(btnW, btnH)
	btnImage.Fill(btnColor)

	// 应用缩放和位置变换
	btnOp := &ebiten.DrawImageOptions{}

	if scale != 1.0 {
		// 计算缩放后的按钮中心点
		centerX := float64(btnX) + float64(btnW)/2
		centerY := float64(btnY) + float64(btnH)/2

		// 先平移到中心点，再缩放，最后平移回正确位置
		btnOp.GeoM.Translate(-float64(btnW)/2, -float64(btnH)/2)
		btnOp.GeoM.Scale(scale, scale)
		btnOp.GeoM.Translate(centerX, centerY)
	} else {
		btnOp.GeoM.Translate(float64(btnX), float64(btnY))
	}
	screen.DrawImage(btnImage, btnOp)

	// 按钮文字居中显示
	text := "START_THE_GAME"
	const charWidth = 6   // 假设每个字符宽度为6像素
	const charHeight = 16 // 假设字体高度为16像素
	textWidth := len(text) * charWidth
	textHeight := charHeight

	// 计算文字的居中位置
	var textX, textY int
	if scale != 1.0 {
		// 缩放时，基于按钮中心点计算文字位置
		centerX := float64(btnX) + float64(btnW)/2
		centerY := float64(btnY) + float64(btnH)/2
		textX = int(centerX - float64(textWidth)/2)
		textY = int(centerY - float64(textHeight)/2)
	} else {
		// 无缩放时，直接基于按钮位置计算
		textX = btnX + (btnW-textWidth)/2
		textY = btnY + (btnH-textHeight)/2
	}

	ebitenutil.DebugPrintAt(screen, text, textX, textY)
}
