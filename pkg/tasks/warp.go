package tasks

import (
	"fmt"

	"github.com/airbusgeo/godal"
)

type Task struct {
	Ds *godal.Dataset
}

func (t *WarpTask) switches() []string {
	var switches []string

	if t.OutputFormat != "" {
		switches = append(switches, "-of", t.OutputFormat)
	}
	if t.TargetExtent != [4]float64{0, 0, 0, 0} {
		switches = append(
			switches, "-te",
			fmt.Sprintf("%f", t.TargetExtent[0]),
			fmt.Sprintf("%f", t.TargetExtent[1]),
			fmt.Sprintf("%f", t.TargetExtent[2]),
			fmt.Sprintf("%f", t.TargetExtent[3]))
	}
	if t.TargetExtentSRS != "" {
		switches = append(switches, "-te_srs", t.TargetExtentSRS)
	}
	if t.OutputType != "" {
		switches = append(switches, "-ot", t.OutputType)
	}
	if t.SourceSRS != "" {
		switches = append(switches, "-s_srs", t.SourceSRS)
	}
	if t.TargetSRS != "" {
		switches = append(switches, "-t_srs", t.TargetSRS)
	}
	if t.Cutline != "" {
		switches = append(switches, "-cutline", t.Cutline)
	}
	if t.Refine != "" {
		switches = append(switches, "-refine_gcps", "tolerance")
	}

	switches = append(switches, "-r", "bilinear")

	return switches
}

type WarpTask struct {
	Task
	OutputFile      string
	OutputFormat    string
	OutputType      string
	SourceSRS       string
	TargetSRS       string
	TargetExtent    [4]float64
	TargetExtentSRS string
	Cutline         string
	Refine          string
}

func (t *WarpTask) Warp(opts ...godal.DatasetWarpOption) (*godal.Dataset, error) {
	return t.Ds.Warp(t.OutputFile, t.switches(), opts...)
}
