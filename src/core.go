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
    "./config"
    "./monster"
    "./player"
    "./renderer"
    wm "./worldmap"
)

const RESOURCESDIR = "resources/"
const SPRITEDIR = "resources/textures/spritesheets/"


type Duration float64


func init() {
	runtime.GOMAXPROCS(runtime.NumCPU() + 1)
	runtime.LockOSThread()

    // X11 multithreading, linux/X11 only
	C.XInitThreads()

    initLua()

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
    defer config.Lua.Close()

    toDraw := make([]sf.Drawer, 0)
	window := sf.NewRenderWindow(sf.VideoMode{config.Conf.ScreenWidth, config.Conf.ScreenHeight, config.Conf.BitDepth}, config.Conf.GameTitle, sf.StyleDefault, config.Conf.ContextSettings)
    config.Conf.Window = window

    // create view
    view := sf.NewViewFromRect(sf.FloatRect{0, 0, 520, 390})




    // read default worldmap and add to toDraw
    worldmap := wm.Read(RESOURCESDIR + "maps/test.dat")
    config.Conf.CurrentMap = worldmap
    for x, v := range worldmap.Tiles {
        for y := range v {
            tileType := worldmap.Tiles[x][y]
            var sprite *sf.Sprite
            switch tileType{
            case 1:
                sprite = sf.NewSprite(config.Conf.Rm.Texture(RESOURCESDIR + "textures/tiles/grass.png"))
            case 2:
                sprite = sf.NewSprite(config.Conf.Rm.Texture(RESOURCESDIR + "textures/tiles/dirt.png"))
            }
            if sprite != nil {
                sprite.SetPosition(sf.Vector2f{float32(x * wm.TILEWIDTH), float32(y * wm.TILEHEIGHT)})
                toDraw = append(toDraw, sprite)
            } else {
                log.Fatal("Unknown tile type: " + fmt.Sprintf("%i", tileType))
            }
        }
    }
    err := config.Lua.DoFile("test.lua")
    if err != nil {
        fmt.Println(err)
    }


    // music
    music, err := sf.NewMusicFromFile(RESOURCESDIR + "sound/test.ogg")
    if err != nil {
        log.Fatal("Could not load sound: ")
    }
    music.SetLoop(true)
    music.Play()


    // default Player
    player1 := player.New("Player", "vik", true)
    player1.SetPosition(100, 200)

    // spawn one test monster
    for i := 0; i < 1; i++ {
        //mon := monster.New("monster", config.Coord(rand.Intn(1500)), config.Coord(rand.Intn(1500)), 500)
        monster.New("monster").SetPosition(200, 200)
    }

    // assign view to main player
    {
        posx, posy := player1.Position()
        view.SetCenter(sf.Vector2f{posx, posy})
        window.SetView(view)
    }


    // set window to inactive for OpenGL
    if !window.SetActive(false) {
		log.Fatal("Could not set window OpenGL context to false")
	}


    font, err := sf.NewFontFromFile("/usr/share/fonts/truetype/ubuntu-font-family/Ubuntu-B.ttf")
    if err != nil {
        log.Fatal("...")
    }
    text := sf.NewText(font)
    text.SetCharacterSize(12)

    // start rendering
	go renderer.Render(window, toDraw, text)

    config.Conf.GameActive = true
    // game loop
    var e sf.Event
	for window.IsOpen() {
        <-config.GameTicker

        // player moving
        if config.Conf.GameActive {
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

        //game event loop
        for e = window.PollEvent(); e != nil; e = window.PollEvent() {
            switch e.(type) {
            case sf.EventClosed:
                window.Close()
            case sf.EventLostFocus:
                config.Conf.GameActive = false
                runtime.GC()
            case sf.EventGainedFocus:
                config.Conf.GameActive = true
            case sf.EventKeyPressed:
                config.TriggerEvent(&config.EventKeyPressed{Event: config.Event{EType: config.EventTypeKeyPressed}, KeyCode: uint16(e.(sf.EventKeyPressed).Code)})
                switch e.(sf.EventKeyPressed).Code {
                case sf.KeySpace:
                    if config.Conf.GameActive {
                        player1.StopAnimation()
                        player1.Jump()
                    }
                case sf.KeyTab:
                    config.Conf.Scrolling = !config.Conf.Scrolling && true
                case sf.KeyEscape:
                    config.Conf.GameActive = !config.Conf.GameActive && true
                case sf.KeyX:
                    for _, entity := range config.Monsters {
                        if mon, ok := entity.(*monster.Monster); ok {
                            mon.Run()
                        }
                    }
                case sf.KeyL:
                    config.EventRegistry = nil
                    config.EventRegistry = make(map[config.EventType][]*config.EventHandler)
                    config.Lua.Close()
                    initLua()
                    err := config.Lua.DoFile("test.lua")
                    if err != nil {
                        fmt.Println(err)
                    }
                }
            }
        }
	}
}

