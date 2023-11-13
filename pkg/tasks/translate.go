package tasks

import (
	"fmt"
	"strconv"

	"github.com/airbusgeo/godal"
)

type TranslateTask struct {
	Task
	ColorModel   string
	Stat         []TileStat
	Bands        []int
	OutputFile   string
	OutputFormat string
	OutputType   string
	OutputSRS    string
	OutSize      [2]int
	SrcWin       [4]int
	ProjWin      [4]float64
	ProjWinSrs   string
	NoData       string
	Scale        bool
	Unscale      bool
	Photometric  string
	AULLR        [4]float64
	scaleFunc    func() []string
}

type TileStat struct {
	Min float64
	Max float64
}

func (t *TranslateTask) scaleRG() []string {
	var switches []string
	if t.Stat != nil {
		for i, b := range t.Bands {
			switches = append(switches, "-scale_"+strconv.Itoa(b), fmt.Sprintf("%.0f", t.Stat[i].Min), fmt.Sprintf("%.0f", t.Stat[i].Max), "0", "255")
			if i == 2 {
				continue
			}
		}
		switches = append(switches, "-scale_3", fmt.Sprintf("%.0f", t.Stat[2].Min), fmt.Sprintf("%.0f", t.Stat[2].Max), "0", "0")
	}
	return switches
}

func (t *TranslateTask) scaleRGBa() []string {
	var switches []string
	if t.Stat != nil {
		for i, b := range t.Bands {
			switches = append(switches, "-scale_"+strconv.Itoa(b), fmt.Sprintf("%.0f", t.Stat[i].Min), fmt.Sprintf("%.0f", t.Stat[i].Max), "0", "255")
		}
	}
	return switches
}

func (t *TranslateTask) scaleAlpha() []string {
	var switches []string
	if t.Scale {
		switches = append(switches, "-scale")
	}
	return switches
}

func (t *TranslateTask) getScaleFunc() func() []string {
	switch t.ColorModel {
	case "", "A":
		return t.scaleAlpha
	case "RG":
		return t.scaleRG
	case "RGBA":
		return t.scaleRGBa
	default:
		return t.scaleAlpha
	}
}

func (t *TranslateTask) switches() []string {
	var switches []string
	scale := t.getScaleFunc()
	switches = append(switches, scale()...)
	if t.AULLR != [4]float64{0, 0, 0, 0} {
		switches = append(switches, "-a_ullr", fmt.Sprintf("%f", t.AULLR[0]), fmt.Sprintf("%f", t.AULLR[1]), fmt.Sprintf("%f", t.AULLR[2]), fmt.Sprintf("%f", t.AULLR[3]))
	}
	if t.ProjWin != [4]float64{0, 0, 0, 0} {
		switches = append(switches, "-projwin", fmt.Sprintf("%f", t.ProjWin[0]), fmt.Sprintf("%f", t.ProjWin[1]), fmt.Sprintf("%f", t.ProjWin[2]), fmt.Sprintf("%f", t.ProjWin[3]))
	}
	if t.ProjWinSrs != "" {
		switches = append(switches, "-projwin_srs", t.ProjWinSrs)
	}
	if t.Photometric != "" {
		switches = append(switches, "-co", "PHOTOMETRIC="+t.Photometric)
	}
	if t.OutputFormat != "" {
		switches = append(switches, "-of", t.OutputFormat)
	}
	if t.OutputType != "" {
		switches = append(switches, "-ot", t.OutputType)
	}
	if t.Bands != nil {
		for _, b := range t.Bands {
			switches = append(switches, "-b", strconv.Itoa(b))
		}
	}
	if t.OutSize != [2]int{0, 0} {
		switches = append(switches, "-outsize", strconv.Itoa(t.OutSize[0]), strconv.Itoa(t.OutSize[1]))

	}
	if t.SrcWin != [4]int{0, 0, 0, 0} {
		switches = append(switches, "-srcwin", strconv.Itoa(t.SrcWin[0]), strconv.Itoa(t.SrcWin[1]), strconv.Itoa(t.SrcWin[2]), strconv.Itoa(t.SrcWin[3]))

	}
	if t.OutputSRS != "" {
		switches = append(switches, "-a_srs", t.OutputSRS)
	}
	if t.NoData != "" {
		switches = append(switches, "-a_nodata", t.NoData)
	}
	if t.Unscale {
		switches = append(switches, "-unscale")
	}
	switches = append(switches, "-r", "bilinear")
	return switches
}

func (t *TranslateTask) Translate(opts ...godal.DatasetTranslateOption) (*godal.Dataset, error) {
	return t.Ds.Translate(t.OutputFile, t.switches(), opts...)
}
