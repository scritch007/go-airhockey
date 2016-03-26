package main

import (
	"encoding/binary"

	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/gl"

	"github.com/scritch007/go-tools"
)

var (
	texture                        *gl.Texture
	program                        *gl.Program
	buffer                         *gl.Buffer
	a_position_location            gl.Attrib
	a_texture_coordinates_location gl.Attrib
	u_texture_unit_location        gl.Uniform
)
var rec = f32.Bytes(binary.LittleEndian,
	-1.0, -1.0, 0.0, 0.0,
	-1.0, 1.0, 0.0, 1.0,
	1.0, -1.0, 1.0, 0.0,
	1.0, 1.0, 1.0, 1.0,
)

func on_surface_created(glctx gl.Context) {
	on_surface_changed(glctx, nil)
	glctx.ClearColor(0, 0, 0, 0)

}

/*
void on_surface_changed() {
	texture = load_png_asset_into_texture("textures/air_hockey_surface.png");
	buffer = create_vbo(sizeof(rect), rect, GL_STATIC_DRAW);
	program = build_program_from_assets("shaders/shader.vsh", "shaders/shader.fsh");

	a_position_location = glGetAttribLocation(program, "a_Position");
	a_texture_coordinates_location = glGetAttribLocation(program, "a_TextureCoordinates");
	u_texture_unit_location = glGetUniformLocation(program, "u_TextureUnit");
}
*/

func on_surface_changed(glctx gl.Context, sz *size.Event) {
	var err error
	texture, err = load_png_asset_into_texture(glctx, "textures/air_hockey_surface.png")
	if err != nil {
		tools.LOG_ERROR.Printf("Failed to create texture ", err)
		return
	}
	buffer = create_vbo(glctx, rec, gl.STATIC_DRAW)
	program, err = build_program_from_assets(glctx, "shaders/shader.vsh", "shaders/shader.fsh")
	if err != nil {
		tools.LOG_ERROR.Printf("Failed to create program ", err)
		return
	}

	a_position_location = glctx.GetAttribLocation(*program, "a_Position")
	a_texture_coordinates_location = glctx.GetAttribLocation(*program, "a_TextureCoordinates")
	u_texture_unit_location = glctx.GetUniformLocation(*program, "u_TextureUnit")
}

func onStop(glctx gl.Context) {

}

/*
void on_draw_frame() {
	glClear(GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT);

	glUseProgram(program);

	glActiveTexture(GL_TEXTURE0);
	glBindTexture(GL_TEXTURE_2D, texture);
	glUniform1i(u_texture_unit_location, 0);

	glBindBuffer(GL_ARRAY_BUFFER, buffer);
	glVertexAttribPointer(a_position_location, 2, GL_FLOAT, GL_FALSE, 4 * sizeof(GL_FLOAT), BUFFER_OFFSET(0));
	glVertexAttribPointer(a_texture_coordinates_location, 2, GL_FLOAT, GL_FALSE, 4 * sizeof(GL_FLOAT), BUFFER_OFFSET(2 * sizeof(GL_FLOAT)));
	glEnableVertexAttribArray(a_position_location);
	glEnableVertexAttribArray(a_texture_coordinates_location);
	glDrawArrays(GL_TRIANGLE_STRIP, 0, 4);

	glBindBuffer(GL_ARRAY_BUFFER, 0);
}
*/
func on_draw_frame(glctx gl.Context, sz *size.Event) {
	glctx.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	glctx.UseProgram(*program)

	glctx.ActiveTexture(gl.TEXTURE0)
	glctx.BindTexture(gl.TEXTURE_2D, *texture)
	glctx.Uniform1i(u_texture_unit_location, 0)

	glctx.BindBuffer(gl.ARRAY_BUFFER, *buffer)
	glctx.VertexAttribPointer(a_position_location, 2, gl.FLOAT, false, 16, 0) // 16 = 4 * 4 (float size)
	glctx.VertexAttribPointer(a_texture_coordinates_location, 2, gl.FLOAT, false, 16, 8)
	glctx.EnableVertexAttribArray(a_position_location)
	glctx.EnableVertexAttribArray(a_texture_coordinates_location)
	glctx.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
}
