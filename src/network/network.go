package network

import (
    "net"
    "log"
    "strconv"
//    "encoding/gob"
    "../config"
)

const (
    PORT = 22342
    HOST = "localhost"
)

var Conn net.Conn

func dial() (net.Conn, error) {
    service := HOST + ":" + strconv.FormatInt(PORT, 10)

    tcpAddr, err := net.ResolveTCPAddr("tcp", service)
    if err != nil {
        return nil, err
    }

    conn, err := net.DialTCP("tcp", nil, tcpAddr)
    if err != nil {
        return nil, err
    }

    return conn, nil
}

func Connect() bool {
    var err error
    Conn, err = dial()
    if err != nil {
        log.Println("Failed to connect to server:", err)
        return false
    }
    config.Conf.Connected = true
    log.Println("Connected to", Conn.RemoteAddr())
    return true
}

func Disconnect() bool {
    if Conn != nil {
        if err := Conn.Close(); err != nil {
            log.Println("Failed to disconnect from server:", err)
            return false
        }
        log.Println("Disconnected from", Conn.RemoteAddr())
        config.Conf.Connected = false
        Conn = nil
        return true
    }
    if config.Conf.Connected {
        config.Conf.Connected = false
    }
    log.Println("Failed to disconnect from server: not connected")
    return false
}
