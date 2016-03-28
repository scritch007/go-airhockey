package main

import (
	"math"

	"github.com/rkusa/gm/mat4"
	"github.com/rkusa/gm/vec3"
)

func mat4x4_perspective(y_fov_in_degrees, aspect, n, f float32) *mat4.Mat4 {
	var m4 mat4.Mat4

	return m4.Perspective(y_fov_in_degrees*math.Pi/180, aspect, n, f)
}

func mat4x4_look_at(eyeX, eyeY, eyeZ,
	centerX, centerY, centerZ,
	upX, upY, upZ float32) *mat4.Mat4 {

	var m4 mat4.Mat4

	eye := vec3.New(eyeX, eyeY, eyeZ)
	center := vec3.New(centerX, centerY, centerZ)
	up := vec3.New(upX, upY, upZ)

	return m4.LookAt(eye, *center, *up)

}

func mat4x4_mul(p, v *mat4.Mat4) *mat4.Mat4 {
	return p.Mul(v)
}
