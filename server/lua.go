package main

import (
    "./config"
    "./event"
    "./monster"
    "./player"
    "./worldmap"
    sf "bitbucket.org/krepa098/gosfml2"
    "fmt"
    "github.com/aarzilli/golua/lua"
    "github.com/stevedonovan/luar"
    "log"
    "reflect"
    "time"
)

func initLua(v **lua.State) {
    *v = luar.Init()

    luar.Register(*v, "entity", luar.Map{
        "Player":  config.LivingEntityPlayer,
        "Monster": config.LivingEntityMonster,
    })
    luar.Register(*v, "player", luar.Map{
        "New":     player.New,
        "Players": config.Players,
    })
    luar.Register(*v, "monster", luar.Map{
        "New":      monster.New,
        "Monsters": config.Monsters,
    })
    luar.Register(*v, "", luar.Map{
        "Print":  fmt.Println,
        "config": config.Conf,
        "Sleep":  func(t uint16) { time.Sleep(time.Duration(t) * time.Second) },
    })
    luar.Register(*v, "", luar.Map{
        "config.CurrentMap": worldmap.Current,
    })
    luar.Register(*v, "event", luar.Map{
        "Trigger": func(fn interface{}) {
            //println("called")
            fmt.Println(reflect.TypeOf(fn))
            st := luar.NewLuaObjectFromValue(config.Lua.State, fn)
            st.Push()
            str := luar.LuaToGo(config.Lua.State, reflect.TypeOf(event.MonsterHurt{}), -1)
            if s, ok := str.(event.MonsterHurt); ok {
                s.Event = event.New(event.TypeMonsterHurt)
                event.Trigger(s)
            }
        },
        "Register":                event.Register,
        "TEST":                    reflect.TypeOf(event.MonsterHurt{}),
        "TEST2":                   reflect.TypeOf(event.MonsterChangedDirection{}),
        "TEST3":                   reflect.TypeOf(event.MonsterChangedDirection{}),
        "Unregister":              event.Unregister,
        "Ticker":                  event.RegisterTicker,
        "Timer":                   event.RegisterTimer,
        "Registry":                event.Registry,
        "PlayerNew":               event.TypePlayerNew,
        "PlayerMove":              event.TypePlayerMove,
        "PlayerCollision":         event.TypePlayerCollision,
        "PlayerJump":              event.TypePlayerJump,
        "PlayerTalk":              event.TypePlayerTalk,
        "PlayerChangedDirection":  event.TypePlayerChangedDirection,
        "PlayerChangedPosition":   event.TypePlayerChangedPosition,
        "PlayerHurt":              event.TypePlayerHurt,
        "PlayerKilled":            event.TypePlayerKilled,
        "PlayerRemoved":           event.TypePlayerRemoved,
        "MonsterNew":              event.TypeMonsterNew,
        "MonsterMove":             event.TypeMonsterMove,
        "MonsterCollision":        event.TypeMonsterCollision,
        "MonsterJump":             event.TypeMonsterJump,
        "MonsterTalk":             event.TypeMonsterTalk,
        "MonsterChangedDirection": event.TypeMonsterChangedDirection,
        "MonsterChangedPosition":  event.TypeMonsterChangedPosition,
        "MonsterHurt":             event.TypeMonsterHurt,
        "MonsterKilled":           event.TypeMonsterKilled,
        "MonsterRemoved":          event.TypeMonsterRemoved,
        "KeyPressed":              event.TypeKeyPressed,
        "CharPressed":             event.TypeCharPressed,
    })
    luar.Register(*v, "key", luar.Map{
        // TODO: extract to func
        "IsPressed": func(k uint16) bool {
            if k < sf.KeyCount {
                if config.Conf.GameActive {
                    return sf.KeyboardIsKeyPressed(sf.KeyCode(k))
                }
            } else {
                v.RaiseError(fmt.Sprintf("unknown keycode: %d", k))
            }
            return false
        },
        "A":         sf.KeyA,
        "B":         sf.KeyB,
        "C":         sf.KeyC,
        "D":         sf.KeyD,
        "E":         sf.KeyE,
        "F":         sf.KeyF,
        "G":         sf.KeyG,
        "H":         sf.KeyH,
        "I":         sf.KeyI,
        "J":         sf.KeyJ,
        "K":         sf.KeyK,
        "L":         sf.KeyL,
        "M":         sf.KeyM,
        "N":         sf.KeyN,
        "O":         sf.KeyO,
        "P":         sf.KeyP,
        "Q":         sf.KeyQ,
        "R":         sf.KeyR,
        "S":         sf.KeyS,
        "T":         sf.KeyT,
        "U":         sf.KeyU,
        "V":         sf.KeyV,
        "W":         sf.KeyW,
        "X":         sf.KeyX,
        "Y":         sf.KeyY,
        "Z":         sf.KeyZ,
        "0":         sf.KeyNum0,
        "1":         sf.KeyNum1,
        "2":         sf.KeyNum2,
        "3":         sf.KeyNum3,
        "4":         sf.KeyNum4,
        "5":         sf.KeyNum5,
        "6":         sf.KeyNum6,
        "7":         sf.KeyNum7,
        "8":         sf.KeyNum8,
        "9":         sf.KeyNum9,
        "Escape":    sf.KeyEscape,
        "LControl":  sf.KeyLControl,
        "LShift":    sf.KeyLShift,
        "LAlt":      sf.KeyLAlt,
        "LSystem":   sf.KeyLSystem,
        "RControl":  sf.KeyRControl,
        "RShift":    sf.KeyRShift,
        "RAlt":      sf.KeyRAlt,
        "RSystem":   sf.KeyRSystem,
        "Menu":      sf.KeyMenu,
        "LBracket":  sf.KeyLBracket,
        "RBracket":  sf.KeyRBracket,
        "SemiColon": sf.KeySemiColon,
        "Comma":     sf.KeyComma,
        "Period":    sf.KeyPeriod,
        "Quote":     sf.KeyQuote,
        "Slash":     sf.KeySlash,
        "BackSlash": sf.KeyBackSlash,
        "Tilde":     sf.KeyTilde,
        "Equal":     sf.KeyEqual,
        "Dash":      sf.KeyDash,
        "Space":     sf.KeySpace,
        "Return":    sf.KeyReturn,
        "Back":      sf.KeyBack,
        "Tab":       sf.KeyTab,
        "PageUp":    sf.KeyPageUp,
        "PageDown":  sf.KeyPageDown,
        "End":       sf.KeyEnd,
        "Home":      sf.KeyHome,
        "Insert":    sf.KeyInsert,
        "Delete":    sf.KeyDelete,
        "Add":       sf.KeyAdd,
        "Subtract":  sf.KeySubtract,
        "Multiply":  sf.KeyMultiply,
        "Divide":    sf.KeyDivide,
        "Left":      sf.KeyLeft,
        "Right":     sf.KeyRight,
        "Up":        sf.KeyUp,
        "Down":      sf.KeyDown,
        "Numpad0":   sf.KeyNumpad0,
        "Numpad1":   sf.KeyNumpad1,
        "Numpad2":   sf.KeyNumpad2,
        "Numpad3":   sf.KeyNumpad3,
        "Numpad4":   sf.KeyNumpad4,
        "Numpad5":   sf.KeyNumpad5,
        "Numpad6":   sf.KeyNumpad6,
        "Numpad7":   sf.KeyNumpad7,
        "Numpad8":   sf.KeyNumpad8,
        "Numpad9":   sf.KeyNumpad9,
        "F1":        sf.KeyF1,
        "F2":        sf.KeyF2,
        "F3":        sf.KeyF3,
        "F4":        sf.KeyF4,
        "F5":        sf.KeyF5,
        "F6":        sf.KeyF6,
        "F7":        sf.KeyF7,
        "F8":        sf.KeyF8,
        "F9":        sf.KeyF9,
        "F10":       sf.KeyF10,
        "F11":       sf.KeyF11,
        "F12":       sf.KeyF12,
        "F13":       sf.KeyF13,
        "F14":       sf.KeyF14,
        "F15":       sf.KeyF15,
        "Pause":     sf.KeyPause,
    })
    // deprecated, + use setfenv
    v.DoString(`
        os.execute   = nil
        os.exit      = nil
        os.remove    = nil
        os.rename    = nil
        os.setlocale = nil
        io           = nil
        loadfile     = nil
        dofile       = nil
        load         = nil
        require      = nil
        package      = nil
        string.dump  = nil
        debug        = nil
        go           = luar
        luar         = nil
    `)

}

func initEventQueue() {
    go func() {
        defer config.Lua.State.Close()
        for e := range event.Queue {
            for _, listener := range event.Registry[e.Type()] {
                if listener.Handler != nil {
                    //println("lock queue")
                    config.Lua.Lock()
                    //println("locked queue")
                    _, err := listener.Handler.Call(e)
                    //println("unlock queue")
                    config.Lua.Unlock()
                    if err != nil {
                        log.Println("error in function call of event type '", e.Type(), "': ", err)
                    }
                }
            }
            //event.WaitChan <- e.Cancelled()
        }
    }()
}

/*
func initCallerQueue() {
    go func() {
        for h := range event.CallerQueue {
            h.Call()
        }
    }()
}*/

// events calling events, howto wait for the first to finish?
// compare monster type by passing empty struct

/*
   types

   event.Timer(ticks, fn)
   - execute just once
   - substitue for sleep

   event.Ticker(ticks, fn)
   - runs from the start
   - can refer to 'self'
   - always same period
   - can't take infinit loops
   - substitue for infinite loops
   - no sleep

   m.SetHandler(fn) + m.Run()
   - runs at load map
   . refers to m
   - always 1 ticks (configurable?)
   - substitue for infinite loops
   - no sleep

   m.NewState(fn)
   - runs on any event
   - infinite loops
   - sleep = varying periods
   - no need for ticks
   - running infinite loop in event will lock game

*/
