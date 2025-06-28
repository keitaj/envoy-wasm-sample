package main

import (
	"github.com/proxy-wasm/proxy-wasm-go-sdk/proxywasm"
	"github.com/proxy-wasm/proxy-wasm-go-sdk/proxywasm/types"
)

func main() {}

func init() {
	proxywasm.SetVMContext(&vmContext{})
}

type vmContext struct {
	types.DefaultVMContext
}

func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &pluginContext{}
}

type pluginContext struct {
	types.DefaultPluginContext
}

func (p *pluginContext) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	proxywasm.LogInfo("plugin started")
	return types.OnPluginStartStatusOK
}

func (*pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &httpAuthContext{contextID: contextID}
}

type httpAuthContext struct {
	types.DefaultHttpContext
	contextID uint32
}

func (ctx *httpAuthContext) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	// Get path
	path, err := proxywasm.GetHttpRequestHeader(":path")
	if err != nil {
		proxywasm.LogErrorf("failed to get path: %v", err)
		path = "/"
	}

	proxywasm.LogInfof("Processing request to path: %s", path)

	// Skip health check
	if path == "/health" {
		proxywasm.LogInfo("Health check endpoint, allowing request")
		return types.ActionContinue
	}

	// Get Authorization header
	authHeader, err := proxywasm.GetHttpRequestHeader("authorization")
	if err != nil || authHeader == "" {
		proxywasm.LogWarn("Missing authorization header")
		return ctx.denyRequest("Missing authorization header")
	}

	// Authentication check
	if authHeader == "Bearer secret-token-123" {
		proxywasm.LogInfo("Valid user token")
		proxywasm.AddHttpRequestHeader("x-auth-user", "user")
		return types.ActionContinue
	}
	if authHeader == "Bearer admin-token-456" {
		proxywasm.LogInfo("Valid admin token")
		proxywasm.AddHttpRequestHeader("x-auth-user", "admin")
		return types.ActionContinue
	}

	proxywasm.LogWarnf("Invalid token: %s", authHeader)
	return ctx.denyRequest("Invalid token")
}

func (ctx *httpAuthContext) OnHttpResponseHeaders(numHeaders int, endOfStream bool) types.Action {
	proxywasm.AddHttpResponseHeader("x-wasm-filter", "go-auth")
	return types.ActionContinue
}

func (ctx *httpAuthContext) denyRequest(reason string) types.Action {
	body := `{"error": "Unauthorized", "message": "` + reason + `"}`
	
	err := proxywasm.SendHttpResponse(401, [][2]string{
		{"content-type", "application/json"},
		{"x-wasm-filter", "go-auth"},
	}, []byte(body), -1)
	
	if err != nil {
		proxywasm.LogErrorf("failed to send response: %v", err)
	}
	
	return types.ActionPause
}