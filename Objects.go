package box

// FileObject : File information describe file objects in Box, with attributes
// like who created the file, when it was last modified, and other information.
// The actual content of the file itself is accessible through the
// /files/{id}/content endpoint. Italicized attributes are not returned by
// default and must be retrieved through the fields parameter.
type FileObject struct {
	Type           string          `json:"type,omitempty"`
	ID             string          `json:"id,omitempty"`
	FileVersion    *FileVersion    `json:"file_version,omitempty"`
	SequenceID     string          `json:"sequence_id,omitempty"`
	Etag           string          `json:"etag,omitempty"`
	Sha1           string          `json:"sha1,omitempty"`
	Name           string          `json:"name,omitempty"`
	Description    string          `json:"description,omitempty"`
	Size           int             `json:"size,omitempty"`
	PathCollection *PathCollection `json:"path_collection,omitempty"`
	CreatedAt      string          `json:"created_at,omitempty"`
	ModifiedAt     string          `json:"modified_at,omitempty"`
	CreatedBy      *User           `json:"created_by,omitempty"`
	ModifiedBy     *User           `json:"modified_by,omitempty"`
	OwnedBy        *User           `json:"owned_by,omitempty"`
	SharedLink     *SharedLink     `json:"shared_link,omitempty"`
	Parent         *Parent         `json:"parent,omitempty"`
	ItemStatus     string          `json:"item_status,omitempty"`
}

// FileVersion : Contains version information of a FileObject.
type FileVersion struct {
	Type       string `json:"type,omitempty"`
	ID         string `json:"id,omitempty"`
	Sha1       string `json:"sha1,omitempty"`
	Name       string `json:"name,omitempty"`
	Size       int    `json:"size,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	ModifiedAt string `json:"modified_at,omitempty"`
	ModifiedBy *User  `json:"modified_by,omitempty"`
	TrashedAt  string `json:"trashed_at,omitempty"`
	PurgedAt   string `json:"purged_at,omitempty"`
}

// FolderObject : A Box Folder object.
type FolderObject struct {
	Type              string             `json:"type,omitempty"`
	ID                string             `json:"id,omitempty"`
	SequenceID        string             `json:"sequence_id,omitempty"`
	Etag              string             `json:"etag,omitempty"`
	Name              string             `json:"name,omitempty"`
	CreatedAt         string             `json:"created_at,omitempty"`
	ModifiedAt        string             `json:"modified_at,omitempty"`
	Description       string             `json:"description,omitempty"`
	Size              int                `json:"size,omitempty"`
	PathCollection    *PathCollection    `json:"path_collection,omitempty"`
	CreatedBy         *User              `json:"created_by,omitempty"`
	ModifiedBy        *User              `json:"modified_by,omitempty"`
	OwnedBy           *User              `json:"owned_by,omitempty"`
	SharedLink        *SharedLink        `json:"shared_link,omitempty"`
	FolderUploadEmail *FolderUploadEmail `json:"folder_upload_email,omitempty"`
	Parent            *Parent            `json:"parent,omitempty"`
	ItemStatus        string             `json:"item_status,omitempty"`
	ItemCollection    *ItemCollection    `json:"item_collection,omitempty"`
	Tags              []string           `json:"tags,omitempty"`
}

// Entries : A more in-depth response containing more information about box objects.
type Entries struct {
	Type              string          `json:"type,omitempty"`
	ID                string          `json:"id,omitempty"`
	SequenceID        interface{}     `json:"sequence_id,omitempty"`
	Etag              string          `json:"etag,omitempty"`
	Name              string          `json:"name,omitempty"`
	Sha1              string          `json:"sha1,omitempty"`
	Description       string          `json:"description,omitempty"`
	Size              int             `json:"size,omitempty"`
	PathCollection    *PathCollection `json:"path_collection,omitempty"`
	CreatedAt         string          `json:"created_at,omitempty"`
	ModifiedAt        string          `json:"modified_at,omitempty"`
	TrashedAt         interface{}     `json:"trashed_at,omitempty"`
	PurgedAt          interface{}     `json:"purged_at,omitempty"`
	ContentCreatedAt  string          `json:"content_created_at,omitempty"`
	ContentModifiedAt string          `json:"content_modified_at,omitempty"`
	CreatedBy         *User           `json:"created_by,omitempty"`
	ModifiedBy        *User           `json:"modified_by,omitempty"`
	OwnedBy           *User           `json:"owned_by,omitempty"`
	SharedLink        *SharedLink     `json:"shared_link,omitempty"`
	Parent            *Parent         `json:"parent,omitempty"`
	ItemStatus        string          `json:"item_status,omitempty"`
}

// PathCollection : The total amount of entries in a given path, as well as the entries themselves.
type PathCollection struct {
	TotalCount int        `json:"total_count,omitempty"`
	Entries    []*Entries `json:"entries,omitempty"`
}

// User : Contains information about a Box user.
type User struct {
	Type  string `json:"type,omitempty"`
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Login string `json:"login,omitempty"`
}

// Permissions : Flags for downloaded files.
type Permissions struct {
	CanDelete              bool `json:"can_delete,omitempty"`
	CanDownload            bool `json:"can_download,omitempty"`
	CanInviteCollaborator  bool `json:"can_invite_collaborator,omitempty"`
	CanRename              bool `json:"can_rename,omitempty"`
	CanSetShareAccess      bool `json:"can_set_share_access,omitempty"`
	CanShare               bool `json:"can_share,omitempty"`
	CanAnnotate            bool `json:"can_annotate,omitempty"`
	CanComment             bool `json:"can_comment,omitempty"`
	CanPreview             bool `json:"can_preview,omitempty"`
	CanUpload              bool `json:"can_upload,omitempty"`
	CanViewAnnotationsAll  bool `json:"can_view_annotations_all,omitempty"`
	CanViewAnnotationsSelf bool `json:"can_view_annotations_self,omitempty"`
}

// SharedLink : A shared link to a downloadable file.
type SharedLink struct {
	URL               string       `json:"url,omitempty"`
	DownloadURL       interface{}  `json:"download_url,omitempty"`
	VanityURL         interface{}  `json:"vanity_url,omitempty"`
	IsPasswordEnabled bool         `json:"is_password_enabled,omitempty"`
	UnsharedAt        interface{}  `json:"unshared_at,omitempty"`
	DownloadCount     int          `json:"download_count,omitempty"`
	PreviewCount      int          `json:"preview_count,omitempty"`
	Access            string       `json:"access,omitempty"`
	Permissions       *Permissions `json:"permissions,omitempty"`
}

// FolderUploadEmail : Access level and email address of upload folder.
type FolderUploadEmail struct {
	Access string `json:"access,omitempty"`
	Email  string `json:"email,omitempty"`
}

// Parent : Parent folder of a returned box object.
type Parent struct {
	Type       string      `json:"type,omitempty"`
	ID         string      `json:"id,omitempty"`
	SequenceID interface{} `json:"sequence_id,omitempty"`
	Etag       interface{} `json:"etag,omitempty"`
	Name       string      `json:"name,omitempty"`
}

// ItemCollection : Total count up to the limit of the number of entries in a folder, as well as the entries themselves.
type ItemCollection struct {
	TotalCount int        `json:"total_count,omitempty"`
	Entries    []*Entries `json:"entries,omitempty"`
	Offset     int        `json:"offset,omitempty"`
	Limit      int        `json:"limit,omitempty"`
}

// Order : Defines how to sort objects.
type Order struct {
	By        string `json:"by,omitempty"`
	Direction string `json:"direction,omitempty"`
}

// EmbeddedFile : An HTML embeddable file.
type EmbeddedFile struct {
	Type              string `json:"type,omitempty"`
	ID                string `json:"id,omitempty"`
	Etag              string `json:"etag,omitempty"`
	ExpiringEmbedLink struct {
		URL string `json:"url,omitempty"`
	} `json:"expiring_embed_link,omitempty"`
}
