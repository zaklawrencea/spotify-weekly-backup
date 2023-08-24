package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/zaklawrencea/spotify-weekly-backup/utils"
)

const SPOTIFY_API = "https://api.spotify.com/v1"

func main() {
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
        "%s/playlists/%s/tracks?fields=items(track(name,uri))",
        SPOTIFY_API,
        utils.DiscoverWeeklyPlaylistID,
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
	if err != nil {
		return nil, err
	}
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


// Add the songs from the current Discover Weekly and add them to the
// backup playlist
func addToPlaylist(token string, songData *utils.SongData) error {
	var songURIs []string
	for _, item := range songData.Items {
		songURIs = append(songURIs, item.Track.URI)
	}

	// Format songURIs
	var discoveredSongs = ""
	for _, s := range songURIs {
		discoveredSongs += ("\"" + s + "\",")
	}

	// Remove trailing comma
	discoveredSongs = strings.TrimSuffix(discoveredSongs, ",")

	// Format JSON body
	var data = strings.NewReader( 
	`{"uris":[` + discoveredSongs + `]}`)

    url := fmt.Sprintf("%s/playlists/%s/tracks", SPOTIFY_API, utils.BackupDiscoverWeeklyPlaylistID)
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

    // Read response
    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return err
    }
	defer resp.Body.Close()

    // Spotify returns 201 with snapshot_id upon success
    if resp.StatusCode != http.StatusCreated {
        return errors.New(fmt.Sprintf("Non OK status code when adding to playlist: %v", string(respBody)))
    }

    log.Infof("Added discover weekly songs to archive playlist ID: %s\n\n", utils.BackupDiscoverWeeklyPlaylistID)

    return nil
}
