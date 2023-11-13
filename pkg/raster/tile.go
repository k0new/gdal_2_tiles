package raster

import (
	"path/filepath"
	"sync"

	"github.com/sirupsen/logrus"

	"grib-tiler/pkg/ds"
	"grib-tiler/pkg/tasks"
	"grib-tiler/pkg/utils"
)

type tile struct {
	tileSrc    ds.Dataset
	bands      []int
	fn         string
	xOffset    int
	yOffset    int
	zoom       int
	x          int // tile width in dataset projection.
	y          int // tile height in dataset projection.
	outSizeX   int // tile out width in px.
	outSizeY   int // tile out height in px.
	nodata     string
	srs        string
	stat       []tasks.TileStat
	colorModel ds.ColorModel
}

// TODO: concurrent error handling.
// tiler is a concurrent safe instance for tile generation.
type tiler struct {
	log   *logrus.Entry
	mutex sync.RWMutex
	wg    sync.WaitGroup
	d     driver
}

// renderAll concurrently creates specified tiles.
func (t *tiler) renderAll(tiles []tile) {
	t.wg.Add(len(tiles))
	for _, v := range tiles {
		go t.render(v)
	}
	t.wg.Wait()
}

// render creates a single tile.
func (t *tiler) render(tile tile) {
	defer t.wg.Done()

	fn := tile.fn + t.d.format()
	// lock for dir creation
	t.mutex.Lock()
	if err := utils.CreateDir(filepath.Dir(fn)); err != nil {
		t.log.Errorf("Ошибка при подготовке директории: %v", err)
		return
	}
	t.mutex.Unlock()
	if err := t.computeStat(&tile); err != nil {
		t.log.Errorf("Ошибка при вычислении min/max значений: %v", err)
		return
	}
	// lock for image creation
	t.mutex.Lock()
	_, err := tile.tileSrc.Translate(tasks.TranslateTask{
		ColorModel:   string(tile.colorModel),
		Bands:        tile.bands,
		OutputFile:   fn,
		OutputFormat: t.d.string(),
		OutputType:   "Byte",
		OutSize:      [2]int{tile.outSizeX, tile.outSizeY},
		SrcWin:       [4]int{tile.xOffset, tile.yOffset, tile.x, tile.y},
		NoData:       tile.nodata,
		Scale:        true,
		OutputSRS:    tile.srs,
		Stat:         tile.stat,
	})
	t.mutex.Unlock()
	if err != nil {
		t.log.Errorf("Ошибка при записи тайла: %v", err)
		return
	}
}

// computeStat retrieves tile min/max statistics.
func (t *tiler) computeStat(tile *tile) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	for _, b := range tile.bands {
		s, err := tile.tileSrc.Get().Bands()[b-1].ComputeStatistics()
		if err != nil {
			return err
		}
		tile.stat = append(tile.stat, tasks.TileStat{
			Min: s.Min,
			Max: s.Max,
		})
	}
	if tile.tileSrc.ColorModel() == ds.RG {
		s, err := tile.tileSrc.Get().Bands()[2].ComputeStatistics()
		if err != nil {
			return err
		}
		tile.stat = append(tile.stat, tasks.TileStat{
			Min: s.Min,
			Max: s.Max,
		})
	}
	return nil
}
