package box

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Config : Basic structure for a Box API JWT.
type Config struct {
	BoxAppSettings struct {
		ClientID     string `json:"clientID"`
		ClientSecret string `json:"clientSecret"`
		AppAuth      struct {
			PublicKeyID string `json:"publicKeyID"`
			PrivateKey  string `json:"privateKey"`
			Passphrase  string `json:"passphrase"`
		} `json:"appAuth"`
	} `json:"boxAppSettings"`
	EnterpriseID string `json:"enterpriseID"`
}

// AccessResponse : Object returned by a successful request to the Box API.
type AccessTokenObject struct {
	AccessToken     string `json:"access_token"`
	ExpiresIn       int    `json:"expires_in"`
	IssuedTokenType string `json:"issued_token_type"`
	RefreshToken    string `json:"refresh_token"`
	RestrictedTo    []struct {
		Scope  string      `json:"scope,omitempty"`
		Object *FileObject `json:"object,omitempty`
	} `json:"restricted_to,omitempty"`
	TokenType string `json:"token_type"`
}

// SDK is the structure for establishing the connection to the Box API.
type SDK struct {
	access *AccessTokenObject
	config *Config
	client *http.Client
}

// NewConfig sets the configuration for the SDK to establish it's connection.
func (sdk *SDK) NewConfig(cfg *Config) {
	sdk.config = cfg
	sdk.client = &http.Client{}
}

// NewConfigFromFile sets the config file to read Box info from.
func (sdk *SDK) NewConfigFromFile(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = json.Unmarshal(content, &sdk.config)
	if err != nil {
		log.Fatal(err)
		return err
	}

	sdk.client = &http.Client{}

	return nil
}

// Request runs an HTTP request to the Box API.
func (sdk *SDK) Request(method string, url string, body io.Reader, headers map[string]string) ([]byte, error) {

	request, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	// Add all user specified headers to the request header.
	if headers != nil {
		for k, v := range headers {
			request.Header.Set(k, v)
		}
	}

	// Check if we have a valid access token object.
	if sdk.access != nil {
		request.Header.Add("Authorization", "Bearer "+sdk.access.AccessToken)
	}

	response, err := sdk.client.Do(request)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	defer response.Body.Close()

	respBytes, err := ioutil.ReadAll(response.Body)
	response.Body.Close()

	log.Println("URL    :", url)
	log.Println("Status :", response.Status)

	return respBytes, nil
}

// RequestAccessToken requests a valid access token from the Box API.
func (sdk *SDK) RequestAccessToken() error {
	if sdk.config == nil {
		return errors.New("No API configuration set")
	}
	// Create a unique 32 character long string.
	rBytes := make([]byte, 32)
	_, err := rand.Read(rBytes)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	jti := base64.URLEncoding.EncodeToString(rBytes)

	// Build the header. This includes the PublicKey as the ID.
	token := jwt.New(jwt.SigningMethodRS512)
	token.Header["keyid"] = sdk.config.BoxAppSettings.AppAuth.PublicKeyID

	// Construct claims.
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = sdk.config.BoxAppSettings.ClientID
	claims["sub"] = sdk.config.EnterpriseID
	claims["box_sub_type"] = "enterprise"
	claims["aud"] = "https://api.box.com/oauth2/token"
	claims["jti"] = jti
	claims["exp"] = time.Now().Add(time.Second * 3).Unix()

	// Decrypt the PrivateKey using its passphrase.
	signedKey, err := jwt.ParseRSAPrivateKeyFromPEMWithPassword(
		[]byte(sdk.config.BoxAppSettings.AppAuth.PrivateKey),
		sdk.config.BoxAppSettings.AppAuth.Passphrase,
	)

	if err != nil {
		log.Fatalln(err)
		return err
	}

	// Build the assertion from the signedKey and claims.
	assertion, err := token.SignedString(signedKey)

	if err != nil {
		log.Fatalln(err)
		return err
	}

	// Build header
	header := make(map[string]string)
	header["Content-Type"] = "application/x-www-form-urlencoded"

	// Build the access token request.
	payload := url.Values{}
	payload.Add("grant_type", "urn:ietf:params:oauth:grant-type:jwt-bearer")
	payload.Add("assertion", assertion)
	payload.Add("client_id", sdk.config.BoxAppSettings.ClientID)
	payload.Add("client_secret", sdk.config.BoxAppSettings.ClientSecret)

	// Post the request to the Box API.
	response, err := sdk.Request("POST", "https://api.box.com/oauth2/token", bytes.NewBufferString(payload.Encode()), header)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	// Set the access token.
	err = json.Unmarshal(response, &sdk.access)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	return nil
}
