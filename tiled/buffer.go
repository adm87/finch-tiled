package tiled

import (
	"log/slog"

	"github.com/adm87/finch-core/finch"
	"github.com/adm87/finch-core/geom"
	"github.com/hajimehoshi/ebiten/v2"
)

func Buffer(ctx finch.Context, tmxID string) *ebiten.Image {
	tmx, exists := GetTmx(tmxID)

	if !exists {
		ctx.Logger().Warn("cannot buffer tilemap, tmx not found:", slog.String("tmxID", tmxID))
		return nil
	}

	width := tmx.Width * tmx.TileWidth
	height := tmx.Height * tmx.TileHeight

	return buffer_tilemap(ctx, tmx, geom.NewRect64(0, 0, float64(width), float64(height)))
}

func BufferRegion(ctx finch.Context, tmxID string, region geom.Rect64) *ebiten.Image {
	tmx, exists := GetTmx(tmxID)

	if !exists {
		ctx.Logger().Warn("cannot buffer tilemap region, tmx not found:", slog.String("tmxID", tmxID))
		return nil
	}

	return buffer_tilemap(ctx, tmx, region)
}

func buffer_tilemap(ctx finch.Context, tmx *TMX, region geom.Rect64) *ebiten.Image {
	if region.Width() <= 0 || region.Height() <= 0 {
		ctx.Logger().Warn("cannot buffer tilemap, region has non-positive width or height", slog.Float64("width", region.Width()), slog.Float64("height", region.Height()))
		return nil
	}

	return nil
}
