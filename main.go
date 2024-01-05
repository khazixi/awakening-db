package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"database/sql"
	"strings"

	_ "embed"

	// tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/gocolly/colly/v2"
	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schema string

const db_name = "awakening.db"
const basestats_URL = "https://serenesforest.net/awakening/characters/base-stats/main-story/"
const basegrowths_URL = "https://serenesforest.net/awakening/characters/growth-rates/base/"
const class_URL = "https://serenesforest.net/awakening/characters/class-sets/"
const class_base_URL = "https://serenesforest.net/awakening/classes/base-stats/"
const skills_URL = "https://serenesforest.net/awakening/miscellaneous/skills/"
const character_assets_URL = "https://serenesforest.net/awakening/characters/maximum-stats/modifiers/"

// WARNING: Does not scrape difficulties well.
// WARNING: Will need if checks to succssfully parse the data
func scrape_base_stats(wg *sync.WaitGroup, db *sql.DB) {
	c := colly.NewCollector()

	c.OnHTML("body", func(h *colly.HTMLElement) {
		h.ForEach("tr", func(i int, h *colly.HTMLElement) {
			row := make([]any, 0)
			h.ForEach("td", func(i int, h *colly.HTMLElement) {
				if i < 14 {
					row = append(row, strings.TrimSpace(h.Text))
				}
			})

			if len(row) == 14 {
				fmt.Println(row, len(row))
				_, err := db.Exec(
					`INSERT INTO basestats VALUES(
        ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
        )`, row...)

				if err != nil {
					log.Println(err)
				}
			}
		})
	})

	c.Visit(basestats_URL)

	wg.Done()
}

func scrape_growth_rates(wg *sync.WaitGroup, db *sql.DB) {
	c := colly.NewCollector()

	c.OnHTML("body", func(h *colly.HTMLElement) {
		h.ForEach("table", func(i int, h *colly.HTMLElement) {
			if h.Index == 3 || h.Index == 4 {
				h.ForEach("tr", func(i int, h *colly.HTMLElement) {
					row := make([]any, 0, 9)
					h.ForEach("td", func(i int, h *colly.HTMLElement) {
						row = append(row, h.Text)
					})

					if len(row) == 9 {
						fmt.Println(row, len(row))
						_, err := db.Exec(
							`INSERT INTO basegrowths VALUES(
                ?, ?, ?, ?, ?, ?, ?, ?, ?
              )`, row...)

						if err != nil {
							log.Println(err)
						}
					}
				})
			}

			if h.Index == 1 {
				h.ForEach("tr", func(i int, h *colly.HTMLElement) {
					row := make([]string, 0)
					h.ForEach("td", func(i int, h *colly.HTMLElement) {
						row = append(row, h.Text)
					})
					fmt.Println(row)
				})
			}

		})
	})

	c.Visit(basegrowths_URL)

	wg.Done()
}

func scrape_class_sets(wg *sync.WaitGroup, db *sql.DB) {
	c := colly.NewCollector()

	c.OnHTML("body", func(h *colly.HTMLElement) {
		h.ForEach("table", func(i int, h *colly.HTMLElement) {
			if i < 2 {
				h.ForEach("tr", func(i int, h *colly.HTMLElement) {
					row := make([]any, 0, 5)
					h.ForEach("td", func(i int, h *colly.HTMLElement) {
						if h.Text == "–" {
							return
						}
						row = append(row, h.Text)
					})
					if len(row) == 3 {
						fmt.Println(row, len(row))
						_, err := db.Exec(
							`INSERT INTO characterclasses VALUES(
                ?, ?, ?, ?, ?
              )`, row[0], row[1], row[2], nil, nil)

						if err != nil {
							log.Println(err)
						}
					} else if len(row) == 4 {
						fmt.Println(row, len(row))
						_, err := db.Exec(
							`INSERT INTO characterclasses VALUES(
                ?, ?, ?, ?, ?
              )`, row[0], row[1], row[2], row[3], nil)

						if err != nil {
							log.Println(err)
						}
					} else if len(row) == 5 {
						fmt.Println(row, len(row))
						_, err := db.Exec(
							`INSERT INTO characterclasses VALUES(
                ?, ?, ?, ?, ?
              )`, row[0], row[1], row[2], row[3], row[4])

						if err != nil {
							log.Println(err)
						}
					}

				})
			}
		})
	})

	c.Visit(class_URL)

	wg.Done()
}

func scrape_base_class(wg *sync.WaitGroup, db *sql.DB) {
	c := colly.NewCollector()

	c.OnHTML("tr", func(h *colly.HTMLElement) {
		row := make([]any, 0)
		h.ForEach("td", func(i int, h *colly.HTMLElement) {
			row = append(row, h.Text)
		})
		fmt.Println(row)
		if len(row) == 10 {
			fmt.Println(row, len(row))
			_, err := db.Exec(
				`INSERT INTO classbase VALUES(
        ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
        )`, row...)

			if err != nil {
				log.Println(err)
			}
		}
	})

	c.Visit(class_base_URL)

	wg.Done()
}

func scrape_skills(wg *sync.WaitGroup, db *sql.DB) {
	c := colly.NewCollector()

	c.OnHTML("tr", func(h *colly.HTMLElement) {
		row := make([]any, 0)
		h.ForEach("td", func(i int, h *colly.HTMLElement) {
			if i == 0 {
				return
			}
			row = append(row, h.Text)
		})

		if len(row) == 5 {
			fmt.Println(row, len(row))
			_, err := db.Exec(
				`INSERT INTO skills VALUES(
        ?, ?, ?, ?, ?
        )`, row...)

			if err != nil {
				log.Println(err)
			}
		}
	})

	c.Visit(skills_URL)

	wg.Done()
}

func scrape_char_assets(wg *sync.WaitGroup, db *sql.DB) {
	c := colly.NewCollector()

	c.OnHTML("body", func(h *colly.HTMLElement) {
		h.ForEach("table", func(i int, h *colly.HTMLElement) {
			if i == 1 {
				h.ForEach("tr", func(i int, h *colly.HTMLElement) {
					row := make([]string, 0)
					h.ForEach("td", func(i int, h *colly.HTMLElement) {
						row = append(row, h.Text)
					})
					fmt.Println(row)
				})
			} else if i == 3 {
				h.ForEach("tr", func(i int, h *colly.HTMLElement) {
					row := make([]string, 0)
					h.ForEach("td", func(i int, h *colly.HTMLElement) {
						row = append(row, h.Text)
					})
					fmt.Println(row)
				})
			} else if i > 4 {
				h.ForEach("tr", func(i int, h *colly.HTMLElement) {
					row := make([]string, 0)
					h.ForEach("td", func(i int, h *colly.HTMLElement) {
						row = append(row, h.Text)
					})
					fmt.Println(row)
				})
			}
		})
	})

	c.Visit(character_assets_URL)

	wg.Done()
}

func main() {
	var options []string
	var dbName string = "awakening"

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Scrape Options").
				Options(
					huh.NewOption("Include Base Stats", "basestats"),
					huh.NewOption("Include Base Growths", "basegrowths"),
					huh.NewOption("Include Class Sets", "classets"),
					huh.NewOption("Include Class Base Stats", "classbase"),
					huh.NewOption("Include Character Assets", "charassets"),
					huh.NewOption("Include Character Skills", "charskills"),
				).
				Value(&options),

			huh.NewInput().
				Title("Enter the name you want to save the database to:").
				Prompt("> ").
				Value(&dbName),
		),
	)

	err := form.Run()

	if err != nil {
		log.Fatal(err)
	}

	wg := sync.WaitGroup{}

	wg.Add(len(options))

	fmt.Println(dbName)

	if len(options) == 0 {
		fmt.Println("Decided not to scrape anything")
	} else {
		file, err := os.Create(dbName + ".db")
		if err != nil {
			log.Fatal(err)
		}
		file.Close()

		db, err := sql.Open("sqlite", dbName+".db")
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec(schema)
		if err != nil {
			log.Fatal(err)
		}

		for _, option := range options {
			switch option {
			case "classbase":
				go scrape_base_class(&wg, db)
			case "classets":
				go scrape_class_sets(&wg, db)
			case "basestats":
				go scrape_base_stats(&wg, db)
			case "basegrowths":
				go scrape_growth_rates(&wg, db)
			case "charassets":
				go scrape_char_assets(&wg, db)
			case "charskills":
				go scrape_skills(&wg, db)
			default:
				fmt.Println("Unknown Option", option)
			}
		}
	}

	wg.Wait()

	// for v := range work {
	// 	fmt.Println(v)
	// }
}
