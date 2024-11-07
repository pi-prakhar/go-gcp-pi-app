package utils

import (
	"io/ioutil"
	"os"
	"testing"
)

// Test function for GetClientId
func TestGetClientId(t *testing.T) {
	// Create a temp file with client ID content
	tmpFile, err := ioutil.TempFile("", "client_id")
	if err != nil {
		t.Fatalf("Unable to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Cleanup

	clientID := "test-client-id"
	if _, err = tmpFile.Write([]byte(clientID)); err != nil {
		t.Fatalf("Unable to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Set environment variable
	os.Setenv("GOOGLE_OAUTH_CLIENT_ID", tmpFile.Name())
	defer os.Unsetenv("GOOGLE_OAUTH_CLIENT_ID")

	// Call function and verify output
	result := GetClientId()
	if result != clientID {
		t.Errorf("Expected client ID %v, got %v", clientID, result)
	}
}

// Test function for GetClientSecret
func TestGetClientSecret(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "client_secret")
	if err != nil {
		t.Fatalf("Unable to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	clientSecret := "test-client-secret"
	if _, err = tmpFile.Write([]byte(clientSecret)); err != nil {
		t.Fatalf("Unable to write to temp file: %v", err)
	}
	tmpFile.Close()

	os.Setenv("GOOGLE_OAUTH_CLIENT_SECRET", tmpFile.Name())
	defer os.Unsetenv("GOOGLE_OAUTH_CLIENT_SECRET")

	result := GetClientSecret()
	if result != clientSecret {
		t.Errorf("Expected client secret %v, got %v", clientSecret, result)
	}
}

// Test function for GetCallbackURL
func TestGetCallbackURL(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "auth_service_host")
	if err != nil {
		t.Fatalf("Unable to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	authServiceHost := "http://localhost:8080"
	if _, err = tmpFile.Write([]byte(authServiceHost)); err != nil {
		t.Fatalf("Unable to write to temp file: %v", err)
	}
	tmpFile.Close()

	os.Setenv("AUTH_SERVICE_HOST", tmpFile.Name())
	defer os.Unsetenv("AUTH_SERVICE_HOST")

	expectedURL := authServiceHost + "/api/v1/auth/google/callback"
	result := GetCallbackURL()
	if result != expectedURL {
		t.Errorf("Expected callback URL %v, got %v", expectedURL, result)
	}
}

// Test function for GenerateRandomString
func TestGenerateRandomString(t *testing.T) {
	length := 10
	result, err := GenerateRandomString(length)
	if err != nil {
		t.Fatalf("Failed to generate random string: %v", err)
	}

	if len(result) == 0 {
		t.Errorf("Expected a non-empty string")
	}
}

// Test function for GetJWTKey
func TestGetJWTKey(t *testing.T) {
	key := "test-key"
	expected := []byte(key)
	result := GetJWTKey(key)

	if string(result) != string(expected) {
		t.Errorf("Expected JWT key %v, got %v", expected, result)
	}
}
