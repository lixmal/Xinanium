package main

// #cgo LDFLAGS: -lX11
// #include <X11/Xlib.h>
import "C"
import (
	sf "bitbucket.org/krepa098/gosfml2"
	"log"
	"runtime"
    "math"
    "fmt"
    "strconv"
    "./config"
    "./monster"
    "./player"
    "./renderer"
    "./event"
    wm "./worldmap"
    _ "net/http/pprof"
    "net/http"
)

const RESOURCESDIR = "resources/"
const SPRITEDIR = "resources/textures/spritesheets/"


type Duration float64


func init() {
	runtime.GOMAXPROCS(runtime.NumCPU() + 1)
	runtime.LockOSThread()

    // X11 multithreading, linux/X11 only
	C.XInitThreads()
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()

}


func i2c(x, y config.Coord) (config.Coord, config.Coord) {
    return (2 * y + x)/2, (2 * y - x)/2
}

func c2i(x, y config.Coord) (config.Coord, config.Coord) {
    return 2 * y - x, (2 * y + x)/2
}

func getTileCoord(x, y config.Coord, h int8) (config.Coord, config.Coord) {
    return config.Coord(math.Floor(float64(x/config.Coord(h)))), config.Coord(math.Floor(float64(y/config.Coord(h)))) // height and width equal
}

func main() {

    var toDraw []sf.Drawer
	window := sf.NewRenderWindow(sf.VideoMode{config.Conf.ScreenWidth, config.Conf.ScreenHeight, config.Conf.BitDepth}, config.Conf.GameTitle, sf.StyleDefault, config.Conf.ContextSettings)
    config.Conf.Window = window

    // create view
    view := sf.NewViewFromRect(sf.FloatRect{0, 0, float32(config.Conf.ScreenWidth) * 0.7, float32(config.Conf.ScreenHeight) * 0.7})




    // read default worldmap and add to toDraw
    worldmap := wm.Read(RESOURCESDIR + "maps/gobmap.dat")
    wm.Current = worldmap
    for x, v := range worldmap.Tiles {
        for y := range v {
            tileType := worldmap.Tiles[x][y]
            var sprite *sf.Sprite
            switch tileType {
            case 0:
                sprite = sf.NewSprite(config.Conf.Rm.Texture(RESOURCESDIR + "textures/tiles/grass.png"))
            case 1:
                sprite = sf.NewSprite(config.Conf.Rm.Texture(RESOURCESDIR + "textures/tiles/dirt.png"))
            case 2:
                sprite = sf.NewSprite(config.Conf.Rm.Texture(RESOURCESDIR + "textures/tiles/water.png"))
            case 3:
                sprite = sf.NewSprite(config.Conf.Rm.Texture(RESOURCESDIR + "textures/tiles/w_br.png"))
            }
            if sprite != nil {
                sprite.SetPosition(sf.Vector2f{float32(x * wm.TILEWIDTH), float32(y * wm.TILEHEIGHT)})
                toDraw = append(toDraw, sprite)
            } else {
                log.Fatal("Unknown tile type: " + fmt.Sprintf("%i", tileType))
            }
        }
    }

    font, err := sf.NewFontFromFile("/usr/share/fonts/truetype/ubuntu-font-family/Ubuntu-B.ttf")
    if err != nil {
        log.Fatal("...")
    }
    text := sf.NewText(font)
    text.SetCharacterSize(12)

    // set window to inactive for OpenGL
    if !window.SetActive(false) {
        log.Fatal("Could not set window OpenGL context to false")
    }

    // start rendering
    go renderer.Render(window, toDraw, text)

    // init Luap
    initLua(&config.Lua.State)

    // start event loop
    initEventQueue()

    config.Conf.GameActive = true


    worldmap.RunScripts()

    var textEntered []rune

/*
    // music
    music, err := sf.NewMusicFromFile(RESOURCESDIR + "sound/test.ogg")
    if err != nil {
        log.Fatal("Could not load sound: ")
    }
    music.SetLoop(true)
    music.Play()
*/

    // default Player
    player1 := player.New("Player", "vik", true)
    player1.SetPosition(100, 200)


    // assign view to main player
    {
        posx, posy := player1.Position()
        view.SetCenter(sf.Vector2f{posx, posy})
        window.SetView(view)
    }

    // spawn one test monster
    for i := 0; i < 1; i++ {
        //mon := monster.New("monster", config.Coord(rand.Intn(1500)), config.Coord(rand.Intn(1500)), 500)
        monster.New("monster").SetPosition(200, 200)
    }






    // game loop
	for window.IsOpen() {
        <-config.GameTicker

        // player moving
        if !config.Conf.TextMode && config.Conf.GameActive {
            var x, y float32
            var pressed bool
            if sf.KeyboardIsKeyPressed(sf.KeyDown) {
                y += 1
                pressed = true
            } else if sf.KeyboardIsKeyPressed(sf.KeyUp) {
                y -= 1
                pressed = true
            }
            if sf.KeyboardIsKeyPressed(sf.KeyRight) {
                x += 1
                pressed = true
            } else if sf.KeyboardIsKeyPressed(sf.KeyLeft) {
                pressed = true
                x -= 1
            }
            if pressed {
                player1.Move(x, y)
            }
        }

        // sfml event loop
        for e := window.PollEvent(); e != nil; e = window.PollEvent() {
            switch e.(type) {
            case sf.EventClosed:
                window.Close()
            case sf.EventLostFocus:
                config.Conf.GameActive = false
                runtime.GC()
            case sf.EventGainedFocus:
                runtime.GC()
                config.Conf.GameActive = true
            case sf.EventTextEntered:
                char := e.(sf.EventTextEntered).Char
                // trigger text entered in any case
                ev := &event.CharPressed{Event: event.New(event.TypeCharPressed), Char: string(char)}
                cancelled := event.Trigger(ev)
                if config.Conf.TextMode && config.Conf.GameActive && !cancelled  {
                    char = rune(ev.Char[0])
                    if strconv.IsPrint(char)  {
                        textEntered = append(textEntered, char)
                    }
                    text.SetString(string(textEntered))
                }
            case sf.EventKeyPressed:
                keyCode := e.(sf.EventKeyPressed).Code
                // TODO: add possibility to manually trigger keys

                // TODO: check other cancelled if no short circuit
                if !event.Trigger(&event.KeyPressed{Event: event.New(event.TypeKeyPressed), Key: uint16(keyCode)}) && !config.Conf.TextMode {
                // add modifier keys
                    switch keyCode {
                    case sf.KeySpace:
                        if config.Conf.GameActive {
                            player1.StopAnimation()
                            player1.Jump()
                        }
                    case sf.KeyTab:
                        config.Conf.Scrolling = !config.Conf.Scrolling && true
                    case sf.KeyH:
                        for _, mon := range config.Monsters {
                            mon.Remove()
                            break
                        }
                    case sf.KeyJ:
                        monster.New("hhh")
                    case sf.KeyX:
                        for _, entity := range config.Monsters {
                            if mon, ok := entity.(*monster.Monster); ok {
                                mon.Run()
                            }
                        }
                    case sf.KeyL:
                        event.Registry = nil
                        event.Registry = make(map[event.Type][]*event.Listener)
                        config.Lua.State.Close()
                        initLua(&config.Lua.State)
                        err := config.Lua.State.DoFile("test.lua")
                        if err != nil {
                            fmt.Println(err)
                        }
                    }
                }

                // only script defined Return and Escape events can be cancelled
                if keyCode == sf.KeyReturn && config.Conf.GameActive {
                    textMode := config.Conf.TextMode
                    if !textMode {
                        config.Conf.TextMode = true
                        err := config.Lua.State.DoString(string(textEntered))
                        if err != nil {
                            log.Println(err)
                        }
                    } else {
                        config.Conf.TextMode = false
                        textEntered = nil
                        text.SetString("")
                    }
                } else if keyCode == sf.KeyEscape {
                    config.Conf.GameActive = !config.Conf.GameActive && true
                }
            }
        }
	}
}

