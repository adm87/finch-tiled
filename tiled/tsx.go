package tiled

// ======================================================
// TSX File
// ======================================================

type TSX struct {
	Attrs      TiledXMLAttrTable `xml:",any,attr"`
	TileOffset *Offset           `xml:"tileoffset"`
	Image      *Image            `xml:"image"`
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
