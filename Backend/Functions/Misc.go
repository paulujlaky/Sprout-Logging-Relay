package Functions

import (
	"encoding/json"
	"math/rand"
	"net/http"
)

// HTTP Requests

type RequestOptions struct {

	Headers map[string]string

}

var HTTPClient = &http.Client{}

func MakeHTTPRequest(Method string, URL string, Options RequestOptions) (*http.Response, error) {

	// Make the request

	RequestInstance, Err := http.NewRequest(Method, URL, nil)

	if Err != nil {

		return nil, Err

	}

	for HeaderKey, HeaderValue := range Options.Headers {

		RequestInstance.Header.Set(HeaderKey, HeaderValue)

	}

	RequestResponse, RequestErr := HTTPClient.Do(RequestInstance)

	if RequestErr != nil {

		return nil, RequestErr // Return the error

	}

	return RequestResponse, nil

}

func GetHTTPRequestJSONResponse(Response *http.Response) (map[string]interface{}, error) {

	var ResponseData map[string]interface{}

	Decoder := json.NewDecoder(Response.Body)

	if Err := Decoder.Decode(&ResponseData); Err != nil {

		return nil, Err

	}

	return ResponseData, nil

}

// Random

func RandomString(Length int) string {

	var Charset string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	var RandomStr string

	for i := 0; i < Length; i++ {

		RandomStr += string(Charset[rand.Intn(len(Charset))])

	}

	return RandomStr

}

// Strings

func PluralizeString(Count int, Str string) string {

	if Count == 1 {

		return Str

	} else {

		return Str + "s"

	}

}