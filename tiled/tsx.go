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

func GetTsx(handle resources.ResourceHandle) (*TSX, bool) {
	sys, ok := resources.GetSystem(tsxSystem).(*TsxResourceSystem)
	if !ok {
		return nil, false
	}

	sys.mu.Lock()
	defer sys.mu.Unlock()

	tsx, exists := sys.tilesets[handle.Key()]
	return tsx, exists
}

func (rs *TsxResourceSystem) ResourceTypes() []string {
	return []string{"tsx"}
}

func (rs *TsxResourceSystem) Type() resources.ResourceSystemType {
	return tsxSystem
}

func (rs *TsxResourceSystem) IsLoaded(handle resources.ResourceHandle) bool {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	_, exists := rs.tilesets[handle.Key()]
	return exists
}

func (rs *TsxResourceSystem) Load(ctx finch.Context, handle resources.ResourceHandle) error {
	metadata, exists := handle.Metadata()
	if !exists {
		return fmt.Errorf("cannot find metadata in manifest: %s", handle.Key())
	}

	_, err := rs.load_tsx(ctx, handle.Key(), metadata)
	if err != nil {
		return err
	}

	ctx.Logger().Info("resource loaded", slog.String("key", handle.Key()))
	return nil
}

func (rs *TsxResourceSystem) Unload(ctx finch.Context, handle resources.ResourceHandle) error {
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

	metadata.Properties = map[string]any{
		"refs": imgRefs,
	}

	return nil
}

func (rs *TsxResourceSystem) GetDependencies(ctx finch.Context, handle resources.ResourceHandle) (imgRefs []resources.ResourceHandle) {
	metadata, exists := handle.Metadata()
	if !exists {
		return nil
	}

	if metadata.Properties == nil {
		return nil
	}

	raw, exists := metadata.Properties["refs"]
	if !exists {
		return nil
	}

	s, ok := raw.([]any)
	if !ok {
		return nil
	}

	imgRefs = make([]resources.ResourceHandle, 0, len(s))
	for _, v := range s {
		str, ok := v.(string)
		if ok {
			imgRefs = append(imgRefs, resources.ResourceHandle(str))
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
