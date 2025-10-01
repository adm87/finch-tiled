package tiled

import (
	"encoding/xml"
	"fmt"
	"path"

	"github.com/adm87/finch-core/finch"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	TMXAssetType = "tmx"
	TSXAssetType = "tsx"
)

func resolveSourcePath(basePath, source string) string {
	resolvedPath := path.Join(path.Dir(basePath), source)
	resolvedPath = path.Clean(resolvedPath)
	return resolvedPath
}

func RegisterTiledAssetImporter() {
	// TMX Asset Support
	finch.RegisterAssetImporter(&finch.AssetManager{
		AssetTypes: []finch.AssetType{TMXAssetType},
		ProcessAssetFile: func(file finch.AssetFile, data []byte) (any, error) {
			var tmx TMX

			if err := xml.Unmarshal(data, &tmx); err != nil {
				return nil, err
			}

			for i := range tmx.Tilesets {
				tmx.Tilesets[i].Attrs[SourceAttr] = AttrString(resolveSourcePath(file.Path(), tmx.Tilesets[i].Source()))
			}

			return &tmx, nil
		},
	})
	// TSX Asset Support
	finch.RegisterAssetImporter(&finch.AssetManager{
		AssetTypes: []finch.AssetType{TSXAssetType},
		ProcessAssetFile: func(file finch.AssetFile, data []byte) (any, error) {
			var tsx TSX

			if err := xml.Unmarshal(data, &tsx); err != nil {
				return nil, err
			}

			tsx.Image.Attrs[SourceAttr] = AttrString(resolveSourcePath(file.Path(), tsx.Image.Source()))

			return &tsx, nil
		},
	})
}

// GetTMX retrieves a TMX asset by its file reference.
func GetTMX(file finch.AssetFile) (*TMX, error) {
	asset, err := finch.GetAsset[*TMX](file)
	if err != nil {
		return nil, err
	}
	return asset, nil
}

// GetTSX retrieves a TSX asset by its file reference.
func GetTSX(file finch.AssetFile) (*TSX, error) {
	asset, err := finch.GetAsset[*TSX](file)
	if err != nil {
		return nil, err
	}
	return asset, nil
}

// GetTSXImg retrieves the image associated with a TSX asset.
func GetTSXImg(file finch.AssetFile) (*ebiten.Image, error) {
	tsx, err := GetTSX(file)
	if err != nil {
		return nil, err
	}

	imgFile := finch.AssetFile(tsx.Image.Source())

	imgAsset, err := imgFile.Get()
	if err != nil {
		return nil, err
	}

	img, ok := imgAsset.(*ebiten.Image)
	if !ok {
		return nil, fmt.Errorf("could not retrieve tsx image from asset file: %s", imgFile.Path())
	}

	return img, nil
}

// MustGetTMX is like GetTMX but panics if the asset cannot be loaded.
func MustGetTMX(file finch.AssetFile) *TMX {
	tmx, err := GetTMX(file)
	if err != nil {
		panic(err)
	}
	return tmx
}

// MustGetTSX is like GetTSX but panics if the asset cannot be loaded.
func MustGetTSX(src string) *TSX {
	tsx, err := GetTSX(finch.AssetFile(src))
	if err != nil {
		panic(err)
	}
	return tsx
}

// MustGetTSXImg is like GetTSXImg but panics if the asset cannot be loaded.
func MustGetTSXImg(src string) *ebiten.Image {
	img, err := GetTSXImg(finch.AssetFile(src))
	if err != nil {
		panic(err)
	}
	return img
}
