package listing

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
	"io/ioutil"
	"log"
	"net/http"
)

func init() {
	http.HandleFunc("/", getServiceListing)
}

func getServiceListing(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	resp, err := getClientWithOAuthContext(ctx).Get("https://appengine.googleapis.com/v1/apps/zen-formations/services")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("HTTP GET returned status %v", resp.Status)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	transformAndDisplay(body, w)
}

func getClientWithOAuthContext(ctx context.Context) *http.Client {
	transport := &oauth2.Transport{
		Source: google.AppEngineTokenSource(ctx, "https://www.googleapis.com/auth/cloud-platform.read-only"),
		Base:   &urlfetch.Transport{Context: ctx},
	}
	return &http.Client{Transport: transport}
}

type AppEngineServices struct {
	Services []AppEngineService `json:"services"`
}

type AppEngineService struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func transformAndDisplay(body []byte, w http.ResponseWriter) {
	var services AppEngineServices
	log.Println("construction du json")
	err := json.Unmarshal(body, &services)
	if err != nil {
		log.Print("err2", err)
	}
	fmt.Fprintf(w, "<ul>")
	for index := range services.Services {
		fmt.Fprintf(w, "<li>%s</li>", services.Services[index].Id)
	}
	fmt.Fprintf(w, "</ul>")
}
