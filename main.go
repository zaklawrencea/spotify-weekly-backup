package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/zaklawrencea/spotify-weekly-backup/utils"
)

const SPOTIFY_API = "https://api.spotify.com/v1"

func main() {
	// var songURIs []string

	token, err := tokenRefresh()
    if err != nil {
        log.Errorf("Unable to refresh token: %v", err)
    }

    songData, err := getSongs(token)
    if err != nil {
        log.Errorf("Unable to fetch discover weekly songs: %v", err)
    }

	err = addToPlaylist(token, songData)
    if err != nil {
        log.Errorf("Unable to add songs to backup playlist: %v", err)
    }

	//createPlaylist(token, discoverWeeklyBackup)
}

// Use the refresh token to generate a new bearer token
func tokenRefresh() (string, error) {
	var data = strings.NewReader( 
		"client_id=" + utils.ClientID + 
		"&client_secret="	+ utils.ClientSecret +
		"&grant_type="	+ utils.GrantType +
		"&refresh_token=" + utils.RefreshToken)

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", data)	
	if err != nil {
		return "", err
	}
    
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
	    return "", err
	}
	defer req.Body.Close()

	// Read and parse the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	    return "", err
	}

	// Get access token from spotify response
	var authResponse utils.AuthResponse
	json.Unmarshal([]byte(body), &authResponse)

    return fmt.Sprintf("Bearer %s",authResponse.Access_Token), nil
}

// Get songs from discover weekly playlist
// func getSongs(token string) ([]string, error) {
func getSongs(token string) (*utils.SongData, error) {
	// Get song name + song URI for tracks in discovery weekly playlist
    url := fmt.Sprintf(
        "%s/%s/tracks?fields=items(track(name,uri))",
        SPOTIFY_API,
        utils.DiscoverWeekly,
    )
    req, err := http.NewRequest(
        http.MethodGet,
        url,
        nil,
    )
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	// Parse response as JSON
	r := bytes.NewReader(body)
	decoder := json.NewDecoder(r)
	
	discoverWeeklySongData := &utils.SongData{}
	err = decoder.Decode(discoverWeeklySongData)
	if err != nil {
		return nil, err
	}

	return discoverWeeklySongData, nil
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
								utils.UserID + "/playlists", data)	
	
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
func addToPlaylist(token string, songData *utils.SongData) error {
	
	var songURIs []string
	for _, item := range songData.Items {
		songURIs = append(songURIs, item.Track.URI)
	}

	// Format songURI 
	var discoveredSongs = ""
	for _, s := range songURIs {
		discoveredSongs += ("\"" + s + "\",")
	}

	// Remove trailing comma
	discoveredSongs = strings.TrimSuffix(discoveredSongs, ",")

	// Format JSON body
	var data = strings.NewReader( 
	`{"uris":[` + discoveredSongs + `]}`)

    url := fmt.Sprintf("%s/playlists/%s/tracks", SPOTIFY_API, utils.DiscoverWeeklyPlaylist)
	req, err := http.NewRequest(
        http.MethodPost,
        url,
        data,
    )	
	if err != nil {
		return err
	}
	
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
    resp, err := client.Do(req)
	if err != nil {
        return err
	}
	defer req.Body.Close()

	// Read response
    if resp.StatusCode != http.StatusOK {
        return errors.New("Non OK status code when adding to playlist")
    }

    log.Infof("Added discover weekly songs to archive playlist ID: %s\n\n", utils.DiscoverWeeklyPlaylist)

    return nil
}
