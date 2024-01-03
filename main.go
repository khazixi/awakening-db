package main

import (
	"fmt"
	"log"
	// "os"
	"strings"

	// tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/gocolly/colly/v2"
)

var schema string

const db_name = "awakening.db"
const basestats_URL = "https://serenesforest.net/awakening/characters/base-stats/main-story/"
const basegrowths_URL = "https://serenesforest.net/awakening/characters/growth-rates/base/"
const class_URL = "https://serenesforest.net/awakening/characters/class-sets/"
const class_base_URL = "https://serenesforest.net/awakening/classes/base-stats/"
const skills_URL = "https://serenesforest.net/awakening/miscellaneous/skills/"
const character_assets_URL = "https://serenesforest.net/awakening/characters/maximum-stats/modifiers/"

type CharacterModel struct {
	Name  string
	Class string
	level int32
	hp    int32
	str   int32
	mag   int32
	skl   int32
	spd   int32
	lck   int32
	def   int32
	res   int32
	mov   int32
}

// WARNING: Does not scrape difficulties well.
// WARNING: Will need if checks to succssfully parse the data
func scrape_base_stats(work chan struct{}) {
	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println(r.URL.Hostname())
	})

	c.OnHTML("body", func(h *colly.HTMLElement) {
		chars := make([][]string, 30)
		h.ForEach("tr", func(i int, h *colly.HTMLElement) {
			row := make([]string, 0)
			h.ForEach("td", func(i int, h *colly.HTMLElement) {
				if i < 14 {
					row = append(row, strings.TrimSpace(h.Text))
				}
			})
			if len(row) != 0 {
				fmt.Println(row, len(row))
			}
			chars = append(chars, row)
		})
	})

	c.Visit(basestats_URL)
  work <- struct{}{}
}

func scrape_growth_rates(work chan struct{}) {
	c := colly.NewCollector()

	c.OnHTML("body", func(h *colly.HTMLElement) {
		h.ForEach("table", func(i int, h *colly.HTMLElement) {
			if h.Index == 3 || h.Index == 4 {
				h.ForEach("tr", func(i int, h *colly.HTMLElement) {
					row := make([]string, 0)
					h.ForEach("td", func(i int, h *colly.HTMLElement) {
						row = append(row, h.Text)
					})
					fmt.Println(row)
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
  work <- struct{}{}
}

func scrape_class_sets(work chan struct{}) {
	c := colly.NewCollector()

	c.OnHTML("body", func(h *colly.HTMLElement) {
		h.ForEach("table", func(i int, h *colly.HTMLElement) {
			if i < 2 {
				h.ForEach("tr", func(i int, h *colly.HTMLElement) {
					row := make([]string, 0)
					h.ForEach("td", func(i int, h *colly.HTMLElement) {
						if h.Text == "–" {
							return
						}
						row = append(row, h.Text)
					})
					fmt.Println(row)
				})
			}
		})
	})

	c.Visit(class_URL)
  work <- struct{}{}
}

func scrape_base_class(work chan struct{}) {
	c := colly.NewCollector()

	c.OnHTML("tr", func(h *colly.HTMLElement) {
		row := make([]string, 0)
		h.ForEach("td", func(i int, h *colly.HTMLElement) {
			row = append(row, h.Text)
		})
		fmt.Println(row)
	})

	c.Visit(class_base_URL)
  work <- struct{}{}
}

func scrape_skills(work chan struct{}) {
	c := colly.NewCollector()

	c.OnHTML("tr", func(h *colly.HTMLElement) {
		row := make([]string, 0)
		h.ForEach("td", func(i int, h *colly.HTMLElement) {
			if i == 0 {
				return
			}
			row = append(row, h.Text)
		})
		fmt.Println(row)
	})

	c.Visit(skills_URL)
  work <- struct{}{}
}

func scrape_char_assets(work chan struct{}) {
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
  work <- struct{}{}
}

func main() {
	options := make([]string, 0)
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
		),
	)

	work := make(chan struct{})

	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}

	if len(options) == 0 {
		fmt.Println("Decided not to scrape anything")
	} else {
		for _, option := range options {

			switch option {
			case "classbase":
				go scrape_base_class(work)
			case "classets":
				go scrape_class_sets(work)
			case "basestats":
				go scrape_base_stats(work)
			case "basegrowths":
				go scrape_growth_rates(work)
      case "charassets":
        go scrape_char_assets(work)
      case "charskills":
        go scrape_skills(work)
			default:
				fmt.Println("Unknown Option", option)
			}
		}
	}

	for v := range work {
		fmt.Println(v)
	}
}
