package main

import (
	"encoding/binary"
	"math"

	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/gl"

	"github.com/rkusa/gm/mat4"
)

/*
typedef struct {
	GLuint texture;
	GLuint buffer;
} Table;

typedef struct {
	vec4 color;
	GLuint buffer;
	int num_points;
} Puck;

typedef struct {
	vec4 color;
	GLuint buffer;
	int num_points;
} Mallet;
*/
type Table struct {
	texture *gl.Texture
	buffer  *gl.Buffer
}

type Puck struct {
	color      [4]float32
	buffer     *gl.Buffer
	num_points int
}

type Mallet struct {
	color      [4]float32
	buffer     *gl.Buffer
	num_points int
}

var table_data = f32.Bytes(binary.LittleEndian,
	0.0, 0.0, 0.5, 0.5,
	-0.5, -0.8, 0.0, 0.9,
	0.5, -0.8, 1.0, 0.9,
	0.5, 0.8, 1.0, 0.1,
	-0.5, 0.8, 0.0, 0.1,
	-0.5, -0.8, 0.0, 0.9,
)

func create_table(glctx gl.Context, texture *gl.Texture) (*Table, error) {
	table := new(Table)
	table.texture = texture
	table.buffer = create_vbo(glctx, table_data, gl.STATIC_DRAW)
	return table, nil
}

/*
void draw_table(const Table* table, const TextureProgram* texture_program, mat4x4 m)
{
	glUseProgram(texture_program->program);

	glActiveTexture(GL_TEXTURE0);
	glBindTexture(GL_TEXTURE_2D, table->texture);
	glUniformMatrix4fv(texture_program->u_mvp_matrix_location, 1, GL_FALSE, (GLfloat*)m);
	glUniform1i(texture_program->u_texture_unit_location, 0);

	glBindBuffer(GL_ARRAY_BUFFER, table->buffer);
	glVertexAttribPointer(texture_program->a_position_location, 2, GL_FLOAT, GL_FALSE, 4 * sizeof(GL_FLOAT), BUFFER_OFFSET(0));
	glVertexAttribPointer(texture_program->a_texture_coordinates_location, 2, GL_FLOAT, GL_FALSE, 4 * sizeof(GL_FLOAT), BUFFER_OFFSET(2 * sizeof(GL_FLOAT)));
	glEnableVertexAttribArray(texture_program->a_position_location);
	glEnableVertexAttribArray(texture_program->a_texture_coordinates_location);
	glDrawArrays(GL_TRIANGLE_FAN, 0, 6);

	glBindBuffer(GL_ARRAY_BUFFER, 0);
}
*/

func draw_table(glctx gl.Context, table *Table, texture_program *TextureProgram, m *mat4.Mat4) {
	glctx.UseProgram(*texture_program.program)

	glctx.ActiveTexture(gl.TEXTURE0)
	glctx.BindTexture(gl.TEXTURE_2D, *table.texture)
	tmp := [16]float32(*m)
	glctx.UniformMatrix4fv(texture_program.u_mvp_matrix_location, tmp[:])
	glctx.Uniform1i(texture_program.u_texture_unit_location, 0)

	glctx.BindBuffer(gl.ARRAY_BUFFER, *table.buffer)
	glctx.VertexAttribPointer(texture_program.a_position_location, 2, gl.FLOAT, false, 4*4, 0)
	glctx.VertexAttribPointer(texture_program.a_texture_coordinates_location, 2, gl.FLOAT, false, 4*4, 2*4)
	glctx.EnableVertexAttribArray(texture_program.a_position_location)
	glctx.EnableVertexAttribArray(texture_program.a_texture_coordinates_location)
	glctx.DrawArrays(gl.TRIANGLE_FAN, 0, 6)
}

/*
static inline int size_of_circle_in_vertices(int num_points) {
	return 1 + (num_points + 1);
}

static inline int size_of_open_cylinder_in_vertices(int num_points) {
	return (num_points + 1) * 2;
}
*/
func size_of_circle_in_vertices(num_points int) int {
	return 1 + (num_points + 1)
}

func size_of_open_cylinder_in_vertices(num_points int) int {
	return (num_points + 1) * 2
}

/*
static inline int gen_circle(float* out, int offset, float center_x, float center_y, float center_z, float radius, int num_points)
{
	out[offset++] = center_x;
	out[offset++] = center_y;
	out[offset++] = center_z;

	int i;
	for (i = 0; i <= num_points; ++i) {
		float angle_in_radians = ((float) i / (float) num_points) * ((float) M_PI * 2.0f);
		out[offset++] = center_x + radius * cos(angle_in_radians);
		out[offset++] = center_y;
		out[offset++] = center_z + radius * sin(angle_in_radians);
	}

	return offset;
}*/

func gen_circle(out []float32, offset int, center_x, center_y, center_z, radius float32, num_points int) int {
	out[offset] = center_x
	offset++
	out[offset] = center_y
	offset++
	out[offset] = center_z
	offset++
	for i := 0; i <= num_points; i++ {
		angle_in_radians := ((float32)(i) / (float32)(num_points)) * 2.0 * (float32)(math.Pi)
		out[offset] = center_x + radius*float32(math.Cos(float64(angle_in_radians)))
		offset++
		out[offset] = center_y
		offset++
		out[offset] = center_z + radius*float32(math.Sin(float64(angle_in_radians)))
		offset++
	}
	return offset
}

/*
static inline int gen_cylinder(float* out, int offset, float center_x, float center_y, float center_z, float height, float radius, int num_points)
{
	const float y_start = center_y - (height / 2.0f);
	const float y_end = center_y + (height / 2.0f);

	int i;
	for (i = 0; i <= num_points; i++) {
		float angle_in_radians = ((float) i / (float) num_points) * ((float) M_PI * 2.0f);

		float x_position = center_x + radius * cos(angle_in_radians);
		float z_position = center_z + radius * sin(angle_in_radians);

		out[offset++] = x_position;
		out[offset++] = y_start;
		out[offset++] = z_position;

		out[offset++] = x_position;
		out[offset++] = y_end;
		out[offset++] = z_position;
	}

	return offset;
}*/
func gen_cylinder(out []float32, offset int, center_x, center_y, center_z, height, radius float32, num_points int) int {
	y_start := center_y - (height / 2.0)
	y_end := center_y + (height / 2.0)

	for i := 0; i <= num_points; i++ {
		angle_in_radians := (float32(i) / float32(num_points)) * (math.Pi * 2.0)
		x_position := center_x + radius*float32(math.Cos(float64(angle_in_radians)))
		z_position := center_z + radius*float32(math.Sin(float64(angle_in_radians)))

		out[offset] = x_position
		offset++
		out[offset] = y_start
		offset++
		out[offset] = z_position
		offset++

		out[offset] = x_position
		offset++
		out[offset] = y_end
		offset++
		out[offset] = z_position
		offset++
	}
	return offset
}

/*
Puck create_puck(float radius, float height, int num_points, vec4 color)
{
	float data[(size_of_circle_in_vertices(num_points) + size_of_open_cylinder_in_vertices(num_points)) * 3];

	int offset = gen_circle(data, 0, 0.0f, height / 2.0f, 0.0f, radius, num_points);
	gen_cylinder(data, offset, 0.0f, 0.0f, 0.0f, height, radius, num_points);

	return (Puck) {{color[0], color[1], color[2], color[3]},
				   create_vbo(sizeof(data), data, GL_STATIC_DRAW),
				   num_points};
}
*/
func create_puck(glctx gl.Context, radius, height float32, num_points int, color [4]float32) (*Puck, error) {
	puck := new(Puck)
	puck.color = color
	puck.num_points = num_points
	data := make([]float32, (size_of_circle_in_vertices(num_points)+size_of_open_cylinder_in_vertices(num_points))*3)
	offset := gen_circle(data, 0, 0.0, height/2.0, 0.0, radius, num_points)
	offset = gen_cylinder(data, offset, 0.0, 0.0, 0.0, height, radius, num_points)
	puck.buffer = create_vbo(glctx, f32.Bytes(binary.LittleEndian, data...), gl.STATIC_DRAW)
	return puck, nil
}

/*
void draw_puck(const Puck* puck, const ColorProgram* color_program, mat4x4 m)
{
	glUseProgram(color_program->program);

	glUniformMatrix4fv(color_program->u_mvp_matrix_location, 1, GL_FALSE, (GLfloat*)m);
	glUniform4fv(color_program->u_color_location, 1, puck->color);

	glBindBuffer(GL_ARRAY_BUFFER, puck->buffer);
	glVertexAttribPointer(color_program->a_position_location, 3, GL_FLOAT, GL_FALSE, 0, BUFFER_OFFSET(0));
	glEnableVertexAttribArray(color_program->a_position_location);

	int circle_vertex_count = size_of_circle_in_vertices(puck->num_points);
	int cylinder_vertex_count = size_of_open_cylinder_in_vertices(puck->num_points);

	glDrawArrays(GL_TRIANGLE_FAN, 0, circle_vertex_count);
	glDrawArrays(GL_TRIANGLE_STRIP, circle_vertex_count, cylinder_vertex_count);
	glBindBuffer(GL_ARRAY_BUFFER, 0);
}
*/
func draw_puck(glctx gl.Context, puck *Puck, color_program *ColorProgram, m *mat4.Mat4) {
	glctx.UseProgram(*color_program.program)

	tmp := [16]float32(*m)
	glctx.UniformMatrix4fv(color_program.u_mvp_matrix_location, tmp[:])
	glctx.Uniform4fv(color_program.u_color_location, puck.color[:])

	glctx.BindBuffer(gl.ARRAY_BUFFER, *puck.buffer)
	glctx.VertexAttribPointer(color_program.a_position_location, 3, gl.FLOAT, false, 0, 0)
	glctx.EnableVertexAttribArray(color_program.a_position_location)

	circle_vertex_count := size_of_circle_in_vertices(puck.num_points)
	cylinder_vertex_count := size_of_open_cylinder_in_vertices(puck.num_points)

	glctx.DrawArrays(gl.TRIANGLE_FAN, 0, circle_vertex_count)
	glctx.DrawArrays(gl.TRIANGLE_STRIP, circle_vertex_count, cylinder_vertex_count)
}

/*
Mallet create_mallet(float radius, float height, int num_points, vec4 color)
{
	float data[(size_of_circle_in_vertices(num_points) * 2 + size_of_open_cylinder_in_vertices(num_points) * 2) * 3];

	float base_height = height * 0.25f;
	float handle_height = height * 0.75f;
	float handle_radius = radius / 3.0f;

	int offset = gen_circle(data, 0, 0.0f, -base_height, 0.0f, radius, num_points);
	offset = gen_circle(data, offset, 0.0f, height * 0.5f, 0.0f, handle_radius, num_points);
	offset = gen_cylinder(data, offset, 0.0f, -base_height - base_height / 2.0f, 0.0f, base_height, radius, num_points);
	gen_cylinder(data, offset, 0.0f, height * 0.5f - handle_height / 2.0f, 0.0f, handle_height, handle_radius, num_points);

	return (Mallet) {{color[0], color[1], color[2], color[3]},
					 create_vbo(sizeof(data), data, GL_STATIC_DRAW),
				     num_points};
}
*/
func create_mallet(glctx gl.Context, radius, height float32, num_points int, color [4]float32) (*Mallet, error) {
	data := make([]float32, (size_of_circle_in_vertices(num_points)*2+size_of_open_cylinder_in_vertices(num_points)*2)*3)
	base_height := height * 0.25
	handle_height := height * 0.75
	handle_radius := radius / 3.0

	offset := gen_circle(data, 0, 0.0, -base_height, 0.0, radius, num_points)
	offset = gen_circle(data, offset, 0.0, height*0.5, 0.0, handle_radius, num_points)
	offset = gen_cylinder(data, offset, 0, -base_height-(base_height/2.0), 0.0, base_height, radius, num_points)
	offset = gen_cylinder(data, offset, 0.0, height*0.5-handle_height/2.0, 0.0, handle_height, handle_radius, num_points)

	mallet := Mallet{
		color:      color,
		buffer:     create_vbo(glctx, f32.Bytes(binary.LittleEndian, data...), gl.STATIC_DRAW),
		num_points: num_points,
	}
	return &mallet, nil
}

/*
void draw_mallet(const Mallet* mallet, const ColorProgram* color_program, mat4x4 m)
{
	glUseProgram(color_program->program);

	glUniformMatrix4fv(color_program->u_mvp_matrix_location, 1, GL_FALSE, (GLfloat*)m);
	glUniform4fv(color_program->u_color_location, 1, mallet->color);

	glBindBuffer(GL_ARRAY_BUFFER, mallet->buffer);
	glVertexAttribPointer(color_program->a_position_location, 3, GL_FLOAT, GL_FALSE, 0, BUFFER_OFFSET(0));
	glEnableVertexAttribArray(color_program->a_position_location);

	int circle_vertex_count = size_of_circle_in_vertices(mallet->num_points);
	int cylinder_vertex_count = size_of_open_cylinder_in_vertices(mallet->num_points);
	int start_vertex = 0;

	glDrawArrays(GL_TRIANGLE_FAN, start_vertex, circle_vertex_count); start_vertex += circle_vertex_count;
	glDrawArrays(GL_TRIANGLE_FAN, start_vertex, circle_vertex_count); start_vertex += circle_vertex_count;
	glDrawArrays(GL_TRIANGLE_STRIP, start_vertex, cylinder_vertex_count); start_vertex += cylinder_vertex_count;
	glDrawArrays(GL_TRIANGLE_STRIP, start_vertex, cylinder_vertex_count);
	glBindBuffer(GL_ARRAY_BUFFER, 0);
}
*/

func draw_mallet(glctx gl.Context, mallet *Mallet, color_program *ColorProgram, m *mat4.Mat4) {
	glctx.UseProgram(*color_program.program)

	tmp := [16]float32(*m)
	glctx.UniformMatrix4fv(color_program.u_mvp_matrix_location, tmp[:])
	glctx.Uniform4fv(color_program.u_color_location, mallet.color[:])

	glctx.BindBuffer(gl.ARRAY_BUFFER, *mallet.buffer)
	glctx.VertexAttribPointer(color_program.a_position_location, 3, gl.FLOAT, false, 0, 0)
	glctx.EnableVertexAttribArray(color_program.a_position_location)

	circle_vertex_count := size_of_circle_in_vertices(mallet.num_points)
	cylinder_vertex_count := size_of_open_cylinder_in_vertices(mallet.num_points)
	start_vertex := 0

	glctx.DrawArrays(gl.TRIANGLE_FAN, start_vertex, circle_vertex_count)
	start_vertex += circle_vertex_count
	glctx.DrawArrays(gl.TRIANGLE_FAN, start_vertex, circle_vertex_count)
	start_vertex += circle_vertex_count
	glctx.DrawArrays(gl.TRIANGLE_STRIP, start_vertex, cylinder_vertex_count)
	start_vertex += cylinder_vertex_count
	glctx.DrawArrays(gl.TRIANGLE_STRIP, start_vertex, cylinder_vertex_count)

}
