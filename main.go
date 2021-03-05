package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/gocolly/colly"
)

type Movie struct {
	Title string `json:"title"`
	// Date string `json:"date"`
	// Summary string `json:"summary"`
	// Metascore string `json:"metascore"`
	// Userscore string `json:"userscore"`
}

func main() {
	listMovies := make([]Movie, 0)

	collector := colly.NewCollector(
		colly.AllowedDomains("metacritic.com", "www.metacritic.com"),
	)

	collector.OnHTML(".clamp-summary-wrap a h3", func(el *colly.HTMLElement) {
		title := el.Text
		movie := Movie{
			Title: title,
		}

		listMovies = append(listMovies, movie)
	})

	collector.OnRequest(func(req *colly.Request) {
		fmt.Println("Visiting", req.URL.String())
	})

	collector.Visit("https://www.metacritic.com/browse/movies/score/metascore/year/filtered")

	writeJSON(listMovies)
}

func writeJSON(data []Movie) {
	file, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Println("Unable to create json file")
	}
	_ = ioutil.WriteFile("metacritic_best_movies_2021.json", file, 0644)
}
