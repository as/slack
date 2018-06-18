package user

type User struct {
	Name   string `json:"name"`
	Id     string `json:"id"`
	TeamID string `json:"team_id"`

	Color    string `json:"color"`
	RealName string `json:"real_name"`
	Deleted  bool   `json:"deleted"`
	Updated  int    `json:"updated"`

	Zone
	Flags

	Profile Profile `json:"profile"`
}

type Zone struct {
	Tz       string `json:"tz"`
	TzLabel  string `json:"tz_label"`
	TzOffset int    `json:"tz_offset"`
}

type Profile struct {
	Team       string `json:"team"`
	Email      string `json:"email"`
	AvatarHash string `json:"avatar_hash"`
	Name
	Status
	Image
}

type Name struct {
	Real        string `json:"real_name"`
	Display     string `json:"display_name"`
	RealNorm    string `json:"real_name_normalized,omitempty"`
	DisplayNorm string `json:"display_name_normalized,omitempty"`
}

type Image struct {
	I24  string `json:"image_24"`
	I32  string `json:"image_32"`
	I48  string `json:"image_48"`
	I72  string `json:"image_72"`
	I192 string `json:"image_192"`
	I512 string `json:"image_512"`
}

type Status struct {
	Text  string `json:"status_text"`
	Emoji string `json:"status_emoji"`
}

type Flags struct {
	PrimaryOwner    bool `json:"is_primary_owner"`
	Owner           bool `json:"is_owner"`
	AppUser         bool `json:"is_app_user"`
	Restricted      bool `json:"is_restricted"`
	UltraRestricted bool `json:"is_ultra_restricted"`
	Admin           bool `json:"is_admin"`
	Bot             bool `json:"is_bot"`
	TwoFA           bool `json:"has_2fa"`
}
