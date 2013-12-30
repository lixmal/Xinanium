package worldmap

import(
//	sf "bitbucket.org/krepa098/gosfml2"
    "log"
    "io/ioutil"
    "encoding/binary"
)
const TILEWIDTH = 40
const TILEHEIGHT = 40
const MAPHEADERSIZE = 1024
const MAPTITLELENGTH = 256

type Coord uint64

type WorldMap struct {
    Tiles [][]byte
    Name string
    Width Coord
    Height Coord
    ServerVersion uint32
}

func Read(filename string) *WorldMap {
    wmap := &WorldMap{}
    file, err := ioutil.ReadFile(filename)

    if err != nil {
        log.Fatal("Could not open map: ", filename,": ", err)
    }

    wmap.Name = string(file[0:MAPTITLELENGTH])

    height, bytes :=  binary.Uvarint(file[MAPTITLELENGTH:MAPTITLELENGTH+8])
    checkBytes(bytes)
    width, bytes  :=  binary.Uvarint(file[MAPTITLELENGTH+8:MAPTITLELENGTH+8*2])
    checkBytes(bytes)
    ver, bytes  :=  binary.Uvarint(file[MAPTITLELENGTH+8*2:MAPTITLELENGTH+8*3])
    checkBytes(bytes)

    wmap.Height = Coord(height)
    wmap.Width = Coord(width)
    wmap.ServerVersion = uint32(ver)
    wmap.Tiles = make([][]byte, height)
    file = file[MAPHEADERSIZE:]
    var j Coord = 0
    for i := Coord(0); i < wmap.Height * wmap.Width; i += wmap.Height {
        wmap.Tiles[j] = file[i:i+wmap.Height]
        j++
    }
    return wmap
}

func checkBytes(bytes int){
    if bytes < 8 {
       // log.Fatal("Bytes read != 8: ", bytes)
    }
}
