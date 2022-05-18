package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/gocolly/colly"
)

type Movie struct {
	Rank         string `json:"rank"`
	Title        string `json:"title"`
	Release_Date string `json:"release_date"`
	Summary      string `json:"summary"`
	Metascore    string `json:"metascore"`
	Userscore    string `json:"userscore"`
	Link         string `json:"link"`
}

func main() {

	c := colly.NewCollector(
		colly.AllowedDomains("metacritic.com", "www.metacritic.com"),
		colly.CacheDir("./metacritic_movies_2022"),
	)

	movies := make([]Movie, 0, 200)

	c.OnHTML("tr", func(e *colly.HTMLElement) {
		e.ForEach(".clamp-summary-wrap", func(_ int, el *colly.HTMLElement) {
			movie := Movie{}
			movie.Rank = el.ChildText("span.numbered")
			movie.Title = el.ChildText("a.title h3")
			movie.Release_Date = el.ChildText(".clamp-details span:first-of-type")
			movie.Summary = el.ChildText(".summary")
			movie.Metascore = el.ChildText(".clamp-metascore a div")
			movie.Userscore = el.ChildText(".clamp-userscore a div")
			movie.Link = "https://www.metacritic.com" + el.ChildAttr("a.title", "href")
			movies = append(movies, movie)
		})
	})

	c.OnRequest(func(req *colly.Request) {
		fmt.Println("Visiting", req.URL.String())
	})

	c.Visit("https://www.metacritic.com/browse/movies/score/metascore/year/filtered")

	writeJSON(movies)
	writeCSV(movies)
}

func writeJSON(data []Movie) {
	file, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Println("Unable to create json file")
	}
	cts := time.Now().Format("2006-Jan-02")
	fn := "metacritic_best_movies_" + cts + ".json"
	_ = ioutil.WriteFile(fn, file, 0644)
}

func writeCSV(data []Movie) {
	cts := time.Now().Format("2006-Jan-02")
	fn := "metacritic_best_movies_2021_" + cts + ".csv"
	csvFile, err := os.Create(fn)
	if err != nil {
		fmt.Println(err)
	}
	defer csvFile.Close()

	w := csv.NewWriter(csvFile)
	var header []string
	header = append(header, "Rank")
	header = append(header, "Title")
	header = append(header, "Release_Date")
	header = append(header, "Metascore")
	header = append(header, "Userscore")
	header = append(header, "Summary")
	header = append(header, "Link")
	w.Write(header)

	for _, mv := range data {
		var row []string
		row = append(row, mv.Rank)
		row = append(row, mv.Title)
		row = append(row, mv.Release_Date)
		row = append(row, mv.Metascore)
		row = append(row, mv.Userscore)
		row = append(row, mv.Summary)
		row = append(row, mv.Link)
		w.Write(row)
	}

	w.Flush()
}
