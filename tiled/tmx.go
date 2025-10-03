package tiled

import "github.com/adm87/finch-core/enum"

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

func (tmx TMX) FindTilesetByTileGID(gid uint32) (*Tileset, bool) {
	for i := len(tmx.Tilesets) - 1; i >= 0; i-- {
		if tmx.Tilesets[i].FirstGID() <= gid {
			return tmx.Tilesets[i], true
		}
	}
	return nil, false
}

func (tmx TMX) GetLayerByName(name string) (*Layer, bool) {
	for _, layer := range tmx.Layers {
		if layer.Name() == name {
			return layer, true
		}
	}
	return nil, false
}

func (tmx TMX) GetObjectGroupByName(name string) (*ObjectGroup, bool) {
	for _, objectGroup := range tmx.ObjectGroups {
		if objectGroup.Name() == name {
			return objectGroup, true
		}
	}
	return nil, false
}
