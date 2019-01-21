package main

import (
    "./config"
    //    "./event"
    //    "./monster"
    "./network"
    "./player"
    //    "./renderer"
    //    wm "./worldmap"
    "image"
    "log"
    "math"

    "azul3d.org/gfx.v1"
    "azul3d.org/gfx/window.v2"
    "azul3d.org/keyboard.v1"
    "azul3d.org/lmath.v1"
    "azul3d.org/tmx.dev"
    //    "net/http"
    //    _ "net/http/pprof"
    "runtime"
)

// TODO: Lock every sprite/window and then test!

type Duration float64

func init() {
    runtime.GOMAXPROCS(runtime.NumCPU() + 1)
    /*
        go func() {
            log.Println(http.ListenAndServe("localhost:6060", nil))
        }()
    */

    log.SetFlags(log.LstdFlags | log.Lshortfile)
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

func gfxLoop(w window.Window, r gfx.Renderer) {
    config.Conf.Window = w
    config.Conf.Renderer = r

    // create view
    //view := sf.NewViewFromRect(sf.FloatRect{0, 0, float32(config.Conf.ScreenWidth) * 0.7, float32(config.Conf.ScreenHeight) * 0.7})

    network.Connect()
    defer network.Disconnect()
    if err := network.Login("vik", "secret"); err != nil {
        return
    }

    // get current map from server
    /*
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

    */
    //var textEntered []rune

    /*
       // music
       music, err := sf.NewMusicFromFile(RESOURCESDIR + "sound/test.ogg")
       if err != nil {
           log.Fatal("Could not load sound: ")
       }
       music.SetLoop(true)
       music.Play()
    */

    /*
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
    */

    // start rendering at last

    config.Conf.GameActive = true

    go func() {

        // Create an event mask for the events we are interested in
        evMask :=
            window.KeyboardStateEvents |
                window.CloseEvents |
                window.LostFocusEvents |
                window.GainedFocusEvents

        // Create a channel of events
        events := make(chan window.Event, 256)

        // Have the window notify our channel whenever events occur
        w.Notify(events, evMask)

        // event loop
        for e := range events {
            switch e := e.(type) {
            case keyboard.StateEvent:
                if e.State == keyboard.Down {
                    switch e.Key {
                    case keyboard.Escape:
                        config.Conf.GameActive = !config.Conf.GameActive && true
                        runtime.GC()
                    case keyboard.Z:
                        config.Players["vik"].(*player.Player).Speed += 30
                        log.Println(config.Players["vik"].(*player.Player).Speed)
                    case keyboard.X:
                        config.Players["vik"].(*player.Player).Speed -= 30
                        log.Println(config.Players["vik"].(*player.Player).Speed)
                    }
                }
            case window.Close:
                // TODO: quit network
                config.Conf.GameActive = false
            case window.LostFocus:
                config.Conf.GameActive = false
                runtime.GC()
            case window.GainedFocus:
                runtime.GC()
                config.Conf.GameActive = true
            }
            /*
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
            */
        }
    }()

    //    wm.ReadOpen("gobmap")
    player1 := player.New("vik", "vik", true)
    config.ToDraw = append(config.ToDraw, player1.GetSprite())
    player1.GetSprite().SetPos(lmath.Vec3{float64(r.Bounds().Dx()) / 2.0, -1, float64(r.Bounds().Dy()) / 2.0})

    camera := gfx.NewCamera()
    camera.SetOrtho(r.Bounds(), 0.01, 1000)
    camera.SetParent(player1.GetSprite().Transform)
    camera.SetPos(lmath.Vec3{-(float64(r.Bounds().Dx()) / 2.0), -2, -(float64(r.Bounds().Dy()) / 2.0)})

    r.Clock().SetMaxFrameRate(config.TICKS)

    tmxMap, layers, err := tmx.LoadFile("resources/maps/default.tmx", nil)
    if err != nil {
        log.Fatal(err)
    }

    go network.InitListener()

    watcher := w.Keyboard()
    for config.Conf.Connected {

        // player moving
        if !config.Conf.TextMode && config.Conf.GameActive {
            var x, y float64
            var pressed bool
            if watcher.Down(keyboard.S) {
                y -= 1
                pressed = true
            } else if watcher.Down(keyboard.W) {
                y += 1
                pressed = true
            }
            if watcher.Down(keyboard.D) {
                x += 1
                pressed = true
            } else if watcher.Down(keyboard.A) {
                pressed = true
                x -= 1
            }
            if pressed {
                player1.Move(x, y)
            }
        }

        //s := float64(r.Bounds().Dy()) / 2.0 // Card is two units wide, so divide by two.
        //player1.GetSprite().SetScale(lmath.Vec3{s, s, s})
        //player1.GetSprite().SetScale(lmath.Vec3{10, 1, 10})
        //fmt.Println(player1.GetSprite().Scale())

        /*
           spr := player1.GetSprite()
           rot := spr.Rot()
           spr.SetRot(lmath.Vec3{
               X: rot.X,
               Y: rot.Y,
               Z: rot.Z + (15 * r.Clock().Dt()),
           })
        */
        r.Clear(image.ZR, gfx.Color{0, 0, 0, 0})
        r.ClearDepth(image.ZR, 1.0)

        for _, v := range config.ToDraw {
            r.Draw(image.ZR, v, camera)
        }

        for _, layer := range tmxMap.Layers {
            objects, ok := layers[layer.Name]
            if ok {
                for _, obj := range objects {
                    r.Draw(image.ZR, obj, camera)
                }
            }
        }

        r.Render()
    }
}

func main() {
    props := window.NewProps()
    props.SetTitle(config.Conf.GameTitle + " - {FPS}")
    props.SetSize(config.Conf.ScreenWidth, config.Conf.ScreenHeight)
    window.Run(gfxLoop, props)
}
