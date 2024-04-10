package global

import (
	. "github.com/pspaces/gospace/space"
)

var Ts Space

func InitTupleSpace() {
	Ts = NewSpace("ts")
}
