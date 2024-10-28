package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
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
        <textarea id="input" placeholder="Enter text here"></textarea>
        <button onclick="encodeBase64()">Encode</button>
        <button onclick="decodeBase64()">Decode</button>
        <div id="output">
            <textarea id="result" readonly rows="5" style="width: 100%; background-color: #444; color: #fff; border: 1px solid #666; padding: 5px;"></textarea>
        </div>
        <button onclick="copyToClipboard('result')">Copy to Clipboard</button>
        <div id="notification" style="display: none;"></div>
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

func (a *App) HandleCryptRSA() string {
	return `
        <div id="rsa" class="crypto-tab-content" style="display:block;">
            <button onclick="generateRSAKeyPair()" style="margin-bottom: 10px;">Generate Key Pair</button>
            <div style="display: flex; flex-direction: column; gap: 10px;">
                <div>
                    <label for="rsa-private-key">Private Key:</label>
                    <div style="display: flex; align-items: center;">
                        <textarea id="rsa-private-key" rows="2" style="flex-grow: 1; resize: none; overflow-y: auto; height: 3em;" readonly></textarea>
                        <button onclick="copyToClipboard('rsa-private-key')" style="margin-left: 5px; height: 3em;">Copy</button>
                    </div>
                </div>
                <div>
					<label for="rsa-public-key">Public Key:</label>
                    <div style="display: flex; align-items: center;">
                        <textarea id="rsa-public-key" rows="2" style="flex-grow: 1; resize: none; overflow-y: auto; height: 3em;" readonly></textarea>
                        <button onclick="copyToClipboard('rsa-public-key')" style="margin-left: 5px; height: 3em;">Copy</button>
                    </div>
                </div>
                <div>
                    <label for="rsa-secret-message">Secret Message:</label>
                    <div style="display: flex; align-items: center;">
                        <textarea id="rsa-secret-message" rows="2" style="flex-grow: 1; resize: none; overflow-y: auto; height: 3em;"></textarea>
                    </div>
                </div>
            </div>
			<button onclick="rsaEncrypt()" style="margin-bottom: 10px;">Encrypt Message</button>
            <div style="display: flex; align-items: center;">
                <textarea id="rsa-signed-message" rows="2" style="flex-grow: 1; resize: none; overflow-y: auto; height: 5em;" readonly></textarea>
                <button onclick="copyToClipboard('rsa-signed-message')" style="margin-left: 5px; height: 3em;">Copy</button>
			</div>
			<div id="notification" style="display: none;"></div>
        </div>
    `
}

func (a *App) HandleRSAGenerateKeyPair() string {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err.Error()
	}

	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	privateKeyStr := string(pem.EncodeToMemory(privateKeyPEM))

	publicKeyPEM := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
	}
	publicKeyStr := string(pem.EncodeToMemory(publicKeyPEM))

	result := map[string]string{
		"privateKey": privateKeyStr,
		"publicKey":  publicKeyStr,
	}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		return err.Error()
	}

	return string(jsonResult)
}

func (a *App) HandleRSAEncrypt(message, publicKeyStr string) string {
	block, _ := pem.Decode([]byte(publicKeyStr))
	if block == nil {
		return "Failed to parse public key"
	}

	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return err.Error()
	}

	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(message))
	if err != nil {
		return err.Error()
	}

	return base64.StdEncoding.EncodeToString(ciphertext)
}

func (a *App) HandleRSADecrypt(ciphertext, privateKeyStr string) string {
	block, _ := pem.Decode([]byte(privateKeyStr))
	if block == nil {
		return "Failed to parse private key"
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return err.Error()
	}

	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return err.Error()
	}

	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertextBytes)
	if err != nil {
		return err.Error()
	}

	return string(plaintext)
}
