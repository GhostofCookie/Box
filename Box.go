package box

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var accessToken string
var configFile string

// BoxRequest : Runs an HTTP request to the box app.
func BoxRequest(method string, url string, body io.Reader, headers map[string]string) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	} else {
		if len(headers) == 0 {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	}
	if accessToken != "" {
		req.Header.Add("Authorization", "Bearer "+accessToken)
	}

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	respBytes, err := ioutil.ReadAll(response.Body)
	response.Body.Close()

	if response.StatusCode != 200 {
		log.Println(" >> URL    :", url)
		log.Println(" >> Status :", response.Status)
	}
	return respBytes, nil
}

// RequestAccessToken : Get valid ACCESS_TOKEN using JWT.
func RequestAccessToken() error {
	name, err := ioutil.ReadFile(configFile)
	var boxConfig BoxJWTRequest

	err = json.Unmarshal(name, &boxConfig)

	if err != nil {
		log.Println(err)
		return err
	}

	// Create a unique 32 character long string.
	rBytes := make([]byte, 32)
	_, err = rand.Read(rBytes)
	if err != nil {
		log.Println(err)
		return err
	}
	jti := base64.URLEncoding.EncodeToString(rBytes)

	// Build the header. This includes the PublicKey as the ID.
	token := jwt.New(jwt.SigningMethodRS512)
	token.Header["keyid"] = boxConfig.BoxAppSettings.AppAuth.PublicKeyID

	// Construct claims.
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = boxConfig.BoxAppSettings.ClientID
	claims["sub"] = boxConfig.EnterpriseID
	claims["box_sub_type"] = "enterprise"
	claims["aud"] = os.Getenv("authURL")
	claims["jti"] = jti
	claims["exp"] = time.Now().Add(time.Second * 10).Unix()

	// Decrypt the PrivateKey using its passphrase.
	signedKey, err := jwt.ParseRSAPrivateKeyFromPEMWithPassword(
		[]byte(boxConfig.BoxAppSettings.AppAuth.PrivateKey),
		boxConfig.BoxAppSettings.AppAuth.Passphrase,
	)

	if err != nil {
		log.Println(err)
		return err
	}

	// Build the assertion from the signedKey and claims.
	assertion, err := token.SignedString(signedKey)

	if err != nil {
		log.Println(err)
		return err
	}

	// Build the access token request.
	payload := url.Values{}
	payload.Add("grant_type", "urn:ietf:params:oauth:grant-type:jwt-bearer")
	payload.Add("assertion", assertion)
	payload.Add("client_id", boxConfig.BoxAppSettings.ClientID)
	payload.Add("client_secret", boxConfig.BoxAppSettings.ClientSecret)

	// Post the request to the Box API.
	response, err := BoxRequest("POST", os.Getenv("authURL"), bytes.NewBufferString(payload.Encode()), nil)
	if err != nil {
		log.Println(err)
		return err
	}

	// Set the access token.
	var js AccessResponse
	err = json.Unmarshal(response, &js)
	if err != nil {
		log.Println(err)
		return err
	}
	accessToken = js.AccessToken

	return nil
}
