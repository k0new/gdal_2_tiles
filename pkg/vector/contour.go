package vector

import (
	"fmt"

	"github.com/airbusgeo/godal"
)

type contour struct {
	parent *godal.Dataset

	ds *godal.Dataset

	layer     godal.Layer
	elevation float64
	tolerance float64
}

func (c *contour) createContourLayer() error {
	layer, err := c.ds.CreateLayer(
		"contour", c.parent.SpatialRef(), godal.GTLineString,
		godal.NewFieldDefinition("ID", godal.FTInt),
		godal.NewFieldDefinition("ELEV", godal.FTReal),
	)
	if err != nil {
		return fmt.Errorf("Ошибка при создании слоев: %v", err)
	}
	c.layer = layer
	return nil
}

func (c *contour) contour(bandIdx int) error {
	err := c.parent.Bands()[bandIdx].ContourGenerate(c.layer, c.elevation, 0, 1)
	if err != nil {
		return fmt.Errorf("Ошибка при генерации контура: %v", err)
	}
	return nil
}

func (c *contour) simplify() error {
	cnt, err := c.layer.FeatureCount()
	if err != nil {
		return fmt.Errorf("Ошибка при получении кол-ва Feature: %v", err)
	}
	for f := 0; f < cnt; f++ {
		feat := c.layer.NextFeature()
		geom := feat.Geometry()
		simple, err := geom.Simplify(c.tolerance)
		if err != nil {
			return fmt.Errorf("Ошибка при упрощении геометрии: %v", err)
		}

		err = feat.SetGeometry(simple)
		if err != nil {
			return fmt.Errorf("Ошибка при обновлении геометрии: %v", err)
		}
		err = c.layer.UpdateFeature(feat)
		if err != nil {
			return fmt.Errorf("Ошибка при обновлении Feature: %v", err)
		}
	}
	return nil
}

func (c *contour) save(fn string) error {
	translate, err := c.ds.VectorTranslate(fn, nil, godal.GeoJSON)
	if err != nil {
		return fmt.Errorf("Ошибка при трансляции в финальный датасет: %v", err)
	}
	if err = translate.Close(); err != nil {
		return fmt.Errorf("Ошибка при закрытии датасета: %v", err)
	}

	return nil
}
