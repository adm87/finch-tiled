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
	TsxSrc        string
	X, Y          float64
	Width, Height float64
	Flags         TiledFlags
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
