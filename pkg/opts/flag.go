package opts

import "github.com/urfave/cli/v2"

const (
	tilerCat   = "Tiler"
	contourCat = "Contour"
)

type Opts struct {
	InputFile       string
	OutputDirectory string
	Transparency    uint
	Equator         string
	IsMultiBand     bool
	TileSize        int
	OutCRS          string
	NoData          string
	Zoom            []int
	Bands           cli.IntSlice
	MakeContours    bool
	Elevation       float64
	SimplifyEpsilon float64
	Cutline         string
	Format          string
}

func (o *Opts) transparency() *cli.UintFlag {
	return &cli.UintFlag{
		Name:        "transparency",
		Category:    tilerCat,
		Usage:       "Процент прозрачности тайлов. [0-100].",
		Value:       0,
		Destination: &o.Transparency,
		Action:      o.transparencyValidator,
	}
}

func (o *Opts) equator() cli.Flag {
	return &cli.StringFlag{
		Name:        "equator",
		Usage:       "Обрезка тайлов/изолиний по экватору. [northern|southern].",
		Destination: &o.Equator,
		Action:      o.equatorValidator,
	}
}

func (o *Opts) multiband() cli.Flag {
	return &cli.BoolFlag{
		Name:        "multiband",
		Category:    tilerCat,
		Aliases:     []string{"m"},
		Usage:       "Генерация тайлов в многоканальном режиме.",
		Destination: &o.IsMultiBand,
	}
}
func (o *Opts) tilesize() cli.Flag {
	return &cli.IntFlag{
		Name:        "tilesize",
		Category:    tilerCat,
		Usage:       "Размер выходных тайлов.",
		Value:       256,
		Destination: &o.TileSize,
	}
}
func (o *Opts) outCRS() cli.Flag {
	return &cli.StringFlag{
		Name:        "out-crs",
		Usage:       "Выходная система координат.",
		Value:       "EPSG:3857",
		Destination: &o.OutCRS,
	}
}

func (o *Opts) nodata() cli.Flag {
	return &cli.StringFlag{
		Name:        "nodata",
		Usage:       "Заполнение nodata пикселей.",
		Value:       "none",
		Destination: &o.NoData,
	}
}

func (o *Opts) zooms() cli.Flag {
	return &cli.StringFlag{
		Name:     "zooms",
		Category: tilerCat,
		Aliases:  []string{"z"},
		Usage:    "Увеличение тайлов в формате 0-2 или 0.",
		Value:    "0-2",
		Action:   o.zoomsValidator,
	}
}

func (o *Opts) bands() cli.Flag {
	return &cli.IntSliceFlag{
		Name:        "bands",
		Usage:       "Генерация тайлов/изолиний для выбранных бэндов в формате 1,2,3 или 1.",
		Required:    true,
		Aliases:     []string{"b"},
		Destination: &o.Bands,
	}
}

func (o *Opts) contours() cli.Flag {
	return &cli.BoolFlag{
		Name:        "contours",
		Category:    contourCat,
		Usage:       "Включение генерации изолиний.",
		Aliases:     []string{"c"},
		Destination: &o.MakeContours,
	}
}

func (o *Opts) cutline() cli.Flag {
	return &cli.StringFlag{
		Name:        "cutline",
		Usage:       "Путь к файлу обрезки растра.",
		TakesFile:   true,
		Destination: &o.Cutline,
	}
}
func (o *Opts) format() cli.Flag {
	return &cli.StringFlag{
		Name:        "format",
		Usage:       "Формат выходных тайлов. [PNG|JPEG].",
		Aliases:     []string{"f"},
		Destination: &o.Format,
		Action:      o.formatValidator,
	}
}

func (o *Opts) elevation() cli.Flag {
	return &cli.Float64Flag{
		Name:        "contours-elev",
		Category:    contourCat,
		Usage:       "Интервал между изолиниями (в высоте).",
		Aliases:     []string{"e"},
		Value:       10.0,
		Destination: &o.Elevation,
	}
}

func (o *Opts) simplify() cli.Flag {
	return &cli.Float64Flag{
		Name:        "contours-simplify",
		Category:    contourCat,
		Usage:       "Tolerance для уменьшения кол-ва точек изолиний.",
		Aliases:     []string{"s"},
		Value:       0.0,
		Destination: &o.SimplifyEpsilon,
	}
}

func (o *Opts) Flags() []cli.Flag {
	return []cli.Flag{
		o.transparency(),
		o.equator(),
		o.multiband(),
		o.tilesize(),
		o.outCRS(),
		o.nodata(),
		o.zooms(),
		o.bands(),
		o.contours(),
		o.elevation(),
		o.simplify(),
		o.cutline(),
		o.format(),
	}
}
