package websocket

import (
	"fmt"
	"sync"
	"time"
)

const (
	CanvasWidth  = 32
	CanvasHeight = 32
	DefaultColor = "#0d1117"
)

type PixelCanvas struct {
	Width      int               `json:"width"`
	Height     int               `json:"height"`
	Pixels     map[string]string `json:"pixels"`
	LastUpdate time.Time         `json:"last_update"`
	mutex      sync.RWMutex
}

func NewPixelCanvas() *PixelCanvas {
	return &PixelCanvas{
		Width:      CanvasWidth,
		Height:     CanvasHeight,
		Pixels:     make(map[string]string),
		LastUpdate: time.Now(),
	}
}

func (c *PixelCanvas) UpdatePixel(x, y int, color string) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if x < 0 || x >= c.Width || y < 0 || y >= c.Height {
		return false
	}

	key := fmt.Sprintf("%d,%d", x, y)
	if color == DefaultColor {
		delete(c.Pixels, key)
	} else {
		c.Pixels[key] = color
	}
	c.LastUpdate = time.Now()
	return true
}

func (c *PixelCanvas) GetPixel(x, y int) string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if x < 0 || x >= c.Width || y < 0 || y >= c.Height {
		return DefaultColor
	}

	key := fmt.Sprintf("%d,%d", x, y)
	if color, exists := c.Pixels[key]; exists {
		return color
	}
	return DefaultColor
}

func (c *PixelCanvas) GetPixels() map[string]string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	pixels := make(map[string]string)
	for k, v := range c.Pixels {
		pixels[k] = v
	}
	return pixels
}

func (c *PixelCanvas) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.Pixels = make(map[string]string)
	c.LastUpdate = time.Now()
}
