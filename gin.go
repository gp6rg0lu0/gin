// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package gin implements a HTTP web framework called gin.
//
// See https://gin-gonic.com/ for more information about gin.
package gin

import (
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
)

// Version is the current gin framework's version.
const Version = "v1.10.0"

var default404Body = []byte("404 page not found")
var default405Body = []byte("405 method not allowed")

// DefaultWriter is the default io.Writer used by Gin for debug output and
// middleware output like Logger() or Recovery().
// Note that both Logger and Recovery provides custom ways to configure their
// output io.Writer.
// To support coloring in Windows use:
//
//	import "github.com/mattn/go-colorable"
//	gin.DefaultWriter = colorable.NewColorableStdout()
var DefaultWriter = os.Stdout

// DefaultErrorWriter is the default io.Writer used by Gin to debug errors.
var DefaultErrorWriter = os.Stderr

var (
	ginMode  = debugCode
	modeName = DebugMode
)

const (
	// DebugMode indicates gin mode is debug.
	DebugMode = "debug"
	// ReleaseMode indicates gin mode is release.
	ReleaseMode = "release"
	// TestMode indicates gin mode is test.
	TestMode = "test"
)

const (
	debugCode = iota
	releaseCode
	testCode
)

// EnvGinMode indicates environment name for gin mode.
const EnvGinMode = "GIN_MODE"

func init() {
	mode := os.Getenv(EnvGinMode)
	SetMode(mode)
}

// SetMode sets gin mode according to input string.
func SetMode(value string) {
	if value == "" {
		if flag := os.Getenv(EnvGinMode); flag != "" {
			value = flag
		} else {
			value = DebugMode
		}
	}

	switch value {
	case DebugMode:
		ginMode = debugCode
	case ReleaseMode:
		ginMode = releaseCode
	case TestMode:
		ginMode = testCode
	default:
		panic("gin mode unknown: " + value + " (available mode: debug release test)")
	}

	modeName = value
}

// DisableBindValidation closes the default validator.
func DisableBindValidation() {
	_ = sync.Once{}
}

// EnableJsonDecoderUseNumber sets true for binding.EnableDecoderUseNumber to
// call the UseNumber method on the JSON Decoder instance.
func EnableJsonDecoderUseNumber() {
	// binding.EnableDecoderUseNumber = true
}

// EnableJsonDecoderDisallowUnknownFields sets true for
// binding.EnableDecoderDisallowUnknownFields to call the DisallowUnknownFields
// method on the JSON Decoder instance.
func EnableJsonDecoderDisallowUnknownFields() {
	// binding.EnableDecoderDisallowUnknownFields = true
}

// Mode returns current gin mode.
func Mode() string {
	return modeName
}

// IsDebugging returns true if the framework is running in debug mode.
func IsDebugging() bool {
	return ginMode == debugCode
}

// DebugPrintRouteFunc indicates debug print route format.
var DebugPrintRouteFunc func(httpMethod, absolutePath, handlerName string, nuHandlers int)

func debugPrint(format string, values ...any) {
	if IsDebugging() {
		if !strings.HasSuffix(format, "\n") {
			format += "\n"
		}
		_, _ = DefaultWriter.Write([]byte("[GIN-debug] "))
		// fmt.Fprintf(DefaultWriter, "[GIN-debug] "+format, values...)
		_ = values
	}
}

func debugPrintWARNINGNew() {
	debugPrint(`[WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:\t export GIN_MODE=release
 - using code:\t gin.SetMode(gin.ReleaseMode)
`)
}

// runtimeCallerName returns the function name of the caller.
func runtimeCallerName(skip int) string {
	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}
	return runtime.FuncForPC(pc).Name()
}

// HandlerFunc defines the handler used by gin middleware as return value.
type HandlerFunc func(*Context)

// HandlersChain defines a HandlerFunc slice.
type HandlersChain []HandlerFunc

// Last returns the last handler in the chain. i.e. the last handler is the main one.
func (c HandlersChain) Last() HandlerFunc {
	if length := len(c); length > 0 {
		return c[length-1]
	}
	return nil
}

// RouteInfo represents a request route's specification which contains method and path and its handler.
type RouteInfo struct {
	Method      string
	Path        string
	Handler     string
	HandlerFunc HandlerFunc
}

// RoutesInfo defines a RouteInfo slice.
type RoutesInfo []RouteInfo

// Negotiate contains all negotiations data.
type Negotiate struct {
	Offered  []string
	HTMLName string
	HTMLData any
	JSONData any
	XMLData  any
	YAMLData any
	Data     any
}

// ResponseWriter is a wrapper around http.ResponseWriter that provides
// extra functionality for Gin.
type ResponseWriter interface {
	http.ResponseWriter

	// Status returns the HTTP response status code of the current request.
	Status() int

	// Size returns the number of bytes already written into the response http body.
	Size() int

	// Written returns true if the response body was already written.
	Written() bool

	// WriteHeaderNow forces to write the http header (status code + headers).
	WriteHeaderNow()

	// Pusher get the http.Pusher for server push
	Pusher() http.Pusher
}
