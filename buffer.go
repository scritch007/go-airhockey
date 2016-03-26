package main

import (
	"golang.org/x/mobile/gl"
)

func create_vbo(glctx gl.Context, data []byte, glType gl.Enum) *gl.Buffer {
	buffer := glctx.CreateBuffer()
	glctx.BindBuffer(gl.ARRAY_BUFFER, buffer)
	glctx.BufferData(gl.ARRAY_BUFFER, rec, glType)

	return &buffer
}
