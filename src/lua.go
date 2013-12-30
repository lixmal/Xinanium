package main
import (
    "fmt"
    "./config"
    "./monster"
    "./player"
    "github.com/stevedonovan/luar"
    sf "bitbucket.org/krepa098/gosfml2"
)

func initLua() {
    config.InitEventQueue()
    config.Lua = luar.Init()
    luar.Register(config.Lua, "entity", luar.Map{
        "Player": config.LivingEntityPlayer,
        "Monster": config.LivingEntityMonster,
    })
    luar.Register(config.Lua, "player", luar.Map{
        "New":  player.New,
        "Players": config.Players,
    })
    luar.Register(config.Lua, "monster", luar.Map{
        "New":  monster.New,
        "Monsters": config.Monsters,
    })
    luar.Register(config.Lua, "", luar.Map{
        "Print":  fmt.Println,
    })
    luar.Register(config.Lua, "event", luar.Map{
        "Register": config.RegisterEvent,
        "Ticker": config.RegisterTicker,
        "Registry": config.EventRegistry,
        "PlayerNew": config.EventTypePlayerNew,
        "PlayerMove": config.EventTypePlayerMove,
        "PlayerCollision": config.EventTypePlayerCollision,
        "PlayerJump": config.EventTypePlayerJump,
        "PlayerTalk": config.EventTypePlayerTalk,
        "PlayerChangedDirection": config.EventTypePlayerChangedDirection,
        "PlayerChangedPosition": config.EventTypePlayerChangedPosition,
        "PlayerHurt": config.EventTypePlayerHurt,
        "PlayerKilled": config.EventTypePlayerKilled,
        "PlayerRemoved": config.EventTypePlayerRemoved,
        "MonsterNew": config.EventTypeMonsterNew,
        "MonsterMove": config.EventTypeMonsterMove,
        "MonsterCollision": config.EventTypeMonsterCollision,
        "MonsterJump": config.EventTypeMonsterJump,
        "MonsterTalk": config.EventTypeMonsterTalk,
        "MonsterChangedDirection": config.EventTypeMonsterChangedDirection,
        "MonsterChangedPosition": config.EventTypeMonsterChangedPosition,
        "MonsterHurt": config.EventTypeMonsterHurt,
        "MonsterKilled": config.EventTypeMonsterKilled,
        "MonsterRemoved": config.EventTypeMonsterRemoved,
        "KeyPressed": config.EventTypeKeyPressed,
        "KeyA": sf.KeyA,
        "KeyB": sf.KeyB,
        "KeyC": sf.KeyC,
        "KeyD": sf.KeyD,
        "KeyE": sf.KeyE,
        "KeyF": sf.KeyF,
        "KeyG": sf.KeyG,
        "KeyH": sf.KeyH,
        "KeyI": sf.KeyI,
        "KeyJ": sf.KeyJ,
        "KeyK": sf.KeyK,
        "KeyL": sf.KeyL,
        "KeyM": sf.KeyM,
        "KeyN": sf.KeyN,
        "KeyO": sf.KeyO,
        "KeyP": sf.KeyP,
        "KeyQ": sf.KeyQ,
        "KeyR": sf.KeyR,
        "KeyS": sf.KeyS,
        "KeyT": sf.KeyT,
        "KeyU": sf.KeyU,
        "KeyV": sf.KeyV,
        "KeyW": sf.KeyW,
        "KeyX": sf.KeyX,
        "KeyY": sf.KeyY,
        "KeyZ": sf.KeyZ,
        "KeyNum0": sf.KeyNum0,
        "KeyNum1": sf.KeyNum1,
        "KeyNum2": sf.KeyNum2,
        "KeyNum3": sf.KeyNum3,
        "KeyNum4": sf.KeyNum4,
        "KeyNum5": sf.KeyNum5,
        "KeyNum6": sf.KeyNum6,
        "KeyNum7": sf.KeyNum7,
        "KeyNum8": sf.KeyNum8,
        "KeyNum9": sf.KeyNum9,
        "KeyEscape": sf.KeyEscape,
        "KeyLControl": sf.KeyLControl,
        "KeyLShift": sf.KeyLShift,
        "KeyLAlt": sf.KeyLAlt,
        "KeyLSystem": sf.KeyLSystem,
        "KeyRControl": sf.KeyRControl,
        "KeyRShift": sf.KeyRShift,
        "KeyRAlt": sf.KeyRAlt,
        "KeyRSystem": sf.KeyRSystem,
        "KeyMenu": sf.KeyMenu,
        "KeyLBracket": sf.KeyLBracket,
        "KeyRBracket": sf.KeyRBracket,
        "KeySemiColon": sf.KeySemiColon,
        "KeyComma": sf.KeyComma,
        "KeyPeriod": sf.KeyPeriod,
        "KeyQuote": sf.KeyQuote,
        "KeySlash": sf.KeySlash,
        "KeyBackSlash": sf.KeyBackSlash,
        "KeyTilde": sf.KeyTilde,
        "KeyEqual": sf.KeyEqual,
        "KeyDash": sf.KeyDash,
        "KeySpace": sf.KeySpace,
        "KeyReturn": sf.KeyReturn,
        "KeyBack": sf.KeyBack,
        "KeyTab": sf.KeyTab,
        "KeyPageUp": sf.KeyPageUp,
        "KeyPageDown": sf.KeyPageDown,
        "KeyEnd": sf.KeyEnd,
        "KeyHome": sf.KeyHome,
        "KeyInsert": sf.KeyInsert,
        "KeyDelete": sf.KeyDelete,
        "KeyAdd": sf.KeyAdd,
        "KeySubtract": sf.KeySubtract,
        "KeyMultiply": sf.KeyMultiply,
        "KeyDivide": sf.KeyDivide,
        "KeyLeft": sf.KeyLeft,
        "KeyRight": sf.KeyRight,
        "KeyUp": sf.KeyUp,
        "KeyDown": sf.KeyDown,
        "KeyNumpad0": sf.KeyNumpad0,
        "KeyNumpad1": sf.KeyNumpad1,
        "KeyNumpad2": sf.KeyNumpad2,
        "KeyNumpad3": sf.KeyNumpad3,
        "KeyNumpad4": sf.KeyNumpad4,
        "KeyNumpad5": sf.KeyNumpad5,
        "KeyNumpad6": sf.KeyNumpad6,
        "KeyNumpad7": sf.KeyNumpad7,
        "KeyNumpad8": sf.KeyNumpad8,
        "KeyNumpad9": sf.KeyNumpad9,
        "KeyF1": sf.KeyF1,
        "KeyF2": sf.KeyF2,
        "KeyF3": sf.KeyF3,
        "KeyF4": sf.KeyF4,
        "KeyF5": sf.KeyF5,
        "KeyF6": sf.KeyF6,
        "KeyF7": sf.KeyF7,
        "KeyF8": sf.KeyF8,
        "KeyF9": sf.KeyF9,
        "KeyF10": sf.KeyF10,
        "KeyF11": sf.KeyF11,
        "KeyF12": sf.KeyF12,
        "KeyF13": sf.KeyF13,
        "KeyF14": sf.KeyF14,
        "KeyF15": sf.KeyF15,
        "KeyPause": sf.KeyPause,
    })
    config.Lua.DoString(`
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
    `)
}
