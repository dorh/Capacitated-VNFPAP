package types

type Graph interface {
	Edges() []Edge
	Clients() []Client
	Servers() []Server
	MaxColor() int
}
