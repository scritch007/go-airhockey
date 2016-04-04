package main

import (
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/gl"

	"github.com/scritch007/gm/mat4"
	"github.com/scritch007/gm/vec3"
	"github.com/scritch007/gm/vec4"
	"github.com/scritch007/go-tools"
)

var (
	puck_height   float32 = 0.02
	puck_radius   float32 = 0.06
	mallet_height float32 = 0.15
	mallet_radius float32 = 0.08

	left_bound  float32 = -0.5
	right_bound float32 = 0.5
	far_bound   float32 = -0.8
	near_bound  float32 = 0.8

	texture                         *gl.Texture
	table                           *Table
	puck                            *Puck
	red_mallet                      *Mallet
	blue_mallet                     *Mallet
	texture_program                 *TextureProgram
	color_program                   *ColorProgram
	projection_matrix               *mat4.Mat4
	view_matrix                     *mat4.Mat4
	view_project_matrix             *mat4.Mat4
	model_view_projection_matrix    *mat4.Mat4
	inverted_view_projection_matrix *mat4.Mat4
	mallet_pressed                  bool

	puck_vector                   *vec3.Vec3
	puck_position                 *vec3.Vec3
	blue_mallet_position          *vec3.Vec3
	previous_blue_mallet_position *vec3.Vec3
)

/*
void on_touch_press(float normalized_x, float normalized_y) {
	Ray ray = convert_normalized_2D_point_to_ray(normalized_x, normalized_y);

	// Now test if this ray intersects with the mallet by creating a
	// bounding sphere that wraps the mallet.
	Sphere mallet_bounding_sphere = (Sphere) {
	   {blue_mallet_position[0],
		blue_mallet_position[1],
		blue_mallet_position[2]},
	mallet_height / 2.0f};

	// If the ray intersects (if the user touched a part of the screen that
	// intersects the mallet's bounding sphere), then set malletPressed =
	// true.
	mallet_pressed = sphere_intersects_ray(mallet_bounding_sphere, ray);
}
*/
func on_touch_press(glctx gl.Context, normalized_x, normalized_y float32) {
	ray := convert_normalize_2D_point_to_ray(glctx, normalized_x, normalized_y)
	mallet_bounding_sphere := Sphere{blue_mallet_position, mallet_height / 2.0}

	mallet_pressed = sphere_intersects_ray(&mallet_bounding_sphere, ray)
}

/*
static Ray convert_normalized_2D_point_to_ray(float normalized_x, float normalized_y) {
	// We'll convert these normalized device coordinates into world-space
	// coordinates. We'll pick a point on the near and far planes, and draw a
	// line between them. To do this transform, we need to first multiply by
	// the inverse matrix, and then we need to undo the perspective divide.
	vec4 near_point_ndc = {normalized_x, normalized_y, -1, 1};
	vec4 far_point_ndc = {normalized_x, normalized_y,  1, 1};

    vec4 near_point_world, far_point_world;
    mat4x4_mul_vec4(near_point_world, inverted_view_projection_matrix, near_point_ndc);
    mat4x4_mul_vec4(far_point_world, inverted_view_projection_matrix, far_point_ndc);

	// Why are we dividing by W? We multiplied our vector by an inverse
	// matrix, so the W value that we end up is actually the *inverse* of
	// what the projection matrix would create. By dividing all 3 components
	// by W, we effectively undo the hardware perspective divide.
    divide_by_w(near_point_world);
    divide_by_w(far_point_world);

	// We don't care about the W value anymore, because our points are now
	// in world coordinates.
	vec3 near_point_ray = {near_point_world[0], near_point_world[1], near_point_world[2]};
	vec3 far_point_ray = {far_point_world[0], far_point_world[1], far_point_world[2]};
	vec3 vector_between;
	vec3_sub(vector_between, far_point_ray, near_point_ray);
	return (Ray) {
		{near_point_ray[0], near_point_ray[1], near_point_ray[2]},
		{vector_between[0], vector_between[1], vector_between[2]}};
}
*/
func convert_normalize_2D_point_to_ray(glctx gl.Context, normalized_x, normalized_y float32) *Ray {
	near_point_ndc := vec4.New(normalized_x, normalized_y, -1, 1)
	far_point_ndc := vec4.New(normalized_x, normalized_y, 1, 1)

	near_point_world := near_point_ndc.Clone().Transform(inverted_view_projection_matrix)
	far_point_world := far_point_ndc.Clone().Transform(inverted_view_projection_matrix)

	//tools.LOG_DEBUG.Printf("npw %s  fpw%s\n", near_point_world, far_point_world)

	divide_by_w(near_point_world)
	divide_by_w(far_point_world)

	near_point_ray := vec3.New(near_point_world[0], near_point_world[1], near_point_world[2])
	far_point_ray := vec3.New(far_point_world[0], far_point_world[1], far_point_world[2])
	vector_between := far_point_ray.Clone().Sub(near_point_ray)
	return &Ray{
		near_point_ray,
		vector_between,
	}
}

/*
static void divide_by_w(vec4 vector) {
	vector[0] /= vector[3];
	vector[1] /= vector[3];
	vector[2] /= vector[3];
}
*/
func divide_by_w(vector *vec4.Vec4) {
	vector[0] /= vector[3]
	vector[1] /= vector[3]
	vector[2] /= vector[3]
}

/*
void on_touch_drag(float normalized_x, float normalized_y) {
	if (mallet_pressed == 0)
		return;

	Ray ray = convert_normalized_2D_point_to_ray(normalized_x, normalized_y);
	// Define a plane representing our air hockey table.
	Plane plane = (Plane) {{0, 0, 0}, {0, 1, 0}};

	// Find out where the touched point intersects the plane
	// representing our table. We'll move the mallet along this plane.
	vec3 touched_point;
	ray_intersection_point(touched_point, ray, plane);

	memcpy(previous_blue_mallet_position, blue_mallet_position, sizeof(blue_mallet_position));

	// Clamp to bounds
	blue_mallet_position[0] = clamp(touched_point[0], left_bound + mallet_radius, right_bound - mallet_radius);
	blue_mallet_position[1] = mallet_height / 2.0f;
	blue_mallet_position[2] = clamp(touched_point[2], 0.0f + mallet_radius, near_bound - mallet_radius);

	// Now test if mallet has struck the puck.
	vec3 mallet_to_puck;
	vec3_sub(mallet_to_puck, puck_position, blue_mallet_position);
	float distance = vec3_len(mallet_to_puck);

	if (distance < (puck_radius + mallet_radius)) {
		// The mallet has struck the puck. Now send the puck flying
		// based on the mallet velocity.
		vec3_sub(puck_vector, blue_mallet_position, previous_blue_mallet_position);
	}
}
*/

func on_touch_drag(glctx gl.Context, normalized_x, normalized_y float32) {
	if !mallet_pressed {
		return
	}
	ray := convert_normalize_2D_point_to_ray(glctx, normalized_x, normalized_y)
	plane := Plane{vec3.New(0, 0, 0), vec3.New(0, 1, 0)}

	touched_point := ray_intersection_point(ray, &plane)
	previous_blue_mallet_position = blue_mallet_position.Clone()

	blue_mallet_position[0] = clamp(touched_point[0], left_bound+mallet_radius, right_bound-mallet_radius)
	blue_mallet_position[1] = mallet_height / 2.0
	blue_mallet_position[2] = clamp(touched_point[2], 0.0+mallet_radius, near_bound-mallet_radius)

	mallet_to_puck := puck_position.Clone().Sub(blue_mallet_position)

	distance := mallet_to_puck.Len()

	if distance < (puck_radius + mallet_radius) {
		puck_vector = blue_mallet_position.Clone().Sub(previous_blue_mallet_position)
	}
}

/*
void on_surface_created() {
	glClearColor(0.0f, 0.0f, 0.0f, 0.0f);
	glEnable(GL_DEPTH_TEST);

	table = create_table(load_png_asset_into_texture("textures/air_hockey_surface.png"));

	vec4 puck_color = {0.8f, 0.8f, 1.0f, 1.0f};
	vec4 red = {1.0f, 0.0f, 0.0f, 1.0f};
	vec4 blue = {0.0f, 0.0f, 1.0f, 1.0f};

	puck = create_puck(puck_radius, puck_height, 32, puck_color);
	red_mallet = create_mallet(mallet_radius, mallet_height, 32, red);
	blue_mallet = create_mallet(mallet_radius, mallet_height, 32, blue);

	blue_mallet_position[0] = 0;
	blue_mallet_position[1] = mallet_height / 2.0f;
	blue_mallet_position[2] = 0.4f;
	puck_position[0] = 0;
	puck_position[1] = puck_height / 2.0f;
	puck_position[2] = 0;
	puck_vector[0] = 0;
	puck_vector[1] = 0;
	puck_vector[2] = 0;

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

	puck, err = create_puck(glctx, puck_radius, puck_height, 32, puck_color)
	if nil != err {
		tools.LOG_ERROR.Printf("Failed to create puck ", err)
		return
	}
	red_mallet, err = create_mallet(glctx, mallet_radius, mallet_height, 32, red)
	if nil != err {
		tools.LOG_ERROR.Printf("Failed to create red mallet ", err)
		return
	}
	blue_mallet, err = create_mallet(glctx, mallet_radius, mallet_height, 32, blue)
	if nil != err {
		tools.LOG_ERROR.Printf("Failed to create blue allet ", err)
		return
	}

	blue_mallet_position = vec3.New(0.0, mallet_height/2.0, 0.4)
	puck_position = vec3.New(0.0, puck_height/2.0, 0.0)
	puck_vector = vec3.New(0, 0, 0)

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

func on_surface_changed(glctx gl.Context, sz *size.Event) bool {

	//tools.LOG_DEBUG.Printf("%s\n", *sz)
	if sz.WidthPx == 0 {
		tools.LOG_ERROR.Println("Invalid Width")
		return false
	}
	if sz.HeightPx == 0 {
		tools.LOG_ERROR.Println("Invalid Height")
		return false
	}
	glctx.Viewport(0, 0, sz.WidthPx, sz.HeightPx)
	projection_matrix = mat4x4_perspective(45, float32(float32(sz.WidthPx)/float32(sz.HeightPx)), 1.0, 10.0)
	view_matrix = mat4x4_look_at(0.0, 1.2, 2.2, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0)

	return true
}

func onStop(glctx gl.Context) {

}

/*
void on_draw_frame() {
    glClear(GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT);

    // Translate the puck by its vector
	vec3_add(puck_position, puck_position, puck_vector);

	// If the puck struck a side, reflect it off that side.
	if (puck_position[0] < left_bound + puck_radius
	 || puck_position[0] > right_bound - puck_radius) {
		puck_vector[0] = -puck_vector[0];
		vec3_scale(puck_vector, puck_vector, 0.9f);
	}
	if (puck_position[2] < far_bound + puck_radius
	 || puck_position[2] > near_bound - puck_radius) {
		puck_vector[2] = -puck_vector[2];
		vec3_scale(puck_vector, puck_vector, 0.9f);
	}

	// Clamp the puck position.
	puck_position[0] = clamp(puck_position[0], left_bound + puck_radius, right_bound - puck_radius);
	puck_position[2] = clamp(puck_position[2], far_bound + puck_radius, near_bound - puck_radius);

	// Friction factor
	vec3_scale(puck_vector, puck_vector, 0.99f);

    mat4x4_mul(view_projection_matrix, projection_matrix, view_matrix);
    mat4x4_invert(inverted_view_projection_matrix, view_projection_matrix);

	position_table_in_scene();
    draw_table(&table, &texture_program, model_view_projection_matrix);

	position_object_in_scene(0.0f, mallet_height / 2.0f, -0.4f);
	draw_mallet(&red_mallet, &color_program, model_view_projection_matrix);

	position_object_in_scene(blue_mallet_position[0], blue_mallet_position[1], blue_mallet_position[2]);
	draw_mallet(&blue_mallet, &color_program, model_view_projection_matrix);

	// Draw the puck.
	position_object_in_scene(puck_position[0], puck_position[1], puck_position[2]);
	draw_puck(&puck, &color_program, model_view_projection_matrix);
}
*/
func on_draw_frame(glctx gl.Context, sz *size.Event) {
	glctx.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	puck_position.Add(puck_vector)

	if (puck_position[0] < left_bound+puck_radius) || (puck_position[0] > right_bound-puck_radius) {
		puck_vector[0] = -puck_vector[0]
		puck_vector.Mul(0.9)
	}
	if (puck_position[2] < far_bound+puck_radius) || (puck_position[2] > near_bound-puck_radius) {
		puck_vector[2] = -puck_vector[2]
		puck_vector.Mul(0.9)
	}

	puck_position[0] = clamp(puck_position[0], left_bound+puck_radius, right_bound-puck_radius)
	puck_position[2] = clamp(puck_position[2], far_bound+puck_radius, near_bound-puck_radius)

	puck_vector.Mul(0.99)

	//tools.LOG_DEBUG.Printf("projection %s", *projection_matrix)
	//tools.LOG_DEBUG.Printf("view %s", *view_matrix)
	view_project_matrix = mat4x4_mul(projection_matrix, view_matrix)
	inverted_view_projection_matrix = view_project_matrix.Clone().Invert()

	position_table_in_scene(glctx)

	//tools.LOG_DEBUG.Printf("model view %s", *model_view_projection_matrix)

	//model_view_projection_matrix = mat4.Identity()
	draw_table(glctx, table, texture_program, model_view_projection_matrix)

	position_object_in_scene(glctx, 0.0, mallet_height/2.0, -0.4)
	//model_view_projection_matrix = mat4.Identity()
	draw_mallet(glctx, red_mallet, color_program, model_view_projection_matrix)

	position_object_in_scene(glctx, blue_mallet_position[0], blue_mallet_position[1], blue_mallet_position[2])
	//model_view_projection_matrix = mat4.Identity()
	draw_mallet(glctx, blue_mallet, color_program, model_view_projection_matrix)

	position_object_in_scene(glctx, puck_position[0], puck_position[1], puck_position[2])
	//model_view_projection_matrix = mat4.Identity()
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
