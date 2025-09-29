package tiled

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
