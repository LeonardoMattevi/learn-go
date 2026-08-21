package camera

type Camera struct {
	X, Y          float64
	ScreenW, ScreenH int
}

// Follow centers the camera on (px, py), clamped to world bounds.
func (c *Camera) Follow(px, py float64, worldPixelW, worldPixelH int) {
	c.X = px - float64(c.ScreenW)/2
	c.Y = py - float64(c.ScreenH)/2
	if c.X < 0 {
		c.X = 0
	}
	if c.Y < 0 {
		c.Y = 0
	}
	if maxX := float64(worldPixelW - c.ScreenW); c.X > maxX {
		c.X = maxX
	}
	if maxY := float64(worldPixelH - c.ScreenH); c.Y > maxY {
		c.Y = maxY
	}
}

// WorldToScreen converts world-pixel coords to screen coords.
func (c Camera) WorldToScreen(wx, wy float64) (float64, float64) {
	return wx - c.X, wy - c.Y
}
