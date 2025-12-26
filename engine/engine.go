// Package engine how mqtt broker react to different messages
package engine

import "github.com/ilgianlu/tagyou/model"

func NewEngine() model.Engine {
	return StandardEngine{}
}
