package tiled

import (
	"encoding/xml"
	"path"

	"github.com/adm87/finch-core/finch"
)

// GetTMX retrieves a TMX asset by its file reference.
func GetTMX(file finch.AssetFile) (*TMX, error) {
	asset, err := finch.GetAsset[*TMX](file)
	if err != nil {
		return nil, err
	}
	return asset, nil
}

// MustGetTMX is like GetTMX but panics if the asset cannot be loaded.
func MustGetTMX(file finch.AssetFile) *TMX {
	tmx, err := GetTMX(file)
	if err != nil {
		panic(err)
	}
	return tmx
}

// ======================================================
// TMX Asset Manager
// ======================================================

func RegisterTMXAssetManager() {
	finch.RegisterAssetManager(&finch.AssetManager{
		Types: []finch.AssetType{"tmx"},
		Allocator: func(file finch.AssetFile, data []byte) (any, error) {
			var tmx TMX

			if err := xml.Unmarshal(data, &tmx); err != nil {
				return nil, err
			}

			for i := range tmx.Tilesets {
				tmxDir := path.Dir(file.Path())

				resolvedPath := path.Join(tmxDir, tmx.Tilesets[i].Source())
				resolvedPath = path.Clean(resolvedPath)

				tmx.Tilesets[i].Attrs[SourceAttr] = AttrString(resolvedPath)
			}

			return &tmx, nil
		},
		Deallocator: func(file finch.AssetFile, data any) error {
			return nil
		},
	})
}
