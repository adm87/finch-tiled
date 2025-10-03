package tiled

import (
	"encoding/xml"
	"fmt"

	"github.com/adm87/finch-core/enum"
	"github.com/adm87/finch-core/geom"
)

// ======================================================
// Miscellaneous Types
// ======================================================

type FlipFlags uint8

const (
	FLIP_NONE       FlipFlags = 0
	FLIP_HORIZONTAL FlipFlags = 1 << iota
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

func (f FlipFlags) FlipHorizontal() bool {
	return (f & FLIP_HORIZONTAL) != 0
}

func (f FlipFlags) FlipVertical() bool {
	return (f & FLIP_VERTICAL) != 0
}

func (f FlipFlags) FlipDiagonal() bool {
	return (f & FLIP_DIAGONAL) != 0
}

func (f FlipFlags) FlipRotatedHex() bool {
	return (f & FLIP_ROTATED_HEX) != 0
}

// ======================================================
// Tile Type
// ======================================================

type Tile struct {
	GID           uint32
	TsxSrc        string
	X, Y          float64
	Width, Height float64
	Flags         FlipFlags
}

type LayerPartitions map[geom.Rect64][]*Tile

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

type TiledXMLAttr any
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
	OrientationAttr  = "orientation"
	PropertyTypeAttr = "propertytype"
	RenderOrderAttr  = "renderorder"
	SourceAttr       = "source"
	SpacingAttr      = "spacing"
	TemplateAttr     = "template"
	TileCountAttr    = "tilecount"
	TileHeightAttr   = "tileheight"
	TileWidthAttr    = "tilewidth"
	TiledVersionAttr = "tiledversion"
	ValueAttr        = "value"
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
	PropertyTypeAttr: func(s string) (TiledXMLAttr, error) { return UnmarshalAttrString(s) },
	ValueAttr:        func(s string) (TiledXMLAttr, error) { return UnmarshalAttrString(s) },
	TemplateAttr:     func(s string) (TiledXMLAttr, error) { return UnmarshalAttrString(s) },
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
// TSX TileOffset Property
// ======================================================

type Offset struct {
	Attrs TiledXMLAttrTable `xml:",any,attr"`
}

func (offset Offset) X() int {
	if x, exists := offset.Attrs[XAttr]; exists {
		if attr, ok := x.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}

func (offset Offset) Y() int {
	if y, exists := offset.Attrs[YAttr]; exists {
		if attr, ok := y.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}

// ======================================================
// TMX Encoding
// ======================================================

type Encoding int

const (
	TMXEncodingCSV Encoding = iota
	TMXEncodingBase64
)

func (e Encoding) String() string {
	switch e {
	case TMXEncodingCSV:
		return "csv"
	case TMXEncodingBase64:
		return "base64"
	default:
		return "unknown"
	}
}

func (e Encoding) IsValid() bool {
	return e >= TMXEncodingCSV && e <= TMXEncodingBase64
}

func (e Encoding) MarshalJSON() ([]byte, error) {
	return enum.MarshalEnum(e)
}

func (e *Encoding) UnmarshalJSON(data []byte) error {
	val, err := enum.UnmarshalEnum[Encoding](data)
	if err != nil {
		return err
	}
	*e = val
	return nil
}

// ======================================================
// TMX Orientation
// ======================================================

type Orientation int

const (
	Orthogonal Orientation = iota
	Isometric
	Staggered
	Hexagonal
)

func (o Orientation) String() string {
	switch o {
	case Orthogonal:
		return "orthogonal"
	case Isometric:
		return "isometric"
	case Staggered:
		return "staggered"
	case Hexagonal:
		return "hexagonal"
	default:
		return "unknown"
	}
}

func (o Orientation) IsValid() bool {
	return o >= Orthogonal && o <= Hexagonal
}

func (o Orientation) MarshalJSON() ([]byte, error) {
	return enum.MarshalEnum(o)
}

func (o *Orientation) UnmarshalJSON(data []byte) error {
	val, err := enum.UnmarshalEnum[Orientation](data)
	if err != nil {
		return err
	}
	*o = val
	return nil
}

// ======================================================
// TMX Render Order
// ======================================================

type RenderOrder int

const (
	TMXRightDown RenderOrder = iota
	TMXRightUp
	TMXLeftDown
	TMXLeftUp
)

func (ro RenderOrder) String() string {
	switch ro {
	case TMXRightDown:
		return "right-down"
	case TMXRightUp:
		return "right-up"
	case TMXLeftDown:
		return "left-down"
	case TMXLeftUp:
		return "left-up"
	default:
		return "unknown"
	}
}

func (ro RenderOrder) IsValid() bool {
	return ro >= TMXRightDown && ro <= TMXLeftUp
}

func (ro RenderOrder) MarshalJSON() ([]byte, error) {
	return enum.MarshalEnum(ro)
}

func (ro *RenderOrder) UnmarshalJSON(data []byte) error {
	val, err := enum.UnmarshalEnum[RenderOrder](data)
	if err != nil {
		return err
	}
	*ro = val
	return nil
}

// ======================================================
// Image Property
// ======================================================

type Image struct {
	Attrs TiledXMLAttrTable `xml:",any,attr"`
}

func (img Image) Source() string {
	if source, exists := img.Attrs[SourceAttr]; exists {
		if attr, ok := source.(AttrString); ok {
			return attr.String()
		}
	}
	return ""
}

func (img Image) Width() int {
	if width, exists := img.Attrs[WidthAttr]; exists {
		if attr, ok := width.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}

func (img Image) Height() int {
	if height, exists := img.Attrs[HeightAttr]; exists {
		if attr, ok := height.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}

// ======================================================
// Layer Data
// ======================================================

type LayerData struct {
	Attrs  TiledXMLAttrTable `xml:",any,attr"`
	Chunks []*DataChunk      `xml:"chunk"`
	Data   string            `xml:",chardata"`
}

func (data LayerData) Encoding() Encoding {
	if encoding, exists := data.Attrs[EncodingAttr]; exists {
		if attr, ok := encoding.(AttrString); ok {
			e, err := enum.Value[Encoding](attr.String())
			if err != nil {
				panic(err)
			}
			return e
		}
	}
	return TMXEncodingCSV
}

// ======================================================
// Data Chunk
// ======================================================

type DataChunk struct {
	Attrs TiledXMLAttrTable `xml:",any,attr"`
	Data  string            `xml:",chardata"`
}

func (chunk DataChunk) X() int {
	if x, exists := chunk.Attrs[XAttr]; exists {
		if attr, ok := x.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}

func (chunk DataChunk) Y() int {
	if y, exists := chunk.Attrs[YAttr]; exists {
		if attr, ok := y.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}

func (chunk DataChunk) Width() int {
	if width, exists := chunk.Attrs[WidthAttr]; exists {
		if attr, ok := width.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}

func (chunk DataChunk) Height() int {
	if height, exists := chunk.Attrs[HeightAttr]; exists {
		if attr, ok := height.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}

// ======================================================
// Layer
// ======================================================

type Layer struct {
	Attrs TiledXMLAttrTable `xml:",any,attr"`
	Data  *LayerData        `xml:"data"`

	// Should these be stored here? Don't serialize them!
	tiles      []*Tile
	partitions LayerPartitions
}

func (layer Layer) ID() int {
	if id, exists := layer.Attrs[IDAttr]; exists {
		if attr, ok := id.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}

func (layer Layer) Name() string {
	if name, exists := layer.Attrs[NameAttr]; exists {
		if attr, ok := name.(AttrString); ok {
			return attr.String()
		}
	}
	return ""
}

func (layer Layer) Width() int {
	if width, exists := layer.Attrs[WidthAttr]; exists {
		if attr, ok := width.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}

func (layer Layer) Height() int {
	if height, exists := layer.Attrs[HeightAttr]; exists {
		if attr, ok := height.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}

func (layer Layer) IsVisible() bool {
	if visible, exists := layer.Attrs[VisibleAttr]; exists {
		if attr, ok := visible.(AttrBool); ok {
			return attr.Bool()
		}
	}
	return true
}

// ======================================================
// Property
// ======================================================

type Property struct {
	Attrs      TiledXMLAttrTable `xml:",any,attr"`
	Properties []*Property       `xml:"property"`
}

func (prop Property) Name() string {
	if name, exists := prop.Attrs[NameAttr]; exists {
		if attr, ok := name.(AttrString); ok {
			return attr.String()
		}
	}
	return ""
}

func (prop Property) Type() string {
	if ptype, exists := prop.Attrs[PropertyTypeAttr]; exists {
		if attr, ok := ptype.(AttrString); ok {
			return attr.String()
		}
	}
	return "string"
}

func (prop Property) Value() string {
	if value, exists := prop.Attrs["value"]; exists {
		if attr, ok := value.(AttrString); ok {
			return attr.String()
		}
	}
	return ""
}

func (prop Property) PropertyType() string {
	if ptype, exists := prop.Attrs[PropertyTypeAttr]; exists {
		if attr, ok := ptype.(AttrString); ok {
			return attr.String()
		}
	}
	return "string"
}

func (prop Property) PropertyOfType(ptype string) (*Property, bool) {
	for _, p := range prop.Properties {
		if p.PropertyType() == ptype {
			return p, true
		}
	}
	return nil, false
}

// ======================================================
// ObjectGroups
// ======================================================

type ObjectGroup struct {
	Attrs   TiledXMLAttrTable `xml:",any,attr"`
	Objects []*Object         `xml:"object"`
}

func (og ObjectGroup) ID() int {
	if id, exists := og.Attrs[IDAttr]; exists {
		if attr, ok := id.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}

func (og ObjectGroup) Name() string {
	if name, exists := og.Attrs[NameAttr]; exists {
		if attr, ok := name.(AttrString); ok {
			return attr.String()
		}
	}
	return ""
}

// ======================================================
// Object
// ======================================================

type Object struct {
	Attrs      TiledXMLAttrTable `xml:",any,attr"`
	Properties []*Property       `xml:"properties>property"`
}

func (obj Object) ID() int {
	if id, exists := obj.Attrs[IDAttr]; exists {
		if attr, ok := id.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}

func (obj Object) X() int {
	if x, exists := obj.Attrs[XAttr]; exists {
		if attr, ok := x.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}

func (obj Object) Y() int {
	if y, exists := obj.Attrs[YAttr]; exists {
		if attr, ok := y.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}

func (obj Object) Width() int {
	if width, exists := obj.Attrs[WidthAttr]; exists {
		if attr, ok := width.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}

func (obj Object) Height() int {
	if height, exists := obj.Attrs[HeightAttr]; exists {
		if attr, ok := height.(AttrInt); ok {
			return attr.Int()
		}
	}
	return 0
}

func (obj Object) Name() string {
	if name, exists := obj.Attrs[NameAttr]; exists {
		if attr, ok := name.(AttrString); ok {
			return attr.String()
		}
	}
	return ""
}

func (obj Object) Template() string {
	if template, exists := obj.Attrs[TemplateAttr]; exists {
		if attr, ok := template.(AttrString); ok {
			return attr.String()
		}
	}
	return ""
}

func (obj Object) PropertyOfType(ptype string) (*Property, bool) {
	for _, prop := range obj.Properties {
		if prop.PropertyType() == ptype {
			return prop, true
		}
	}
	return nil, false
}

func (obj Object) HasTemplate() bool {
	return obj.Template() != ""
}

// ======================================================
// Tileset
// ======================================================

type Tileset struct {
	Attrs TiledXMLAttrTable `xml:",any,attr"`
}

func (ts Tileset) FirstGID() uint32 {
	if firstGID, exists := ts.Attrs[FirstGIDAttr]; exists {
		if attr, ok := firstGID.(AttrInt); ok {
			return uint32(attr.Int())
		}
	}
	return 0
}

func (ts Tileset) Source() string {
	if source, exists := ts.Attrs[SourceAttr]; exists {
		if attr, ok := source.(AttrString); ok {
			return attr.String()
		}
	}
	return ""
}
