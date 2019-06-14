package types

type Vertex interface {
	Vertex() int
}

type Location interface {
	Longitude() float64
	Latitude() float64
	Distance(l Location) float64
}

type Node interface{
	Location
	Vertex
}

type Client interface {
	Node
	Color() int
}

type Server interface {
	Location
	Vertex
	Colors() []int
	Capacity() int
}
