package vkclient

type Audio struct {
	Artist string `json:"artist"`
	Title  string `json:"title"`
}

type Attachment struct {
	Type  string `json:"type"`
	Audio Audio  `json:"audio"`
}

type Item struct {
	Id          int          `json:"id"`
	Date        int          `json:"date"`
	OwnerId     int          `json:"owner_id"`
	FromId      int          `json:"from_id"`
	PostType    string       `json:"post_type"`
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
}

type Response struct {
	Items      []Item `json:"items"`
	NextFrom   string `json:"next_from"`
	Count      int    `json:"count"`
	TotalCount int    `json:"total_count"`
}

type VKModel struct {
	Response `json:"response"`
}
