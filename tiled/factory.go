package tiled

type TiledObjectFactory[T any] struct {
	FromTemplate func(instance *Object, template *TX, tmx *TMX) T
	FromObject   func(obj *Object, tmx *TMX) T
}
