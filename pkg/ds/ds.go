package ds

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/airbusgeo/godal"

	"grib-tiler/pkg/projection"
	"grib-tiler/pkg/tasks"
)

// Wrapper around *godal.Dataset
type dataset struct {
	*godal.Dataset
	cutline      string
	bandsMap     bandsMap
	noData       string
	baseFilename string
	colorModel   ColorModel
	proj         projection.Projection
}

type Dataset interface {
	Get() *godal.Dataset
	Prepare() error
	Bounds() [4]float64
	Size() (int, int)
	GetGeoTransform() [6]float64
	Filename() string
	BaseFilename() string
	Projection() projection.Projection
	NewVector(name string, driver godal.DriverName) (Dataset, error)
	GetBands() bandsMap
	SetCutLine(file string)
	GetCutLine() string
	SetColorModel(bandsCnt int) error
	ColorModel() ColorModel
	Translate(task tasks.TranslateTask, opts ...godal.DatasetTranslateOption) (Dataset, error)
	Warp(task tasks.WarpTask, opts ...godal.DatasetWarpOption) (Dataset, error)
	VRT(task tasks.BuildVRTTask, opts ...godal.BuildVRTOption) (Dataset, error)
	Clear() error
}

func New(ds *godal.Dataset, bands map[int]int, noData string, proj projection.Projection) Dataset {
	return &dataset{Dataset: ds, bandsMap: bands, noData: noData, baseFilename: filepath.Base(ds.Description()), proj: proj}
}

type ColorModel string

const (
	A    ColorModel = "A"    // Alfa channel only.
	RG   ColorModel = "RG"   // Red Green.
	RGB  ColorModel = "RGB"  // Red Green Blue.
	RGBA ColorModel = "RGBA" // Red Green Blue Alfa.
)

func (ds *dataset) ColorModel() ColorModel {
	return ds.colorModel
}

func (ds *dataset) Projection() projection.Projection {
	return ds.proj
}

func (ds *dataset) SetColorModel(bandsCnt int) error {
	switch bandsCnt {
	case 1:
		ds.colorModel = A
	case 2:
		ds.colorModel = RG
	case 3:
		ds.colorModel = RGB
	case 4:
		ds.colorModel = RGBA
	default:
		return fmt.Errorf("Цветовая модель для %d каналов не распознана.", bandsCnt)
	}
	return nil
}

// Get returns *godal.Dataset instance.
func (ds *dataset) Get() *godal.Dataset {
	return ds.Dataset
}

func (ds *dataset) GetBands() bandsMap {
	return ds.bandsMap
}

func (ds *dataset) Bounds() [4]float64 {
	b, _ := ds.Get().Bounds()
	return b
}

// Size returns dataset's width and height.
func (ds *dataset) Size() (int, int) {
	return ds.Structure().SizeX, ds.Structure().SizeY
}

func (ds *dataset) GetGeoTransform() [6]float64 {
	gt, _ := ds.GeoTransform()
	return gt
}

func (ds *dataset) Filename() string {
	return filepath.Base(ds.Description())
}

func (ds *dataset) Translate(task tasks.TranslateTask, opts ...godal.DatasetTranslateOption) (Dataset, error) {
	task.Ds = ds.Dataset
	d, err := task.Translate(opts...)
	if err != nil {
		return nil, err
	}
	return &dataset{
		Dataset:  d,
		cutline:  ds.cutline,
		bandsMap: ds.bandsMap,
	}, nil
}

func (ds *dataset) Warp(task tasks.WarpTask, opts ...godal.DatasetWarpOption) (Dataset, error) {
	task.Ds = ds.Dataset
	d, err := task.Warp(opts...)
	if err != nil {
		return nil, err
	}
	return &dataset{
		Dataset:  d,
		cutline:  ds.cutline,
		bandsMap: ds.bandsMap,
	}, nil
}

func (ds *dataset) VRT(task tasks.BuildVRTTask, opts ...godal.BuildVRTOption) (Dataset, error) {
	task.Ds = ds.Dataset
	d, err := task.BuildVRT(opts...)
	if err != nil {
		return nil, err
	}
	return &dataset{
		Dataset:  d,
		cutline:  ds.cutline,
		bandsMap: ds.bandsMap,
	}, nil
}

func (ds *dataset) SetCutLine(file string) {
	ds.cutline = file
}

func (ds *dataset) GetCutLine() string {
	return ds.cutline
}

func (ds *dataset) set(t Dataset) {
	ds.Dataset = t.Get()
}

// Prepare warps selected bands to a new dataset with EPSG:4326.
func (ds *dataset) Prepare() error {
	bands := ds.bandsMap.Values()
	if ds.colorModel == RG {
		bands = append(bands, bands[len(bands)-1])
	}
	tr, err := ds.Translate(tasks.TranslateTask{
		Bands:        bands,
		OutputFile:   ds.Description() + "_tr.tif",
		OutputFormat: "GTiff",
		OutputSRS:    "EPSG:4326",
		NoData:       "0",
		Scale:        true,
		AULLR:        [4]float64{-180, 90, 180, -90},
	})
	if err != nil {
		return err
	}
	warp, err := tr.Warp(tasks.WarpTask{
		OutputFile:   ds.Description() + ds.proj.String() + ".tif",
		SourceSRS:    tr.Get().Projection(),
		TargetSRS:    ds.proj.String(),
		OutputType:   "Byte",
		TargetExtent: ds.proj.Extent(),
		Cutline:      ds.cutline,
	})
	if err != nil {
		return err
	}
	ds.set(warp)
	if err = tr.Clear(); err != nil {
		return err
	}
	return nil
}

func (ds *dataset) BaseFilename() string {
	return ds.baseFilename
}

func (ds *dataset) Clear() error {
	path := ds.Description()
	if err := ds.Close(); err != nil {
		return err
	}
	return os.Remove(path)
}

type bandsMap map[int]int

func (m bandsMap) Keys() []int {
	var k []int
	for i := range m {
		k = append(k, i)
	}
	return k
}

func (m bandsMap) Values() []int {
	var v []int
	for _, k := range m {
		v = append(v, k)
	}

	return v
}

func (ds *dataset) NewVector(name string, driver godal.DriverName) (Dataset, error) {
	d, err := godal.CreateVector(driver, name)
	if err != nil {
		return nil, err
	}
	return &dataset{
		Dataset:      d,
		cutline:      ds.cutline,
		bandsMap:     ds.bandsMap,
		noData:       ds.noData,
		baseFilename: ds.baseFilename,
		colorModel:   ds.colorModel,
		proj:         ds.proj,
	}, nil
}
