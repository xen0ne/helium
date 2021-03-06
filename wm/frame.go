package wm

import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/xen0ne/helium/consts"
	"github.com/xen0ne/helium/grid"
)

// Frame is an interface for a frame because I'm not quite sure what to do
type Frame interface {
	Contains(xproto.Window) bool
	FrameId() xproto.Window
	Focus()
	Unfocus()
	Close()
	Tag() int
	SetTag(int)
	Map()
	Unmap()
	State() consts.State
	UpdateBar()
	ResizeRel(int, consts.Direction)
	Resize(int, int)
	SnapToGrid(g *grid.Grid)
}
