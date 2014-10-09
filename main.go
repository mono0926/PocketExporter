package main

import (
	"github.com/Shaked/getpocket/auth"
	"log"
	"net/http"
	"fmt"
	"github.com/mono0926/getpocket/commands"
)

type Page struct {
	auth *auth.Auth
}

var (
	redirectURI     = "https://localhost:10443/authcheck"
	consumerKey     = "32988-d9f4685756fba39521660fb2"
	ssl_certificate = "/Users/mono/go/src/workspace/PocketExporter/ssl/server.crt"
	ssl_key         = "/Users/mono/go/src/workspace/PocketExporter/ssl/server.key"
)

func main() {

	a, e := auth.Factory(consumerKey, redirectURI)
	if nil != e {
		log.Fatal(e)
	}
	log.Printf("Listen on 10443")
	p := &Page{auth: a}
	http.HandleFunc("/auth", p.Auth)
	http.HandleFunc("/authcheck", p.AuthCheck)

	err := http.ListenAndServeTLS(":10443", ssl_certificate, ssl_key, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (p *Page) Auth(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Referer(), "Example Handler, Should connect and request app permissions")
	requestToken, err := p.auth.Connect()
	if nil != err {
		fmt.Fprintf(w, "Token error %s (%d)", err.Error(), err.ErrorCode())
	}

	p.auth.RequestPermissions(requestToken, w, r)
}

func (p *Page) AuthCheck(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Referer(), "GetPocket connection check, should get the username and access token")
	requestToken := r.URL.Query().Get("requestToken")
	if "" == requestToken {
		fmt.Fprintf(w, "Request token is invalid")
		return
	}
	user, err := p.auth.User(requestToken)
	if nil != err {
		fmt.Fprintf(w, "%s (%d)", err.Error(), err.ErrorCode())
		return
	}
	fmt.Fprintf(w, "%#v", user)



	retrieve := commands.NewRetrieve()
	retrieve.SetState("all")
	retrieve.SetTag("golang")
	c := commands.New(user, consumerKey)
	resp, e := c.Exec(retrieve)
	if nil != e {
		fmt.Errorf("ERROR%s\n", e)
		return
	}
	list := resp.(*commands.RetrieveResponse).List
	for k, item := range list {
		fmt.Fprintln(w, "key: ", k)
		fmt.Fprintln(w, "title: ", item.GivenTitle)
		fmt.Fprintln(w, "url: ", item.GivenURL)
		fmt.Fprintln(w)
	}
}
