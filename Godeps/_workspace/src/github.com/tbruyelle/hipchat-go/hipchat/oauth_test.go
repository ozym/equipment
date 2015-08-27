package hipchat

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	setup()
	defer teardown()

	clientID := "client-abcdef"
	clientSecret := "secret-12345"

	mux.HandleFunc("/oauth/token", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() != "/oauth/token" {
			t.Errorf("Incorrect URL = %v, want %v", r.URL, "/oauth/token")
		}

		if m := "POST"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}

		if r.Header.Get("Authorization") != "Basic Y2xpZW50LWFiY2RlZjpzZWNyZXQtMTIzNDU=" {
			t.Errorf("Incorrect authorization header")
		}

		if r.FormValue("grant_type") != "client_credentials" {
			t.Errorf("grant_type should be 'client_credentials'")
		}

		if r.FormValue("scope") != "send_notification view_room" {
			t.Errorf("scope should be 'send_notification view_room'")
		}

		fmt.Fprintf(w, `
		{
            "access_token": "GeneratedAuthToken",
            "expires_in": 3599,
            "group_id": 123456,
            "group_name": "TestGroup",
            "scope": "send_notification view_room",
            "token_type": "bearer"
        }
        `)
	})
	want := &OAuthAccessToken{
		AccessToken: "GeneratedAuthToken",
		ExpiresIn:   3599,
		GroupID:     123456,
		GroupName:   "TestGroup",
		Scope:       "send_notification view_room",
		TokenType:   "bearer",
	}

	credentials := ClientCredentials{ClientID: clientID, ClientSecret: clientSecret}

	token, _, err := client.GenerateToken(credentials, []string{ScopeSendNotification, ScopeViewRoom})
	if err != nil {
		t.Fatalf("Client.GetAccessToken returns an error %v", err)
	}
	if !reflect.DeepEqual(want, token) {
		t.Errorf("Client.GetAccessToken returned %+v, want %+v", token, want)
	}
	if client.authToken != want.AccessToken {
		t.Errorf("Client.authToken = %s, want %s", client.authToken, want.AccessToken)
	}
}
