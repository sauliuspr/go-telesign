package tsauth

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func compute(canonicalString, authMethod, secret string) string {
	defer timeTrack(time.Now(), "compute")
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

func canonicalString(r *http.Request) string {
	defer timeTrack(time.Now(), "cannonicalString")
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
		fmt.Printf("String: \n%s\n", ret)
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

// Sign Telesign Request with CusteomerID and apiKey and return the request for futher processing
//
func Sign(req *http.Request, customerID, apiKey string) *http.Request {
	defer timeTrack(time.Now(), "tsauth.Sign")
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
			panic(err)
		}
		//rClone.Header.Add("X-Ts-Nonce", fmt.Sprintf("%X", b))
	}
	sig := compute(canonicalString(rClone), "sha256", apiKey)
	rClone.Header.Add("Authorization", fmt.Sprintf("TSA %s:%s", customerID, sig))

	if bodyBytes != nil {
		rClone.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	return rClone
}
