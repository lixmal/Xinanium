package network

// actions
const (
    PLAYER_MOVE = iota
    PLAYER_POS
    GET_PLAYER
    GET_PLAYER_TEX
    PLAYER_LOGIN
    PLAYER_TEX
    NO_PLAYER_TEX
    LOGIN_OK
    LOGIN_FAIL
)

const (
    PROTO_VERSION          = 1
    IDENTCODE       uint32 = 0x58696E<<8 | PROTO_VERSION
    LOGIN_MAXLENGTH        = 128
    PW_MAXLENGTH           = 128
    PACKET_MAXSIZE         = 512
)
