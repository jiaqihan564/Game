package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
)

// PlayScreen 游戏运行界面
type PlayScreen struct {
	playerX, playerY  float64       // 玩家位置坐标
	gridSize          int           // 网格单元格大小
	mainChar          *ebiten.Image // 玩家角色图像
	inventoryLoaded   bool          // 背包是否需要加载
	selectedItemIndex int           // 当前选中的背包物品索引
	backgroundImage   *ebiten.Image // 背景图片（新增字段）
	gridData          [][]int       // 网格数据数组（新增字段）
	inventorySize     int           // 背包大小
	items             []*item       // 背包物品数组
}

type item struct {
	id    int           // 物品ID
	name  string        // 物品名称
	count int           // 数量
	image *ebiten.Image // 物品图片
}

func NewPlayScreen() *PlayScreen {
	p := &PlayScreen{
		gridSize: 32,
		playerX:  0,
		playerY:  0,
	}

	// 初始化网格数据数组
	screenW, screenH := 800, 600 // 默认窗口大小
	gridWidth := screenW / p.gridSize
	gridHeight := screenH / p.gridSize

	// 创建二维数组并初始化为默认值（例如0表示空地）
	p.initGridData(gridHeight, gridWidth)

	// 预加载主角图片
	var err error
	p.mainChar, _, err = ebitenutil.NewImageFromFile("photos/zhu.png")
	if err != nil {
		panic(err)
	}

	// 预加载背景图片
	p.backgroundImage, _, err = ebitenutil.NewImageFromFile("photos/playBeijing.png")
	if err != nil {
		panic(err)
	}

	// 初始化背包物品
	p.initBag()

	return p
}

func (p *PlayScreen) Update() {
	// 事件监听
	if ebiten.IsKeyPressed(ebiten.KeyJ) {
		p.MovePlayer(-1, 0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyL) {
		p.MovePlayer(1, 0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyI) {
		p.MovePlayer(0, -1)
	}
	if ebiten.IsKeyPressed(ebiten.KeyK) {
		p.MovePlayer(0, 1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		// 打开或关闭背包
		p.inventoryLoaded = !p.inventoryLoaded
	}

	// 背包物品选择逻辑
	if p.inventoryLoaded {
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
			p.selectedItemIndex--
			if p.selectedItemIndex < 0 {
				p.selectedItemIndex = len(p.items) - 1
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
			p.selectedItemIndex++
			if p.selectedItemIndex >= len(p.items) {
				p.selectedItemIndex = 0
			}
		}
	}
}

// MovePlayer 处理角色移动逻辑
func (p *PlayScreen) MovePlayer(dx, dy float64) {
	// 基于实际帧时间计算速度（像素/秒）
	speed := 96.0                              // 每秒移动96像素
	delta := 1.0 / float64(ebiten.ActualTPS()) // 获取真实帧时间

	// 更新位置
	p.playerX += dx * speed * delta
	p.playerY += dy * speed * delta

	// 限制角色在屏幕范围内
	screenW, screenH := ebiten.WindowSize()
	if p.playerX < 0 {
		p.playerX = 0
	}
	if p.playerY < 0 {
		p.playerY = 0
	}
	if p.playerX > float64(screenW-32) {
		p.playerX = float64(screenW - 32)
	}
	if p.playerY > float64(screenH-32) {
		p.playerY = float64(screenH - 32)
	}
}

// Draw 绘制网格
func (p *PlayScreen) Draw(screen *ebiten.Image) {
	// 绘制背景
	p.DrawBackground(screen)
	p.DrawGrid(screen)
	p.DrawPlayer(screen)

	// 判断是否需要渲染背包
	if p.inventoryLoaded {
		p.DrawInventory(screen)
	}
}

// DrawPlayer 绘制玩家
func (p *PlayScreen) DrawPlayer(screen *ebiten.Image) {
	// 绘制角色（使用浮点坐标）
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(
		32.0/float64(p.mainChar.Bounds().Dx()),
		32.0/float64(p.mainChar.Bounds().Dy()),
	)
	op.GeoM.Translate(p.playerX, p.playerY) // 直接使用浮点数
	screen.DrawImage(p.mainChar, op)
}

// DrawGrid 绘制网格
func (p *PlayScreen) DrawGrid(screen *ebiten.Image) {
	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()

	// 使用 vector 绘制垂直线
	for x := 0; x <= w; x += p.gridSize {
		vector.DrawFilledRect(screen, float32(x), 0, 1, float32(h), color.White, false)
	}
	// 使用 vector 绘制水平线
	for y := 0; y <= h; y += p.gridSize {
		vector.DrawFilledRect(screen, 0, float32(y), float32(w), 1, color.White, false)
	}
}

// DrawInventory 绘制背包
func (p *PlayScreen) DrawInventory(screen *ebiten.Image) {
	// 绘制背包背景
	inventoryWidth := 300
	inventoryHeight := 400
	inventoryX := (screen.Bounds().Dx() - inventoryWidth) / 2
	inventoryY := (screen.Bounds().Dy() - inventoryHeight) / 2

	// 使用 vector 绘制半透明背景（带圆角效果的模拟）
	vector.DrawFilledRect(screen, float32(inventoryX), float32(inventoryY), float32(inventoryWidth), float32(inventoryHeight), color.RGBA{0, 0, 0, 200}, false)

	// 绘制边框（双层边框增强视觉效果）
	vector.StrokeRect(screen, float32(inventoryX), float32(inventoryY), float32(inventoryWidth), float32(inventoryHeight), 2, color.RGBA{100, 100, 150, 255}, false)
	vector.StrokeRect(screen, float32(inventoryX-1), float32(inventoryY-1), float32(inventoryWidth+2), float32(inventoryHeight+2), 1, color.RGBA{200, 200, 255, 255}, false)

	// 绘制标题
	titleX := inventoryX + (inventoryWidth-len("Inventory")*6)/2 // 居中标题
	titleY := inventoryY + 15
	ebitenutil.DebugPrintAt(screen, "Inventory", titleX, titleY)

	// 绘制分隔线
	vector.DrawFilledRect(screen, float32(inventoryX+10), float32(titleY+20), float32(inventoryWidth-20), 1, color.RGBA{100, 100, 150, 255}, false)

	// 绘制物品列表
	itemStartX := inventoryX + 20
	itemStartY := titleY + 40
	itemWidth := 48
	itemHeight := 48
	itemSpacing := 15

	// 计算每行最多显示的物品数量
	itemsPerRow := (inventoryWidth - 40) / (itemWidth + itemSpacing)

	for i, item := range p.items {
		if i >= p.inventorySize { // 限制显示数量不超过背包容量
			break
		}

		// 计算物品位置
		row := i / itemsPerRow
		col := i % itemsPerRow
		x := itemStartX + col*(itemWidth+itemSpacing)
		y := itemStartY + row*(itemHeight+itemSpacing)

		// 确保物品不会超出背包边界
		if x+itemWidth > inventoryX+inventoryWidth-20 {
			continue
		}

		// 绘制物品背景阴影
		vector.DrawFilledRect(screen, float32(x+2), float32(y+2), float32(itemWidth), float32(itemHeight), color.RGBA{0, 0, 0, 100}, false)

		// 如果是选中的物品，绘制高亮背景
		if i == p.selectedItemIndex {
			// 绘制选中效果的发光边框
			for j := 0; j < 3; j++ {
				vector.StrokeRect(screen, float32(x-j), float32(y-j), float32(itemWidth+2*j), float32(itemHeight+2*j), 1, color.RGBA{100, 150, 255, uint8(150 - j*50)}, false)
			}
		}

		// 绘制物品背景框
		vector.DrawFilledRect(screen, float32(x), float32(y), float32(itemWidth), float32(itemHeight), color.RGBA{50, 50, 70, 200}, false)
		vector.StrokeRect(screen, float32(x), float32(y), float32(itemWidth), float32(itemHeight), 1, color.RGBA{150, 150, 200, 255}, false)

		// 绘制物品图片
		if item.image != nil {
			op := &ebiten.DrawImageOptions{}
			// 缩放图片以适应物品框，保持一定边距
			imgBounds := item.image.Bounds()
			scaleX := float64(itemWidth-8) / float64(imgBounds.Dx())
			scaleY := float64(itemHeight-8) / float64(imgBounds.Dy())
			finalScale := math.Min(scaleX, scaleY) // 保持图片比例
			op.GeoM.Scale(finalScale, finalScale)
			op.GeoM.Translate(float64(x)+(float64(itemWidth)-float64(imgBounds.Dx())*finalScale)/2,
				float64(y)+(float64(itemHeight)-float64(imgBounds.Dy())*finalScale)/2)
			screen.DrawImage(item.image, op)
		}

		// 绘制物品数量（右下角）
		if item.count > 1 {
			countText := fmt.Sprintf("%d", item.count)
			textWidth := len(countText) * 6
			// 绘制数量背景
			vector.DrawFilledRect(screen, float32(x+itemWidth-textWidth-4), float32(y+itemHeight-16), float32(textWidth+4), 16, color.RGBA{0, 0, 0, 200}, false)
			ebitenutil.DebugPrintAt(screen, countText, x+itemWidth-textWidth-2, y+itemHeight-14)
		}

		// 绘制物品名称（鼠标悬停时显示）
		cursorX, cursorY := ebiten.CursorPosition()
		if cursorX >= x && cursorX <= x+itemWidth && cursorY >= y && cursorY <= y+itemHeight {
			// 绘制名称背景
			nameWidth := len(item.name)*6 + 4
			nameX := x + (itemWidth-nameWidth)/2
			vector.DrawFilledRect(screen, float32(nameX), float32(y-20), float32(nameWidth), 16, color.RGBA{0, 0, 0, 200}, false)
			ebitenutil.DebugPrintAt(screen, item.name, nameX+2, y-18)
		}
	}

	// 绘制关闭按钮
	closeBtnSize := 20
	closeBtnX := inventoryX + inventoryWidth - closeBtnSize - 10
	closeBtnY := inventoryY + 10

	// 关闭按钮背景
	vector.DrawFilledRect(screen, float32(closeBtnX), float32(closeBtnY), float32(closeBtnSize), float32(closeBtnSize), color.RGBA{200, 50, 50, 200}, false)
	vector.StrokeRect(screen, float32(closeBtnX), float32(closeBtnY), float32(closeBtnSize), float32(closeBtnSize), 1, color.RGBA{255, 100, 100, 255}, false)

	// 关闭按钮 "X"
	ebitenutil.DebugPrintAt(screen, "X", closeBtnX+7, closeBtnY+2)

	// 检查关闭按钮点击
	if p.inventoryLoaded {
		cursorX, cursorY := ebiten.CursorPosition()
		if cursorX >= closeBtnX && cursorX <= closeBtnX+closeBtnSize &&
			cursorY >= closeBtnY && cursorY <= closeBtnY+closeBtnSize &&
			inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			p.inventoryLoaded = false
		}
	}
}

// DrawBackground 绘制背景图片并使其铺满整个屏幕
func (p *PlayScreen) DrawBackground(screen *ebiten.Image) {
	// 直接使用已加载的背景图片，无需重复加载
	bgImage := p.backgroundImage

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

// 初始化地图资源
func (p *PlayScreen) initGridData(height int, width int) {
	// 创建二维数组并初始化为默认值（例如0表示空地）
	p.gridData = make([][]int, height)
	for i := range p.gridData {
		p.gridData[i] = make([]int, width)
		for j := range p.gridData[i] {
			p.gridData[i][j] = 0 // 初始化为0，表示空地
		}
	}
}

// 初始化背包物品 金币 50000
func (p *PlayScreen) initBag() {
	p.inventorySize = 5
	p.items = append(p.items, &item{
		count: 50000,
		id:    0,
		image: itemImages["1001"],
		name:  "金币",
	})
}
