package impl

import "Capacitated/graph/types"

func NewClient(vertex_id int, color int, longitude float64, latitude float64) types.Client {
	lv := NewLocationVertex(vertex_id, longitude, latitude)
	return &client{lv, color}
}

func NewServer(vertex_id int, longitude float64, latitude float64, capacity int) types.Server {
	lv := NewLocationVertex(vertex_id, longitude, latitude)
	return &server {lv, []int{}, capacity}
}

type client struct {
	LocationVertex
	color int
}

func (c *client) Color() int {
	return c.color
}

type server struct {
	LocationVertex
	colors []int
	capacity int
}


func (s *server) Colors() []int {
	return s.colors
}

func (s *server) Capacity() int{
	return s.capacity
}
