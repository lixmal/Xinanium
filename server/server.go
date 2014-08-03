package main

import (
	"net"
	//   "os"
	wm "./worldmap"
	"encoding/gob"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
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
	PORT             = 22342
	LADDRESS         = "0.0.0.0"
	TIMEOUT          = 60 * 1000
	IDENTCODE uint32 = 0x58696E4C
)

// actions
const (
	PLAYER_MOVE = iota
	GET_PLAYER
	GET_PLAYER_TEX
	PLAYER_LOGIN
)

const (
	SPRITEDIR       = filepath.Dir(filepath.FromSlash("resources/textures/spritesheets")) + "/"
	SPRITEEXTENSION = ".png"
)

type code uint64
type Dir struct {
	X, Y float32
}

func main() {
	service := LADDRESS + ":" + strconv.FormatInt(PORT, 10)
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	checkErr(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkErr(err)
	log.Println("server listening on port", PORT)

	worldmap := wm.ReadOpen("default")

	// init Luap
	initLua(&config.Lua.State)

	// start event loop
	initEventQueue()

	// run map scripts
	worldmap.RunScripts()

	var connchans []chan code
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		//    conn.SetTimeout(TIMEOUT)
		connchan := make(chan code, 2)
		connchans = append(connchans, connchan)
		log.Println("accepted incoming connection from", conn.RemoteAddr())
		go handleClient(conn, connchan)
	}
}

func handleClient(conn net.Conn, connchan chan code) {
	defer conn.Close()
	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)

	// wait for player login cmd
	var action code
	if err := readConn(decoder, &action); err != nil {
		disconnect(encoder, conn, err.Error())
		return
	}
	if action != PLAYER_LOGIN {
		disconnect(encoder, conn, "not logged in")
		return
	}

	// read login name
	var login string
	if err := readConn(decoder, &login); err != nil {
		disconnect(encoder, conn, err.Error())
		return
	}

	// read password
	var pw string
	if err := readConn(decoder, &pw); err != nil {
		disconnect(encoder, conn, err.Error())
		return
	}

	// check login & pw
	if login != "vik" || pw != "secret" {
		log.Println("login by", login, " failed")
		disconnect(encoder, conn, "wrong login")
		return
	}
	log.Println("login by", login, "successful")

	// wait for player tex cmd
	if err := readConn(decoder, &action); err != nil {
		disconnect(encoder, conn, err.Error())
		return
	}
	if action != GET_PLAYER_TEX {
		disconnect(encoder, conn, "not logged in")
		return
	}

	// TODO: read player from db

	// read player texture from file
	playertex, err := ioutil.ReadFile(SPRITEDIR + login + SPRITEEXTENSION)
	if err != nil {
		disconnect(encoder, conn, "failed to read player texture")
		return
	}
	// send player texture to client
	if err := sendConn(encoder, playertex); err != nil {
		disconnect(encoder, conn, err.Error())
		return
	}
	playertex = nil

	// send current map
	if err := sendConn(encoder, worldmap.Current); err != nil {
		disconnect(encoder, conn, err.Error())
		return
	}

	// TODO: map player to connection
	for {
		if err := readConn(decoder, &action); err != nil {
			disconnect(encoder, conn, err.Error())
			return
		}

		switch action {
		case PLAYER_MOVE:
			var dir Dir
			if err := readConn(decoder, &dir); err != nil {
				disconnect(encoder, conn, err.Error())
				return
			}
			log.Println(dir)
		}
	}
	/*    for {
	          n, err := conn.Read(buf[0:])
	          if err != nil {
	              return
	          }

	          s := string(buf[0:n])
	          // decode request
	          if s[0:2] == CD {
	              chdir(conn, s[3:])
	          } else if s[0:3] == DIR {
	              dirList(conn)
	          } else if s[0:3] == PWD {
	              pwd(conn)
	          }
	      }
	*/
}

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

func sendConn(encoder *gob.Encoder, value interface{}) error {
	err := encoder.Encode(IDENTCODE)
	if err != nil {
		log.Println(err)
		return err
	}
	err = encoder.Encode(value)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func disconnect(encoder *gob.Encoder, conn net.Conn, sendErr string) {
	var ident uint32

	err := encoder.Encode(ident)
	if err != nil {
		log.Println(err)
		return
	}

	err = encoder.Encode(sendErr)
	if err != nil {
		log.Println(err)
	}
	// TODO: actually disconnect
	log.Println("disconnected from", conn.RemoteAddr())
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
