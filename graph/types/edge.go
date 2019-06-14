package types

type Edge interface {
	Client() Client
	Server() Server
	Edge() int
}
