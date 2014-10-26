package network
// TODO: fix preamble stuck to each packet

import (
    "net"
    "log"
    "encoding/gob"
    "../config"
    "errors"
    "time"
)

const (
    PORT = "22342"
    HOST = "localhost"
    PROTO_VERSION    = 1
    IDENTCODE uint32 = 0x58696E << 8 | PROTO_VERSION
    NETWORKRETRYTIMEOUT = time.Second * 5
)

var Conn net.Conn
var encoder *gob.Encoder
var decoder *gob.Decoder

func dial() (net.Conn, error) {
    udpAddr := HOST + ":" + PORT

    conn, err := net.Dial("udp", udpAddr)
    if err != nil {
        return nil, err
    }

    return conn, nil
}

func Connect() bool {
    var err error
    for Conn, err = dial(); err != nil; {
        log.Println("Failed to connect to server:", err)
        time.Sleep(NETWORKRETRYTIMEOUT)
    }
    config.Conf.Connected = true
    encoder = gob.NewEncoder(Conn)
    decoder = gob.NewDecoder(Conn)
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
        encoder = nil
        decoder = nil
        Conn = nil
        return true
    }
    if config.Conf.Connected {
        config.Conf.Connected = false
        encoder = nil
        decoder = nil
    }
    log.Println("Failed to disconnect from server: not connected")
    return false
}

func Send(val interface{}) error {
    if !config.Conf.Connected {
        return errors.New("Not connected to server")
    }
    err := encoder.Encode(IDENTCODE)
    if err != nil {
        log.Println(err)
        Disconnect()
        return err
    }
    err = encoder.Encode(val)
    if err != nil {
        log.Println(err)
        Disconnect()
        return err
    }
    return nil
}


func Read(value interface{}) error {
    if !config.Conf.Connected {
        return errors.New("Not connected to server")
    }

    var ident uint32

    err := decoder.Decode(&ident)
    if err != nil {
        log.Println(err)
        Disconnect()
        return errors.New("value expected, received something else")
    }
    if ident != IDENTCODE {
        Disconnect()
        return errors.New("ident code expected, received something else")
    }

    err = decoder.Decode(value)
    if err != nil {
        log.Println(err)
        Disconnect()
        return errors.New("value expected, received something else")
    }

    return nil
}
