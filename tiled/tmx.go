package tiled

import (
	"encoding/xml"
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
	TMXMetadataTSXRefs = "tsx_refs"
)

// ======================================================
// TMX Encoding
// ======================================================

type TMXEncoding string

const (
	TMXEncodingCSV    TMXEncoding = "csv"
	TMXEncodingBase64 TMXEncoding = "base64"
)

func (e TMXEncoding) String() string {
	return string(e)
}

func (e TMXEncoding) IsValid() bool {
	switch e {
	case TMXEncodingCSV, TMXEncodingBase64:
		return true
	default:
		return false
	}
}

// ======================================================
// TMX Infinite
// ======================================================

type TMXInfinite bool

func (i TMXInfinite) MarshalXMLAttr(name string) (xml.Attr, error) {
	if i {
		return xml.Attr{Name: xml.Name{Local: name}, Value: "1"}, nil
	}
	return xml.Attr{Name: xml.Name{Local: name}, Value: "0"}, nil
}

func (i *TMXInfinite) UnmarshalXMLAttr(attr xml.Attr) error {
	if attr.Value == "1" {
		*i = true
	} else {
		*i = false
	}
	return nil
}

// ======================================================
// TMX Orientation
// ======================================================

const (
	TMXOrthogonal TMXOrientation = "orthogonal"
	TMXIsometric  TMXOrientation = "isometric"
	TMXStaggered  TMXOrientation = "staggered"
	TMXHexagonal  TMXOrientation = "hexagonal"
)

type TMXOrientation string

func (o TMXOrientation) String() string {
	return string(o)
}

func (o TMXOrientation) IsValid() bool {
	switch o {
	case TMXOrthogonal, TMXIsometric, TMXStaggered, TMXHexagonal:
		return true
	default:
		return false
	}
}

// ======================================================
// TMX Render Order
// ======================================================

const (
	TMXRightDown TMXRenderOrder = "right-down"
	TMXRightUp   TMXRenderOrder = "right-up"
	TMXLeftDown  TMXRenderOrder = "left-down"
	TMXLeftUp    TMXRenderOrder = "left-up"
)

type TMXRenderOrder string

func (ro TMXRenderOrder) String() string {
	return string(ro)
}

func (ro TMXRenderOrder) IsValid() bool {
	switch ro {
	case TMXRightDown, TMXRightUp, TMXLeftDown, TMXLeftUp:
		return true
	default:
		return false
	}
}

// ======================================================
// TMX File Structure
// ======================================================

// TMX represents a deserialized Tiled tmx file.
type TMX struct {
	Version      string         `xml:"version,attr"`
	TiledVersion string         `xml:"tiledversion,attr"`
	Width        int            `xml:"width,attr"`
	Height       int            `xml:"height,attr"`
	TileWidth    int            `xml:"tilewidth,attr"`
	TileHeight   int            `xml:"tileheight,attr"`
	NextLayerID  int            `xml:"nextlayerid,attr"`
	NextObjectID int            `xml:"nextobjectid,attr"`
	Orientation  TMXOrientation `xml:"orientation,attr"`
	Infinite     TMXInfinite    `xml:"infinite,attr"`
	RenderOrder  TMXRenderOrder `xml:"renderorder,attr"`
	Tilesets     []TMXTileset   `xml:"tileset"`
	Layers       []TMXLayer     `xml:"layer"`
}

type TMXTileset struct {
	FirstGID int    `xml:"firstgid,attr"`
	Source   string `xml:"source,attr"`
}

type TMXLayer struct {
	ID     int     `xml:"id,attr"`
	Name   string  `xml:"name,attr"`
	Width  int     `xml:"width,attr"`
	Height int     `xml:"height,attr"`
	Data   TMXData `xml:"data"`
}

type TMXData struct {
	Encoding TMXEncoding `xml:"encoding,attr"`
	Chunks   []TMXChunk  `xml:"chunk"`
	Data     string      `xml:",chardata"`
}

type TMXChunk struct {
	X      int    `xml:"x,attr"`
	Y      int    `xml:"y,attr"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`
	Data   string `xml:",chardata"`
}

func (tmx *TMX) IsValid() error {
	if !tmx.Orientation.IsValid() {
		return fmt.Errorf("invalid map orientation: %s", tmx.Orientation)
	}
	if !tmx.RenderOrder.IsValid() {
		return fmt.Errorf("invalid map render order: %s", tmx.RenderOrder)
	}
	if tmx.TileWidth <= 0 || tmx.TileHeight <= 0 {
		return fmt.Errorf("invalid tile size: %dx%d", tmx.TileWidth, tmx.TileHeight)
	}
	if tmx.Width <= 0 || tmx.Height <= 0 {
		return fmt.Errorf("invalid map size: %dx%d", tmx.Width, tmx.Height)
	}
	for _, ts := range tmx.Tilesets {
		if ts.FirstGID <= 0 {
			return fmt.Errorf("invalid tileset firstgid: %d", ts.FirstGID)
		}
		if ts.Source == "" {
			return fmt.Errorf("tileset source is empty")
		}
	}
	for _, layer := range tmx.Layers {
		if layer.ID <= 0 {
			return fmt.Errorf("invalid layer id: %d", layer.ID)
		}
		if !layer.Data.Encoding.IsValid() {
			return fmt.Errorf("invalid layer data encoding: %s", layer.Data.Encoding)
		}
	}
	// TASK: Finish filling out validation checks.
	return nil
}

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

func (rs *TmxResourceSystem) ResourceTypes() []string {
	return []string{"tmx"}
}

func (rs *TmxResourceSystem) Type() resources.ResourceSystemType {
	return tmxSystem
}

func (rs *TmxResourceSystem) IsLoaded(key string) bool {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	_, exists := rs.tilemaps[key]
	return exists
}

func (rs *TmxResourceSystem) Load(ctx finch.Context, key string, metadata *resources.Metadata) error {
	_, err := rs.load_tmx(ctx, key, metadata)
	if err != nil {
		return err
	}

	ctx.Logger().Info("resource loaded", slog.String("key", key))
	return nil
}

func (rs *TmxResourceSystem) Unload(ctx finch.Context, key string) error {
	return nil
}

func (rs *TmxResourceSystem) GenerateMetadata(ctx finch.Context, key string, metadata *resources.Metadata) error {
	tmx, err := rs.load_tmx(ctx, key, metadata)
	if err != nil {
		return err
	}

	var tsxRefs []string
	for _, ts := range tmx.Tilesets {
		b := filepath.Base(ts.Source)
		tsxRefs = append(tsxRefs, strings.TrimSuffix(b, filepath.Ext(b)))
	}

	metadata.Extras = map[string]any{
		TMXMetadataTSXRefs: tsxRefs,
	}

	return nil
}

func (rs *TmxResourceSystem) GetDependencies(ctx finch.Context, key string, metadata *resources.Metadata) (tsxRefs []string) {
	if metadata.Extras == nil {
		return nil
	}

	raw, ok := metadata.Extras[TMXMetadataTSXRefs]
	if !ok {
		return nil
	}

	s, ok := raw.([]any)
	if !ok {
		return nil
	}

	tsxRefs = make([]string, 0, len(s))
	for _, v := range s {
		str, ok := v.(string)
		if ok {
			tsxRefs = append(tsxRefs, str)
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

	if err := tmx.IsValid(); err != nil {
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
