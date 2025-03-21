package Models

import (
	"elucid503/SproutLoggingRelay/Functions"
	"errors"
)

type SproutURLSchema struct {

	// General

	BaseAPIURL string

	// Accounts

	GetMyAccount string

	// Other

}

var SproutAPIURLs = SproutURLSchema{

	BaseAPIURL: "https://sprout.software/API",

	GetMyAccount: "/Accounts/Me",
}

type SproutAccount struct {
	UID string `json:"UID"`

	Username string `json:"Username"`
	Email    string `json:"Email"`

	Avatar string `json:"Avatar"`

	Flags []any `json:"Flags"`
}

func (Account SproutAccount) FromJSON(JSONData *map[string]interface{}) SproutAccount {

	// JSON is already in the correct format, so we just need to cast it to the correct types

	return SproutAccount{

		UID:      (*JSONData)["UID"].(string),
		Username: (*JSONData)["Username"].(string),
		Email:    (*JSONData)["Email"].(string),
		Avatar:   (*JSONData)["Avatar"].(string),
		Flags:    (*JSONData)["Flags"].([]any),
	}

}

func GetSproutAccountByToken(Token string) (SproutAccount, error) {

	// Make a request to the Sprout API to get the account information

	var URL string = SproutAPIURLs.BaseAPIURL + SproutAPIURLs.GetMyAccount

	// Make the request

	RequestOptions := Functions.RequestOptions{

		Headers: map[string]string{

			"Authorization": "Bearer " + Token,
		},
	}

	RequestResponse, RequestError := Functions.MakeHTTPRequest("GET", URL, RequestOptions)

	if (RequestError != nil) || (RequestResponse.StatusCode != 200) {

		return SproutAccount{}, RequestError // Returning the error

	} else {

		JSONResp, Err := Functions.GetHTTPRequestJSONResponse(RequestResponse)

		if Err != nil {

			return SproutAccount{}, Err // Returning the error

		}

		var Success bool = JSONResp["Success"].(bool)

		if !Success {

			return SproutAccount{}, errors.New("Failed to get account") // Returning the error

		}

		return SproutAccount{}.FromJSON(&JSONResp), nil

	}

}
