package tiled

import (
	"image"
	"log/slog"

	"github.com/adm87/finch-core/finch"
	"github.com/adm87/finch-core/geom"
	"github.com/adm87/finch-resources/images"
	"github.com/adm87/finch-resources/resources"
	"github.com/hajimehoshi/ebiten/v2"
)

// Buffer attempts to retrieve the specified .tmx resource and, if found, creates a new image the size of the tilemap and draws the entire tilemap onto it.
// It's recommended to cache the returned image for performance instead of buffering it every frame.
func Buffer(ctx finch.Context, tmxHandle resources.ResourceHandle) *ebiten.Image {
	tmx, exists := GetTmx(tmxHandle)

	if !exists {
		ctx.Logger().Warn("tmx resource not found", slog.String("tmx", tmxHandle.Key()))
		return nil
	}

	return draw_region(ctx, tmx, geom.NewRect64(0, 0, float64(tmx.Width()*tmx.TileWidth()), float64(tmx.Height()*tmx.TileHeight())))
}

// BufferRegion attempts to retrieve the specified .tmx resource and, if found, creates a new image the size of the region and draws the specified region of the tilemap onto it.
// It's recommended to cache the returned image for performance instead of buffering it every frame.
func BufferRegion(ctx finch.Context, tmxHandle resources.ResourceHandle, region geom.Rect64) *ebiten.Image {
	tmx, exists := GetTmx(tmxHandle)

	if !exists {
		ctx.Logger().Warn("tmx resource not found", slog.String("tmx", tmxHandle.Key()))
		return nil
	}

	return draw_region(ctx, tmx, region)
}

// Draw attempts to retrieve the specified .tmx resource and, if found, draws the entire tilemap onto the provided image.
func Draw(ctx finch.Context, img *ebiten.Image, tmxHandle resources.ResourceHandle) {
	tmx, exists := GetTmx(tmxHandle)

	if !exists {
		ctx.Logger().Warn("tmx resource not found", slog.String("tmx", tmxHandle.Key()))
		return
	}

	draw_region_to(ctx, img, tmx, geom.NewRect64(0, 0, float64(tmx.Width()*tmx.TileWidth()), float64(tmx.Height()*tmx.TileHeight())))
}

func DrawRegion(ctx finch.Context, img *ebiten.Image, tmxHandle resources.ResourceHandle, region geom.Rect64) {
	tmx, exists := GetTmx(tmxHandle)

	if !exists {
		ctx.Logger().Warn("tmx resource not found", slog.String("tmx", tmxHandle.Key()))
		return
	}

	draw_region_to(ctx, img, tmx, region)
}

func draw_region(ctx finch.Context, data *TMX, region geom.Rect64) *ebiten.Image {
	img := ebiten.NewImage(int(region.Width()), int(region.Height()))
	draw_region_to(ctx, img, data, region)
	return img
}

func draw_region_to(ctx finch.Context, img *ebiten.Image, tmx *TMX, region geom.Rect64) {
	tileWidth := tmx.TileWidth()
	tileHeight := tmx.TileHeight()

	for _, layer := range tmx.Layers {
		if !layer.Visible() {
			continue
		}

		// TASK: Support infinite maps
		//       Check if tmx is infinite and retreive data from chunks intsecting with region
		//       Only data positions within from each chunk that is visible by the region should be returned

		tileIndices, err := DecodeData(layer.Data.Data, layer.Data.Encoding())
		if err != nil {
			ctx.Logger().Error("failed to decode layer data", slog.String("layer", layer.Name()), slog.String("error", err.Error()))
			return
		}

		for i, tileIndex := range tileIndices {
			if tileIndex == 0 {
				continue
			}

			x := (i % layer.Width()) * tileWidth
			y := (i / layer.Width()) * tileHeight

			if !is_tile_in_region(tmx, region, x, y, tileWidth, tileHeight) {
				continue
			}

			gid, hFlip, vFlip, dFlip, _ := DecodeTile(tileIndex)
			tileset, exists := tmx.FindTilesetByTileGID(gid)

			if !exists {
				ctx.Logger().Warn("no tileset found for tile GID", slog.Any("gid", gid), slog.String("layer", layer.Name()))
				continue
			}

			tsxKey := resources.KeyFromPath(tileset.Source())
			tsx, exists := GetTsx(resources.ResourceHandle(tsxKey))

			if !exists {
				ctx.Logger().Error("failed to find tileset for tilemap", slog.String("tileset", tsxKey))
				return
			}

			tsxImgKey := resources.KeyFromPath(tsx.Image.Source())
			tsxImg, exists := images.GetImage(resources.ResourceHandle(tsxImgKey))

			if !exists {
				ctx.Logger().Error("failed to find tileset image", slog.String("image", tsxImgKey), slog.String("tileset", tsx.Name()))
				return
			}

			tile := gid - tileset.FirstGID()
			minx, miny := region.Min()

			blit_tile(ctx, img, tsxImg, tmx, tsx, tile, hFlip, vFlip, dFlip, x, y, -minx, -miny)
		}
	}
}

func is_tile_in_region(tmx *TMX, region geom.Rect64, x, y, tw, th int) bool {
	minx, miny := region.Min()
	maxx, maxy := region.Max()

	if float64(x+tw) < minx || float64(x) > maxx {
		return false
	}
	if float64(y+th) < miny || float64(y) > maxy {
		return false
	}

	return true
}

func blit_tile(ctx finch.Context, img *ebiten.Image, tilesetImg *ebiten.Image, tmx *TMX, tsx *TSX, tile uint32, hFlip, vFlip, dFlip bool, x, y int, dx, dy float64) {
	tilesPerRow := tsx.Image.Width() / tsx.TileWidth()
	tileX := (int(tile) % tilesPerRow) * tsx.TileWidth()
	tileY := (int(tile) / tilesPerRow) * tsx.TileHeight()

	sub := tilesetImg.SubImage(image.Rect(tileX, tileY, tileX+tsx.TileWidth(), tileY+tsx.TileHeight())).(*ebiten.Image)

	op := &ebiten.DrawImageOptions{}

	// Tiled anchors tiles at the bottom-left of their cell
	// See: https://doc.mapeditor.org/en/stable/reference/tmx-map-format/
	op.GeoM.Translate(0, float64(-tsx.TileHeight()+tmx.TileHeight()))

	if hFlip {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(float64(tsx.TileWidth()), 0)
	}
	if vFlip {
		op.GeoM.Scale(1, -1)
		op.GeoM.Translate(0, float64(tsx.TileHeight()))
	}
	if dFlip {
		op.GeoM.Rotate(-3.14159265 / 2) // -90 degrees in radians
		op.GeoM.Scale(1, -1)
		op.GeoM.Translate(float64(tsx.TileHeight()), 0)
	}

	x = x + tsx.TileOffset.X() + int(dx)
	y = y + tsx.TileOffset.Y() + int(dy)

	op.GeoM.Translate(float64(x), float64(y))

	img.DrawImage(sub, op)
}
