package engine

import "github.com/ilgianlu/tagyou/model"

func NewEngine() model.Engine {
	return StandardEngine{}
}
