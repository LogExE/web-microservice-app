package boxes

type NewBoxEvent struct {
	BoxID   int64  `json:"box_id"`
	Content string `json:"content"`
	Author  string `json:"author"`
}

type BoxLikedEvent struct {
	BoxID  int64  `json:"box_id"`
	Author string `json:"author"`
	Likes  int64  `json:"likes"`
}
