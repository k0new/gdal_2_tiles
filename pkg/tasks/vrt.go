package tasks

import (
	"strconv"

	"github.com/airbusgeo/godal"
)

type BuildVRTTask struct {
	Task
	Bands      []int
	OutputFile string
	OutputType string
}

func (t *BuildVRTTask) switches() []string {
	var switches []string
	if t.OutputType != "" {
		switches = append(switches, "-ot", t.OutputType)
	}
	if t.Bands != nil {
		for _, b := range t.Bands {
			switches = append(switches, "-b", strconv.Itoa(b))
		}
	}
	switches = append(switches, "-r", "bilinear")

	return switches
}

func (t *BuildVRTTask) BuildVRT(opts ...godal.BuildVRTOption) (*godal.Dataset, error) {
	return godal.BuildVRT(t.OutputFile, []string{t.Ds.Description()}, t.switches(), opts...)
}
