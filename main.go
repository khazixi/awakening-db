package main

import (
	"fmt"
	"os"
	"strings"

	// tea "github.com/charmbracelet/bubbletea"
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
func scrape_base_stats() {
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
}

func scrape_growth_rates() {
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
}

func scrape_class_sets() {
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
}

func scrape_base_class() {
  c := colly.NewCollector()

  c.OnHTML("tr", func(h *colly.HTMLElement) {
    row := make([]string, 0)
    h.ForEach("td", func(i int, h *colly.HTMLElement) {
      row = append(row, h.Text)
    })
    fmt.Println(row)
  })

  c.Visit(class_base_URL)
}

func scrape_skills() {
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
}

func scrape_char_assets() {
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
}

func main() {
	os.Create(db_name)
	// db, err := sql.Open("sqlite3", db_name)
	//
	// defer db.Close()

	// if err != nil {
	// 	fmt.Println("Failed to Print Value")
	// }

	// go scrape_base_stats()
	// go scrape_growth_rates()
  // go scrape_class_sets()
  // go scrape_base_class()
  // go scrape_skills()
  scrape_char_assets()

	// db.Exec(schema)
}
