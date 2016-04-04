package main

import (
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/gl"

	_ "github.com/scritch007/go-tools"
)

func main() {

	app.Main(func(a app.App) {
		var glctx gl.Context
		var sz size.Event
		init := false
		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					glctx, _ = e.DrawContext.(gl.Context)
					on_surface_created(glctx)
					a.Send(paint.Event{})
				case lifecycle.CrossOff:
					onStop(glctx)
					glctx = nil
				}
			case size.Event:
				if glctx != nil {
					tmp := on_surface_changed(glctx, &e)
					if tmp {
						sz = e
						init = true
					}
				}
			case paint.Event:
				if !init || e.External {
					// As we are actively painting as fast as
					// we can (usually 60 FPS), skip any paint
					// events sent by the system.
					continue
				}

				on_draw_frame(glctx, &sz)
				a.Publish()
				// Drive the animation by preparing to paint the next frame
				// after this one is shown.
				a.Send(paint.Event{})
			case touch.Event:
				normalized_x := (e.X/float32(sz.WidthPx))*2.0 - 1.0
				normalized_y := -((e.Y/float32(sz.HeightPx))*2.0 - 1.0)
				//normalized_x := e.X
				//normalized_y := e.Y
				//tools.LOG_DEBUG.Printf("Received %f %f\n", normalized_x, normalized_y)
				if e.Type == touch.TypeBegin {
					on_touch_press(glctx, normalized_x, normalized_y)
				} else if e.Type == touch.TypeMove {
					on_touch_drag(glctx, normalized_x, normalized_y)
				}
			}
		}
	})
}
