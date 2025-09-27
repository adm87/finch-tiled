package tiled

import (
	"encoding/xml"
	"path"

	"github.com/adm87/finch-core/finch"
	"github.com/adm87/finch-core/images"
	"github.com/hajimehoshi/ebiten/v2"
)

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
	img, err := images.Get(finch.AssetFile(tsx.Image.Source()))
	if err != nil {
		return nil, err
	}
	return img, nil
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

// ======================================================
// TSX Asset Manager
// ======================================================

func RegisterTSXAssetManager() {
	finch.RegisterAssetManager(&finch.AssetManager{
		Types: []finch.AssetType{"tsx"},
		Allocator: func(file finch.AssetFile, data []byte) (any, error) {
			var tsx TSX

			if err := xml.Unmarshal(data, &tsx); err != nil {
				return nil, err
			}

			tsxDir := path.Dir(file.Path())

			resolvedPath := path.Join(tsxDir, tsx.Image.Source())
			resolvedPath = path.Clean(resolvedPath)

			tsx.Image.Attrs[SourceAttr] = AttrString(resolvedPath)

			return &tsx, nil
		},
		Deallocator: func(file finch.AssetFile, data any) error {
			return nil
		},
	})
}
