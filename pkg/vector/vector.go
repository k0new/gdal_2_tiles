package vector

import (
	"math/rand"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/airbusgeo/godal"
	"github.com/sirupsen/logrus"
	"github.com/twpayne/go-geos"

	"grib-tiler/pkg/ds"
	"grib-tiler/pkg/opts"
	"grib-tiler/pkg/utils"
)

type Vector struct {
	ds     ds.Dataset
	log    *logrus.Entry
	opts   *opts.Opts
	wg     *sync.WaitGroup
	outDir string
}

func New(ds ds.Dataset, opts *opts.Opts, log *logrus.Logger, wg *sync.WaitGroup) *Vector {
	return &Vector{
		ds:     ds,
		log:    log.WithFields(logrus.Fields{"task": "isoline", "crs": opts.OutCRS}),
		opts:   opts,
		outDir: filepath.Join(opts.OutputDirectory, ds.BaseFilename(), "isoline", ds.Projection().Name()),
		wg:     wg,
	}
}

func (v *Vector) Generate() {
	bands := v.ds.GetBands()
	if err := utils.CreateDir(v.outDir); err != nil {
		v.log.Errorf("Ошибка при создании директории: %v", err)
		return
	}
	rand.NewSource(123132)
	for i, b := range bands {
		tgt, err := v.ds.NewVector(strconv.Itoa(rand.Int())+".gpkg", godal.GeoPackage)
		if err != nil {
			v.log.Errorf("Ошибка при создании векторного датасета: %v", err)
			return
		}
		c := contour{
			parent:    v.ds.Get(),
			ds:        tgt.Get(),
			elevation: v.opts.Elevation,
			tolerance: v.opts.SimplifyEpsilon,
		}
		if err = c.createContourLayer(); err != nil {
			v.log.Error(err)
			return
		}
		if err = c.contour(i - 1); err != nil {
			v.log.Error(err)
			return
		}

		if v.opts.SimplifyEpsilon != 0 {
			if err = c.simplify(); err != nil {
				v.log.Error(err)
				return
			}
		}
		conOut := filepath.Join(v.outDir, strconv.Itoa(b)+"_"+v.opts.OutCRS+".geojson")
		if err = c.save(conOut); err != nil {
			v.log.Error(err)
			return
		}
		err = tgt.Clear()
		if err != nil {
			v.log.Errorf("Ошибка при очистке датасета: %v", err)
			return
		}
	}

	v.log.Debug("Генерация изолиний...ОК")
}

func (v *Vector) godalSimplify(l godal.Layer) {
	cnt, err := l.FeatureCount()
	if err != nil {
		v.log.Errorf("Ошибка при получении кол-ва Feature: %v", err)
		return
	}
	for f := 0; f < cnt; f++ {
		feat := l.NextFeature()
		geom := feat.Geometry()
		simple, err := geom.Simplify(v.opts.SimplifyEpsilon)
		if err != nil {
			v.log.Errorf("Ошибка при упрощении геометрии: %v", err)
			return
		}

		err = feat.SetGeometry(simple)
		if err != nil {
			v.log.Errorf("Ошибка при обновлении геометрии: %v", err)
			return
		}
		err = l.UpdateFeature(feat)
		if err != nil {
			v.log.Errorf("Ошибка при обновлении Feature: %v", err)
			return
		}
	}
}

func (v *Vector) simplif(l godal.Layer) {
	cnt, err := l.FeatureCount()
	if err != nil {
		v.log.Errorf("Ошибка при получении кол-ва Feature: %v", err)
		return
	}
	for f := 0; f < cnt; f++ {
		feat := l.NextFeature()
		geomJSON, err := feat.Geometry().GeoJSON()
		if err != nil {
			v.log.Errorf("Ошибка при получении геометрии: %v", err)
			return
		}
		geom, err := geos.NewGeomFromGeoJSON(geomJSON)
		if err != nil {
			v.log.Errorf("Ошибка при создании GEOS-геометрии: %v", err)
			return
		}
		simplified := geom.Simplify(v.opts.SimplifyEpsilon)
		godalGeom, err := godal.NewGeometryFromGeoJSON(simplified.ToGeoJSON(0))
		if err != nil {
			v.log.Errorf("Ошибка при создании GODAL-геометрии: %v", err)
			return
		}
		err = feat.SetGeometry(godalGeom)
		if err != nil {
			v.log.Errorf("Ошибка при обновлении геометрии: %v", err)
			return
		}
		err = l.UpdateFeature(feat)
		if err != nil {
			v.log.Errorf("Ошибка при обновлении Feature: %v", err)
			return
		}
	}
}
