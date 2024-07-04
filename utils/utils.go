package utils

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

type ServiceAccount struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
}

func GetPrivateKey(key string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(key))
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey.(*rsa.PrivateKey), nil
}

func GenerateJWT(sa ServiceAccount, scope string, privateKey *rsa.PrivateKey) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"iss":   sa.ClientEmail,
		"scope": scope,
		"aud":   sa.TokenURI,
		"exp":   now.Add(time.Hour).Unix(),
		"iat":   now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}

func GetAccessToken(sa ServiceAccount, jwt string) (string, error) {

	data := map[string]string{
		"grant_type": "urn:ietf:params:oauth:grant-type:jwt-bearer",
		"assertion":  jwt,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", sa.TokenURI, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tokenResponse map[string]interface{}

	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return "", err
	}

	accessToken, ok := tokenResponse["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("no access token found in response")
	}

	return accessToken, nil
}
