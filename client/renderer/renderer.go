package renderer

import (
    "time"
    "github.com/veandco/go-sdl2/sdl"
//    "fmt"
//    "../config"
)

type Duration float64

//var Elapsed =  make(chan Duration)
var Elapsed Duration

//var ToDraw []sf.Drawer

func Render(window *sdl.Window) {
    renderer := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
    defer renderer.Destroy()

    for start := time.Now(); true; start = time.Now() {
        renderer.Clear()
        renderer.SetDrawColor(255, 255, 255, 255)
        /*
        for _, v := range ToDraw {
            //channel or lock
            window.Draw(v, states)
        }
        for _, entity := range config.Monsters {
            window.Draw(entity.GetSprite(), states)
        }
        for _, entity := range config.Players {
            window.Draw(entity.GetSprite(), states)
        }
        */

       // text.SetString(fmt.Sprintf("%.0f fps", 1/ float64(Duration(time.Since(start)) / Duration(time.Second))))

        //window.Draw(text, states)
        /*
        select {
            case Elapsed <- Duration(time.Since(start)) / Duration(time.Second):
            default:
        }
        */

        renderer.Present()
        Elapsed = Duration(time.Since(start)) / Duration(time.Second)
    }
}
