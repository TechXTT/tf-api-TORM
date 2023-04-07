package main

// every 3599 seconds, the token expires and we need to refresh it with a new one from the API
import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Body struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token"`
	GrantType    string `json:"grant_type"`
}

type Response struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println("[ERROR] Could not load .env file: ", err)
		panic(err)
	}
	timer := time.NewTicker(1 * time.Second)
	access_token := ""
	client_id := os.Getenv("GMAIL_CLIENT_ID")
	client_secret := os.Getenv("GMAIL_CLIENT_SECRET")
	refresh_token := os.Getenv("GMAIL_REFRESH_TOKEN")

	body := Body{
		ClientID:     client_id,
		ClientSecret: client_secret,
		RefreshToken: refresh_token,
		GrantType:    "refresh_token",
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		fmt.Println("[ERROR] Could not marshal body: ", err)
		panic(err)
	}

	var response Response

	for {
		select {
		case <-timer.C:

			resp, err := http.Post("https://www.googleapis.com/oauth2/v4/token", "application/json", bytes.NewBuffer(jsonBody))
			if err != nil {
				fmt.Println("[ERROR] Could not post to API: ", err)
				panic(err)
			}

			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				fmt.Println("[ERROR] Could not get token: ", resp.StatusCode)
				panic(err)
			}

			err = json.NewDecoder(resp.Body).Decode(&response)
			if err != nil {
				print(err)
				panic(err)
			}

			access_token = response.AccessToken
			timer = time.NewTicker(time.Duration(response.ExpiresIn) * time.Second)

			f, err := os.OpenFile(".env", os.O_RDWR, 0644)
			if err != nil {
				fmt.Println("[ERROR] Could not open .env file: ", err)
				panic(err)
			}
			defer f.Close()

			fileInfo, err := f.Stat()
			if err != nil {
				fmt.Println("[ERROR] Could not get file info: ", err)
				panic(err)
			}

			fileSize := fileInfo.Size()

			_, err = f.Seek(0, 2)
			if err != nil {
				fmt.Println("[ERROR] Could not seek to end of file: ", err)
				panic(err)
			}

			for i := fileSize - 1; i >= 0; i-- {
				_, err = f.Seek(i, 0)
				if err != nil {
					fmt.Println("[ERROR] Could not seek to end of file: ", err)
					panic(err)
				}

				b := make([]byte, 1)
				_, err = f.Read(b)
				if err != nil {
					fmt.Println("[ERROR] Could not read file: ", err)
					panic(err)
				}
				if b[0] == '\n' {
					_, err = f.Seek(i+1, 0)
					if err != nil {
						fmt.Println("[ERROR] Could not seek to end of file: ", err)
						panic(err)
					}
					break
				}

			}

			err = f.Truncate(fileSize)
			if err != nil {
				fmt.Println("[ERROR] Could not truncate file: ", err)
				panic(err)
			}

			_, err = f.WriteString("GMAIL_ACCESS_TOKEN=" + access_token + "")
			if err != nil {
				fmt.Println("[ERROR] Could not write to file: ", err)
				panic(err)
			}

			f.Sync()
			fmt.Println("Access token refreshed: ", access_token)
		}
	}
}
