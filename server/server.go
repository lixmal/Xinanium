package main

import (
    "net"
 //   "os"
    "log"
    "strconv"
    "encoding/gob"
)

const (
    PORT      = 22342
    LADDRESS  = "0.0.0.0"
    TIMEOUT   = 60 * 1000
    IDENTCODE = 0x58696E4C
)

const (
    MOVE = iota
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

    var connchans []chan code
    for {
        conn, err := listener.Accept()
        if err != nil {
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

    // map player to connection
    for {
        var ident uint32

        //
        err := decoder.Decode(&ident)
        if err != nil {
            log.Println(err)
            return
        }
        if ident != IDENTCODE {
            continue
        }

        var what code
        err = decoder.Decode(&what)
        if err != nil {
            log.Println(err)
            return
        }
        //
        err = decoder.Decode(&ident)
        if err != nil {
            log.Println(err)
            return
        }
        if ident != IDENTCODE {
            continue
        }

        switch what {
            case MOVE:
                var dir Dir
                err := decoder.Decode(&dir)
                if err != nil {
                    log.Println(err)
                    return
                }
                log.Println(dir)
        }
    }
    _ = encoder
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

func checkErr(err error) {
    if err != nil  {
        log.Fatal(err)
    }
}
