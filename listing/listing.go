package listing

import (
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
	"io/ioutil"
	"net/http"
)

func init() {
	http.HandleFunc("/", getServiceListing)
}

func getServiceListing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!\n")
	ctx := appengine.NewContext(r)
	transport := &oauth2.Transport{
		Source: google.AppEngineTokenSource(ctx, "https://www.googleapis.com/auth/cloud-platform.read-only"),
		Base:   &urlfetch.Transport{Context: ctx},
	}
	client := &http.Client{Transport: transport}
	resp, err := client.Get("https://appengine.googleapis.com/v1/apps/zen-formations/services")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Fprintf(w, "HTTP GET returned status %v", resp.Status)
	fmt.Fprintf(w, "\n%s", body)

}
