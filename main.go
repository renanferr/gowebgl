package gowebgl

import (
	"fmt"
	"os"
	"syscall/js"
	"math"

	"gonum.org/v1/gonum/mat"
	"github.com/renanferr/gowebgl/mat4"
)

// Shaders sources
type Shaders struct {
	vertex   string
	fragment string
}

type AttribLocationsInfo struct {
	vertexPosition js.Value
}

type UniformsLocationsInfo struct {
	projectionMatrix 	js.Value
	modelViewMatrix 	js.Value
}

type WebGLContext struct {
	shaderProgram 		js.Value
	canvasCtx 			js.Value
	attribLocations 	AttribLocationsInfo
	uniformLocations 	UniformLocationsInfo
}

func (gl *WebGLContext) GetContext() js.Value {
	return gl.canvasCtx
}

// InitShaders initializes Shaders interface
func InitShaders(v string, f string) Shaders {
	s := Shaders{
		vertex:   v,
		fragment: f,
	}
	return s
}


// Init initializes WebGL
func Init(canvas js.Value, shaders Shaders) WebGLContext {
	canvasCtx := canvas.Call("getContext", "webgl")
	// if !canvasCtx {
	// 	panic("Unable to initialize WebGL. Your browser or machine may not support it.")
	// }

	fmt.Println("canvasCtx", canvasCtx)
	// Set clear color to black, fully opaque
	canvasCtx.Call("clearColor", 0.0, 0.0, 0.0, 1.0)

	// Clear the color buffer with specified clear color
	canvasCtx.Call("clear", canvasCtx.Get("COLOR_BUFFER_BIT"))

	shaderProgram := initShaderProgram(canvasCtx, shaders)

	attribLocations := AttribLocationsInfo{
		vertexPosition: canvasCtx.Call("getAttribLocation", shaderProgram, "aVertexPosition"),
	}
	uniformLocations := UniformsLocationsInfo{
		projectionMatrix: canvasCtx.Call("getUniformLocation", shaderProgram, "uProjectionMatrix"),
		modelViewMatrix: canvasCtx.Call("getUniformLocation", shaderProgram, "uModelViewMatrix"),
	}

	gl := WebGLContext{
        shaderProgram,
        attribLocations,
		uniformLocations,
		canvasCtx,
	}
	
	return gl
}

func initShaderProgram(ctx js.Value, shaders Shaders) js.Value {
	vertexShader := loadShader(ctx, ctx.Get("VERTEX_SHADER"), shaders.vertex)
	fragmentShader := loadShader(ctx, ctx.Get("FRAGMENT_SHADER"), shaders.fragment)

	// Create the shader program
	shaderProgram := ctx.Call("createProgram")
	ctx.Call("attachShader", shaderProgram, vertexShader)
	ctx.Call("attachShader", shaderProgram, fragmentShader)
	ctx.Call("linkProgram", shaderProgram)

	// If creating the shader program failed, alert
	LINK_STATUS := ctx.Get("LINK_STATUS")
	fmt.Println("LINK_STATUS", LINK_STATUS)
	success := ctx.Call("getProgramParameter", shaderProgram, LINK_STATUS).Bool()
	fmt.Println("shaderProgram", success)
	if !success {
		programInfoLog := ctx.Call("getProgramInfoLog", shaderProgram)
		panic(fmt.Sprintf("Unable to initialize the shader program: %v", programInfoLog))
		os.Exit(0)
	}

	return shaderProgram
}

func loadShader(ctx js.Value, sType js.Value, source string) js.Value {
	shader := ctx.Call("createShader", sType)

	// Send the source to the shader object
	ctx.Call("shaderSource", shader, source)

	// Compile the shader program
	ctx.Call("compileShader", shader)

	// See if it compiled successfully
	COMPILE_STATUS := ctx.Get("COMPILE_STATUS")
	fmt.Println("COMPILE_STATUS", COMPILE_STATUS)
	success := ctx.Call("getShaderParameter", shader, COMPILE_STATUS).Bool()
	fmt.Println("shader", sType, success)

	if !success {
		infoLog := ctx.Call("getShaderInfoLog", shader).String()
		err := fmt.Sprintf("An error occurred compiling the shaders: %v", infoLog)
		panic(err)
		ctx.Call("deleteShader", shader)
		os.Exit(0)
	}

	return shader
}


func DrawScene(gl WebGLContext, buffers Buffers) {
	ctx := gl.GetContext()
	ctx.Call("clearColor", 0.0, 0.0, 0.0, 1.0) 	// Clear to black, fully opaque
    ctx.Call("clearDepth", 1.0)                	// Clear everything
    ctx.Call("enable", ctx.Get("DEPTH_TEST"))	// Enable depth testing
    ctx.Call("depthFunc", ctx.Get("LEQUAL"))	// Near things obscure far things
  
    // Clear the canvas before start drawing on it.
    ctx.Call("clear", (ctx.Get("COLOR_BUFFER_BIT") | ctx.Get("DEPTH_BUFFER_BIT")));
  
    // Create a perspective matrix, a special matrix that is
    // used to simulate the distortion of perspective in a camera.
    // Our field of view is 45 degrees, with a width/height
    // ratio that matches the display size of the canvas
    // and we only want to see objects between 0.1 units
    // and 100 units away from the camera.
  
	fieldOfView := 45 * math.Pi / 180

	aspect := ctx.Get("canvas").Get("clientWidth").Float() / ctx.Get("canvas").Get("clientHeight").Float()

	zNear := 0.1
	zFar := 100.0
	projectionMatrix := mat4.Create()
  
    // note: glmatrix.js always has the first argument
    // as the destination to receive the result.
    mat4.Perspective(projectionMatrix,
                     fieldOfView,
                     aspect,
                     zNear,
                     zFar)
  
    // Set the drawing position to the "identity" point, which is
    // the center of the scene.
	modelViewMatrix := mat4.Create()
    // Now move the drawing position a bit to where we want to
    // start drawing the square.
  
    mat4.translate(modelViewMatrix,     // destination matrix
                   modelViewMatrix,     // matrix to translate
                   [-0.0, 0.0, -6.0]);  // amount to translate
  
    // Tell WebGL how to pull out the positions from the position
    // buffer into the vertexPosition attribute.
    {
      const numComponents = 2;  // pull out 2 values per iteration
      const type = ctx.FLOAT;    // the data in the buffer is 32bit floats
      const normalize = false;  // don't normalize
      const stride = 0;         // how many bytes to get from one set of values to the next
                                // 0 = use type and numComponents above
      const offset = 0;         // how many bytes inside the buffer to start from
      ctx.bindBuffer(ctx.ARRAY_BUFFER, buffers.position);
      ctx.vertexAttribPointer(
          programInfo.attribLocations.vertexPosition,
          numComponents,
          type,
          normalize,
          stride,
          offset);
      ctx.enableVertexAttribArray(
          programInfo.attribLocations.vertexPosition);
    }
  
    // Tell WebGL to use our program when drawing
  
    ctx.useProgram(programInfo.program);
  
    // Set the shader uniforms
  
    ctx.uniformMatrix4fv(
        programInfo.uniformLocations.projectionMatrix,
        false,
        projectionMatrix);
    ctx.uniformMatrix4fv(
        programInfo.uniformLocations.modelViewMatrix,
        false,
        modelViewMatrix);
  
    {
      const offset = 0;
      const vertexCount = 4;
      ctx.drawArrays(ctx.TRIANGLE_STRIP, offset, vertexCount);
    }
  }