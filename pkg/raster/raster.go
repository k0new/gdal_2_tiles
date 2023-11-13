package raster

import (
	"fmt"
	"math"
	"math/rand"
	"path/filepath"
	"sync"

	"github.com/sirupsen/logrus"

	"grib-tiler/pkg/ds"
	"grib-tiler/pkg/opts"
	"grib-tiler/pkg/utils"
)

type Raster struct {
	ds     ds.Dataset
	log    *logrus.Entry
	opts   *opts.Opts
	wg     *sync.WaitGroup
	outDir string
}

func New(ds ds.Dataset, opts *opts.Opts, log *logrus.Logger, wg *sync.WaitGroup) *Raster {
	return &Raster{
		ds:     ds,
		log:    log.WithFields(logrus.Fields{"task": "tiler"}),
		opts:   opts,
		wg:     wg,
		outDir: filepath.Join(opts.OutputDirectory, ds.BaseFilename(), "tiles", ds.Projection().Name()),
	}
}

func (r *Raster) Generate() {
	r.log.Debug("Генерация тайлов...")
	if r.opts.Cutline != "" {
		r.ds.SetCutLine(r.opts.Cutline)
	}
	if err := utils.CreateDir(r.outDir); err != nil {
		r.log.Errorf("Ошибка при создании директории: %v", err)
		return
	}
	r.log.Debugf("Количество каналов - %d", len(r.ds.GetBands().Keys()))
	var tiles []tile
	if r.opts.IsMultiBand {
		r.log.Debug("Текущий режим - мультиканальный")
		bands := r.ds.GetBands().Keys()
		if r.ds.ColorModel() == ds.RG {
			bands = append(bands, bands[len(bands)-1])
		}
		tiles = r.prepareTiles(bands...)
	} else {
		r.log.Debug("Текущий режим - одноканальный")
		for _, b := range r.ds.GetBands().Keys() {
			tiles = r.prepareTiles(b)
		}
	}

	t := tiler{
		log:   r.log,
		mutex: sync.RWMutex{},
		wg:    sync.WaitGroup{},
		d:     driver(r.opts.Format),
	}

	t.renderAll(tiles)

	r.log.Debug("Генерация тайлов...ОК")
}

func (r *Raster) prepareTiles(bands ...int) []tile {
	var tiles []tile
	rand.NewSource(1314512)

	warp := r.ds
	//zoom := r.getZoom()
	var tileBase string
	if len(bands) == 1 {
		tileBase = fmt.Sprintf("%s/%d", r.outDir, r.ds.GetBands()[bands[0]])
	} else {
		tileBase = fmt.Sprintf("%s/%d", r.outDir, rand.Int())
	}

	for _, z := range r.opts.Zoom {
		r.log.Debug("Вычисление размеров и оффсетов...")

		outSizeX, outSizeY := r.opts.TileSize, r.opts.TileSize
		if z == full {
			outSizeX, outSizeY = warp.Size()
			r.log.Debugf("Выставлен размера тайла %dx%d для увеличения 0", outSizeX, outSizeY)
		}
		tileX, tileY := calcTileSize(warp, z)
		dsSizeX, dsSizeY := warp.Size()
		numXTiles := int(math.Ceil(float64(dsSizeX) / float64(tileX)))
		numYTiles := int(math.Ceil(float64(dsSizeY) / float64(tileY)))
		for x := 0; x < numXTiles; x++ {
			for y := 0; y < numYTiles; y++ {
				xOffset := x * tileX
				yOffset := y * tileY
				tileFilename := fmt.Sprintf("%s/%d/%d/%d", tileBase, z, x, y)

				tX, tY := reCalcTileSize(warp, xOffset, tileX, yOffset, tileY)
				if tX <= 5 || tY <= 5 {
					continue
				}
				t := tile{
					tileSrc:    warp,
					fn:         tileFilename,
					xOffset:    xOffset,
					yOffset:    yOffset,
					zoom:       z,
					x:          tX,
					y:          tY,
					outSizeX:   outSizeX,
					outSizeY:   outSizeY,
					bands:      bands,
					nodata:     r.opts.NoData,
					colorModel: r.ds.ColorModel(),
				}
				tiles = append(tiles, t)
			}
		}
	}
	return tiles
}

type driver string

const (
	png  driver = "PNG"
	jpeg driver = "JPEG"
)

func (d driver) string() string {
	return string(d)
}

func (d driver) format() string {
	switch d {
	case png:
		return ".png"
	case jpeg:
		return ".jpg"
	default:
		return ".png"
	}
}

const (
	full         = iota // Full dataset projection.
	fourTiles           // 2x2 dataset projection.
	sixteenTiles        // 4x4 dataset projection.
)

// calcTileSize() calculates tile resolution depends on specified Zoom.
func calcTileSize(ds ds.Dataset, zoom int) (int, int) {
	w, h := ds.Size()
	maxDim := w
	if h > maxDim {
		maxDim = h
	}
	switch zoom {
	case full:
		return w, h
	case fourTiles:
		s := maxDim / 2
		return s, s
	case sixteenTiles:
		s := maxDim / 4
		return s, s
	default:
		return w, h
	}
}

func reCalcTileSize(ds ds.Dataset, xOff, tileX, yOff, tileY int) (int, int) {
	tx, ty := tileX, tileY
	w, h := ds.Size()
	if w-xOff < tileX {
		tx = w - xOff
	}
	if h-yOff < tileY {
		ty = h - yOff
	}
	return tx, ty
}

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}
