package dashgram

type TrackEventRequest struct {
	Updates []any  `json:"updates"`
	Origin  string `json:"origin,omitempty"`
}

type InvitedByRequest struct {
	UserID    int    `json:"user_id"`
	InvitedBy int    `json:"invited_by"`
	Origin    string `json:"origin,omitempty"`
}
