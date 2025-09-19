package tiled

import (
	"encoding/xml"
	"fmt"
	"log/slog"
	"sync"

	"github.com/adm87/finch-core/finch"
	"github.com/adm87/finch-core/types"
	"github.com/adm87/finch-resources/resources"
)

// ======================================================
// TMX Resource System
// ======================================================

var tmxSystem = resources.NewResourceSystemKey[*TmxResourceSystem]()

type TmxResourceSystem struct {
	tilemaps map[string]*TMX
	loading  types.HashSet[string]
	mu       sync.Mutex
}

func NewTmxResourceSystem() *TmxResourceSystem {
	return &TmxResourceSystem{
		tilemaps: make(map[string]*TMX),
		loading:  make(types.HashSet[string]),
		mu:       sync.Mutex{},
	}
}

func GetTmx(handle resources.ResourceHandle) (*TMX, bool) {
	sys, ok := resources.GetSystem(tmxSystem).(*TmxResourceSystem)
	if !ok {
		return nil, false
	}

	sys.mu.Lock()
	defer sys.mu.Unlock()

	tmx, exists := sys.tilemaps[handle.Key()]
	return tmx, exists
}

func (rs *TmxResourceSystem) ResourceTypes() []string {
	return []string{"tmx"}
}

func (rs *TmxResourceSystem) Type() resources.ResourceSystemType {
	return tmxSystem
}

func (rs *TmxResourceSystem) IsLoaded(handle resources.ResourceHandle) bool {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	_, exists := rs.tilemaps[handle.Key()]
	return exists
}

func (rs *TmxResourceSystem) Load(ctx finch.Context, handle resources.ResourceHandle) error {
	metadata, exists := handle.Metadata()
	if !exists {
		return fmt.Errorf("cannot find metadata in manifest: %s", handle.Key())
	}

	_, err := rs.load_tmx(ctx, handle.Key(), metadata)
	if err != nil {
		return err
	}

	ctx.Logger().Info("resource loaded", slog.String("key", handle.Key()))
	return nil
}

func (rs *TmxResourceSystem) Unload(ctx finch.Context, handle resources.ResourceHandle) error {
	return nil
}

func (rs *TmxResourceSystem) GenerateMetadata(ctx finch.Context, key string, metadata *resources.Metadata) error {
	tmx, err := rs.load_tmx(ctx, key, metadata)
	if err != nil {
		return err
	}

	var tsxRefs []string
	for _, ts := range tmx.Tilesets {
		tsxRefs = append(tsxRefs, resources.KeyFromPath(ts.Source()))
	}

	metadata.Properties = map[string]any{
		"refs": tsxRefs,
	}

	return nil
}

func (rs *TmxResourceSystem) GetDependencies(ctx finch.Context, handle resources.ResourceHandle) (tsxRefs []resources.ResourceHandle) {
	metadata, exists := handle.Metadata()
	if !exists {
		return nil
	}

	if metadata.Properties == nil {
		return nil
	}

	raw, ok := metadata.Properties["refs"]
	if !ok {
		return nil
	}

	s, ok := raw.([]any)
	if !ok {
		return nil
	}

	tsxRefs = make([]resources.ResourceHandle, 0, len(s))
	for _, v := range s {
		str, ok := v.(string)
		if ok {
			tsxRefs = append(tsxRefs, resources.ResourceHandle(str))
		}
	}

	return
}

func (rs *TmxResourceSystem) load_tmx(ctx finch.Context, key string, metadata *resources.Metadata) (*TMX, error) {
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

	var tmx TMX
	if err := xml.Unmarshal(data, &tmx); err != nil {
		return nil, err
	}

	rs.mu.Lock()
	rs.tilemaps[key] = &tmx
	rs.mu.Unlock()

	return &tmx, nil
}

func (rs *TmxResourceSystem) try_load(key string) error {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	if _, exists := rs.tilemaps[key]; exists {
		return fmt.Errorf("tmx resource is already loaded: %s", key)
	}
	if rs.loading.Contains(key) {
		return fmt.Errorf("tmx resource is already loading: %s", key)
	}

	rs.loading.Add(key)
	return nil
}
