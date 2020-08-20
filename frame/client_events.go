package frame

import (
	"fmt"
	"log"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/mousebind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xprop"
	"github.com/xen0ne/helium/wm"
)

func (f *Frame) addClientEvents() {
	err := f.client.Listen(xproto.EventMaskPropertyChange |
		xproto.EventMaskStructureNotify)

	if err != nil {
		log.Println(err)
	}
	// tell us if the client is killed
	f.cDestroyNotify().Connect(wm.X, f.client.Id)
	f.cPropertyNotify().Connect(wm.X, f.client.Id)
	// c.cbMapNotify().Connect(wm.X, c.Id())
	// c.cbConfigureRequest().Connect(wm.X, c.Id())
	// c.cbClientMessage().Connect(wm.X, c.Id())

	// add focous on click
	err = mousebind.ButtonPressFun(
		func(X *xgbutil.XUtil, ev xevent.ButtonPressEvent) {
			xproto.AllowEvents(X.Conn(), xproto.AllowReplayPointer, 0)
			f.Focus()
		}).Connect(wm.X, f.client.Id, "1", false, true)
	if err != nil {
		fmt.Println(err)
	}

	mousebind.Drag(wm.X, f.client.Id, f.client.Id, "Mod1-1", true,
		f.moveDragBegin, f.moveDragStep, f.moveDragEnd)
}

func (f *Frame) cPropertyNotify() xevent.PropertyNotifyFun {
	fn := func(X *xgbutil.XUtil, ev xevent.PropertyNotifyEvent) {
		name, err := xprop.AtomName(wm.X, ev.Atom)
		if err != nil {
			log.Printf("Could not get property atom name for '%s' because: %s.", ev, err)
			return
		}
		f.handleProperty(name)
	}

	return xevent.PropertyNotifyFun(fn)
}

func (f *Frame) handleProperty(p string) {
	switch p {
	case "_NET_WM_VISIBLE_NAME":
		fallthrough
	case "_NET_WM_NAME":
		fallthrough
	case "WM_NAME":
		f.UpdateBar()
	default:
		log.Printf("Dont know how to handle property '%s'\n", p)
	}
}

func (f *Frame) cDestroyNotify() xevent.DestroyNotifyFun {
	fn := func(X *xgbutil.XUtil, ev xevent.DestroyNotifyEvent) {
		log.Printf("destroy notify for %x\n", ev.Window)

		// get the last destroy notify
		X.Sync()
		xevent.Read(X, false)

		for i, ee := range xevent.Peek(X) {
			if dn, ok := ee.Event.(xproto.DestroyNotifyEvent); ok {
				if dn.Window == ev.Window {
					log.Printf("another destroy for %x\n", ev.Window)
					xevent.DequeueAt(X, i)
				}
			}
		}

		flast := wm.FoucusQ[0].FrameId() == f.FrameId()
		wm.FoucusQ = wm.RemoveFrame(f, wm.FoucusQ)

		if len(wm.FoucusQ) > 0 && flast {
			wm.FoucusQ[0].Focus()
		}

		f.bar.Destroy()
		f.Destroy()
		wm.ManagedFrames = wm.RemoveFrame(f, wm.ManagedFrames)
	}
	return xevent.DestroyNotifyFun(fn)
}