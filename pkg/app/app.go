package app

import (
	"sync"

	"github.com/sirupsen/logrus"

	"grib-tiler/pkg/opts"
)

type App struct {
	log  *logrus.Logger
	opts *opts.Opts
	wg   *sync.WaitGroup
}

func NewApp(opts *opts.Opts, log *logrus.Logger) *App {
	return &App{opts: opts, log: log}
}

func (a *App) Run() {
	a.wg = &sync.WaitGroup{}

	a.runReader(a.opts.InputFile)

	a.log.Info("Завершение...")
}
