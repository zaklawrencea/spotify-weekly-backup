# Spotify Weekly Backup
Ever listen through your Spotify Discover Weekly playlist, find one or two new bangers, but forget to save them before the playlist is refreshed? Me too...

This app is designed to run as a weekly cronjob that will save (append) the contents of a provided playlist to that of another given playlist. In my case: from the Discover Weekly playlist to a Backup playlist.

# Setup & Local development
Please note this is a developmental app and isn't particularly designed to be used by others. However, if you want to use the code, consider setting up a Spotify App and following the [Spotify Authorization guide](https://developer.spotify.com/documentation/general/guides/authorization-guide/#authorization-code-flow), specifically the [Authorization Code Flow](https://developer.spotify.com/documentation/web-api/tutorials/code-flow), to generate a refresh token.

Afterward, fill in the various variables in `utils/secrets.go`.
