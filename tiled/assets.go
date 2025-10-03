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
	TXAssetType  = "tx"
)

func resolveSourcePath(basePath, source string) string {
	resolvedPath := path.Join(path.Dir(basePath), source)
	resolvedPath = path.Clean(resolvedPath)
	return resolvedPath
}

func RegisterTiledAssetImporters() {
	// TMX Asset Support
	finch.RegisterAssetImporter(&finch.AssetImporter{
		AssetTypes: []finch.AssetType{TMXAssetType},
		ProcessAssetFile: func(file finch.AssetFile, data []byte) (any, error) {
			var tmx TMX

			if err := xml.Unmarshal(data, &tmx); err != nil {
				return nil, err
			}

			for i := range tmx.Tilesets {
				if _, exists := tmx.Tilesets[i].Attrs[SourceAttr]; exists {
					tmx.Tilesets[i].Attrs[SourceAttr] = AttrString(resolveSourcePath(file.Path(), tmx.Tilesets[i].Source()))
				}
			}

			for i := range tmx.ObjectGroups {
				for j := range tmx.ObjectGroups[i].Objects {
					if _, exists := tmx.ObjectGroups[i].Objects[j].Attrs[TemplateAttr]; !exists {
						continue
					}
					tmx.ObjectGroups[i].Objects[j].Attrs[TemplateAttr] = AttrString(resolveSourcePath(file.Path(), tmx.ObjectGroups[i].Objects[j].Template()))
				}
			}

			return &tmx, nil
		},
	})
	// TSX Asset Support
	finch.RegisterAssetImporter(&finch.AssetImporter{
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
	// TX Asset Support
	finch.RegisterAssetImporter(&finch.AssetImporter{
		AssetTypes: []finch.AssetType{TXAssetType},
		ProcessAssetFile: func(file finch.AssetFile, data []byte) (any, error) {
			var tx TX

			if err := xml.Unmarshal(data, &tx); err != nil {
				return nil, err
			}

			if tx.Tileset != nil {
				if _, exists := tx.Tileset.Attrs[SourceAttr]; exists {
					tx.Tileset.Attrs[SourceAttr] = AttrString(resolveSourcePath(file.Path(), tx.Tileset.Source()))
				}
			}

			return &tx, nil
		},
	})
}

// GetTX retrieves a TX asset by its file reference.
func GetTX(file finch.AssetFile) (*TX, error) {
	asset, err := finch.GetAsset[*TX](file)
	if err != nil {
		return nil, err
	}
	return asset, nil
}

// GetTXTSX retrieves the TSX asset referenced by a TX asset.
func GetTXTSX(file finch.AssetFile) (*TSX, error) {
	tx, err := GetTX(file)
	if err != nil {
		return nil, err
	}
	if tx.Tileset == nil {
		return nil, fmt.Errorf("tx does not contain a tileset: %s", file.Path())
	}
	tsxFile := finch.AssetFile(tx.Tileset.Source())

	tsx, err := GetTSX(tsxFile)
	if err != nil {
		return nil, err
	}
	return tsx, nil
}

// GetTXImg retrieves the image associated with a TX asset.
//
// Images are retrieved from the TSX asset referenced by the TX.
func GetTXImg(file finch.AssetFile) (*ebiten.Image, error) {
	tsx, err := GetTXTSX(file)
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
		return nil, fmt.Errorf("could not retrieve tx image from asset file: %s", imgFile.Path())
	}

	return img, nil
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

// MustGetTX is like GetTX but panics if the asset cannot be found.
func MustGetTX(file finch.AssetFile) *TX {
	tx, err := GetTX(file)
	if err != nil {
		panic(err)
	}
	return tx
}

// MustGetTXTSX is like GetTXTSX but panics if the asset cannot be found.
func MustGetTXTSX(file finch.AssetFile) *TSX {
	tsx, err := GetTXTSX(file)
	if err != nil {
		panic(err)
	}
	return tsx
}

// MustGetTXImg is like GetTXImg but panics if the asset cannot be found.
func MustGetTXImg(file finch.AssetFile) *ebiten.Image {
	img, err := GetTXImg(file)
	if err != nil {
		panic(err)
	}
	return img
}

// MustGetTMX is like GetTMX but panics if the asset cannot be found.
func MustGetTMX(file finch.AssetFile) *TMX {
	tmx, err := GetTMX(file)
	if err != nil {
		panic(err)
	}
	return tmx
}

// MustGetTSX is like GetTSX but panics if the asset cannot be found.
func MustGetTSX(src string) *TSX {
	tsx, err := GetTSX(finch.AssetFile(src))
	if err != nil {
		panic(err)
	}
	return tsx
}

// MustGetTSXImg is like GetTSXImg but panics if the asset cannot be found.
func MustGetTSXImg(src string) *ebiten.Image {
	img, err := GetTSXImg(finch.AssetFile(src))
	if err != nil {
		panic(err)
	}
	return img
}
