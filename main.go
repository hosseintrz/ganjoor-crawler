package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/gocql/gocql"
	"github.com/hosseintrz/ganjoor_crawler/db"
	"github.com/hosseintrz/ganjoor_crawler/model"
	"github.com/hosseintrz/ganjoor_crawler/persistence"
)

func main() {
	err := db.InitDatabase()
	if err != nil {
		log.Fatalf("error connecting to database: %v\n", err)
	}

	cql, err := db.GetDB()
	if err != nil {
		log.Fatalf("error getting cassandra instance: %v\n", err)
	}
	defer cql.Close()

	err = collectPoems(cql)
	if err != nil {
		log.Fatalf("error collecting poems: %v\n", err)
	}

	// err = collectPoets()
	// if err != nil{
	// 	log.Fatalf("error collecting poets: %v\n", err)
	// }
}

func collectPoems(cql *gocql.Session) error {
	c := colly.NewCollector()

	c.OnHTML(".poem", func(e *colly.HTMLElement) {
		fmt.Println("found .poem elemnet")

		poetName := e.ChildTexts("#garticle #page-hierarchy a")[0]
		title := e.ChildText("#garticle #page-hierarchy h2 a")

		var lines []string
		e.DOM.Find("#garticle .b").Each(func(_ int, s *goquery.Selection) {
			line := s.Text()
			lines = append(lines, line)
		})

		err := validateData(poetName, title, lines)
		if err != nil {
			log.Printf("error validating data: %v\n", err)
			return
		}

		poet, err := persistence.GetPoetByName(cql, poetName)
		if err != nil {
			poet = &model.Poet{
				ID:   gocql.TimeUUID(),
				Name: poetName,
			}
			err = persistence.InsertPoet(cql, poet)
			if err != nil {
				log.Printf("error inserting poet: %v\n", err)
				return
			}
		}

		content := formatPoem(lines)

		poem := model.Poem{
			ID:       gocql.TimeUUID(),
			Title:    title,
			Content:  content,
			PoetID:   poet.ID,
			PoetName: poet.Name,
		}

		err = persistence.InsertPoem(cql, poem)
		if err != nil {
			log.Printf("error inserting poem: %v\n", err)
			return
		}

		fmt.Printf("Poet: %s\nTitle: %s\nContent:\n%v\n\n", poem.PoetName, poem.Title, poem.Content)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	fmt.Println("start visiting")
	err := c.Visit("https://ganjoor.net/rahi/ghazalha4/sh1")
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func formatPoem(lines []string) string {
	result := ""
	for _, line := range lines {
		formattedLine := strings.Replace(line, "\n", "    ----    ", 1)
		result += formattedLine + "\n"
	}
	return result
}

func validateData(poetName, title string, lines []string) error {
	valid := len(poetName) > 0 && len(title) > 0
	if !valid {
		return fmt.Errorf("invalid data")
	}

	return nil
}
