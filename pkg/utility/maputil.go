package utility

import (
	"github.com/gookit/goutil/maputil"
)

type MapData maputil.Data

func (d *MapData) MustOk() maputil.Data {
	if *d == nil {
		*d = make(MapData)
	}
	return maputil.Data(*d)
}

func (d MapData) StdMap() map[string]any {
	return d
}
