package graw

type redditorResponse struct {
	Redditor
	err string `json:"error"`
}

type Redditor struct {
	HasMail bool `json:"has_mail"`
	Name string `json:"name"`
	Created float32 `json:"created"`
	HideFromRobots bool `json:"hide_from_robots"`
	GoldCredits int `json:"gold_credits"`
	CreatedUTC float32 `json:"created_utc"`
	HasModMail bool `json:"has_mod_mail"`
	LinkKarma int `json:"link_karma"`
	CommentKarma int `json:"comment_karma"`
	Over18 bool `json:"over_18"`
	IsGold bool `json:"is_gold"`
	IsMod bool `json:"is_mod"`
	GoldExpiration int `json:"gold_expiration"`
	HasVerifiedEmail bool `json:"has_verified_email"`
	ID string `json:"id"`
	InboxCount int `json:"inbox_count"`
}
