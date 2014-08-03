package main

import (
	"./config"
	"./event"
	"./monster"
	"./network"
	"./player"
	"./renderer"
	wm "./worldmap"
	sf "bitbucket.org/krepa098/gosfml2"
	"fmt"
	"log"
	"math"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"strconv"
)

// TODO: Lock every sprite/window and then test!

type Duration float64

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU() + 1)
	runtime.LockOSThread()

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

}

func i2c(x, y config.Coord) (config.Coord, config.Coord) {
	return (2*y + x) / 2, (2*y - x) / 2
}

func c2i(x, y config.Coord) (config.Coord, config.Coord) {
	return 2*y - x, (2*y + x) / 2
}

func getTileCoord(x, y config.Coord, h int8) (config.Coord, config.Coord) {
	return config.Coord(math.Floor(float64(x / config.Coord(h)))), config.Coord(math.Floor(float64(y / config.Coord(h)))) // height and width equal
}

func main() {

	window := sf.NewRenderWindow(sf.VideoMode{config.Conf.ScreenWidth, config.Conf.ScreenHeight, config.Conf.BitDepth}, config.Conf.GameTitle, sf.StyleDefault, config.Conf.ContextSettings)
	config.Conf.Window = window

	// create view
	view := sf.NewViewFromRect(sf.FloatRect{0, 0, float32(config.Conf.ScreenWidth) * 0.7, float32(config.Conf.ScreenHeight) * 0.7})

	network.Connect()
	defer network.Disconnect()

	if err := network.Send(config.PLAYER_LOGIN); err != nil {
		log.Println(err)
	}
	if err := network.Send("vik"); err != nil {
		log.Println(err)
	}
	if err := network.Send("secret"); err != nil {
		log.Println(err)
	}

	// get current map from server

	// load fonts, move this to RM
	font, err := sf.NewFontFromFile("/usr/share/fonts/truetype/ubuntu-font-family/Ubuntu-B.ttf")
	if err != nil {
		log.Fatal("...")
	}
	text, err := sf.NewText(font)
	text.SetCharacterSize(12)
	_ = err

	// set window to inactive for OpenGL
	if !window.SetActive(false) {
		log.Fatal("Could not set window OpenGL context to false")
	}

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
	if player1 == nil {
		log.Fatal("could not create Player")
	}
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

	// start rendering at last
	go renderer.Render(window, text)

	config.Conf.GameActive = true

	// game loop
	for window.IsOpen() && config.Conf.Connected {
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
			switch eT := e.(type) {
			case sf.EventClosed:
				window.Close()
			case sf.EventLostFocus:
				config.Conf.GameActive = false
				runtime.GC()
			case sf.EventGainedFocus:
				runtime.GC()
				config.Conf.GameActive = true
			case sf.EventTextEntered:
				char := eT.Char
				// trigger text entered in any case
				ev := &event.CharPressed{Event: event.New(event.TypeCharPressed), Char: string(char)}
				cancelled := event.Trigger(ev)
				if config.Conf.TextMode && config.Conf.GameActive && !cancelled {
					char = rune(ev.Char[0])
					if strconv.IsPrint(char) {
						textEntered = append(textEntered, char)
					}
					text.SetString(string(textEntered))
				}
			case sf.EventKeyPressed:
				keyCode := eT.Code
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
