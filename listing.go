package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"google.golang.org/api/appengine/v1"
)

func main() {
	http.HandleFunc("/api/services", getServiceListingAPI)
	http.HandleFunc("/", getServiceListingHTML)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func getServiceListing(w http.ResponseWriter, r *http.Request) (appEngineServicesStruct, error) {
	appName := os.Getenv("GOOGLE_CLOUD_PROJECT")
	log.Println("appName:"+appName)
	// add CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	ctx := context.Background()
	appengineService, err := appengine.NewService(ctx)
	var appEngineServices appEngineServicesStruct
	response, err := appengine.NewAppsServicesService(appengineService).List(appName).Do()
	if err != nil {
		log.Print("err calling NewAppsServicesService.List", err)
		return appEngineServices, err
	}
	log.Println("Loading JSON")
	appEngineServices = transform(response.Services)
	return appEngineServices, nil
}

func getServiceListingAPI(w http.ResponseWriter, r *http.Request) {
	appEngineServices, err := getServiceListing(w, r)
	if err != nil {
		return
	}
	data, err := json.Marshal(appEngineServices)
	if err != nil {
		log.Print("err mashal", err)
	}
	fmt.Fprintf(w, "%s", data)
}

func getServiceListingHTML(w http.ResponseWriter, r *http.Request) {
	appEngineServices, err := getServiceListing(w, r)
	if err != nil {
		return
	}
	transformAndDisplay(appEngineServices, w)
}

type appEngineServicesStruct struct {
	Services []appEngineServiceStruct `json:"services"`
}

type appEngineServiceStruct struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	URL   string `json:"url"`
	Title string `json:"title"`
}

func formatName(name string) string {
	// Function replacing words (assuming lower case input)
	replace := func(word string) string {
		switch word {
		case "zen-" /*, "formation-"*/ :
			return ""
		}

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

func transform(services []*appengine.Service) appEngineServicesStruct {
	appEngineServices := appEngineServicesStruct{}
	for index := range services {
		appEngineService := appEngineServiceStruct{
			ID: services[index].Id,
			Name: services[index].Name,
			URL: fmt.Sprintf("https://%s-dot-zen-formations.appspot.com/", services[index].Id),
			Title: formatName(services[index].Id),
		}
		appEngineServices.Services = append(appEngineServices.Services, appEngineService)
	}
	return appEngineServices
}

func transformAndDisplay(appEngineServices appEngineServicesStruct, w http.ResponseWriter) {
	tmpl := template.Must(template.ParseFiles("template/listing.gohtml"))
	tc := make(map[string]interface{})
	tc["Services"] = appEngineServices.Services
	if err := tmpl.Execute(w, tc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
