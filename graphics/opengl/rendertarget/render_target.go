package rendertarget

// #cgo LDFLAGS: -framework OpenGL
//
// #include <OpenGL/gl.h>
import "C"
import (
	"github.com/hajimehoshi/go-ebiten/graphics/opengl/texture"
	gtexture "github.com/hajimehoshi/go-ebiten/graphics/texture"
	"github.com/hajimehoshi/go-ebiten/graphics/rendertarget"
)

type Framebuffer C.GLuint

func createFramebuffer(nativeTexture C.GLuint) C.GLuint {
	framebuffer := C.GLuint(0)
	C.glGenFramebuffers(1, &framebuffer)

	origFramebuffer := C.GLint(0)
	C.glGetIntegerv(C.GL_FRAMEBUFFER_BINDING, &origFramebuffer)

	C.glBindFramebuffer(C.GL_FRAMEBUFFER, framebuffer)
	defer C.glBindFramebuffer(C.GL_FRAMEBUFFER, C.GLuint(origFramebuffer))

	C.glFramebufferTexture2D(C.GL_FRAMEBUFFER, C.GL_COLOR_ATTACHMENT0,
		C.GL_TEXTURE_2D, nativeTexture, 0)
	if C.glCheckFramebufferStatus(C.GL_FRAMEBUFFER) !=
		C.GL_FRAMEBUFFER_COMPLETE {
		panic("creating framebuffer failed")
	}

	return framebuffer
}

func New(width, height int, filter texture.Filter) (
	*rendertarget.RenderTarget, *gtexture.Texture, error) {
	tex, err := texture.New(width, height, filter)
	if err != nil {
		return nil, nil, err
	}
	f := func(native interface{}) interface{}{
		return createFramebuffer(C.GLuint(native.(texture.Native)))
	}
	framebuffer := tex.CreateFramebuffer(f)
	return rendertarget.NewWithFramebuffer(tex,
		Framebuffer(framebuffer.(C.GLuint))), tex, nil
}

func NewWithFramebuffer(width, height int, framebuffer Framebuffer,
	filter texture.Filter) (*rendertarget.RenderTarget, error) {
	tex, err := texture.New(width, height, filter)
	if err != nil {
		return nil, err
	}
	return rendertarget.NewWithFramebuffer(tex, framebuffer), nil
}