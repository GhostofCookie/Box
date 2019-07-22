package box

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

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

// UploadFile : Creates an Access Token to the Box API, then uploads a given name to the specified folder.
func UploadFile(file interface{}, newName string, folderID string) (*PathCollection, error) {
	RequestAccessToken()

	t := reflect.TypeOf(file)

	var name string
	if t.Name() == "string" {
		name = file.(string)
	} else {
		name = newName
	}
	if newName == "" && name != "" {
		newName = name
	}

	var contents []byte
	if t.Name() == "" {
		contents = file.([]byte)
	} else {
		f, err := os.Open(name)
		if err != nil {
			log.Println(err)
		}
		defer f.Close()

		contents, err = ioutil.ReadAll(f)
		if err != nil {
			log.Println(err)
		}
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", name)
	if err != nil {
		log.Println(err)
	}
	part.Write(contents)

	err = writer.WriteField("filename", name)
	if err != nil {
		log.Println(err)
	}

	err = writer.Close()
	if err != nil {
		log.Println(err)
	}

	headers := make(map[string]string)
	headers["Content-Type"] = writer.FormDataContentType()
	headers["Content-Length"] = string(body.Len())

	response, err := BoxRequest("POST",
		"https://upload.box.com/api/2.0/files/content?attributes={%22name%22:%22"+newName+"%22,%20%22parent%22:{%22id%22:%22"+folderID+"%22}}",
		body, headers)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	fileObject := &PathCollection{}
	json.Unmarshal(response, &fileObject)

	return fileObject, nil
}

// GetFileInfo : Returns information about the file with 'ID' fileID.
func GetFileInfo(fileID string) (*FileObject, error) {
	RequestAccessToken()
	response, err := BoxRequest("GET", "https://api.box.com/2.0/files/"+fileID, nil, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	fileObject := &FileObject{}
	json.Unmarshal(response, &fileObject)

	return fileObject, nil
}

// GetEmbedLink : Returns information about the file with 'ID' fileID.
func GetEmbedLink(fileID string) (*EmbeddedFile, error) {
	RequestAccessToken()
	response, err := BoxRequest("GET", "https://api.box.com/2.0/files/"+fileID+"?fields=expiring_embed_link", nil, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	fileObject := &EmbeddedFile{}
	json.Unmarshal(response, &fileObject)

	return fileObject, nil
}

// DownloadFile : Downloads a file with 'ID' fileID.
func DownloadFile(fileID string, location string) error {
	RequestAccessToken()
	response, err := BoxRequest("GET", "https://api.box.com/2.0/files/"+fileID+"/content", nil, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	fInfo, err := GetFileInfo(fileID)
	file, err := os.Create(location + fInfo.Name)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	_, err = file.Write(response)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// DeleteFile : Deletes a file in a specific folder with 'ID" fileID.
func DeleteFile(fileID string, etag string) error {
	RequestAccessToken()
	headers := make(map[string]string)
	headers["If-Match"] = etag
	_, err := BoxRequest("DELETE", "https://api.box.com/2.0/files/"+fileID, nil, headers)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Folder Functions

// CreateFolder : Creates a new folder under the parent folder that has 'ID' parentFolderID.
func CreateFolder(name string, parentFolderID string) (*FolderObject, error) {
	RequestAccessToken()
	body := strings.NewReader(`{"name":"` + name + `", "parent": {"id": "` + parentFolderID + `"}}`)

	response, err := BoxRequest("POST", "https://api.box.com/2.0/folders", body, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	folderObject := &FolderObject{}
	json.Unmarshal(response, &folderObject)

	return folderObject, nil
}

// GetFolderItems : Returns all the items contained inside the folder with 'ID' folderID.
func GetFolderItems(folderID string, limit int, offset int) (*ItemCollection, error) {
	RequestAccessToken()

	response, err := BoxRequest("GET", "https://api.box.com/2.0/folders/"+folderID+"/items?limit="+strconv.Itoa(limit)+"&offset="+strconv.Itoa(offset), nil, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	items := &ItemCollection{}
	json.Unmarshal(response, &items)

	return items, nil
}

// DeleteFolder : Deletes the folder with 'ID' folderID.
func DeleteFolder(folderID string) {
	RequestAccessToken()
	_, err := BoxRequest("DELETE", "https://api.box.com/2.0/folders/"+folderID+"?recursive=true", nil, nil)
	if err != nil {
		log.Println(err)
		return
	}
}
