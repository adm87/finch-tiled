package tiled

import (
	"encoding/xml"
	"path"

	"github.com/adm87/finch-core/finch"
)

// ======================================================
// TMX Asset Manager
// ======================================================

func RegisterTMXAssetManager() {
	finch.RegisterAssetManager(&finch.AssetManager{
		Types:       []finch.AssetType{"tmx"},
		Allocator:   allocate_tmx,
		Deallocator: deallocate_tmx,
	})
}

func GetTMX(file finch.AssetFile) (*TMX, error) {
	asset, err := finch.GetAsset[*TMX](file)
	if err != nil {
		return nil, err
	}
	return asset, nil
}

func MustGetTMX(file finch.AssetFile) *TMX {
	return finch.MustGetAsset[*TMX](file)
}

func allocate_tmx(file finch.AssetFile, data []byte) (any, error) {
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
}

func deallocate_tmx(file finch.AssetFile, data any) error {
	return nil
}
