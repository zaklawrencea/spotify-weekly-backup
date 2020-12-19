package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"strings"

	secrets "./mylib"
)

// type Song struct {
// 	Href string
// }

type AuthResponse struct {

	// Maps to JSON properties in spotify response
	Access_Token string
	Token_Type string
	Expires_In string
	Scope string
}

func main() {

	var token = tokenRefresh()
	fmt.Println(token)
	//getPlaylist()
}

func tokenRefresh() string {

	// Setting the body and header for the POST request
	var data = strings.NewReader( 
		"client_id=" + secrets.ClientID + 
		"&client_secret="	+ secrets.ClientSecret +
		"&grant_type="	+ secrets.GrantType +
		"&refresh_token=" + secrets.RefreshToken)

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", data)	
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	
	// Error handling
	if err != nil {
		panic(err)
	}

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	defer req.Body.Close()

	// Error handling
	if err != nil {
		panic(err)
	}

	// Read and parse the response body
	body, err := ioutil.ReadAll(resp.Body)

	// Error handling
	if err != nil {
		panic(err)
	}

	// Get access token from spotify response
	var authResponse AuthResponse
	json.Unmarshal([]byte(body), &authResponse)

	
	return string(authResponse.Access_Token)
}

func getPlaylist() {

}
