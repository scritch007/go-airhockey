package main

import (
	"golang.org/x/mobile/gl"
)

func load_texture(glctx gl.Context, width, height int, glType gl.Enum, pixels []byte) *gl.Texture {
	t := glctx.CreateTexture()
	glctx.BindTexture(gl.TEXTURE_2D, t)
	glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	glctx.TexImage2D(gl.TEXTURE_2D, 0, width, height, glType, gl.UNSIGNED_BYTE, pixels)
	glctx.GenerateMipmap(gl.TEXTURE_2D)
	return &t

}
