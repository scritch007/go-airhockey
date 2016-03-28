package main

import (
	"golang.org/x/mobile/gl"
)

/*

typedef struct {
	GLuint program;

	GLint a_position_location;
	GLint a_texture_coordinates_location;
	GLint u_mvp_matrix_location;
	GLint u_texture_unit_location;
} TextureProgram;

typedef struct {
	GLuint program;

	GLint a_position_location;
	GLint u_mvp_matrix_location;
	GLint u_color_location;
} ColorProgram;


TextureProgram get_texture_program(GLuint program)
{
	return (TextureProgram) {
			program,
			glGetAttribLocation(program, "a_Position"),
			glGetAttribLocation(program, "a_TextureCoordinates"),
			glGetUniformLocation(program, "u_MvpMatrix"),
			glGetUniformLocation(program, "u_TextureUnit")};
}

ColorProgram get_color_program(GLuint program)
{
	return (ColorProgram) {
			program,
			glGetAttribLocation(program, "a_Position"),
			glGetUniformLocation(program, "u_MvpMatrix"),
			glGetUniformLocation(program, "u_Color")};
}
*/
type TextureProgram struct {
	program                        *gl.Program
	a_position_location            gl.Attrib
	a_texture_coordinates_location gl.Attrib
	u_mvp_matrix_location          gl.Uniform
	u_texture_unit_location        gl.Uniform
}
type ColorProgram struct {
	program               *gl.Program
	a_position_location   gl.Attrib
	u_mvp_matrix_location gl.Uniform
	u_color_location      gl.Uniform
}

func get_texture_program(glctx gl.Context, program *gl.Program) *TextureProgram {
	return &TextureProgram{
		program,
		glctx.GetAttribLocation(*program, "a_Position"),
		glctx.GetAttribLocation(*program, "a_TextureCoordinates"),
		glctx.GetUniformLocation(*program, "u_MvpMatrix"),
		glctx.GetUniformLocation(*program, "u_TextureUnit"),
	}
}

func get_color_program(glctx gl.Context, program *gl.Program) *ColorProgram {
	return &ColorProgram{
		program,
		glctx.GetAttribLocation(*program, "a_Position"),
		glctx.GetUniformLocation(*program, "u_MvpMatrix"),
		glctx.GetUniformLocation(*program, "u_Color"),
	}
}
