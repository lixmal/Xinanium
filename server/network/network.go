package network

import (
    "net"
    //   "os"
    //    wm "./worldmap"
    "bytes"
    "encoding/binary"
    "encoding/gob"
    "log"
    //     "path/filepath"
    commonnet "../../common/network"
    "../config"
    "sync"
    "time"
)

const (
    TIMEOUT      = 10 * time.Second
    ERRTHRESHOLD = 5
)

type Client struct {
    net.PacketConn
    net.Addr
    ConnChan chan []byte
    Player   config.LivingEntity
}

var Clients = make(map[string]*Client)
var ClientsMutex sync.RWMutex

func (c *Client) Write(value interface{}) error {
    // TODO: check keepalive and disconnect if no _valid_ packets arrived
    // using a 0 slice here so packets stay small
    buf := make([]byte, 0, commonnet.PACKET_MAXSIZE)
    buf0 := bytes.NewBuffer(buf)
    e := gob.NewEncoder(buf0)

    if err := e.Encode(value); err != nil {
        return err
    }
    if _, err := c.WriteTo(buf[:buf0.Len()], c.Addr); err != nil {
        return err
    }
    return nil
}

func (c *Client) WriteAction(action uint16, value interface{}) error {
    // TODO: check keepalive and disconnect if no _valid_ packets arrived
    buf0 := make([]byte, 2, commonnet.PACKET_MAXSIZE)
    binary.BigEndian.PutUint16(buf0[:2], action)

    var n int
    if value != nil {
        buf := bytes.NewBuffer(buf0[2:])
        e := gob.NewEncoder(buf)

        if err := e.Encode(value); err != nil {
            return err
        }
        n = buf.Len()
    }
    if _, err := c.WriteTo(buf0[:2+n], c.Addr); err != nil {
        return err
    }
    return nil
}

func (c *Client) WriteRaw(value []byte) error {
    // TODO: check keepalive and disconnect if no _valid_ packets arrived
    var buf [commonnet.PACKET_MAXSIZE]byte

    // TODO: glue slices instead of copy
    n := copy(buf[:], value)

    if _, err := c.WriteTo(buf[:n], c.Addr); err != nil {
        return err
    }
    return nil
}

func (c *Client) Disconnect(err string) {
    // TODO: write err to client

    ClientsMutex.Lock()
    delete(Clients, c.Addr.String())
    ClientsMutex.Unlock()

    if c.Player != nil {
        c.Player.Remove()
        // TODO: remove player/client mapping and resolve circular referencing
        if c.Player.Name() != "" {
            log.Println("Disconnected " + c.Player.Name() + " from " + c.Addr.String() + ": " + err)
        } else {
            log.Println("Disconnected from " + c.Addr.String() + ": " + err)
        }
    }
}
