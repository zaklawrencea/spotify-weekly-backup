package utils

type SongData struct {
	Items []Items
}

type Items struct {
	Track Track
}

type Track struct {
	URI string
	Name string
}

type AuthResponse struct {
	// Maps to JSON properties in spotify response
	Access_Token string
	Token_Type string
	Expires_In string
	Scope string
}
