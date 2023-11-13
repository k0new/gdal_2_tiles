package app

import (
	"github.com/airbusgeo/godal"

	"grib-tiler/pkg/ds"
	"grib-tiler/pkg/equator"
	"grib-tiler/pkg/projection"
	"grib-tiler/pkg/raster"
	"grib-tiler/pkg/vector"
)

func (a *App) getProj() projection.Projection {
	switch a.opts.OutCRS {
	case "EPSG:3575":
		return projection.New3575()
	case "EPSG:3857":
		return projection.New3857()
	case "EPSG:4326":
		return projection.New4326()
	default:
		return projection.New4326()
	}
}

func (a *App) runReader(path string) {
	godal.RegisterAll()
	dataset, err := godal.Open(path)
	if err != nil {
		a.log.Error(err)
		return
	}
	bands := createBandsMap(a.opts.Bands.Value())

	b := ds.New(dataset, bands, a.opts.NoData, a.getProj())

	if a.opts.IsMultiBand {
		if err = b.SetColorModel(len(bands)); err != nil {
			a.log.Error(err)
			return
		}
	}
	if a.opts.Equator != "" {
		equator.New(b, a.opts, a.log, nil).Generate()
	}
	if err = b.Prepare(); err != nil {
		a.log.Errorf("Ошибка подготовки датасета: %v", err)
		return
	}

	if a.opts.MakeContours {
		vector.New(b, a.opts, a.log, nil).Generate()
	}
	raster.New(b, a.opts, a.log, nil).Generate()

	if err = b.Clear(); err != nil {
		a.log.Errorf("Ошибка при очистке датасета: %v", err)
		return
	}
}

func createBandsMap(bands []int) map[int]int {
	m := make(map[int]int, len(bands))
	for i, b := range bands {
		m[i+1] = b
	}
	return m
}
