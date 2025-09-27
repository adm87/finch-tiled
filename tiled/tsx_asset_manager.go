package tiled

import (
	"encoding/xml"
	"path"

	"github.com/adm87/finch-core/finch"
)

// ======================================================
// TSX Asset Manager
// ======================================================

func RegisterTSXAssetManager() {
	finch.RegisterAssetManager(&finch.AssetManager{
		Types:       []finch.AssetType{"tsx"},
		Allocator:   allocate_tsx,
		Deallocator: deallocate_tsx,
	})
}

func GetTSX(file finch.AssetFile) (*TSX, error) {
	asset, err := finch.GetAsset[*TSX](file)
	if err != nil {
		return nil, err
	}
	return asset, nil
}

func MustGetTSX(file finch.AssetFile) *TSX {
	return finch.MustGetAsset[*TSX](file)
}

func allocate_tsx(file finch.AssetFile, data []byte) (any, error) {
	var tsx TSX

	if err := xml.Unmarshal(data, &tsx); err != nil {
		return nil, err
	}

	tsxDir := path.Dir(file.Path())

	resolvedPath := path.Join(tsxDir, tsx.Image.Source())
	resolvedPath = path.Clean(resolvedPath)

	tsx.Image.Attrs[SourceAttr] = AttrString(resolvedPath)

	return &tsx, nil
}

func deallocate_tsx(file finch.AssetFile, data any) error {
	return nil
}
