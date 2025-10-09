package tiled

import (
	"github.com/adm87/finch-core/enum"
	"github.com/adm87/finch-core/geom"
)

// ======================================================
// TMX File
// ======================================================

// TMX represents a deserialized Tiled tmx file.
type TMX struct {
	Attrs        TiledXMLAttrTable `xml:",any,attr"`
	ObjectGroups []*ObjectGroup    `xml:"objectgroup"`
	Tilesets     []*Tileset        `xml:"tileset"`
	Layers       []*Layer          `xml:"layer"`
}

func (tmx TMX) Orientation() Orientation {
	if orientation, exists := tmx.Attrs[OrientationAttr]; exists {
		if attr, ok := orientation.(AttrString); ok {
			e, err := enum.Value[Orientation](attr.String())
			if err != nil {
				panic(err)
			}
			return e
		}
	}
	return Orthogonal
}

func (tmx TMX) RenderOrder() RenderOrder {
	if renderOrder, exists := tmx.Attrs[RenderOrderAttr]; exists {
		if attr, ok := renderOrder.(AttrString); ok {
			e, err := enum.Value[RenderOrder](attr.String())
			if err != nil {
				panic(err)
			}
			return e
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

func (tmx TMX) LayerByName(name string) *Layer {
	for _, layer := range tmx.Layers {
		if layer.Name() == name {
			return layer
		}
	}
	return nil
}

func (tmx TMX) LayerByProperty(ptype string, pvalue any) *Layer {
	for _, layer := range tmx.Layers {
		if prop, exists := layer.PropertyOfType(ptype); exists {
			if prop.Value() == pvalue {
				return layer
			}
		}
	}
	return nil
}

func (tmx TMX) ObjectGroupByName(name string) *ObjectGroup {
	for _, og := range tmx.ObjectGroups {
		if og.Name() == name {
			return og
		}
	}
	return nil
}

func (tmx TMX) ObjectGroupByProperty(ptype string, pvalue any) *ObjectGroup {
	for _, og := range tmx.ObjectGroups {
		if prop, exists := og.PropertyOfType(ptype); exists {
			if prop.Value() == pvalue {
				return og
			}
		}
	}
	return nil
}

func (tmx TMX) Bounds() geom.Rect64 {
	bounds := geom.Rect64{}

	if len(tmx.Layers) == 0 {
		return bounds
	}

	if tmx.IsInfinite() {
		for _, layer := range tmx.Layers {
			bounds = bounds.Union(layer.Bounds())
		}
	} else {
		bounds = geom.NewRect64(0, 0, float64(tmx.Width()), float64(tmx.Height()))
	}

	bounds.X *= float64(tmx.TileWidth())
	bounds.Y *= float64(tmx.TileHeight())
	bounds.Width *= float64(tmx.TileWidth())
	bounds.Height *= float64(tmx.TileHeight())

	return bounds
}
