package main

import (
	"net"
	//   "os"
//	wm "./worldmap"
	"encoding/gob"
    "bytes"
    "encoding/binary"
	"io/ioutil"
	"log"
// 	"path/filepath"
    "time"
    "sync"
)

type Player struct {
	Handle string
	Name   string
	// sprite missing here
	Speed      uint16
	JumpHeight uint8
	Health     int16
	Invincible bool
	Invisible  bool
	Walking    bool
	InAir      bool
	Dead       bool
	Floating   bool
	dir        *Dir
	entityType uint16
	centric    bool
}


const (
	LISTEN_ADDRESS   = ""
	PORT             = "22342"
	TIMEOUT          = 10 * time.Second
    PROTO_VERSION    = 1
	IDENTCODE uint32 = 0x58696E << 8 | PROTO_VERSION
    ERRTHRESHOLD     = 5
    PACKET_MAXSIZE   = 512
)

// actions
const (
	PLAYER_MOVE = iota
	GET_PLAYER
	GET_PLAYER_TEX
	PLAYER_LOGIN
    NO_PLAYER_TEX
)

const (
	SPRITEDIR       = "resources/textures/spritesheets/"
    // SPRITEDIR       = filepath.Dir(filepath.FromSlash("resources/textures/spritesheets")) + "/"
	SPRITEEXTENSION = ".png"
)

type Dir struct {
	X, Y float32
}
type Conn struct {
    net.PacketConn
    net.Addr
    // TODO: change to player object
    handle string
}

var clients = make(map[net.Addr]chan []byte)
var clientsMutex sync.RWMutex


func main() {
	service := LISTEN_ADDRESS + ":" + PORT

	conn, err := net.ListenPacket("udp", service)
	if err != nil {
		log.Fatal(err)
	}
    defer conn.Close()
	log.Println("Listening on", conn.LocalAddr())

	//worldmap := wm.ReadOpen("default")

	// init Luap
	// initLua(&config.Lua.State)

	// start event loop
	// initEventQueue()

	// run map scripts
	// worldmap.RunScripts()


    // handle incoming packets
    // either move to exisiting client goroutine
    // or negotiate a new client
    // TODO: put "timeout" somewhere in here
	for {
        var buf [PACKET_MAXSIZE]byte
        n, addr, err := conn.ReadFrom(buf[:])
		if err != nil {
            // TODO: add debug output and remove this
			log.Println(err)
			continue
		}
        // minimum of identcode + action
        if n < 4+2 {
            continue
        }
        // 4 bytes from identcode
        if binary.BigEndian.Uint32(buf[:4]) != IDENTCODE {
            continue
        }

        clientsMutex.RLock()
        clientchan, ok := clients[addr]
        clientsMutex.RUnlock()
        if ok {
            clientchan<- buf[4:n]
        } else {
            log.Println("Accepted incoming connection from", addr)
            go handleClient(&Conn{conn, addr, ""}, buf[4:n])
        }
	}
}

func handleClient(conn *Conn, data []byte) {

    action := binary.BigEndian.Uint16(data[:2])

	if action != PLAYER_LOGIN {
        conn.disconnect("not logged in")
		return
	}

    // TODO: session handling in tcp, rest in udp

    // Ignore non session initiating packets. Once session is established: send errors on unknown packets

    // TODO: decode login+pw from data at this point
    {
        login, pw := "vik", "secret"

        // TODO: read player from db
        // check login & pw
        if login != "vik" || pw != "secret" {
            log.Println("Login by", login, "from", conn.Addr, "failed")
            conn.disconnect("wrong login")
            return
        }
        // TODO: sent OK response
        conn.handle = login
        log.Println("Login by", login, "from", conn.Addr, "successful")
    }
    data = nil

    // authenticated: create chan and add to clients
    connchan := make(chan []byte)
    clientsMutex.Lock()
    clients[conn.Addr] = connchan
    clientsMutex.Unlock()


    // errors don't break connection at once from this point
    var errCnt uint8
    for data = range connchan {
        if errCnt > ERRTHRESHOLD {
            conn.disconnect("too many errors")
            return
        }

        // read action
        action = binary.BigEndian.Uint16(data[:2])

        d := gob.NewDecoder(bytes.NewBuffer(data))

		switch action {
		case PLAYER_MOVE:
			var dir Dir
            if d.Decode(dir) != nil {
               errCnt++
               continue
            }
			log.Println(dir)
        case GET_PLAYER_TEX:
            // TODO: check here if in correct game phase
            // TODO: sanitize player name for path or use ID from DB
            // read player texture from file
            playertex, err := ioutil.ReadFile(SPRITEDIR + conn.handle + SPRITEEXTENSION)
            if err != nil {
                // TODO: not sure what todo, maybe just ignore
                conn.disconnect("failed to read player texture")
                return
            }
            // send player texture to client
            if err := conn.write(playertex); err != nil {
                conn.disconnect("failed to send player texture")
                return
            }
        // TODO: send current map
        //if err := sendConn(encoder, worldmap.Current); err != nil {
        //	disconnect(encoder, conn, err.Error())
        //	return
        default:
            // unknown action
            continue
		}
        // all went fine: remove err
        errCnt = 0
	}
}


/*
func readConn(decoder *gob.Decoder, value interface{}) error {
	var ident uint32

	err := decoder.Decode(&ident)
	if err != nil {
		log.Println(err)
		return errors.New("value expected, received something else")
	}
	if ident != IDENTCODE {
		return errors.New("ident code expected, received something else")
	}

	err = decoder.Decode(value)
	if err != nil {
		log.Println(err)
		return errors.New("value expected, received something else")
	}

	return nil
}

func (c *Conn) read(value interface{}) bool {
	var ident uint32

    // TODO: send errors to client instead of false!?
	err := c.decoder.Decode(&ident)
	if err != nil {
		log.Println(err)
		return false
	}
	if ident != IDENTCODE {
		return false
	}

	err = c.decoder.Decode(value)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
*/

func (c *Conn) write(value interface{}) error {
    // TODO: check keepalive and disconnect if no _valid_ packets arrived
    var buf0 [PACKET_MAXSIZE]byte
    binary.BigEndian.PutUint32(buf0[:4], IDENTCODE)
    buf := bytes.NewBuffer(buf0[4:])
    e := gob.NewEncoder(buf)

    if err := e.Encode(value); err != nil {
        return err
    }
    if _, err := c.WriteTo(buf0[:], c.Addr); err != nil {
        return err
    }
    return nil
}

func (c *Conn) disconnect(err string) {
    // TODO: write err to client

    clientsMutex.Lock()
    if clientchan, ok := clients[c.Addr]; ok {
        close(clientchan)
    }
    delete(clients, c.Addr)
    clientsMutex.Unlock()

    // TODO: remove player/client mapping and resolve circle referencing
    log.Println("Disconnected " + c.handle + " from " + c.Addr.String() + ": " + err)
}

