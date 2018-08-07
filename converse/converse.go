package converse

type Reply struct {
	Ok      bool `json:"ok"`
	Channel `json:"channel"`
}

type Channel struct {
	Id             string     `json:"id"`
	Name           string     `json:"name"`
	NameNormalized string     `json:"name_normalized"`
	PreviousNames  []struct{} `json:"previous_names"`
	Creator        string     `json:"creator"`
	Created        int        `json:"created"`
	Unlinked       int        `json:"unlinked"`
	NumMembers     int        `json:"num_members"`
	Locale         string     `json:"locale"`
	LastRead       string     `json:"last_read"`
	Topic          Tag        `json:"topic"`
	Purpose        Tag        `json:"purpose"`
	Is
}

type Tag struct {
	Value   string `json:"value"`
	Creator string `json:"creator"`
	LastSet int    `json:"last_set"`
}

type Is struct {
	Archived         bool `json:"is_archived"`
	Channel          bool `json:"is_channel"`
	ExtShared        bool `json:"is_ext_shared"`
	General          bool `json:"is_general"`
	Group            bool `json:"is_group"`
	Im               bool `json:"is_im"`
	Member           bool `json:"is_member"`
	OrgShared        bool `json:"is_org_shared"`
	PendingExtShared bool `json:"is_pending_ext_shared"`
	Mpim             bool `json:"is_mpim"`
	Private          bool `json:"is_private"`
	ReadOnly         bool `json:"is_read_only"`
	Shared           bool `json:"is_shared"`
}
