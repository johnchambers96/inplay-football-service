package sportMonks

import (
	"encoding/json"
	client "inplay-football-service/internal"
	"io"
	"os"
)

type Participants struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	ShortCode string `json:"short_code"`
}

type Data struct {
	Id           int            `json:"id"`
	Name         string         `json:"name"`
	Participants []Participants `json:"participants"`
}

type GetFixturesResponse struct {
	Data     []Data `json:"data"`
	Timezone string `json:"timezone"`
}

func GetTodaysFixtures() (GetFixturesResponse, error) {
	var data GetFixturesResponse
	client := client.Create(os.Getenv("SPORT_MONKS_BASE_URL"))
	endpoint := "/v3/football/fixtures/date/2023-11-04?include=participants"
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": os.Getenv("SPORT_MONKS_AUTH_TOKEN"),
	}

	res, err := client.MakeRequest("GET", endpoint, nil, headers)

	if err != nil {
		return data, err
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return data, err
	}

	json.Unmarshal(body, &data)
	defer res.Body.Close()

	return data, nil
}
