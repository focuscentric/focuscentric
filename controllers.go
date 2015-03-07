package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
)

type pageData struct {
	Title             string
	SubTitle          string
	CurrentEpisode    *Episode
	CurrentProduction *Production
	Productions       []*Production
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	prod, err := GetFeatured()
	if err != nil {
		log.Printf("error on db: %s", err)
		http.Redirect(w, r, "/error", http.StatusExpectationFailed)
		return
	}
	d := &pageData{Title: "Focus Centric - Formations video techniques", CurrentProduction: prod}
	if err := render(w, "index.html", d); err != nil {
		log.Println(err)
	}
}

func collectionsHandler(w http.ResponseWriter, r *http.Request) {
	id := getID(r.URL.Path, "/collections/")
	productions, err := GetCollection(slugToCategory(id))
	if err != nil {
		log.Printf("error on collectionHandler: %s", err)
		http.Redirect(w, r, "/error", http.StatusExpectationFailed)
		return
	}
	d := &pageData{Title: "Formations: " + id, SubTitle: id, Productions: productions}
	if err := render(w, "collections.html", d); err != nil {
		log.Println(err)
	}
}

func productionHandler(w http.ResponseWriter, r *http.Request) {
	id := getID(r.URL.Path, "/production/")
	production, err := GetProduction(-1, id)
	if err != nil {
		log.Printf("error on productionHandler: %s", err)
		http.Redirect(w, r, "/error", http.StatusExpectationFailed)
		return
	}
	d := &pageData{
		Title:             production.Title,
		SubTitle:          categoryToSlug(production.Category),
		CurrentProduction: production,
	}
	if err := render(w, "production.html", d); err != nil {
		log.Println(err)
	}
}

func episodeHandler(w http.ResponseWriter, r *http.Request) {
	slug := getID(r.URL.Path, "/episode/")
	productionID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Print("no production id supplied in query string")
		http.Redirect(w, r, "/error", http.StatusExpectationFailed)
		return
	}
	production, err := GetProduction(productionID, "")
	if err != nil {
		log.Printf("error on episodeHandler: %s", err)
		http.Redirect(w, r, "/error", http.StatusExpectationFailed)
		return
	}
	var current *Episode
	for _, e := range production.Episodes {
		if e.Slug == slug {
			current = e
			break
		}
	}
	if current == nil {
		log.Print("unable to find the current epidose")
		http.Redirect(w, r, "/error", http.StatusNotFound)
		return
	}
	d := &pageData{
		Title:             current.Title,
		CurrentEpisode:    current,
		CurrentProduction: production,
	}
	if err := render(w, "episode.html", d); err != nil {
		log.Println(err)
	}
}

func getID(url string, controller string) string {
	return url[len(controller):]
}

func slugToCategory(slug string) string {
	category := ""
	switch strings.ToUpper(slug) {
	case "JAVASCRIPT-NODEJS":
		category = "JavaScript / NodeJS"
	case "NET":
		category = ".NET"
	case "MOBILE":
		category = "iOS, Android, Windows Phone"
	case "PYTHON":
		category = "Python"
	}
	return category
}

func categoryToSlug(category string) string {
	slug := ""
	switch category {
	case "JavaScript / NodeJS":
		slug = "JAVASCRIPT-NODEJS"
	case ".NET":
		slug = "NET"
	case "iOS, Android, Windows Phone":
		slug = "MOBILE"
	case "Python":
		slug = "PYTHON"
	}
	return strings.ToLower(slug)
}
