package module

import (
	"github.com/adm87/finch-core/finch"
	"github.com/adm87/finch-resources/resources"
	"github.com/adm87/finch-tiled/tiled"
)

func Register(ctx finch.Context) {
	resources.RegisterSystem(tiled.NewTmxResourceSystem())
	resources.RegisterSystem(tiled.NewTsxResourceSystem())

	ctx.Logger().Info("tiled module registered")
}
