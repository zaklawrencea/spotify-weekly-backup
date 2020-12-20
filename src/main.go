package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"strings"
	"bytes"

	secrets "./mylib"
)

type SongData struct {
	Items []Items
}

type Items struct {
	Track Tracks
}

type Tracks struct {
	ID string
	Name string
}

type AuthResponse struct {

	// Maps to JSON properties in spotify response
	Access_Token string
	Token_Type string
	Expires_In string
	Scope string
}

func main() {
	var token = tokenRefresh()
	getSongs(token)
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

func getSongs(token string) {

	// Get songs from discover weekly playlist

	// Format bearer token for headers
	var bearer = "Bearer " + token

	// Create GET request
	// Get song name + song ID for tracks in discovery weekly playlist
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/playlists/"+secrets.DiscoverWeekly+"/tracks?fields=items(track(name,id))", nil)
	req.Header.Set("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)

	// Error handling
	if err != nil {
		panic(err)
	}

	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	// Parse response as JSON
	r := bytes.NewReader(body)
	decoder := json.NewDecoder(r)
	
	val := &SongData{}

	decodeErr := decoder.Decode(val)

	// Error handling
	if decodeErr != nil {
		panic(err)
	}

	// Print song name + id
	for _,i := range val.Items {
		fmt.Println(i.Track.ID)
		fmt.Println(i.Track.Name)
	}
}
