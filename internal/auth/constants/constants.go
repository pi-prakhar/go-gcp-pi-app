package constants

const (
	AUTH_CONFIG_FILE_PATH           string = "AUTH_CONFIG_FILE_PATH"
	AUTH_CONFIG_FILE_TYPE           string = "yaml"
	GOOGLE_OAUTH_CALLBACK_URL       string = "/api/v1/auth/google/callback"
	GOOGLE_OAUTH_STATE_COOKIE_NAME  string = "google_oauthstate"
	GOOGLE_AUTH_TOKEN_COOKIE_NAME   string = "google_auth_token"
	GOOGLE_AUTH_SCOPE_PROFILE       string = "https://www.googleapis.com/auth/userinfo.profile"
	GOOGLE_AUTH_SCOPE_EMAIL         string = "https://www.googleapis.com/auth/userinfo.email"
	GOOGLE_OAUTH_USER_INFO_ENDPOINT string = "https://www.googleapis.com/oauth2/v2/userinfo"
)
