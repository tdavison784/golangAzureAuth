package main

import (
	"context"
	"fmt"
	"github.com/coreos/go-oidc"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/oauth2"
	"log"
	"net/http"
)

// we are creating an oauth2.Config type that contains all needed information in order to connect
// to our oauth Provider. In this case, that is Azure Entra ID.
// Place in your ClientID, ClientSecret, and tenant values as well as any custom scopes as needed.
var (
	azureADOAuthConfig = &oauth2.Config{
		ClientID:     "xxx-xxx-xxx",
		ClientSecret: "xxx-xxx-xxx",
		RedirectURL:  "http://localhost:8080/callback",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://login.microsoftonline.com/{tenant_id_here}/oauth2/v2.0/authorize",
			TokenURL: "https://login.microsoftonline.com/{tenant_id_here}/oauth2/v2.0/token",
		},
		Scopes: []string{"openid", "email", "profile", "offline_access"},
	}
	oauthStateString = "testing123"
	oidcProvider     = "https://login.microsoftonline.com/{tenant_id_here}/v2.0"
)

func main() {
	router := httprouter.New()
	router.GET("/", HomeHandler)
	router.GET("/login", LoginHandler)
	router.GET("/callback", CallbackHandler)
	router.GET("/protected", ProtectedHandler)
	log.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", router))

}

func HomeHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Welcome to the home page!")
}

func LoginHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	url := azureADOAuthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func CallbackHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	state := r.URL.Query().Get("state")
	if state != oauthStateString {
		fmt.Println("State did not match")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.URL.Query().Get("code")
	token, err := azureADOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Println("Error exchanging code:", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// You can store and use the token for making authorized requests
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token.AccessToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})

	http.Redirect(w, r, "/protected", http.StatusSeeOther)
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cookie, err := r.Cookie("access_token")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// You should validate the access token here before allowing access to protected resources
	accessToken := cookie.Value
	if err := ValidateToken(accessToken); err != nil {
		fmt.Println("Error validating access token", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Fprintf(w, "Protected API Endpoint")
}

func ValidateToken(accessToken string) error {
	ctx := context.Background()
	fmt.Println(oidcProvider)
	provider, err := oidc.NewProvider(ctx, oidcProvider)
	if err != nil {
		return err
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: azureADOAuthConfig.ClientID, SkipIssuerCheck: true})
	idToken, err := verifier.Verify(ctx, accessToken)
	if err != nil {
		return err
	}
	if err != nil {
		fmt.Println("Error putting claim into claims", err)
	}
	var claims struct {
		IpAddr string `json:"ipaddr"`
		UPN    string `json:"upn"`
	}
	idToken.Claims(&claims)
	fmt.Println(claims)

	// You can further validate claims if needed
	// For example:
	// if idToken.Subject != "expected_subject" {
	//     return errors.New("unexpected subject")
	// }

	return nil
}
