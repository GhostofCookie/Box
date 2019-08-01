package box

// BoxJWTRequest : Basic structure for a Box API JWT.
type BoxJWTRequest struct {
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
type AccessResponse struct {
	AccessToken  string        `json:"access_token"`
	ExpiresIn    int           `json:"expires_in"`
	RestrictedTo []interface{} `json:"restricted_to"`
	TokenType    string        `json:"token_type"`
}

// FileObject : File information describe file objects in Box, with attributes
// like who created the file, when it was last modified, and other information.
// The actual content of the file itself is accessible through the
// /files/{id}/content endpoint. Italicized attributes are not returned by
// default and must be retrieved through the fields parameter.
type FileObject struct {
	Type           string         `json:"type"`
	ID             string         `json:"id"`
	FileVersion    FileVersion    `json:"file_version"`
	SequenceID     string         `json:"sequence_id"`
	Etag           string         `json:"etag"`
	Sha1           string         `json:"sha1"`
	Name           string         `json:"name"`
	Description    string         `json:"description"`
	Size           int            `json:"size"`
	PathCollection PathCollection `json:"path_collection"`
	CreatedAt      string         `json:"created_at"`
	ModifiedAt     string         `json:"modified_at"`
	CreatedBy      User           `json:"created_by"`
	ModifiedBy     User           `json:"modified_by"`
	OwnedBy        User           `json:"owned_by"`
	SharedLink     SharedLink     `json:"shared_link"`
	Parent         Parent         `json:"parent"`
	ItemStatus     string         `json:"item_status"`
}

type FileVersionObject struct {
	Type       string `json:"type"`
	ID         string `json:"id"`
	Sha1       string `json:"sha1"`
	Name       string `json:"name"`
	Size       int    `json:"size"`
	CreatedAt  string `json:"created_at"`
	ModifiedAt string `json:"modified_at"`
	ModifiedBy User
	TrashedAt  string `json:"trashed_at"`
	PurgedAt   string `json:"purged_at"`
}

// FolderObject : A Box Folder object.
type FolderObject struct {
	Type              string            `json:"type"`
	ID                string            `json:"id"`
	SequenceID        string            `json:"sequence_id"`
	Etag              string            `json:"etag"`
	Name              string            `json:"name"`
	CreatedAt         string            `json:"created_at"`
	ModifiedAt        string            `json:"modified_at"`
	Description       string            `json:"description"`
	Size              int               `json:"size"`
	PathCollection    PathCollection    `json:"path_collection,omitempty"`
	CreatedBy         User              `json:"created_by,omitempty"`
	ModifiedBy        User              `json:"modified_by,omitempty"`
	OwnedBy           User              `json:"owned_by,omitempty"`
	SharedLink        SharedLink        `json:"shared_link,omitempty"`
	FolderUploadEmail FolderUploadEmail `json:"folder_upload_email,omitempty"`
	Parent            Parent            `json:"parent,omitempty"`
	ItemStatus        string            `json:"item_status"`
	ItemCollection    ItemCollection    `json:"item_collection,omitempty"`
	Tags              []string          `json:"tags"`
}

// FileVersion : Contains version information of a FileObject.
type FileVersion struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	Sha1 string `json:"sha1"`
}

// Entries : A more in-depth response containing more information about box objects.
type Entries struct {
	EntriesMini
	Sha1              string         `json:"sha1 "`
	Description       string         `json:"description"`
	Size              int            `json:"size"`
	PathCollection    PathCollection `json:"path_collection,omitempty"`
	CreatedAt         string         `json:"created_at"`
	ModifiedAt        string         `json:"modified_at"`
	TrashedAt         interface{}    `json:"trashed_at,omitempty"`
	PurgedAt          interface{}    `json:"purged_at,omitempty"`
	ContentCreatedAt  string         `json:"content_created_at"`
	ContentModifiedAt string         `json:"content_modified_at"`
	CreatedBy         User           `json:"created_by,omitempty"`
	ModifiedBy        User           `json:"modified_by,omitempty"`
	OwnedBy           User           `json:"owned_by,omitempty"`
	SharedLink        SharedLink     `json:"shared_link,omitempty"`
	Parent            Parent         `json:"parent,omitempty"`
	ItemStatus        string         `json:"item_status"`
}

// EntriesMini : Basic structure for response carrying info about box objects.
type EntriesMini struct {
	Type       string      `json:"type"`
	ID         string      `json:"id"`
	SequenceID interface{} `json:"sequence_id,omitempty"`
	Etag       string      `json:"etag,omitempty"`
	Name       string      `json:"name"`
}

// PathCollection : The total amount of entries in a given path, as well as the entries themselves.
type PathCollection struct {
	TotalCount int       `json:"total_count"`
	Entries    []Entries `json:"entries"`
}

// User : Contains information about a Box user.
type User struct {
	Type  string `json:"type"`
	ID    string `json:"id"`
	Name  string `json:"name"`
	Login string `json:"login"`
}

// Permissions : Flags for downloaded files.
type Permissions struct {
	CanDownload bool `json:"can_download"`
	CanPreview  bool `json:"can_preview"`
}

// SharedLink : A shared link to a downloadable file.
type SharedLink struct {
	URL               string      `json:"url"`
	DownloadURL       interface{} `json:"download_url,omitempty"`
	VanityURL         interface{} `json:"vanity_url,omitempty"`
	IsPasswordEnabled bool        `json:"is_password_enabled"`
	UnsharedAt        interface{} `json:"unshared_at,omitempty"`
	DownloadCount     int         `json:"download_count"`
	PreviewCount      int         `json:"preview_count"`
	Access            string      `json:"access"`
	Permissions       Permissions `json:"permissions,omitempty"`
}

// FolderUploadEmail : Access level and email address of upload folder.
type FolderUploadEmail struct {
	Access string `json:"access"`
	Email  string `json:"email"`
}

// Parent : Parent folder of a returned box object.
type Parent struct {
	Type       string      `json:"type"`
	ID         string      `json:"id"`
	SequenceID interface{} `json:"sequence_id,omitempty"`
	Etag       interface{} `json:"etag,omitempty"`
	Name       string      `json:"name"`
}

// ItemCollection : Total count up to the limit of the number of entries in a folder, as well as the entries themselves.
type ItemCollection struct {
	TotalCount int           `json:"total_count"`
	Entries    []EntriesMini `json:"entries"`
	Offset     int           `json:"offset"`
	Limit      int           `json:"limit"`
}

// Order : Defines how to sort objects.
type Order struct {
	By        string `json:"by"`
	Direction string `json:"direction"`
}

// EmbeddedFile : An HTML embeddable file.
type EmbeddedFile struct {
	Type              string `json:"type"`
	ID                string `json:"id"`
	Etag              string `json:"etag"`
	ExpiringEmbedLink struct {
		URL string `json:"url"`
	} `json:"expiring_embed_link"`
}
