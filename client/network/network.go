package network

// TODO: fix gob preamble stuck to each packet

import (
    commonconfig "../../common/config"
    commonnet "../../common/network"
    "../config"
    "azul3d.org/lmath.v1"
    "bytes"
    "encoding/binary"
    "encoding/gob"
    "errors"
    "log"
    "net"
    "time"
)

const (
    NETWORKRETRYTIMEOUT = time.Second * 5
    NETWORKTIMEOUT      = time.Second * 5
    NETWORKREADTIMEOUT  = time.Second * 30
)

var Conn net.Conn
var encoder *gob.Encoder
var decoder *gob.Decoder

func dial() (net.Conn, error) {
    udpAddr := net.JoinHostPort(config.Network.Server.Host, config.Network.Server.Port)

    conn, err := net.DialTimeout("udp", udpAddr, NETWORKTIMEOUT)
    if err != nil {
        return nil, err
    }

    return conn, nil
}

func Login(handle, password string) error {
    if err := Write(
        commonnet.PLAYER_LOGIN,
        commonconfig.Credentials{Handle: handle, Password: []byte(password)},
    ); err != nil {
        log.Println(err)
        return err
    }
    Conn.SetReadDeadline(time.Now().Add(NETWORKTIMEOUT))
    action, _, err := ReadAction()
    if err != nil {
        log.Println(err)
        return err
    }
    if action != commonnet.LOGIN_OK {
        log.Println("Login failed")
        return err
    }
    return nil
}

func Connect() bool {
    var err error
    for Conn, err = dial(); err != nil; Conn, err = dial() {
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

func SendR(val interface{}) error {
    if !config.Conf.Connected {
        return errors.New("Not connected to server")
    }
    err := encoder.Encode(commonnet.IDENTCODE)
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

func Write(action uint16, value interface{}) error {
    // TODO: check keepalive and disconnect if no _valid_ packets arrived
    buf0 := make([]byte, 6, commonnet.PACKET_MAXSIZE)
    binary.BigEndian.PutUint32(buf0[:4], commonnet.IDENTCODE)
    binary.BigEndian.PutUint16(buf0[4:6], action)

    var n int
    if value != nil {
        buf := bytes.NewBuffer(buf0[6:])
        e := gob.NewEncoder(buf)

        if err := e.Encode(value); err != nil {
            return err
        }
        n = buf.Len()
    }
    if _, err := Conn.Write(buf0[:6+n]); err != nil {
        return err
    }
    return nil
}

func WriteRaw(value []byte) error {
    // TODO: check keepalive and disconnect if no _valid_ packets arrived
    buf := make([]byte, commonnet.PACKET_MAXSIZE)
    binary.BigEndian.PutUint32(buf[:4], commonnet.IDENTCODE)

    // TODO: glue slices instead of copy
    n := copy(buf[4:], value)

    if _, err := Conn.Write(buf[:n+4]); err != nil {
        return err
    }
    return nil
}

func ReadR(value interface{}) error {
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
    if ident != commonnet.IDENTCODE {
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

func Read(value interface{}) error {
    buf, err := ReadRaw()
    if err != nil {
        Disconnect()
        return err
    }

    if err = DecodeBuf(value, buf); err != nil {
        Disconnect()
        return err
    }

    return nil
}

func ReadAction() (uint16, []byte, error) {
    buf, err := ReadRaw()
    if err != nil {
        Disconnect()
        return 0, nil, err
    }

    return binary.BigEndian.Uint16(buf[:2]), buf[2:], nil
}

func DecodeBuf(value interface{}, buf []byte) error {
    d := gob.NewDecoder(bytes.NewBuffer(buf))
    return d.Decode(value)
}

func ReadRaw() ([]byte, error) {
    buf := make([]byte, commonnet.PACKET_MAXSIZE)
    _, err := Conn.Read(buf)
    if err != nil {
        return nil, err
    }
    return buf, nil
}

func InitListener() {
    for Conn != nil {
        Conn.SetReadDeadline(time.Now().Add(NETWORKREADTIMEOUT))
        action, buf, err := ReadAction()
        if err != nil {
            log.Println(err)
            continue
        }
        switch action {
        case commonnet.PLAYER_POS:
            var v lmath.Vec3
            if err := DecodeBuf(&v, buf); err != nil {
                log.Println(err)
                Disconnect()
            }
            if player, ok := config.Players["vik"]; ok {
                player.SetPosition(&v, true)
            }
        }
    }
}
