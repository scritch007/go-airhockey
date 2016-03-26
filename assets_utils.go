package main

import (
	"fmt"
	"image"
	_ "image/png"
	"io/ioutil"

	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/gl"

	"github.com/scritch007/go-tools"
)

/*
GLuint load_png_asset_into_texture(const char* relative_path) {
	assert(relative_path != NULL);

	const FileData png_file = get_asset_data(relative_path);
	const RawImageData raw_image_data = get_raw_image_data_from_png(png_file.data, png_file.data_length);
	const GLuint texture_object_id = load_texture(
		raw_image_data.width, raw_image_data.height, raw_image_data.gl_color_format, raw_image_data.data);

	release_raw_image_data(&raw_image_data);
	release_asset_data(&png_file);

	return texture_object_id;
}*/
func load_png_asset_into_texture(glctx gl.Context, texture_path string) (*gl.Texture, error) {
	a, err := asset.Open(texture_path)
	if err != nil {
		return nil, tools.LogError(fmt.Sprintf("Failed to open texture %s: %s\n", texture_path, err))
	}
	defer a.Close()

	m, _, err := image.Decode(a)
	if err != nil {
		return nil, tools.LogError(fmt.Sprintf("Failed to Decode texture %s %s\n", texture_path, err))
	}
	rgbImage, ok := m.(*image.RGBA)
	if !ok {
		return nil, tools.LogError("This wasn't an RGBA image")
	}

	tools.LOG_DEBUG.Printf("Dx = > %d\n", m.Bounds().Dx())
	tools.LOG_DEBUG.Printf("Dy = > %d\n", m.Bounds().Dy())

	return load_texture(glctx, m.Bounds().Dx(), m.Bounds().Dy(), gl.RGBA, rgbImage.Pix), nil

}

func build_program_from_assets(glctx gl.Context, vertex_shader_path string, fragment_shader_path string) (*gl.Program, error) {
	v, err := asset.Open(vertex_shader_path)
	if err != nil {
		return nil, tools.LogError(fmt.Sprintf("Failed to open vertex shader %s %s\n", vertex_shader_path, err))
	}
	defer v.Close()

	vertex_shader, err := ioutil.ReadAll(v)
	if err != nil {
		return nil, tools.LogError(fmt.Sprintf("Failed to read vertex shader %s %s\n", vertex_shader_path, err))
	}
	s, err := asset.Open(fragment_shader_path)
	if err != nil {
		return nil, tools.LogError(fmt.Sprintf("Failed to open fragment shader %s %s\n", fragment_shader_path, err))
	}
	defer s.Close()

	fragment_shader, err := ioutil.ReadAll(s)
	if err != nil {
		return nil, tools.LogError(fmt.Sprintf("Failed to read fragment shader %s %s\n", fragment_shader_path, err))
	}

	program, err := glutil.CreateProgram(glctx, string(vertex_shader), string(fragment_shader))
	if err != nil {
		return nil, tools.LogError(fmt.Sprintf("Failed to create program %s\n", err))
	}

	return &program, nil
}
