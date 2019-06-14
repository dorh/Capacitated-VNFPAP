package impl

import (
	"Capacitated/graph/types"
	"math"
)

type LocationVertex interface {
	types.Location
	types.Vertex
}

func NewLocationVertex(vertex_id int, longitude float64, latitude float64) LocationVertex {
	return &locationVertex{vertex_id, longitude, latitude}
}

type locationVertex struct {
	vertex_id int
	longitude float64
	latitude  float64
}

func (lv *locationVertex) Vertex() int {
	return lv.vertex_id
}

func (lv *locationVertex) Longitude() float64 {
	return lv.longitude
}

func (lv *locationVertex) Latitude() float64 {
	return lv.latitude
}

func (lv *locationVertex) Distance(l types.Location) float64 {
	return math.Sqrt(math.Pow(lv.latitude-l.Latitude(), 2.0) + math.Pow(lv.longitude-l.Longitude(), 2.0))
}
