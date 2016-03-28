package main

import (
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/gl"

	"github.com/rkusa/gm/mat4"
	"github.com/rkusa/gm/vec3"
	"github.com/scritch007/go-tools"
)

var (
	puck_height                  float32 = 0.02
	mallet_height                float32 = 0.15
	texture                      *gl.Texture
	table                        *Table
	puck                         *Puck
	red_mallet                   *Mallet
	blue_mallet                  *Mallet
	texture_program              *TextureProgram
	color_program                *ColorProgram
	projection_matrix            *mat4.Mat4
	view_matrix                  *mat4.Mat4
	view_project_matrix          *mat4.Mat4
	model_view_projection_matrix *mat4.Mat4
)

/*
void on_surface_created() {
	glClearColor(0.0f, 0.0f, 0.0f, 0.0f);
	glEnable(GL_DEPTH_TEST);

	table = create_table(load_png_asset_into_texture("textures/air_hockey_surface.png"));

	vec4 puck_color = {0.8f, 0.8f, 1.0f, 1.0f};
	vec4 red = {1.0f, 0.0f, 0.0f, 1.0f};
	vec4 blue = {0.0f, 0.0f, 1.0f, 1.0f};

	puck = create_puck(0.06f, puck_height, 32, puck_color);
	red_mallet = create_mallet(0.08f, mallet_height, 32, red);
	blue_mallet = create_mallet(0.08f, mallet_height, 32, blue);

	texture_program = get_texture_program(build_program_from_assets("shaders/texture_shader.vsh", "shaders/texture_shader.fsh"));
	color_program = get_color_program(build_program_from_assets("shaders/color_shader.vsh", "shaders/color_shader.fsh"));
}
*/

func on_surface_created(glctx gl.Context) {
	var err error
	glctx.ClearColor(0, 0, 0, 0)
	glctx.Enable(gl.DEPTH_TEST)

	texture, err = load_png_asset_into_texture(glctx, "textures/air_hockey_surface.png")
	if err != nil {
		tools.LOG_ERROR.Printf("Failed to create texture ", err)
		return
	}
	table, err = create_table(glctx, texture)
	if err != nil {
		tools.LOG_ERROR.Printf("Failed to create table ", err)
		return
	}
	puck_color := [4]float32{0.8, 0.8, 1.0, 1.0}
	red := [4]float32{1.0, 0.0, 0.0, 1.0}
	blue := [4]float32{0.0, 0.0, 1.0, 1.0}

	puck, err = create_puck(glctx, 0.6, puck_height, 32, puck_color)
	if nil != err {
		tools.LOG_ERROR.Printf("Failed to create puck ", err)
		return
	}
	red_mallet, err = create_mallet(glctx, 0.08, mallet_height, 32, red)
	if nil != err {
		tools.LOG_ERROR.Printf("Failed to create red mallet ", err)
		return
	}
	blue_mallet, err = create_mallet(glctx, 0.08, mallet_height, 32, blue)
	if nil != err {
		tools.LOG_ERROR.Printf("Failed to create blue allet ", err)
		return
	}
	program, err := build_program_from_assets(glctx, "shaders/texture_shader.vsh", "shaders/texture_shader.fsh")
	if nil != err {
		tools.LOG_ERROR.Printf("Failed to create texture Program ", err)
		return
	}
	texture_program = get_texture_program(glctx, program)

	program, err = build_program_from_assets(glctx, "shaders/color_shader.vsh", "shaders/color_shader.fsh")
	if nil != err {
		tools.LOG_ERROR.Printf("Failed to create color Program ", err)
		return
	}
	color_program = get_color_program(glctx, program)
}

/*
void on_surface_changed(int width, int height) {
	glViewport(0, 0, width, height);
	mat4x4_perspective(projection_matrix, 45, (float) width / (float) height, 1.0f, 10.0f);
	mat4x4_look_at(view_matrix, 0.0f, 1.2f, 2.2f, 0.0f, 0.0f, 0.0f, 0.0f, 1.0f, 0.0f);
}
*/

func on_surface_changed(glctx gl.Context, sz *size.Event) {
	glctx.Viewport(0, 0, sz.WidthPx, sz.HeightPx)
	projection_matrix = mat4x4_perspective(45, float32(float32(sz.WidthPx)/float32(sz.HeightPx)), 1.0, 10.0)
	view_matrix = mat4x4_look_at(0.0, 1.2, 2.2, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0)

}

func onStop(glctx gl.Context) {

}

/*
void on_draw_frame() {
    glClear(GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT);
    mat4x4_mul(view_projection_matrix, projection_matrix, view_matrix);

	position_table_in_scene();
    draw_table(&table, &texture_program, model_view_projection_matrix);

	position_object_in_scene(0.0f, mallet_height / 2.0f, -0.4f);
	draw_mallet(&red_mallet, &color_program, model_view_projection_matrix);

	position_object_in_scene(0.0f, mallet_height / 2.0f, 0.4f);
	draw_mallet(&blue_mallet, &color_program, model_view_projection_matrix);

	// Draw the puck.
	position_object_in_scene(0.0f, puck_height / 2.0f, 0.0f);
	draw_puck(&puck, &color_program, model_view_projection_matrix);
}


*/
func on_draw_frame(glctx gl.Context, sz *size.Event) {
	glctx.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	view_project_matrix = mat4x4_mul(projection_matrix, view_matrix)

	position_table_in_scene(glctx)
	draw_table(glctx, table, texture_program, model_view_projection_matrix)

	position_object_in_scene(glctx, 0.0, mallet_height/2.0, -0.4)
	draw_mallet(glctx, red_mallet, color_program, model_view_projection_matrix)

	position_object_in_scene(glctx, 0.0, mallet_height/2.0, 0.4)
	draw_mallet(glctx, blue_mallet, color_program, model_view_projection_matrix)

	position_object_in_scene(glctx, 0.0, puck_height/2.0, 0.0)
	draw_puck(glctx, puck, color_program, model_view_projection_matrix)

}

/*
static void position_table_in_scene() {
	// The table is defined in terms of X & Y coordinates, so we rotate it
	// 90 degrees to lie flat on the XZ plane.
	mat4x4 rotated_model_matrix;
	mat4x4_identity(model_matrix);
	mat4x4_rotate_X(rotated_model_matrix, model_matrix, deg_to_radf(-90.0f));
	mat4x4_mul(model_view_projection_matrix, view_projection_matrix, rotated_model_matrix);
}
*/

func position_table_in_scene(glctx gl.Context) {
	model_matrix := mat4.Identity()
	rotated_model_matrix := model_matrix.Rotation(vec3.New(deg_to_radf(-90.0), 0.0, 0.0))
	model_view_projection_matrix = mat4x4_mul(view_project_matrix, rotated_model_matrix)
}

/*
static void position_object_in_scene(float x, float y, float z) {
	mat4x4_identity(model_matrix);
	mat4x4_translate_in_place(model_matrix, x, y, z);
	mat4x4_mul(model_view_projection_matrix, view_projection_matrix, model_matrix);
}
*/

func position_object_in_scene(glctx gl.Context, x, y, z float32) {
	model_matrix := mat4.Identity()
	model_matrix.Translate(vec3.New(x, y, z))
	model_view_projection_matrix = mat4x4_mul(view_project_matrix, model_matrix)
}
