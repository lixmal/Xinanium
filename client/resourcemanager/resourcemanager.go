package resourcemanager

import (
    "azul3d.org/gfx.v1"
    "log"
    "image"
    _ "image/png"
    "bytes"
    "os"
)

var glslVert = []byte(`
    #version 120
    attribute vec3 Vertex;
    attribute vec2 TexCoord0;
    uniform mat4 MVP;
    varying vec2 tc0;
    void main()
    {
        tc0 = TexCoord0;
        gl_Position = MVP * vec4(Vertex, 1.0);
    }
`)
var glslFrag = []byte(`
    #version 120
    varying vec2 tc0;
    uniform sampler2D Texture0;
    uniform bool BinaryAlpha;
    void main()
    {
        gl_FragColor = texture2D(Texture0, tc0);
        if(BinaryAlpha && gl_FragColor.a < 0.5) {
            discard;
        }
    }
`)

var Textures map[string]*gfx.Texture

// load texture from memory, else read from disk and store
func Texture(path string) *gfx.Texture {
    tex, ok := Textures[path]
    if !ok {
        r, err := os.Open(path)
        if err != nil {
            log.Fatal("Could not load image '", path, "':", err)
        }

        img, _, err := image.Decode(r)
        if err != nil {
            log.Fatal("Could not load image '", path, "':", err)
        }

        tex = gfx.NewTexture()
        tex.Source = img

        Textures[path] = tex
    }

    return tex
}

// from file
func Sprite(path string) *gfx.Object {
    return loadSprite(Texture(path))
}

// from memory
func SpriteFromMemory(mem []byte) *gfx.Object {
    r := bytes.NewBuffer(mem)
    img, _, err := image.Decode(r)
    if err != nil {
        log.Fatal("Could not load image from memory: ", err)
    }

    tex := gfx.NewTexture()
    tex.MinFilter = gfx.LinearMipmapLinear
    tex.MagFilter = gfx.Linear
    tex.Format = gfx.DXT1RGBA
    tex.Source = img

    return loadSprite(tex)
}

func loadSprite (tex *gfx.Texture) *gfx.Object {
    sprite := gfx.NewObject()

    sprite.Textures = []*gfx.Texture{tex}

    mesh := gfx.NewMesh()
    mesh.Vertices = []gfx.Vec3{
        // Bottom-left triangle.
        {-1, 0, -1},
        {1, 0, -1},
        {-1, 0, 1},
        // Top-right triangle.
        {-1, 0, 1},
        {1, 0, -1},
        {1, 0, 1},
    }
    mesh.TexCoords = []gfx.TexCoordSet{
        {
            Slice: []gfx.TexCoord{
                {0, 1},
                {1, 1},
                {0, 0},
                {0, 0},
                {1, 1},
                {1, 0},
            },
        },
    }
    sprite.Meshes = []*gfx.Mesh{mesh}

    shader := gfx.NewShader("Sprite")
    shader.GLSLVert = glslVert
    shader.GLSLFrag = glslFrag
    sprite.Shader = shader

    sprite.AlphaMode = gfx.AlphaToCoverage

    return sprite
}
