package unifi

import "fmt"

// Group represents a UniFi group.
type Group struct {
	ID      string   `json:"unique_id,omitempty"`
	Name    string   `json:"name,omitempty"`
	UpID    string   `json:"up_id,omitempty"`
	UpIDs   []string `json:"up_ids,omitempty"`
	System  string   `json:"system_name,omitempty"`
	Created string   `json:"create_time,omitempty"`
}

// Role represents a UniFi role.
type Role struct {
	ID        string `json:"unique_id,omitempty"`
	Name      string `json:"name,omitempty"`
	IsSystem  bool   `json:"system_role,omitempty"`
	SystemKey string `json:"system_key,omitempty"`
	Level     int64  `json:"level,omitempty"`
}

// Permission represents a UniFi permission.
type Permission map[string][]string

// Account represents an authorized account.
type Account struct {
	ID                 string       `json:"unique_id,omitempty"`
	FirstName          string       `json:"first_name,omitempty"`
	LastName           string       `json:"last_name,omitempty"`
	FullName           string       `json:"full_name,omitempty"`
	Email              string       `json:"email,omitempty"`
	EmailStatus        string       `json:"email_status,omitempty"`
	Phone              string       `json:"phone,omitempty"`
	AvatarRelativePath string       `json:"avatar_relative_path,omitempty"`
	Status             string       `json:"status,omitempty"`
	EmployeeNumber     string       `json:"employee_number,omitempty"`
	CreatedTS          int64        `json:"create_time,omitempty"`
	Username           string       `json:"username,omitempty"`
	IsLocal            bool         `json:"local_account_exist,omitempty"`
	PasswordRevision   int64        `json:"password_revision,omitempty"`
	SSOAccount         string       `json:"sso_account,omitempty"`
	SSOUUID            string       `json:"sso_uuid,omitempty"`
	SSOUsername        string       `json:"sso_username,omitempty"`
	SSOPicture         string       `json:"sso_picture,omitempty"`
	UIDSSOID           string       `json:"uid_sso_id,omitempty"`
	UIDSSOAccount      string       `json:"uid_sso_account,omitempty"`
	Groups             []Group      `json:"groups,omitempty"`
	Roles              []Role       `json:"roles,omitempty"`
	Permissions        []Permission `json:"permissions,omitempty"`
	Scopes             []string     `json:"scopes,omitempty"`
	CloudAccessGranted bool         `json:"cloud_access_granted,omitempty"`
	Updated            int64        `json:"update_time,omitempty"`
	Avatar             string       `json:"avatar,omitempty"`
	NFCToken           string       `json:"nfc_token,omitempty"`
	NFCDisplayID       string       `json:"nfc_display_id,omitempty"`
	NFCCardType        string       `json:"nfc_card_type,omitempty"`
	NFCCardStatus      string       `json:"nfc_card_status,omitempty"`
	DisplayID          string       `json:"id,omitempty"`
	IsOwner            bool         `json:"isOwner,omitempty"`
	IsSuperAdmin       bool         `json:"isSuperAdmin,omitempty"`
}

func (a *Account) String() string {
	if a == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%s %s %s %v", a.ID, a.FullName, a.Email, a.Roles)
}
