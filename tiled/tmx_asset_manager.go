package tiled

import (
	"encoding/xml"
	"path"

	"github.com/adm87/finch-core/finch"
)

func RegisterTMXAssetManager() {
	finch.RegisterAssetManager(&finch.AssetManager{
		AssetTypes: []finch.AssetType{"tmx"},
		ProcessAssetFile: func(file finch.AssetFile, data []byte) (any, error) {
			var tmx TMX

			if err := xml.Unmarshal(data, &tmx); err != nil {
				return nil, err
			}

			// Resolve the relative paths within the TMX file to be absolute
			// based on the location of the TMX file itself.
			for i := range tmx.Tilesets {
				tmxDir := path.Dir(file.Path())

				resolvedPath := path.Join(tmxDir, tmx.Tilesets[i].Source())
				resolvedPath = path.Clean(resolvedPath)

				tmx.Tilesets[i].Attrs[SourceAttr] = AttrString(resolvedPath)
			}

			return &tmx, nil
		},
		CleanupAssetFile: func(file finch.AssetFile, data any) error {
			// Nothing special needs to be done to clean up a TMX asset.
			return nil
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

// MustGetTMX is like GetTMX but panics if the asset cannot be loaded.
func MustGetTMX(file finch.AssetFile) *TMX {
	tmx, err := GetTMX(file)
	if err != nil {
		panic(err)
	}
	return tmx
}
