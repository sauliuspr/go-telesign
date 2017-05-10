package auth

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Compute HMAC 256 of Cannonical String
//
func Compute(canonicalString, authMethod, secret string) string {
	data, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		fmt.Println("ERROR: in Base64 ", err)
	}
	mac := hmac.New(sha256.New, []byte(data))
	if authMethod == "sha1" {
		mac = hmac.New(sha1.New, []byte(data))
		fmt.Println("INFO: Faling back to sha1")
	}
	mac.Write([]byte(canonicalString))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// CanonicalString retrieves string from http Request
//
func CanonicalString(r *http.Request) string {
	contentType := ""
	uri := r.URL.EscapedPath()
	if uri == "" {
		uri = "/"
	}

	if r.URL.RawQuery != "" {
		uri = uri + "?" + r.URL.RawQuery
	}
	header := r.Header

	if r.Method == "POST" || r.Method == "PUT" {
		contentType = "application/x-www-form-urlencoded"
		body := new(bytes.Buffer)
		body.ReadFrom(r.Body)
		ret := strings.Join([]string{
			strings.ToUpper(r.Method),
			contentType,
			"",
			"x-ts-auth-method:" + header.Get("x-ts-auth-method"),
			"x-ts-date:" + header.Get("x-ts-date"),
			fmt.Sprintf("%s", body),
			r.URL.Path,
		}, "\n")
		return ret
	}
	return strings.Join([]string{
		strings.ToUpper(r.Method),
		contentType,
		"",
		"x-ts-auth-method:" + header.Get("x-ts-auth-method"),
		"x-ts-date:" + header.Get("x-ts-date"),
		r.URL.Path,
	}, "\n")
}

// GenerateTelesignHeaders with CusteomerID and apiKey
//
func GenerateTelesignHeaders(req *http.Request, customerID, apiKey string) *http.Request {

	// If Request Body exist preserve it
	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(req.Body)
		req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	// Clone the request
	rClone := new(http.Request)
	*rClone = *req
	rClone.Header = make(http.Header, len(req.Header))
	for idx, header := range req.Header {
		rClone.Header[idx] = append([]string(nil), header...)
	}

	// Adding Telesign headers
	if rClone.Header["X-Ts-Date"] == nil {
		rClone.Header.Add("X-Ts-Date", time.Now().Format(time.RFC1123Z))
	}
	if rClone.Header["X-Ts-Auth-Method"] == nil {
		rClone.Header.Add("X-Ts-Auth-Method", "HMAC-SHA256")
	}
	if rClone.Header["X-Ts-Nonce"] == nil {

		// Generate Random String
		b := make([]byte, 16)
		if _, err := rand.Read(b); err != nil {
			fmt.Println("Could not generate a random string", err)
			return nil
		}

	}
	sig := Compute(CanonicalString(rClone), "sha256", apiKey)
	rClone.Header.Add("Authorization", fmt.Sprintf("TSA %s:%s", customerID, sig))

	if bodyBytes != nil {
		rClone.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	return rClone
}
