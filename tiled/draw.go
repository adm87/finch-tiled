package tiled

import (
	"fmt"
	"image"
	"log/slog"
	"strconv"
	"strings"

	"github.com/adm87/finch-core/finch"
	"github.com/adm87/finch-core/fsys"
	"github.com/adm87/finch-core/geom"
	"github.com/hajimehoshi/ebiten/v2"
)

// TASK: Implement support for all encoding/compression types Tiled supports.
//       - Probably a good idea to support as many features of Tiled as possible - this goes beyond just encoding/compression.

// TASK: Implement support for object layers.
//       - This will be needed to define dynamic collision areas, spawn points, and other interactive elements in a game.

// TASK: Implement support for isometric and staggered maps.
//       - This early in development, it's really just a nice to have - but would be useful for certain types of games.

// TASK: Implement support for dynamically modifying tilemaps (e.g., changing tiles at runtime).
//       - Another nice to have, but could be useful for games that feature destructible environments or tile-based puzzles.

const (
	ErrWhileDrawingLayer = "tiled: error while drawing layer"
)

type DrawMode int

const (
	DrawModeNormal DrawMode = iota
	DrawModeRegional
	DrawModeScene
)

var identity = &ebiten.GeoM{}
var op = &ebiten.DrawImageOptions{}

// Draw attempts to render the entire TMX map onto the provided image.
// If the map is larger than the image, only the top-left portion will be drawn.
func Draw(ctx finch.Context, img *ebiten.Image, tmx *TMX) {
	region := geom.NewRect64(0, 0, float64(img.Bounds().Dx()), float64(img.Bounds().Dy()))
	for i := range tmx.Layers {
		if err := draw_map_layer(DrawModeNormal, img, tmx.Layers[i], tmx.Tilesets, &region, identity, tmx.TileWidth(), tmx.TileHeight(), tmx.IsInfinite()); err != nil {
			ctx.Logger().Error(ErrWhileDrawingLayer, slog.String("layer", tmx.Layers[i].Name()), slog.Any("error", err))
		}
	}
}

// DrawLayer attempts to render a specific layer of the TMX map onto the provided image.
// If the map is larger than the image, only the top-left portion will be drawn.
func DrawLayer(ctx finch.Context, img *ebiten.Image, tmx *TMX, layerName string) {
	layer, ok := tmx.GetLayerByName(layerName)
	if !ok {
		ctx.Logger().Warn("tiled: layer not found", slog.String("layer", layerName))
		return
	}
	region := geom.NewRect64(0, 0, float64(img.Bounds().Dx()), float64(img.Bounds().Dy()))
	if err := draw_map_layer(DrawModeNormal, img, layer, tmx.Tilesets, &region, identity, tmx.TileWidth(), tmx.TileHeight(), tmx.IsInfinite()); err != nil {
		ctx.Logger().Error(ErrWhileDrawingLayer, slog.String("layer", layer.Name()), slog.Any("error", err))
	}
}

// DrawRegion renders only the specified region of the TMX map onto the provided image.
func DrawRegion(ctx finch.Context, img *ebiten.Image, tmx *TMX, region geom.Rect64) {
	for i := range tmx.Layers {
		if err := draw_map_layer(DrawModeRegional, img, tmx.Layers[i], tmx.Tilesets, &region, identity, tmx.TileWidth(), tmx.TileHeight(), tmx.IsInfinite()); err != nil {
			ctx.Logger().Error(ErrWhileDrawingLayer, slog.String("layer", tmx.Layers[i].Name()), slog.Any("error", err))
		}
	}
}

// DrawLayerRegion renders only the specified region of a specific layer of the TMX map onto the provided image.
func DrawLayerRegion(ctx finch.Context, img *ebiten.Image, tmx *TMX, layerName string, region geom.Rect64) {
	layer, ok := tmx.GetLayerByName(layerName)
	if !ok {
		ctx.Logger().Warn("tiled: layer not found", slog.String("layer", layerName))
		return
	}
	if err := draw_map_layer(DrawModeRegional, img, layer, tmx.Tilesets, &region, identity, tmx.TileWidth(), tmx.TileHeight(), tmx.IsInfinite()); err != nil {
		ctx.Logger().Error(ErrWhileDrawingLayer, slog.String("layer", layer.Name()), slog.Any("error", err))
	}
}

// DrawScene renders the TMX map as seen through a camera, using the provided viewport and view matrix.
// This is typically used for rendering the map in a game scene where the camera can move and zoom.
func DrawScene(ctx finch.Context, img *ebiten.Image, tmx *TMX, viewport geom.Rect64, viewMatrix ebiten.GeoM) {
	for i := range tmx.Layers {
		if err := draw_map_layer(DrawModeScene, img, tmx.Layers[i], tmx.Tilesets, &viewport, &viewMatrix, tmx.TileWidth(), tmx.TileHeight(), tmx.IsInfinite()); err != nil {
			ctx.Logger().Error(ErrWhileDrawingLayer, slog.String("layer", tmx.Layers[i].Name()), slog.Any("error", err))
		}
	}
}

// DrawSceneLayer renders a specific layer of the TMX map as seen through a camera, using the provided viewport and view matrix.
// This is typically used for rendering the map in a game scene where the camera can move and zoom.
func DrawSceneLayer(ctx finch.Context, img *ebiten.Image, tmx *TMX, layerName string, viewport geom.Rect64, viewMatrix ebiten.GeoM) {
	layer, ok := tmx.GetLayerByName(layerName)
	if !ok {
		ctx.Logger().Warn("tiled: layer not found", slog.String("layer", layerName))
		return
	}
	if err := draw_map_layer(DrawModeScene, img, layer, tmx.Tilesets, &viewport, &viewMatrix, tmx.TileWidth(), tmx.TileHeight(), tmx.IsInfinite()); err != nil {
		ctx.Logger().Error(ErrWhileDrawingLayer, slog.String("layer", layer.Name()), slog.Any("error", err))
	}
}

func draw_map_layer(mode DrawMode, destImg *ebiten.Image, layer *TMXLayer, tilesets []*TMXTileset, region *geom.Rect64, view *ebiten.GeoM, cellWidth, cellHeight int, isInfinite bool) error {
	if !layer.IsVisible() || len(tilesets) == 0 {
		return nil
	}

	layerWidth := layer.Width() * cellWidth
	layerHeight := layer.Height() * cellHeight

	if err := process_tiles(layer, tilesets, region, layerWidth, layerHeight, cellWidth, cellHeight, isInfinite); err != nil {
		return err
	}

	tiles := collect_tiles(layer, region, cellWidth, cellHeight, isInfinite)

	for i := range tiles {
		op.GeoM.Reset()

		// The order of operations is important here.
		// See: https://doc.mapeditor.org/en/stable/reference/global-tile-ids/#tile-flipping
		if tiles[i].Flags&FLIP_DIAGONAL != 0 {
			op.GeoM.Rotate(fsys.HalfPi)
			op.GeoM.Scale(-1, 1)
			op.GeoM.Translate(float64(tiles[i].Height-tiles[i].Width), 0)
		}
		if tiles[i].Flags&FLIP_HORIZONTAL != 0 {
			op.GeoM.Scale(-1, 1)
			op.GeoM.Translate(float64(tiles[i].Width), 0)
		}
		if tiles[i].Flags&FLIP_VERTICAL != 0 {
			op.GeoM.Scale(1, -1)
			op.GeoM.Translate(0, float64(tiles[i].Height))
		}

		switch mode {
		case DrawModeNormal:
			op.GeoM.Translate(tiles[i].X, tiles[i].Y)
		case DrawModeRegional:
			minx, miny := region.Min()
			op.GeoM.Translate(tiles[i].X-minx, tiles[i].Y-miny)
		case DrawModeScene:
			op.GeoM.Translate(tiles[i].X, tiles[i].Y)
			op.GeoM.Concat(*view)
		default:
			panic("unhandled draw mode")
		}

		srcImg, err := GetTSXImg(finch.AssetFile(tiles[i].TsxSrc))
		if err != nil {
			return err
		}

		tilesPerRow := float64(srcImg.Bounds().Dx()) / tiles[i].Width
		tileX := (int(tiles[i].GID) % int(tilesPerRow)) * int(tiles[i].Width)
		tileY := (int(tiles[i].GID) / int(tilesPerRow)) * int(tiles[i].Height)

		destImg.DrawImage(srcImg.SubImage(image.Rect(tileX, tileY, tileX+int(tiles[i].Width), tileY+int(tiles[i].Height))).(*ebiten.Image), op)
	}

	return nil
}

func process_tiles(layer *TMXLayer, tilesets []*TMXTileset, region *geom.Rect64, layerWidth, layerHeight, cellWidth, cellHeight int, isInfinite bool) error {
	if isInfinite {
		return process_chunks(layer, tilesets, region, layerWidth, layerHeight, cellWidth, cellHeight)
	}

	// Already processed
	if layer.tiles != nil {
		return nil
	}

	tiles, err := decode_tiles(layer.Data.Data, tilesets, 0, 0, layerWidth, layerHeight, cellWidth, cellHeight)
	if err != nil {
		return err
	}

	layer.tiles = tiles
	return nil
}

func process_chunks(layer *TMXLayer, tilesets []*TMXTileset, region *geom.Rect64, layerWidth, layerHeight, cellWidth, cellHeight int) error {
	if layer.Data == nil || len(layer.Data.Chunks) == 0 {
		return nil
	}

	if layer.partitions == nil {
		layer.partitions = make(LayerPartitions)
	}

	minx, miny := region.Min()
	maxx, maxy := region.Max()

	for _, chunk := range layer.Data.Chunks {
		chunkX := float64(chunk.X() * cellWidth)
		chunkY := float64(chunk.Y() * cellHeight)
		chunkW := float64(chunk.Width() * cellWidth)
		chunkH := float64(chunk.Height() * cellHeight)

		cminx, cminy := chunkX, chunkY
		cmaxx, cmaxy := cminx+chunkW, cminy+chunkH

		if cmaxx < minx || cminx > maxx || cmaxy < miny || cminy > maxy {
			continue
		}

		chunkRect := geom.NewRect64(cminx, cminy, cmaxx-cminx, cmaxy-cminy)
		if _, exists := layer.partitions[chunkRect]; exists || !region.Intersects(&chunkRect) {
			continue
		}

		tiles, err := decode_tiles(chunk.Data, tilesets, int(chunkX), int(chunkY), int(chunkW), int(chunkH), cellWidth, cellHeight)
		if err != nil {
			return err
		}

		layer.partitions[chunkRect] = tiles
	}

	return nil
}

func decode_tiles(data string, tilesets []*TMXTileset, localStartX, localStartY, layerWidth, layerHeight, cellWidth, cellHeight int) ([]*Tile, error) {
	parsedData, err := parse_csv_data(data)
	if err != nil {
		return nil, err
	}

	var tiles []*Tile

	cellPerRow := layerWidth / cellWidth

	for i := range parsedData {
		gid := parsedData[i] & TILE_ID_MASK
		if gid == 0 {
			continue // Empty tile
		}

		var flags TiledFlags
		if (parsedData[i] & TILE_FLIP_HORIZONTAL) != 0 {
			flags |= FLIP_HORIZONTAL
		}
		if (parsedData[i] & TILE_FLIP_VERTICAL) != 0 {
			flags |= FLIP_VERTICAL
		}
		if (parsedData[i] & TILE_FLIP_DIAGONAL) != 0 {
			flags |= FLIP_DIAGONAL
			// According to Tiled docs, diagonal flip swaps horizontal and vertical flips
			// See: https://doc.mapeditor.org/en/stable/reference/global-tile-ids/#tile-flipping
			if flags&(FLIP_HORIZONTAL|FLIP_VERTICAL) != 0 {
				flags ^= FLIP_HORIZONTAL | FLIP_VERTICAL
			}
		}
		if (parsedData[i] & TILE_FLIP_ROTATED_HEX) != 0 {
			flags |= FLIP_ROTATED_HEX
		}

		var tileset *TMXTileset
		for j := len(tilesets) - 1; j >= 0; j-- {
			if gid >= tilesets[j].FirstGID() {
				tileset = tilesets[j]
				break
			}
		}

		if tileset == nil {
			return nil, fmt.Errorf("no tileset found for GID %d", gid)
		}

		tsx, err := GetTSX(finch.AssetFile(tileset.Source()))
		if err != nil {
			return nil, err
		}

		x := float64(localStartX + ((i % cellPerRow) * cellWidth))
		y := float64(localStartY + ((i / cellPerRow) * cellHeight))

		if tsx.TileOffset != nil {
			x += float64(tsx.TileOffset.X())
			y += float64(tsx.TileOffset.Y())
		}

		// Tiled anchors tiles at the bottom-left of their cell.
		// Adjust the Y position to offset the tile by the difference between the cell and tile's heights.
		// See: https://doc.mapeditor.org/en/stable/reference/tmx-map-format/
		y += float64(cellHeight) - float64(tsx.TileHeight())

		tiles = append(tiles, &Tile{
			Flags:  flags,
			GID:    gid - tileset.FirstGID(),
			TsxSrc: tileset.Source(),
			X:      x,
			Y:      y,
			Width:  float64(tsx.TileWidth()),
			Height: float64(tsx.TileHeight()),
		})
	}

	return tiles, nil
}

func parse_csv_data(dataStr string) ([]uint32, error) {
	var data []uint32
	for _, s := range strings.Split(dataStr, ",") {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		tileIndex, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("invalid CSV layer data: %w", err)
		}
		data = append(data, uint32(tileIndex))
	}
	return data, nil
}

func collect_tiles(layer *TMXLayer, region *geom.Rect64, cellWidth, cellHeight int, isInfinite bool) []*Tile {
	if layer.tiles == nil && layer.partitions == nil {
		return nil
	}

	tiles := layer.tiles
	if isInfinite {
		for chunkRect, chunkTiles := range layer.partitions {
			if region.Intersects(&chunkRect) {
				tiles = append(tiles, chunkTiles...)
			}
		}
	}

	var result []*Tile

	minx, miny := region.Min()
	maxx, maxy := region.Max()

	for i := range tiles {
		tminx := tiles[i].X
		tminy := tiles[i].Y
		tmaxx := tiles[i].X + float64(tiles[i].Width)
		tmaxy := tiles[i].Y + float64(tiles[i].Height)

		if tmaxx < minx || tminx > maxx || tmaxy < miny || tminy > maxy {
			continue
		}

		result = append(result, tiles[i])
	}

	return result
}
