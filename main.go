package main

import (
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

type DBMsg struct {
	command string
	data    []any
}

func manager(wg *sync.WaitGroup, dbch chan DBMsg) {
	wg.Wait()
	close(dbch)
}

// WARNING: Does not scrape difficulties well.
// WARNING: Will need if checks to succssfully parse the data
func scrape_base_stats(wg *sync.WaitGroup, ch chan DBMsg) {
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
				ch <- DBMsg{
					command: `INSERT INTO basestats VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
					data:    row,
				}
			}
		})
	})

	c.Visit(basestats_URL)

	wg.Done()
}

func scrape_growth_rates(wg *sync.WaitGroup, dbch chan DBMsg) {
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
						dbch <- DBMsg{
							command: `INSERT INTO basegrowths VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)`,
							data:    row,
						}
					}
				})
			}

			if h.Index == 1 {
				h.ForEach("tr", func(i int, h *colly.HTMLElement) {
					row := make([]any, 0)
					h.ForEach("td", func(i int, h *colly.HTMLElement) {
						row = append(row, h.Text)
					})

					if len(row) == 9 {
						dbch <- DBMsg{
							command: `INSERT INTO growth_rate_modifiers VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)`,
							data:    row,
						}
					}
				})
			}

		})
	})

	c.Visit(basegrowths_URL)

	wg.Done()
}

func scrape_class_sets(wg *sync.WaitGroup, dbch chan DBMsg) {
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
						dbch <- DBMsg{
							command: `INSERT INTO characterclasses VALUES(?, ?, ?, null, null)`,
							data:    row,
						}

					} else if len(row) == 4 {

						dbch <- DBMsg{
							command: `INSERT INTO characterclasses VALUES(?, ?, ?, ?, null)`,
							data:    row,
						}
					} else if len(row) == 5 {
						dbch <- DBMsg{
							command: `INSERT INTO characterclasses VALUES(?, ?, ?, ?, ?)`,
							data:    row,
						}
					}
				})
			}
		})
	})

	c.Visit(class_URL)

	wg.Done()
}

func scrape_base_class(wg *sync.WaitGroup, dbch chan DBMsg) {
	c := colly.NewCollector()

	c.OnHTML("tr", func(h *colly.HTMLElement) {
		row := make([]any, 0)
		h.ForEach("td", func(i int, h *colly.HTMLElement) {
			row = append(row, h.Text)
		})
		if len(row) == 10 {
			dbch <- DBMsg{
				command: `INSERT INTO classbase VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				data:    row,
			}
		}
	})

	c.Visit(class_base_URL)

	wg.Done()
}

func scrape_skills(wg *sync.WaitGroup, dbch chan DBMsg) {
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
			dbch <- DBMsg{
				command: `INSERT INTO skills VALUES(?, ?, ?, ?, ?)`,
				data:    row,
			}
		}
	})

	c.Visit(skills_URL)

	wg.Done()
}

// TODO: Implement the SQL for this function
func scrape_char_assets(wg *sync.WaitGroup, dbch chan DBMsg) {
	c := colly.NewCollector()

	c.OnHTML("body", func(h *colly.HTMLElement) {
		h.ForEach("table", func(i int, h *colly.HTMLElement) {
			if i == 1 {
				h.ForEach("tr", func(i int, h *colly.HTMLElement) {
					row := make([]string, 0)
					h.ForEach("td", func(i int, h *colly.HTMLElement) {
						row = append(row, h.Text)
					})
				})
			} else if i == 3 {
				h.ForEach("tr", func(i int, h *colly.HTMLElement) {
					row := make([]string, 0)
					h.ForEach("td", func(i int, h *colly.HTMLElement) {
						row = append(row, h.Text)
					})
				})
			} else if i > 4 {
				h.ForEach("tr", func(i int, h *colly.HTMLElement) {
					row := make([]string, 0)
					h.ForEach("td", func(i int, h *colly.HTMLElement) {
						row = append(row, h.Text)
					})
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

	dbch := make(chan DBMsg)

	wg := sync.WaitGroup{}

	wg.Add(len(options))

	if len(options) == 0 {
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

		go manager(&wg, dbch)

		for _, option := range options {
			switch option {
			case "classbase":
				go scrape_base_class(&wg, dbch)
			case "classets":
				go scrape_class_sets(&wg, dbch)
			case "basestats":
				go scrape_base_stats(&wg, dbch)
			case "basegrowths":
				go scrape_growth_rates(&wg, dbch)
			case "charassets":
				go scrape_char_assets(&wg, dbch)
			case "charskills":
				go scrape_skills(&wg, dbch)
			}
		}

		for dbmsg := range dbch {
			_, err := db.Exec(dbmsg.command, dbmsg.data...)
			if err != nil {
				log.Println(err)
			}
		}

	}
}
