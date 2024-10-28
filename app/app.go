package main

import (
	"context"
	"encoding/base64"
	"fmt"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// HandleGreet handles the HTMX request
func (a *App) HandleGreet() string {
	return "<p>App version: 0.1-beta</p>"
}

// HandleHome returns the home content
func (a *App) HandleHome() string {
	return `
        <h1>Hi, thanks for using this app.</h1>
        <button onclick="window.go.main.App.HandleGreet().then(result => document.getElementById('greeting').innerHTML = result)">
            Show Details
        </button>
        <div id="greeting"></div>
    `
}

// HandleBase64 returns the Base64 encoding/decoding UI
func (a *App) HandleBase64() string {
	return `
        <h2>Base64 Encode/Decode</h2>
        <textarea id="input" placeholder="Enter text here"></textarea>
        <button onclick="encodeBase64()">Encode</button>
        <button onclick="decodeBase64()">Decode</button>
        <div id="output">
            <textarea id="result" readonly rows="5" style="width: 100%; background-color: #444; color: #fff; border: 1px solid #666; padding: 5px;"></textarea>
        </div>
        <button onclick="copyToClipboard()">Copy to Clipboard</button>
    `
}

// HandleBase64Encode encodes the input to Base64
func (a *App) HandleBase64Encode(input string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(input))
	return encoded
}

// HandleBase64Decode decodes the input from Base64
func (a *App) HandleBase64Decode(input string) string {
	decoded, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return fmt.Sprintf("Error: %s", err.Error())
	}
	return string(decoded)
}
