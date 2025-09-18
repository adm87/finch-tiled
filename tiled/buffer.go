package tiled

import "github.com/hajimehoshi/ebiten/v2"

func Buffer() *ebiten.Image {
	return BufferRegion()
}

func BufferVar(img *ebiten.Image) {
	BufferRegionVar(img)
}

func BufferRegion() *ebiten.Image {
	return nil
}

func BufferRegionVar(img *ebiten.Image) {
	_ = img
}
