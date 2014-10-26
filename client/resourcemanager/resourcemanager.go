package resourcemanager

import (
	"azul3d.org/gfx.v1"
	"log"
    "image"
    _ "image/png"
    "bytes"
    "os"
)

var SpriteShader *gfx.Shader

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

var textures = make(map[string]*gfx.Texture)

func init() {
    shader := gfx.NewShader("Sprite")
    shader.GLSLVert = glslVert
    shader.GLSLFrag = glslFrag
    if shader == nil {
        log.Fatal("Failed to load default shader")
    }
    SpriteShader = shader
}

// load texture from memory, else read from disk and store
func Texture(path string) *gfx.Texture {
	tex, ok := textures[path]
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

		textures[path] = tex
	}

	return tex
}


// from file
func Sprite(path string, mesh []*gfx.Mesh) *gfx.Object {
    return loadSprite(Texture(path), mesh)
}

// from memory, not sure if should cache in mem, rather on disk
func SpriteFromMemory(mem *[]byte, mesh []*gfx.Mesh) *gfx.Object {
    r := bytes.NewBuffer(*mem)
    img, _, err := image.Decode(r)
    if err != nil {
        log.Fatal("Could not load image from memory: ", err)
    }

    tex := gfx.NewTexture()
    tex.MinFilter = gfx.LinearMipmapLinear
    tex.MagFilter = gfx.Linear
    tex.Format = gfx.DXT1RGBA
    tex.Source = img

    return loadSprite(tex, mesh)
}


func loadSprite (tex *gfx.Texture, mesh []*gfx.Mesh) *gfx.Object {

    sprite := gfx.NewObject()

    sprite.Textures = []*gfx.Texture{tex}

    //imgbnd := tex.Source.Bounds()
    //aspect := float32(imgbnd.Dx()) / float32(imgbnd.Dy())
    //var height float32 = 40.0
    sprite.Shader = SpriteShader

    sprite.AlphaMode = gfx.AlphaToCoverage

    sprite.Meshes = mesh

    return sprite
}

func TexCoords(u, v, s, t float32) []gfx.TexCoord {
    return []gfx.TexCoord{
        // Left triangle.
        {u, v},
        {u, t},
        {s, t},
        // Right triangle.
        {u, v},
        {s, t},
        {s, v},
    }
}

func Mesh (w, h, wPart, hPart float32) []*gfx.Mesh {
    mesh := gfx.NewMesh()
    mesh.Vertices = []gfx.Vec3{
        // Left triangle.
        {-w, 0, h}, // Left-Top
        {-w, 0, -h}, // Left-Bottom
        {w, 0, -h}, // Right-Bottom

        // Right triangle.
        {-w, 0, h}, // Left-Top
        {w, 0, -h}, // Right-Bottom
        {w, 0, h}, // Right-Top
    }
    mesh.TexCoords = []gfx.TexCoordSet{
        {
            Slice: TexCoords(0, 0, wPart, hPart),
        },
    }
    return []*gfx.Mesh{mesh}
}

