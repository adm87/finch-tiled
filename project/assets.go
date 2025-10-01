package project

import (
	"encoding/json"

	"github.com/adm87/finch-core/finch"
)

func RegisterAssetImporter() {
	finch.RegisterAssetImporter(&finch.AssetManager{
		AssetTypes: []finch.AssetType{"tiled-project"},
		ProcessAssetFile: func(file finch.AssetFile, data []byte) (any, error) {
			project := &TiledProject{}
			if err := json.Unmarshal(data, project); err != nil {
				return nil, err
			}
			return project, nil
		},
	})
}
