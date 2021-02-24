# spotify-weekly-backup
Never lose your discover weekly playlists again.

# Setup
The application uses a refresh token to generate an access token in order to perform actions on behalf of the user. This token must be generated yourself and the relevant parameters added to variables inside the src/mylib/secrets.go file. Steps for generating this token can be found in the Spotify Authorization guide: https://developer.spotify.com/documentation/general/guides/authorization-guide/#authorization-code-flow

