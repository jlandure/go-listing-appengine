package listing

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func init() {
	http.HandleFunc("/api/services", getServiceListingApi)
	http.HandleFunc("/", getServiceListingHtml)
}

func getServiceListing(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	ctx := appengine.NewContext(r)
	resp, err := getClientWithOAuthContext(ctx).Get("https://appengine.googleapis.com/v1/apps/zen-formations/services")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}
	log.Printf("HTTP GET returned status %v", resp.Status)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, nil
}

func getServiceListingApi(w http.ResponseWriter, r *http.Request) {
	body, err := getServiceListing(w, r)
	if err != nil {
		return
	}
	var services AppEngineServices
	log.Println("construction du json")
	err = json.Unmarshal(body, &services)
	services = transform(services)
	data, err := json.Marshal(services)
	if err != nil {
		log.Print("err mashal", err)
	}
	fmt.Fprintf(w, "%s", data)
}

func getServiceListingHtml(w http.ResponseWriter, r *http.Request) {
	body, err := getServiceListing(w, r)
	if err != nil {
		return
	}
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
	Id    string `json:"id"`
	Name  string `json:"name"`
	Url   string `json:"url"`
	Title string `json:"title"`
}

func formatName(name string) string {
	// Function replacing words (assuming lower case input)
	replace := func(word string) string {
		// switch word {
		// case "-", "zen":
		// 	return ""
		// }

		if word[len(word)-1] == '-' {
			word = word[:len(word)-1] + " "
		}
		return strings.Title(word)
	}

	r := regexp.MustCompile(`(\w+-|\w+)`)
	formattedName := r.ReplaceAllStringFunc(name, replace)

	log.Println(formattedName)
	return formattedName
}

func transform(services AppEngineServices) AppEngineServices {
	for index := range services.Services {
		services.Services[index].Title = formatName(services.Services[index].Id)
		log.Printf("%s", services.Services[index].Title)
		services.Services[index].Url = fmt.Sprintf("https://%s-dot-zen-formations.appspot.com/", services.Services[index].Id)
		// 	fmt.Fprintf(w, "<li>%s</li>", services.Services[index].Id)
	}
	return services
}

func transformAndDisplay(body []byte, w http.ResponseWriter) {
	var services AppEngineServices
	log.Println("construction du json")
	err := json.Unmarshal(body, &services)
	if err != nil {
		log.Print("err2", err)
	}
	services = transform(services)
	tmpl := template.Must(template.ParseFiles("template/listing.gohtml"))
	tc := make(map[string]interface{})
	tc["Services"] = services.Services

	if err := tmpl.Execute(w, tc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
