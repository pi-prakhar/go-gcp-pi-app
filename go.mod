module github.com/pi-prakhar/go-gcp-pi-app

go 1.23.1

require (
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/pi-prakhar/go-gcp-auth v0.0.0-20241005170138-90cb0ecef115
	golang.org/x/oauth2 v0.23.0
)

require cloud.google.com/go/compute/metadata v0.3.0 // indirect
