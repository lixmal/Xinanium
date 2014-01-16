package renderer

import (
    "time"
    sf "bitbucket.org/krepa098/gosfml2"
//    "fmt"
    "../config"
)

type Duration float64

//var Elapsed =  make(chan Duration)
var Elapsed Duration

func Render(window *sf.RenderWindow, toDraw []sf.Drawer, text *sf.Text) {
    states := sf.RenderStatesDefault()
    //window.SetFramerateLimit(60)
    window.SetVSyncEnabled(true)
    bgcolor := sf.ColorBlack()

    for start := time.Now(); window.IsOpen(); start = time.Now() {
        window.Clear(bgcolor) // check if needed

        for _, v := range toDraw {
            //channel or lock
            window.Draw(v, states)
        }

        for _, entity := range config.Monsters {
            window.Draw(entity.GetSprite(), states)
        }
        for _, entity := range config.Players {
            window.Draw(entity.GetSprite(), states)
        }

       // text.SetString(fmt.Sprintf("%.0f fps", 1/ float64(Duration(time.Since(start)) / Duration(time.Second))))

        window.Draw(text, states)
/*        select {
            case Elapsed <- Duration(time.Since(start)) / Duration(time.Second):
            default:
        }
        */

        window.Display()
        Elapsed = Duration(time.Since(start)) / Duration(time.Second)
    }
}
