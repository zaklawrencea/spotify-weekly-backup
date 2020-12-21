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
	var songIDs []string

	token := tokenRefresh()
	token = "Bearer " + token

	songIDs = getSongs(token)
	fmt.Println(songIDs)

	addToPlaylist (token, songIDs)
	//createPlaylist(token, discoverWeeklyBackup)
}

// Use the refresh token to generate a new bearer token
func tokenRefresh() (string) {

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

// Get songs from discover weekly playlist
func getSongs(token string) ([]string) {

	// Create GET request
	// Get song name + song ID for tracks in discovery weekly playlist
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/playlists/" 	+
								secrets.DiscoverWeekly 							+
								"/tracks?fields=items(track(name,id))", nil)
	

	//var bearer = "Bearer " + token
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	// Error handling
	if err != nil {
		panic(err)
	}

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

	// Add song ids to a slice
	var foundSongIDs []string
	for _,i := range val.Items {
		//fmt.Println(i.Track.ID)
		//fmt.Println(i.Track.Name)
		foundSongIDs = append(foundSongIDs, i.Track.ID)
		
	}

	return foundSongIDs
}

// TODO - Currently unused
// Create a playlist - return Spotify URI of playlist
func createPlaylist(token string) {

	// Create POST request
	// Setting the body and header for the POST request
	var discoverWeeklyBackup = "Discover Weekly Backup"
	var data = strings.NewReader( 
	`{"name":"` + discoverWeeklyBackup + `","public":false}`)

	req, err := http.NewRequest("POST", "https://api.spotify.com/v1/users/" + 
								secrets.UserID + "/playlists", data)	
	
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	// Error handling
	if err != nil {
		panic(err)
	}

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	defer req.Body.Close()

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	// Error handling
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))
}

// Add the songs from the current Discover Weekly and add them to the
// backup playlist
func addToPlaylist(token string, songIDs []string) {
	
}