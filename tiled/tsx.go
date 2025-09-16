package tiled

import (
	"encoding/xml"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
	"sync"

	"github.com/adm87/finch-core/finch"
	"github.com/adm87/finch-core/types"
	"github.com/adm87/finch-resources/resources"
)

const (
	TSXMetadataImageRefs = "img_refs"
)

// ======================================================
// TSX File Structure
// ======================================================

type TSX struct {
	Version      string   `xml:"version,attr"`
	TiledVersion string   `xml:"tiledversion,attr"`
	Name         string   `xml:"name,attr"`
	TileWidth    int      `xml:"tilewidth,attr"`
	TileHeight   int      `xml:"tileheight,attr"`
	Spacing      int      `xml:"spacing,attr"`
	TileCount    int      `xml:"tilecount,attr"`
	Columns      int      `xml:"columns,attr"`
	Image        TSXImage `xml:"image"`
}

type TSXImage struct {
	Source string `xml:"source,attr"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`
}

func (tsx *TSX) IsValid() error {
	if tsx.Image.Source == "" {
		return errors.New("tsx image source is required")
	}
	// TASK: Finish filling out validation checks.
	return nil
}

// ======================================================
// TSX Resource System
// ======================================================

var tsxSystem = resources.NewResourceSystemKey[*TsxResourceSystem]()

type TsxResourceSystem struct {
	tilesets map[string]*TSX
	loading  types.HashSet[string]
	mu       sync.Mutex
}

func NewTsxResourceSystem() *TsxResourceSystem {
	return &TsxResourceSystem{
		tilesets: make(map[string]*TSX),
		loading:  make(types.HashSet[string]),
		mu:       sync.Mutex{},
	}
}

func (rs *TsxResourceSystem) ResourceTypes() []string {
	return []string{"tsx"}
}

func (rs *TsxResourceSystem) Type() resources.ResourceSystemType {
	return tsxSystem
}

func (rs *TsxResourceSystem) IsLoaded(key string) bool {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	_, exists := rs.tilesets[key]
	return exists
}

func (rs *TsxResourceSystem) Load(ctx finch.Context, key string, metadata *resources.Metadata) error {
	_, err := rs.load_tsx(ctx, key, metadata)
	if err != nil {
		return err
	}

	ctx.Logger().Info("resource loaded", slog.String("key", key))
	return nil
}

func (rs *TsxResourceSystem) Unload(ctx finch.Context, key string) error {
	return errors.New("not implemented")
}

func (rs *TsxResourceSystem) GenerateMetadata(ctx finch.Context, key string, metadata *resources.Metadata) error {
	tsx, err := rs.load_tsx(ctx, key, metadata)
	if err != nil {
		return err
	}

	var imgRefs []string
	if tsx.Image.Source != "" {
		b := filepath.Base(tsx.Image.Source)
		imgRefs = append(imgRefs, strings.TrimSuffix(b, filepath.Ext(b)))
	}

	metadata.Extras = map[string]any{
		TSXMetadataImageRefs: imgRefs,
	}

	return nil
}

func (rs *TsxResourceSystem) GetDependencies(ctx finch.Context, key string, metadata *resources.Metadata) (imgRefs []string) {
	if metadata.Extras == nil {
		return nil
	}

	raw, exists := metadata.Extras[TSXMetadataImageRefs]
	if !exists {
		return nil
	}

	s, ok := raw.([]any)
	if !ok {
		return nil
	}

	imgRefs = make([]string, 0, len(s))
	for _, v := range s {
		str, ok := v.(string)
		if ok {
			imgRefs = append(imgRefs, str)
		}
	}

	return
}

func (rs *TsxResourceSystem) load_tsx(ctx finch.Context, key string, metadata *resources.Metadata) (*TSX, error) {
	if err := rs.try_load(key); err != nil {
		return nil, err
	}

	defer func() {
		rs.mu.Lock()
		rs.loading.Remove(key)
		rs.mu.Unlock()
	}()

	data, err := resources.LoadData(ctx, key, metadata)
	if err != nil {
		return nil, err
	}

	var tsx TSX
	if err := xml.Unmarshal(data, &tsx); err != nil {
		return nil, err
	}

	if err := tsx.IsValid(); err != nil {
		return nil, err
	}

	rs.mu.Lock()
	rs.tilesets[key] = &tsx
	rs.mu.Unlock()
	return &tsx, nil
}

func (rs *TsxResourceSystem) try_load(key string) error {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	if _, exists := rs.tilesets[key]; exists {
		return fmt.Errorf("tsx resource is already loaded: %s", key)
	}
	if rs.loading.Contains(key) {
		return fmt.Errorf("tsx resource is already loading: %s", key)
	}

	rs.loading.Add(key)
	return nil
}
