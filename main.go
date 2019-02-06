package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"io/ioutil"
	"os"
	"os/user"
	"path"
)

func getToken(config *oauth2.Config, file string) *oauth2.Token {
	tok, err := tokenFromFile(file)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(file, tok)
	}
	return tok
}

func configFromFile(file string) *oauth2.Config {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	conf, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		panic(err)
	}
	return conf
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func saveToken(path string, token *oauth2.Token) {
	f, err := os.OpenFile(path, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Println("Go to the following link in your browser then type the authorization code:")
	fmt.Println(authURL)
	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		panic(err)
	}
	tok, err := config.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		panic(err)
	}
	return tok
}

func main() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	dir := path.Join(usr.HomeDir, ".google-auth")
	conf := configFromFile(path.Join(dir, "client-secret.json"))
	tok := getToken(conf, path.Join(dir, "check-gmail.json"))
	client := conf.Client(oauth2.NoContext, tok)
	srv, err := gmail.New(client)
	if err != nil {
		panic(err)
	}
	res, err := srv.Users.Messages.List("me").Q("is:unread").Fields("messages(id)").Do()
	if err != nil {
		panic(err)
	}
	for i, m := range res.Messages {
		if i > 0 {
			fmt.Println()
		}
		m, err = srv.Users.Messages.Get("me", m.Id).Format("metadata").MetadataHeaders(
			"From", "Subject",
		).Do()
		if err != nil {
			panic(err)
		}
		headers := make(map[string]string)
		for _, h := range m.Payload.Headers {
			headers[h.Name] = h.Value
		}
		fmt.Printf("From: %s\n", headers["From"])
		fmt.Printf("Subject: %s\n", headers["Subject"])
	}
}
