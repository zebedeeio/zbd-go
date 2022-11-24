package zebedee

import "net/http"

type PublicUser struct {
	Name                   string            `json:"name"`
	Image                  string            `json:"image"`
	Username               string            `json:"username"`
	PublicBio              string            `json:"publicBio"`
	PublicStaticCharge     string            `json:"publicStaticCharge"`
	IsPublicPayPageEnabled bool              `json:"isPublicPayPageEnabled"`
	Social                 map[string]string `json:"social"`
}

func GetPublicGamertagData(gamertag string) (*PublicUser, error) {
	client := &Client{
		BaseURL:    "https://api.zebedee.io/public/v1",
		APIKey:     "blank",
		HttpClient: &http.Client{},
	}

	var pu PublicUser
	err := client.MakeRequest("GET", "/user/"+gamertag, nil, &pu)
	return &pu, err
}
