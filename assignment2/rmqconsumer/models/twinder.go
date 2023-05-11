package models

type SwipePayload struct {
	Swiper    string `json:"swiper"`
	Swipee    string `json:"swipee"`
	Comment   string `json:"comment"`
	Direction string `json:"direction"`
}

type TwinderUserStats struct {
	UserId   string `json:"usedId"`
	Likes    int    `json:"likes"`
	Dislikes int    `json:"dislikes"`
}

type TwinderMatches struct {
	UserId  string   `json:"usedId"`
	Matches []string `json:"matches"`
}

type AllTwinderUserStats struct {
	UsersStats []TwinderUserStats `json:"users_stats"`
}
