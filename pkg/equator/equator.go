package equator

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/sirupsen/logrus"

	"grib-tiler/pkg/ds"
	"grib-tiler/pkg/opts"
	"grib-tiler/pkg/utils"
)

const (
	northern = "northern"
	southern = "southern"
)

type Equator struct {
	ds     ds.Dataset
	log    *logrus.Entry
	opts   *opts.Opts
	outDir string
	wg     *sync.WaitGroup
}

func New(ds ds.Dataset, opts *opts.Opts, log *logrus.Logger, wg *sync.WaitGroup) *Equator {
	return &Equator{
		ds:     ds,
		log:    log.WithFields(logrus.Fields{"task": "equator", "type": opts.Equator}),
		opts:   opts,
		outDir: filepath.Join(opts.OutputDirectory, ds.Filename(), "equator"),
		wg:     wg,
	}
}

func (e *Equator) Generate() {
	//defer e.wg.Done()
	e.log.Debug("Генерация контура...")

	ex := e.ds.Bounds()
	filename := filepath.Join(e.outDir, e.opts.Equator+".geojson")
	if err := utils.CreateDir(e.outDir); err != nil {
		e.log.Errorf("Ошибка при создании директории: %v", err)
		return
	}
	switch e.opts.Equator {
	case northern:
		ex[1] = 0.0
		ex[3] = float64(int(ex[3]))
	case southern:
		ex[1] = float64(int(ex[1]))
		ex[3] = 0.0
	}
	// Create a bounding box geometry
	geometry := orb.Bound{
		Min: orb.Point{ex[0], ex[1]},
		Max: orb.Point{ex[2], ex[3]},
	}

	// Create a feature with the bounding box geometry
	feature := geojson.NewFeature(geometry)

	// Create a feature collection with the single feature
	featureCollection := geojson.NewFeatureCollection()
	featureCollection.Append(feature)

	// Convert the feature collection to GeoJSON string
	extentGeojson, err := json.Marshal(featureCollection)
	if err != nil {
		e.log.Errorf("Ошибка маршалинга контура: %v", err)
		return
	}
	e.log.Debug("Генерация контура...ОК")
	e.log.Debug("Запись файла контура...")
	err = os.WriteFile(filename, extentGeojson, os.ModePerm)
	if err != nil {
		e.log.Errorf("Ошибка записи файла контура: %v", err)
		return
	}
	e.ds.SetCutLine(filename)
	e.log.Debug("Запись файла контура...ОК")
}
