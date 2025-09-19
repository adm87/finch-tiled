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

// Buffer attempts to retrieve the specified .tmx resource and, if found, creates a new image and draws the entire tilemap onto it.
// It's recommended to cache the returned image for performance instead of buffering it every frame.
//
// Note: This function does not support infinite tilemaps. Use BufferRegion() instead.
func Buffer(ctx finch.Context, tmxHandle resources.ResourceHandle) *ebiten.Image {
	tmx, exists := GetTmx(tmxHandle)

	if !exists {
		ctx.Logger().Warn("tmx resource not found", slog.String("tmx", tmxHandle.Key()))
		return nil
	}

	if !exists {
		ctx.Logger().Warn("tmx resource not found", slog.String("tmx", tmxHandle.Key()))
		return nil
	}

	if tmx.Infinite() {
		ctx.Logger().Error("cannot buffer infinite tilemap. Use BufferRegion() instead.")
		return nil
	}

	return buffer_region(ctx, tmx, geom.NewRect64(0, 0, float64(tmx.Width()*tmx.TileWidth()), float64(tmx.Height()*tmx.TileHeight())))
}

// BufferRegion attempts to retrieve the specified .tmx resource and, if found, creates a new image and draws the specified region of the tilemap onto it.
// It's recommended to cache the returned image for performance instead of buffering it every frame.
func BufferRegion(ctx finch.Context, tmxHandle resources.ResourceHandle, region geom.Rect64) *ebiten.Image {
	panic("not implemented")
}

// BufferVar draws the entire tilemap onto the provided image.
// If the image size does not match the tilemap size, the tilemap may be clipped or not fully drawn.
//
// Note: This function does not support infinite tilemaps. Use BufferRegionVar() instead.
func BufferVar(ctx finch.Context, img *ebiten.Image, tmxHandle resources.ResourceHandle) {
	tmx, exists := GetTmx(tmxHandle)

	if !exists {
		ctx.Logger().Warn("tmx resource not found", slog.String("tmx", tmxHandle.Key()))
		return
	}

	if !exists {
		ctx.Logger().Warn("tmx resource not found", slog.String("tmx", tmxHandle.Key()))
		return
	}

	if tmx.Infinite() {
		ctx.Logger().Error("cannot buffer infinite tilemap. Use BufferRegionVar() instead.")
		return
	}

	check_image_size(ctx, img, tmx.Width(), tmx.Height(), tmxHandle.Key())
	buffer_region_var(ctx, img, tmx, geom.NewRect64(0, 0, float64(tmx.Width()*tmx.TileWidth()), float64(tmx.Height()*tmx.TileHeight())))
}

// BufferRegionVar draws the specified region of the tilemap onto the provided image.
// If the image size does not match the tilemap size, the tilemap may be clipped or not fully drawn.
func BufferRegionVar(ctx finch.Context, img *ebiten.Image, tmxHandle resources.ResourceHandle, region geom.Rect64) {
	panic("not implemented")
}

func buffer_region(ctx finch.Context, data *TMX, region geom.Rect64) *ebiten.Image {
	img := ebiten.NewImage(int(region.Width()), int(region.Height()))
	buffer_region_var(ctx, img, data, region)
	return img
}

func buffer_region_var(ctx finch.Context, img *ebiten.Image, tmx *TMX, region geom.Rect64) {
	for _, layer := range tmx.Layers {
		if !layer.Visible() {
			continue
		}

		tileIndices, err := DecodeData(layer.Data.Data, layer.Data.Encoding())
		if err != nil {
			ctx.Logger().Error("failed to decode layer data", slog.String("layer", layer.Name()), slog.String("error", err.Error()))
			return
		}

		for i, tileIndex := range tileIndices {
			if tileIndex == 0 {
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
			blit_tile(ctx, img, tsxImg, tmx, tsx, tileset, tile, hFlip, vFlip, dFlip, region, i)
		}
	}
}

func blit_tile(ctx finch.Context, img *ebiten.Image, tsxImg *ebiten.Image, tmx *TMX, tsx *TSX, tileset *TMXTileset, tile uint32, hFlip, vFlip, dFlip bool, region geom.Rect64, pos int) {
	tilesPerRow := tsx.Image.Width() / tsx.TileWidth()
	tileX := (int(tile) % tilesPerRow) * tsx.TileWidth()
	tileY := (int(tile) / tilesPerRow) * tsx.TileHeight()

	sub := tsxImg.SubImage(image.Rect(tileX, tileY, tileX+tsx.TileWidth(), tileY+tsx.TileHeight())).(*ebiten.Image)

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

	x := ((pos % tmx.Width()) * tmx.TileWidth()) + tsx.TileOffset.X()
	y := ((pos / tmx.Width()) * tmx.TileHeight()) + tsx.TileOffset.Y()

	op.GeoM.Translate(float64(x), float64(y))

	img.DrawImage(sub, op)
}

func check_image_size(ctx finch.Context, img *ebiten.Image, expectedWidth, expectedHeight int, key string) {
	imgW := img.Bounds().Dx()
	imgH := img.Bounds().Dy()

	if imgW != expectedWidth || imgH != expectedHeight {
		ctx.Logger().Warn("tiled buffer image size mismatch and may result in the tilemap being clipped or wasted space:",
			slog.String("tmx", key),
			slog.Int("expectedWidth", expectedWidth),
			slog.Int("expectedHeight", expectedHeight),
			slog.Int("actualWidth", imgW),
			slog.Int("actualHeight", imgH))
	}
}
