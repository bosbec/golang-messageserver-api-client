package client

import (
	"bytes"
	"code.google.com/p/go-uuid/uuid"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	url      string
	accessId string
	key      string
}

func New(url, accessId, key string) *Client {
	client := &Client{url, accessId, key}

	return client
}

func (this *Client) SendSms(request *SendSmsRequest) error {
	requestUri := strings.ToLower(url.QueryEscape(this.url))
	requestHttpMethod := "POST"
	unixTimestamp := strconv.FormatInt(time.Now().Unix(), 10)
	nonce := strings.Replace(uuid.New(), "-", "", -1)

	content := toJson(request)

	contentHash := computeMd5Hash(content)

	signature := computeSignature(
		this.accessId,
		requestHttpMethod,
		requestUri,
		unixTimestamp,
		nonce,
		contentHash)

	key, _ := base64.StdEncoding.DecodeString(this.key)

	signatureHash := computeHmacHash(signature, key)

	a := []string{this.accessId, signatureHash, nonce, unixTimestamp}

	auth := strings.Join(a, ":")

	return performHttpRequest(requestHttpMethod, this.url, content, auth)
}

func computeSignature(
	accessId string,
	httpMethod string,
	url string,
	timestamp string,
	nonce string,
	content string) []byte {
	var data bytes.Buffer

	data.WriteString(accessId)
	data.WriteString(httpMethod)
	data.WriteString(url)
	data.WriteString(timestamp)
	data.WriteString(nonce)
	data.WriteString(content)

	signature := data.Bytes()

	return signature
}

func performHttpRequest(
	httpMethod string,
	url string,
	content []byte,
	authorization string) error {
	request, err := http.NewRequest(httpMethod, url, bytes.NewBuffer(content))

	if err != nil {
		return err
	}

	request.Header.Set("Authorization", "mr "+authorization)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	return nil
}

func computeMd5Hash(content []byte) string {
	hasher := md5.New()

	hasher.Write(content)

	hash := base64.StdEncoding.EncodeToString(hasher.Sum(nil))

	return hash
}

func computeHmacHash(content []byte, key []byte) string {
	hasher := hmac.New(sha256.New, key)

	hasher.Write(content)

	hash := base64.StdEncoding.EncodeToString(hasher.Sum(nil))

	return hash
}

func toJson(request interface{}) []byte {
	json, _ := json.MarshalIndent(request, "", "  ")

	content := []byte(json)

	return content
}
