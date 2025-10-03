package tiled

// ======================================================
// Tiled XML Attribute Table
// ======================================================

type TX struct {
	Attrs   TiledXMLAttrTable `xml:",any,attr"`
	Tileset *Tileset          `xml:"tileset"`
	Object  *Object           `xml:"object"`
}
