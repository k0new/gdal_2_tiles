package projection

type Projection interface {
	Extent() [4]float64
	String() string
	Name() string
}

type EPSG3575 string
type EPSG3857 string
type EPSG4326 string

func New3857() EPSG3857 {
	return "EPSG:3857"
}

func New3575() EPSG3575 {
	return "EPSG:3575"
}

func New4326() EPSG4326 {
	return "EPSG:4326"
}

func (e EPSG3575) String() string {
	return string(e)
}

func (e EPSG3857) String() string {
	return string(e)
}

func (e EPSG4326) String() string {
	return string(e)
}

func (EPSG3575) Extent() [4]float64 {
	return [4]float64{-4886873.23, -4859208.99, 4886873.23, 4859208.99}
}

func (EPSG3857) Extent() [4]float64 {
	return [4]float64{-20037508.34, -20048966.1, 20037508.34, 20048966.1}
}

func (EPSG4326) Extent() [4]float64 {
	return [4]float64{-180, -90, 180, 90}
}

func (EPSG3857) Name() string {
	return "mercator"
}

func (EPSG3575) Name() string {
	return "north_pole"
}

func (EPSG4326) Name() string {
	return "world"
}
