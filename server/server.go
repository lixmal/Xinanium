package main

import (
    "net"
 //   "os"
    "log"
    "strconv"
    "encoding/gob"
)

const (
    PORT = 22342
    LADDRESS = "0.0.0.0"
    TIMEOUT = 60 * 1000

)

func main() {
    service := LADDRESS + ":" + strconv.FormatInt(PORT, 10)
    tcpAddr, err := net.ResolveTCPAddr("tcp", service)
    checkErr(err)

    listener, err := net.ListenTCP("tcp", tcpAddr)
    checkErr(err)

    for {
        conn, err := listener.Accept()
        if err != nil {
            continue
        }
    //    conn.SetTimeout(TIMEOUT)
        go handleClient(conn)
    }
}

func handleClient(conn net.Conn) {
    defer conn.Close()
    encoder := gob.NewEncoder(conn)
    decoder := gob.NewDecoder(conn)
    var buf [512]byte
    _ = buf
    _ = encoder
    _ = decoder
    // map player to connection
    for {

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

func checkErr(err error) {
    if err != nil  {
        log.Fatal(err)
    }
}
