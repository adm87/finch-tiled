package tiled

import (
	"encoding/xml"
	"fmt"
	"path"

	"github.com/adm87/finch-core/finch"
	"github.com/hajimehoshi/ebiten/v2"
)

func RegisterTSXAssetManager() {
	finch.RegisterAssetManager(&finch.AssetManager{
		AssetTypes: []finch.AssetType{"tsx"},
		ProcessAssetFile: func(file finch.AssetFile, data []byte) (any, error) {
			var tsx TSX

			if err := xml.Unmarshal(data, &tsx); err != nil {
				return nil, err
			}

			// Resolve the relative path of the image within the TSX file to be absolute
			// based on the location of the TSX file itself.

			tsxDir := path.Dir(file.Path())

			resolvedPath := path.Join(tsxDir, tsx.Image.Source())
			resolvedPath = path.Clean(resolvedPath)

			tsx.Image.Attrs[SourceAttr] = AttrString(resolvedPath)

			return &tsx, nil
		},
		CleanupAssetFile: func(file finch.AssetFile, data any) error {
			// Nothing special needs to be done to clean up a TSX asset.
			return nil
		},
	})
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
