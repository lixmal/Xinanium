package event
import(
    "sort"
    "log"
    "time"
    "github.com/stevedonovan/luar"
    "../config"
)

type Listener struct {
    priority uint16
    Handler *luar.LuaObject
}

type Eventer interface {
    Type() Type
    SetCancelled(bool)
    Cancelled() bool
}

// sorting
type byPriority []*Listener

func (p byPriority) Len() int           { return len(p) }
func (p byPriority) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p byPriority) Less(i, j int) bool { return p[i].priority < p[j].priority }


// Event Types

type Type uint16
const (
    TypePlayerNew = iota
    TypePlayerMove
    TypePlayerCollision
    TypePlayerJump
    TypePlayerTalk
    TypePlayerChangedDirection
    TypePlayerChangedPosition
    TypePlayerHurt
    TypePlayerKilled
    TypePlayerRemoved
    TypeMonsterNew
    TypeMonsterMove
    TypeMonsterCollision
    TypeMonsterJump
    TypeMonsterTalk
    TypeMonsterChangedDirection
    TypeMonsterChangedPosition
    TypeMonsterHurt
    TypeMonsterKilled
    TypeMonsterRemoved
    TypeKeyPressed
    TypeCharPressed
    TypeTicker
)

func New(kind Type) *Event {
    return &Event{kind, false}
}

type Event struct {
    kind Type
    cancelled bool
}

func (e *Event) SetCancelled(c bool) {
    e.cancelled = c
}

func (e *Event) Cancelled() bool {
    return e.cancelled
}

func (e *Event) Type() Type {
    return e.kind
}


type PlayerNew struct {
    *Event
    Player config.LivingEntity
}

type PlayerMove struct {
    *Event
    Player config.LivingEntity
    NewX float32
    NewY float32
}

type PlayerCollision struct {
    *Event
    Player config.LivingEntity
    What config.LivingEntity
    X float32
    Y float32
}

type PlayerJump struct {
    *Event
    Player config.LivingEntity
}

type PlayerTalk struct {
    *Event
    Player config.LivingEntity
    Text string
}

type PlayerChangedDirection struct {
    *Event
    Player config.LivingEntity
    NewDirX float32
    NewDirY float32
}

type PlayerChangedPosition struct {
    *Event
    Player config.LivingEntity
    NewX float32
    NewY float32
}

type PlayerHurt struct {
    *Event
    Player config.LivingEntity
    Damager config.Entity
    Damage int16
}

type PlayerKilled struct {
    *Event
    Player config.LivingEntity
    Killer config.Entity
    HurtEvent *PlayerHurt
}

type PlayerRemoved struct {
    *Event
    Player config.LivingEntity
}

type MonsterNew struct {
    *Event
    Monster config.LivingEntity
}

type MonsterMove struct {
    *Event
    Monster config.LivingEntity
    NewX float32
    NewY float32
}

type MonsterCollision struct {
    *Event
    Monster config.LivingEntity
    What config.LivingEntity
    X float32
    Y float32
}

type MonsterJump struct {
    *Event
    Monster config.LivingEntity
}

type MonsterTalk struct {
    *Event
    Monster config.LivingEntity
    Text string
}

type MonsterChangedDirection struct {
    *Event
    Monster config.LivingEntity
    NewDirX float32
    NewDirY float32
}

type MonsterChangedPosition struct {
    *Event
    Monster config.LivingEntity
    NewX float32
    NewY float32
}

type MonsterHurt struct {
    *Event
    Monster config.LivingEntity
    Damager config.Entity
    Damage int16
}

type MonsterKilled struct {
    *Event
    Monster config.LivingEntity
    Killer config.Entity
    HurtEvent *MonsterHurt
}

type MonsterRemoved struct {
    *Event
    Monster config.LivingEntity
}

type KeyPressed struct {
    *Event
    Key uint16
}

type CharPressed struct {
    *Event
    Char string
}

type Ticker struct {
    *Event
    control chan bool
}

func (e *Ticker) Stop() {
    select {
    case e.control <- true:
    default:
    }
}
func (e *Ticker) PauseResume() {
    select {
    case e.control <- false:
    default:
    }
}


func Register(eType uint16, fn interface{}, priority uint16) bool {
    e := Type(eType)
    //lock
    if f, ok := fn.(*luar.LuaObject); ok && f.Type == "function" {
        Registry[e] = append(Registry[e], &Listener{Handler: f, priority: priority})
        sort.Sort(byPriority(Registry[e]))
        return true
    }
    // TODO: add lua error on unknown func
    log.Println("Couldn't add event listener: ", fn)
    return false
}

func Unregister(eType uint16, fn interface{}) bool {
    e := Type(eType)
    //lock
    if f, ok := fn.(*luar.LuaObject); ok {
        slc := Registry[e]
        for i, v := range slc {
            log.Printf("%+#v == %+#v\n", v.Handler, f)
            if (v.Handler == f) {
                Registry[e] = slc[:i]
                Registry[e] = append(Registry[e], slc[i+1:]...)
                f.Close()
                return true
            }
        }
    }
    log.Println("Couldn't remove event: ", fn)
    return false
}

func RegisterTicker(ticks uint16, fn interface{}, priority uint16) bool {
    bl := Register(TypeTicker * ticks + 1000, fn, priority)
    if bl {
        go makeTicker(ticks)
    }
    return bl
}

// TODO: sync with game ticker:
func RegisterTimer(ticks uint16, fn interface{}) bool {
    if f, ok := fn.(*luar.LuaObject); ok && f.Type == "function" {
        go func() {
            time.Sleep(time.Second / config.TICKS * time.Duration(ticks))
            config.Lua.Lock()
            _, err := f.Call()
            f.Close()
            config.Lua.Unlock()
            if err != nil {
                log.Println("error in function call", f)
            }
        }()
        return true
    }
    log.Println("Couldn't add timer listener: ", fn)
    return false
}


// stop stops all for now, shutdown ticker!!
// think about calling events based on gameticker to sync
func makeTicker(ticks uint16) {
    ticker := time.Tick(time.Second / config.TICKS * time.Duration(ticks))
    controlChan := make(chan bool, 1)
    active := true
    tickid := TypeTicker * ticks + 1000

    for _ = range ticker {
        select {
        case v := <-controlChan:
            if v == true {
                return
            }
            active = !active && true
        default:
        }
        if active {
            Trigger(&Ticker{Event: New(Type(tickid)), control: controlChan})
        }
    }
}

// TODO: auto adjust buffer to map
func Trigger(event Eventer) bool {
    log.Println("event into queue", event.Type())
    Queue <- event
    /*select {
        case bl := <-WaitChan:
            return bl
        default:
    }
    */
    log.Println("event wait response", event.Type())
    //bl := <-WaitChan
    //log.Println("event answer of", event.Type(), " is", bl)
    //return bl
    //log.Println("Event queue is full")
    return false
}



var WaitChan = make(chan bool, 100)
var Queue = make(chan Eventer, 100)
var Registry = make(map[Type][]*Listener)
