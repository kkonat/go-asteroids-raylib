package main

import (
	"embed"
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
)

//go:embed shaders
var shadersFS embed.FS

func load_shaders(v_shader_name string, f_shader_name string) rl.Shader {
	vsBytes, err := shadersFS.ReadFile(v_shader_name)
	if err != nil {
		log.Fatalf("failed to read vertex shader: %v", err)
	}
	fsBytes, err := shadersFS.ReadFile(f_shader_name)
	if err != nil {
		log.Fatalf("failed to read fragment shader: %v", err)
	}
	// Convert the byte slices to strings
	vertexShaderCode := string(vsBytes)
	fragmentShaderCode := string(fsBytes)
	shader := rl.LoadShaderFromMemory(vertexShaderCode, fragmentShaderCode)
	return shader
}
