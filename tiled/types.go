package tiled

import (
	"encoding/xml"
	"fmt"

	"github.com/adm87/finch-core/geom"
)

// ======================================================
// Miscellaneous Types
// ======================================================

type TiledFlags uint8

const (
	FLIP_NONE       TiledFlags = 0
	FLIP_HORIZONTAL TiledFlags = 1 << iota
	FLIP_VERTICAL
	FLIP_DIAGONAL
	FLIP_ROTATED_HEX
)

const (
	// These constants represent the bit flags used by Tiled to encode tile transformations.
	// See: https://doc.mapeditor.org/en/stable/reference/global-tile-ids/#tile-flipping

	TILE_FLIP_HORIZONTAL  = 0x80000000
	TILE_FLIP_VERTICAL    = 0x40000000
	TILE_FLIP_DIAGONAL    = 0x20000000
	TILE_FLIP_ROTATED_HEX = 0x10000000
	TILE_ID_MASK          = 0x1FFFFFFF
)

type Tile struct {
	GID           uint32
	TsxKey        string
	X, Y          float64
	Width, Height float64
	Flags         TiledFlags
}

type LayerPartitions map[geom.Rect64][]*Tile

type TiledXMLAttr any

// ======================================================
// String Attribute
// ======================================================

type AttrString string

func UnmarshalAttrString(s string) (AttrString, error) {
	return AttrString(s), nil
}

func (s AttrString) String() string {
	return string(s)
}

// ======================================================
// Integer Attribute
// ======================================================

type AttrInt int

func UnmarshalAttrInt(s string) (AttrInt, error) {
	var v int
	_, err := fmt.Sscanf(s, "%d", &v)
	if err != nil {
		return 0, fmt.Errorf("invalid integer attribute: %s", s)
	}
	return AttrInt(v), nil
}

func (i AttrInt) Int() int {
	return int(i)
}

// ======================================================
// Boolean Attribute
// ======================================================

type AttrBool bool

func UnmarshalAttrBool(s string) (AttrBool, error) {
	var b AttrBool
	if s == "1" || s == "true" {
		b = AttrBool(true)
	} else if s == "0" || s == "false" {
		b = AttrBool(false)
	} else {
		return false, fmt.Errorf("invalid boolean attribute: %s", s)
	}
	return b, nil
}

func (b AttrBool) Bool() bool {
	return bool(b)
}

// ======================================================
// Tiled XML Attribute Table
// ======================================================

type TiledXMLAttrTable map[string]TiledXMLAttr

const (
	ColumnsAttr      = "columns"
	EncodingAttr     = "encoding"
	FirstGIDAttr     = "firstgid"
	HeightAttr       = "height"
	IDAttr           = "id"
	InfiniteAttr     = "infinite"
	LockedAttr       = "locked"
	NameAttr         = "name"
	NextLayerIDAttr  = "nextlayerid"
	NextObjectIDAttr = "nextobjectid"
	object           = "object"
	objectgroup      = "objectgroup"
	OrientationAttr  = "orientation"
	RenderOrderAttr  = "renderorder"
	SourceAttr       = "source"
	SpacingAttr      = "spacing"
	TileCountAttr    = "tilecount"
	TileHeightAttr   = "tileheight"
	TileWidthAttr    = "tilewidth"
	TiledVersionAttr = "tiledversion"
	VersionAttr      = "version"
	VisibleAttr      = "visible"
	WidthAttr        = "width"
	XAttr            = "x"
	YAttr            = "y"
)

var attr_unmarshallers = map[string]func(s string) (TiledXMLAttr, error){
	RenderOrderAttr:  func(s string) (TiledXMLAttr, error) { return UnmarshalAttrString(s) },
	OrientationAttr:  func(s string) (TiledXMLAttr, error) { return UnmarshalAttrString(s) },
	VersionAttr:      func(s string) (TiledXMLAttr, error) { return UnmarshalAttrString(s) },
	TiledVersionAttr: func(s string) (TiledXMLAttr, error) { return UnmarshalAttrString(s) },
	NameAttr:         func(s string) (TiledXMLAttr, error) { return UnmarshalAttrString(s) },
	SourceAttr:       func(s string) (TiledXMLAttr, error) { return UnmarshalAttrString(s) },
	EncodingAttr:     func(s string) (TiledXMLAttr, error) { return UnmarshalAttrString(s) },
	InfiniteAttr:     func(s string) (TiledXMLAttr, error) { return UnmarshalAttrBool(s) },
	VisibleAttr:      func(s string) (TiledXMLAttr, error) { return UnmarshalAttrBool(s) },
	LockedAttr:       func(s string) (TiledXMLAttr, error) { return UnmarshalAttrBool(s) },
	WidthAttr:        func(s string) (TiledXMLAttr, error) { return UnmarshalAttrInt(s) },
	HeightAttr:       func(s string) (TiledXMLAttr, error) { return UnmarshalAttrInt(s) },
	TileWidthAttr:    func(s string) (TiledXMLAttr, error) { return UnmarshalAttrInt(s) },
	TileHeightAttr:   func(s string) (TiledXMLAttr, error) { return UnmarshalAttrInt(s) },
	SpacingAttr:      func(s string) (TiledXMLAttr, error) { return UnmarshalAttrInt(s) },
	TileCountAttr:    func(s string) (TiledXMLAttr, error) { return UnmarshalAttrInt(s) },
	ColumnsAttr:      func(s string) (TiledXMLAttr, error) { return UnmarshalAttrInt(s) },
	FirstGIDAttr:     func(s string) (TiledXMLAttr, error) { return UnmarshalAttrInt(s) },
	IDAttr:           func(s string) (TiledXMLAttr, error) { return UnmarshalAttrInt(s) },
	XAttr:            func(s string) (TiledXMLAttr, error) { return UnmarshalAttrInt(s) },
	YAttr:            func(s string) (TiledXMLAttr, error) { return UnmarshalAttrInt(s) },
	NextLayerIDAttr:  func(s string) (TiledXMLAttr, error) { return UnmarshalAttrInt(s) },
	NextObjectIDAttr: func(s string) (TiledXMLAttr, error) { return UnmarshalAttrInt(s) },
}

func (m *TiledXMLAttrTable) UnmarshalXMLAttr(attr xml.Attr) error {
	unmarshal, ok := attr_unmarshallers[attr.Name.Local]

	if !ok {
		println("TiledXMLAttrTable:UnmarshalXMLAttr - unknown attribute:", attr.Name.Local)
		return nil
	}

	if *m == nil {
		*m = make(map[string]TiledXMLAttr)
	}

	parsed, err := unmarshal(attr.Value)

	if err != nil {
		return err
	}

	(*m)[attr.Name.Local] = parsed
	return nil
}

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
// TMX File
// ======================================================

// TMX represents a deserialized Tiled tmx file.
type TMX struct {
	Attrs        TiledXMLAttrTable `xml:",any,attr"`
	ObjectGroups []*TMXObjectGroup `xml:"objectgroup"`
	Tilesets     []*TMXTileset     `xml:"tileset"`
	Layers       []*TMXLayer       `xml:"layer"`
}

func (tmx TMX) Orientation() TMXOrientation {
	if orientation, exists := tmx.Attrs[OrientationAttr]; exists {
		if attr, ok := orientation.(AttrString); ok {
			return TMXOrientation(attr.String())
		}
	}
	return TMXOrthogonal
}

func (tmx TMX) RenderOrder() TMXRenderOrder {
	if renderOrder, exists := tmx.Attrs[RenderOrderAttr]; exists {
		if attr, ok := renderOrder.(AttrString); ok {
			return TMXRenderOrder(attr.String())
		}
	}
	return TMXRightDown
}

func (tmx TMX) Version() string {
	if version, exists := tmx.Attrs[VersionAttr]; exists {
		if attr, ok := version.(AttrString); ok {
			return attr.String()
		}
	}
	return "unknown"
}

func (tmx TMX) TiledVersion() string {
	if tiledVersion, exists := tmx.Attrs[TiledVersionAttr]; exists {
		if attr, ok := tiledVersion.(AttrString); ok {
			return attr.String()
		}
	}
	return "unknown"
}

func (tmx TMX) Width() int {
	if width, exists := tmx.Attrs[WidthAttr]; exists {
		if attr, ok := width.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}

func (tmx TMX) Height() int {
	if height, exists := tmx.Attrs[HeightAttr]; exists {
		if attr, ok := height.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}

func (tmx TMX) TileWidth() int {
	if tileWidth, exists := tmx.Attrs[TileWidthAttr]; exists {
		if attr, ok := tileWidth.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}

func (tmx TMX) TileHeight() int {
	if tileHeight, exists := tmx.Attrs[TileHeightAttr]; exists {
		if attr, ok := tileHeight.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}

func (tmx TMX) IsInfinite() bool {
	if infinite, exists := tmx.Attrs[InfiniteAttr]; exists {
		if attr, ok := infinite.(AttrBool); ok {
			return attr.Bool()
		}
	}
	return false
}

func (tmx TMX) NextLayerID() int {
	if nextLayerID, exists := tmx.Attrs[NextLayerIDAttr]; exists {
		if attr, ok := nextLayerID.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}

func (tmx TMX) NextObjectID() int {
	if nextObjectID, exists := tmx.Attrs[NextObjectIDAttr]; exists {
		if attr, ok := nextObjectID.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}

func (tmx TMX) FindTilesetByTileGID(gid uint32) (*TMXTileset, bool) {
	for i := len(tmx.Tilesets) - 1; i >= 0; i-- {
		if tmx.Tilesets[i].FirstGID() <= gid {
			return tmx.Tilesets[i], true
		}
	}
	return nil, false
}

// ======================================================
// TMX Tileset Property
// ======================================================

type TMXTileset struct {
	Attrs TiledXMLAttrTable `xml:",any,attr"`
}

func (ts TMXTileset) FirstGID() uint32 {
	if firstGID, exists := ts.Attrs[FirstGIDAttr]; exists {
		if attr, ok := firstGID.(AttrInt); ok {
			return uint32(attr.Int())
		}
	}
	return 0
}

func (ts TMXTileset) Source() string {
	if source, exists := ts.Attrs[SourceAttr]; exists {
		if attr, ok := source.(AttrString); ok {
			return attr.String()
		}
	}
	return ""
}

// ======================================================
// TMX ObjectGroups Property
// ======================================================

type TMXObjectGroup struct {
	Attrs   TiledXMLAttrTable `xml:",any,attr"`
	Objects []*TMXObject      `xml:"object"`
}

func (og TMXObjectGroup) ID() int {
	if id, exists := og.Attrs[IDAttr]; exists {
		if attr, ok := id.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}

func (og TMXObjectGroup) Name() string {
	if name, exists := og.Attrs[NameAttr]; exists {
		if attr, ok := name.(AttrString); ok {
			return attr.String()
		}
	}
	return ""
}

// ======================================================
// TMX Object Property
// ======================================================

type TMXObject struct {
	Attrs TiledXMLAttrTable `xml:",any,attr"`
}

func (obj TMXObject) ID() int {
	if id, exists := obj.Attrs[IDAttr]; exists {
		if attr, ok := id.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}

func (obj TMXObject) X() int {
	if x, exists := obj.Attrs[XAttr]; exists {
		if attr, ok := x.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}

func (obj TMXObject) Y() int {
	if y, exists := obj.Attrs[YAttr]; exists {
		if attr, ok := y.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}

func (obj TMXObject) Width() int {
	if width, exists := obj.Attrs[WidthAttr]; exists {
		if attr, ok := width.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}

func (obj TMXObject) Height() int {
	if height, exists := obj.Attrs[HeightAttr]; exists {
		if attr, ok := height.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}

// ======================================================
// TMX Layer Property
// ======================================================

type TMXLayer struct {
	Attrs TiledXMLAttrTable `xml:",any,attr"`
	Data  *TMXLayerData     `xml:"data"`

	// Should these be stored here? Don't serialize them!
	tiles      []*Tile
	partitions LayerPartitions
}

func (layer TMXLayer) ID() int {
	if id, exists := layer.Attrs[IDAttr]; exists {
		if attr, ok := id.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}

func (layer TMXLayer) Name() string {
	if name, exists := layer.Attrs[NameAttr]; exists {
		if attr, ok := name.(AttrString); ok {
			return attr.String()
		}
	}
	return ""
}

func (layer TMXLayer) Width() int {
	if width, exists := layer.Attrs[WidthAttr]; exists {
		if attr, ok := width.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}

func (layer TMXLayer) Height() int {
	if height, exists := layer.Attrs[HeightAttr]; exists {
		if attr, ok := height.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}

func (layer TMXLayer) IsVisible() bool {
	if visible, exists := layer.Attrs[VisibleAttr]; exists {
		if attr, ok := visible.(AttrBool); ok {
			return attr.Bool()
		}
	}
	return true
}

// ======================================================
// TMX Layer Data Property
// ======================================================

type TMXLayerData struct {
	Attrs  TiledXMLAttrTable `xml:",any,attr"`
	Chunks []*TMXDataChunk   `xml:"chunk"`
	Data   string            `xml:",chardata"`
}

func (data TMXLayerData) Encoding() TMXEncoding {
	if encoding, exists := data.Attrs[EncodingAttr]; exists {
		if attr, ok := encoding.(AttrString); ok {
			return TMXEncoding(attr.String())
		}
	}
	return TMXEncodingCSV
}

// ======================================================
// TMX Data Chunk Property
// ======================================================

type TMXDataChunk struct {
	Attrs TiledXMLAttrTable `xml:",any,attr"`
	Data  string            `xml:",chardata"`
}

func (chunk TMXDataChunk) X() int {
	if x, exists := chunk.Attrs[XAttr]; exists {
		if attr, ok := x.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}

func (chunk TMXDataChunk) Y() int {
	if y, exists := chunk.Attrs[YAttr]; exists {
		if attr, ok := y.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}

func (chunk TMXDataChunk) Width() int {
	if width, exists := chunk.Attrs[WidthAttr]; exists {
		if attr, ok := width.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}

func (chunk TMXDataChunk) Height() int {
	if height, exists := chunk.Attrs[HeightAttr]; exists {
		if attr, ok := height.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}

// ======================================================
// TSX File
// ======================================================

type TSX struct {
	Attrs      TiledXMLAttrTable `xml:",any,attr"`
	TileOffset *TSXTileOffset    `xml:"tileoffset"`
	Image      *TSXImage         `xml:"image"`
}

func (tsx TSX) Version() string {
	if version, exists := tsx.Attrs[VersionAttr]; exists {
		if attr, ok := version.(AttrString); ok {
			return attr.String()
		}
	}
	return ""
}

func (tsx TSX) TiledVersion() string {
	if tiledVersion, exists := tsx.Attrs[TiledVersionAttr]; exists {
		if attr, ok := tiledVersion.(AttrString); ok {
			return attr.String()
		}
	}
	return ""
}

func (tsx TSX) Name() string {
	if name, exists := tsx.Attrs[NameAttr]; exists {
		if attr, ok := name.(AttrString); ok {
			return attr.String()
		}
	}
	return ""
}

func (tsx TSX) TileWidth() int {
	if tileWidth, exists := tsx.Attrs[TileWidthAttr]; exists {
		if attr, ok := tileWidth.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}

func (tsx TSX) TileHeight() int {
	if tileHeight, exists := tsx.Attrs[TileHeightAttr]; exists {
		if attr, ok := tileHeight.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}

func (tsx TSX) Spacing() int {
	if spacing, exists := tsx.Attrs[SpacingAttr]; exists {
		if attr, ok := spacing.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}

func (tsx TSX) TileCount() int {
	if tileCount, exists := tsx.Attrs[TileCountAttr]; exists {
		if attr, ok := tileCount.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}

func (tsx TSX) Columns() int {
	if columns, exists := tsx.Attrs[ColumnsAttr]; exists {
		if attr, ok := columns.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}

// ======================================================
// TSX TileOffset Property
// ======================================================

type TSXTileOffset struct {
	Attrs TiledXMLAttrTable `xml:",any,attr"`
}

func (offset TSXTileOffset) X() int {
	if x, exists := offset.Attrs[XAttr]; exists {
		if attr, ok := x.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}

func (offset TSXTileOffset) Y() int {
	if y, exists := offset.Attrs[YAttr]; exists {
		if attr, ok := y.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}

// ======================================================
// TSX Image Property
// ======================================================

type TSXImage struct {
	Attrs TiledXMLAttrTable `xml:",any,attr"`
}

func (img TSXImage) Source() string {
	if source, exists := img.Attrs[SourceAttr]; exists {
		if attr, ok := source.(AttrString); ok {
			return attr.String()
		}
	}
	return ""
}

func (img TSXImage) Width() int {
	if width, exists := img.Attrs[WidthAttr]; exists {
		if attr, ok := width.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}

func (img TSXImage) Height() int {
	if height, exists := img.Attrs[HeightAttr]; exists {
		if attr, ok := height.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}
