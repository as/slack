package converse

type Channel struct {
	Name           string   `json:"name"`
	Id             string   `json:"id"`
	LastRead       string   `json:"last_read"`
	NumMembers     int      `json:"num_members"`
	Topic          Tag      `json:"topic"`
	Purpose        Tag      `json:"purpose"`
	Creator        string   `json:"creator"`
	Created        int      `json:"created"`
	NameNormalized string   `json:"name_normalized"`
	PreviousNames  []string `json:"previous_names"`
	Unlinked       int      `json:"unlinked"`
	Locale         string   `json:"locale"`
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
