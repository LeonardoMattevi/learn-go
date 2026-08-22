package sprite

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Direction indices matching entity.Direction constants (Down=0, Up=1, Left=2, Right=3).
const (
	DirDown  = 0
	DirUp    = 1
	DirLeft  = 2
	DirRight = 3
)

// BuildPlayer returns a directional animated sprite for the player (stickman).
// 4 directions × 4 walk frames, 12×16 px per frame.
func BuildPlayer() *Sprite {
	const (
		fw, fh       = 12, 16
		dirsCount    = 4
		framesPerDir = 4
	)
	frames := make([]*ebiten.Image, dirsCount*framesPerDir)
	skin := color.RGBA{R: 220, G: 180, B: 120, A: 255}
	robe := color.RGBA{R: 30, G: 100, B: 200, A: 255}
	staff := color.RGBA{R: 160, G: 110, B: 40, A: 255}

	for dir := 0; dir < dirsCount; dir++ {
		for f := 0; f < framesPerDir; f++ {
			img := ebiten.NewImage(fw, fh)

			// Head (circle approximated as filled rect 4×4 centered at x=6, y=2)
			drawCircle(img, 6, 3, 2.5, skin)

			// Body
			vector.DrawFilledRect(img, 4, 6, 4, 5, robe, false)

			// Staff (right side, static)
			vector.DrawFilledRect(img, 10, 1, 1, 10, staff, false)

			// Legs — animate based on frame
			legPhase := f % 2 // 0 or 1
			switch dir {
			case DirDown, DirUp:
				// legs split left/right
				if legPhase == 0 {
					vector.DrawFilledRect(img, 4, 11, 2, 5, robe, false) // left leg forward
					vector.DrawFilledRect(img, 7, 11, 2, 4, robe, false) // right leg back
				} else {
					vector.DrawFilledRect(img, 4, 11, 2, 4, robe, false)
					vector.DrawFilledRect(img, 7, 11, 2, 5, robe, false)
				}
			case DirLeft, DirRight:
				// legs stagger forward/back
				if legPhase == 0 {
					vector.DrawFilledRect(img, 4, 11, 2, 5, robe, false)
					vector.DrawFilledRect(img, 7, 11, 2, 3, robe, false)
				} else {
					vector.DrawFilledRect(img, 4, 11, 2, 3, robe, false)
					vector.DrawFilledRect(img, 7, 11, 2, 5, robe, false)
				}
			}

			// Direction indicator: a small dot on the face side
			switch dir {
			case DirDown:
				vector.DrawFilledRect(img, 5, 5, 2, 1, color.RGBA{0, 0, 0, 180}, false)
			case DirUp:
				// eyes on back of head — just hood color, no dot
			case DirLeft:
				vector.DrawFilledRect(img, 3, 3, 1, 1, color.RGBA{0, 0, 0, 180}, false)
			case DirRight:
				vector.DrawFilledRect(img, 8, 3, 1, 1, color.RGBA{0, 0, 0, 180}, false)
			}

			frames[dir*framesPerDir+f] = img
		}
	}
	return New(frames, dirsCount, framesPerDir)
}

// BuildWalker returns a 2-frame animated sprite for the Walker enemy (12×12 px).
func BuildWalker() *Sprite {
	const (
		fw, fh       = 12, 12
		dirsCount    = 1
		framesPerDir = 2
	)
	body := color.RGBA{R: 200, G: 40, B: 40, A: 255}
	eye := color.RGBA{R: 255, G: 220, B: 0, A: 255}
	frames := make([]*ebiten.Image, framesPerDir)

	for f := 0; f < framesPerDir; f++ {
		img := ebiten.NewImage(fw, fh)
		// Body
		vector.DrawFilledRect(img, 1, 1, 10, 10, body, false)
		// Eye (moves slightly between frames to show life)
		ex := float32(4 + f)
		vector.DrawFilledRect(img, ex, 3, 3, 3, eye, false)
		// Feet (alternate)
		if f == 0 {
			vector.DrawFilledRect(img, 2, 10, 3, 2, color.RGBA{140, 20, 20, 255}, false)
			vector.DrawFilledRect(img, 7, 10, 3, 1, color.RGBA{140, 20, 20, 255}, false)
		} else {
			vector.DrawFilledRect(img, 2, 10, 3, 1, color.RGBA{140, 20, 20, 255}, false)
			vector.DrawFilledRect(img, 7, 10, 3, 2, color.RGBA{140, 20, 20, 255}, false)
		}
		frames[f] = img
	}
	return New(frames, dirsCount, framesPerDir)
}

// BuildShooter returns a 2-frame sprite for the Shooter enemy (12×12 px).
func BuildShooter() *Sprite {
	const (
		fw, fh       = 12, 12
		dirsCount    = 1
		framesPerDir = 2
	)
	body := color.RGBA{R: 180, G: 90, B: 10, A: 255}
	eye := color.RGBA{R: 255, G: 80, B: 0, A: 255}
	frames := make([]*ebiten.Image, framesPerDir)

	for f := 0; f < framesPerDir; f++ {
		img := ebiten.NewImage(fw, fh)
		// Diamond/rhombus body
		drawDiamond(img, 6, 6, 5, body)
		// Pulsing eye in center
		eyeSize := float32(2 + f)
		vector.DrawFilledRect(img, float32(6)-eyeSize/2, float32(6)-eyeSize/2, eyeSize, eyeSize, eye, false)
		frames[f] = img
	}
	return New(frames, dirsCount, framesPerDir)
}

// BuildTileFloor returns a 1×1 tile image for floor (16×16 px).
func BuildTileFloor() *ebiten.Image {
	img := ebiten.NewImage(16, 16)
	img.Fill(color.RGBA{R: 40, G: 35, B: 30, A: 255})
	// Subtle grid lines
	lineCol := color.RGBA{R: 50, G: 45, B: 38, A: 255}
	vector.DrawFilledRect(img, 0, 0, 16, 1, lineCol, false)
	vector.DrawFilledRect(img, 0, 0, 1, 16, lineCol, false)
	return img
}

// BuildTileWall returns a 1×1 tile image for wall (16×16 px).
func BuildTileWall() *ebiten.Image {
	img := ebiten.NewImage(16, 16)
	img.Fill(color.RGBA{R: 70, G: 65, B: 80, A: 255})
	// Brick pattern: two horizontal mortar lines
	mortar := color.RGBA{R: 50, G: 45, B: 55, A: 255}
	vector.DrawFilledRect(img, 0, 7, 16, 2, mortar, false)
	// Vertical mortar offset per row
	vector.DrawFilledRect(img, 7, 0, 2, 7, mortar, false)
	vector.DrawFilledRect(img, 3, 9, 2, 7, mortar, false)
	vector.DrawFilledRect(img, 11, 9, 2, 7, mortar, false)
	return img
}

// BuildProjectilePlayer returns a small glowing orb image (6×6 px).
func BuildProjectilePlayer() *ebiten.Image {
	img := ebiten.NewImage(6, 6)
	vector.DrawFilledRect(img, 1, 1, 4, 4, color.RGBA{R: 100, G: 200, B: 255, A: 255}, false)
	vector.DrawFilledRect(img, 2, 2, 2, 2, color.RGBA{R: 220, G: 240, B: 255, A: 255}, false)
	return img
}

// BuildProjectileEnemy returns a small red orb image (6×6 px).
func BuildProjectileEnemy() *ebiten.Image {
	img := ebiten.NewImage(6, 6)
	vector.DrawFilledRect(img, 1, 1, 4, 4, color.RGBA{R: 255, G: 80, B: 30, A: 255}, false)
	vector.DrawFilledRect(img, 2, 2, 2, 2, color.RGBA{R: 255, G: 200, B: 180, A: 255}, false)
	return img
}

// drawCircle fills a circle approximation using scanlines.
func drawCircle(img *ebiten.Image, cx, cy, r float32, col color.RGBA) {
	for y := cy - r; y <= cy+r; y++ {
		dx := float32(math.Sqrt(float64(r*r - (y-cy)*(y-cy))))
		vector.DrawFilledRect(img, cx-dx, y, dx*2, 1, col, false)
	}
}

// drawDiamond fills a rotated square (diamond) shape.
func drawDiamond(img *ebiten.Image, cx, cy, r float32, col color.RGBA) {
	for y := cy - r; y <= cy+r; y++ {
		dx := r - float32(math.Abs(float64(y-cy)))
		if dx > 0 {
			vector.DrawFilledRect(img, cx-dx, y, dx*2, 1, col, false)
		}
	}
}
