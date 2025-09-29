package tiled

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

func (tmx TMX) GetLayerByName(name string) (*TMXLayer, bool) {
	for _, layer := range tmx.Layers {
		if layer.Name() == name {
			return layer, true
		}
	}
	return nil, false
}

func (tmx TMX) GetObjectGroupByName(name string) (*TMXObjectGroup, bool) {
	for _, objectGroup := range tmx.ObjectGroups {
		if objectGroup.Name() == name {
			return objectGroup, true
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
