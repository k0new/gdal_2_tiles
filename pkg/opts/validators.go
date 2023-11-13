package opts

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
)

func (*Opts) transparencyValidator(_ *cli.Context, t uint) error {
	if t < 1 || t > 100 {
		return fmt.Errorf("Значение прозрачности %d не входит в промежуток [0-100]", t)
	}
	return nil
}

func (*Opts) equatorValidator(_ *cli.Context, e string) error {
	if e != "southern" && e != "northern" {
		return fmt.Errorf("%q не входит в доступные варианты образки по экватору - [northern|southern]", e)
	}
	return nil
}

func (*Opts) formatValidator(_ *cli.Context, f string) error {
	if f != "PNG" && f != "JPEG" {
		return fmt.Errorf("%q не входит в доступные форматы изображения - [PNG|JPEG]", f)
	}
	return nil
}

func (o *Opts) zoomsValidator(_ *cli.Context, z string) error {
	zoomS := strings.Split(strings.TrimSpace(z), "-")
	var zoom []int
	if len(zoomS) == 2 {
		min, err := strconv.Atoi(zoomS[0])
		if err != nil {
			return fmt.Errorf("Ошибка конвертации значения увеличения: %v", err)
		}
		if zoomS[0] == "0" {
			min = 0
		}
		max, err := strconv.Atoi(zoomS[1])
		if err != nil {
			return fmt.Errorf("Ошибка конвертации значения увеличения: %v", err)
		}
		zoom = makeRange(min, max)
	} else if len(zoomS) == 1 {
		s, err := strconv.Atoi(zoomS[0])
		if err != nil {
			return fmt.Errorf("Ошибка конвертации значения увеличения: %v", err)
		}
		if s == 0 {
			zoom = []int{0}
		}
		zoom = []int{s}
	}
	o.Zoom = zoom
	return nil
}

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}
