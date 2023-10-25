package main

import (
	"database/sql"
	_ "embed"
	"fmt"
	"os"

	// tea "github.com/charmbracelet/bubbletea"
	"github.com/gocolly/colly/v2"
	_ "github.com/mattn/go-sqlite3"
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

type CharacterModel struct {
  Name string
  Class string
  level int32
  hp int32
  str int32
  mag int32
  skl int32
  spd int32
  lck int32
  def int32
  res int32
  mov int32
}

func scrape_base_stats(db *sql.DB, c *colly.Collector) {
	c.OnHTML("body", func(e *colly.HTMLElement) {
    // characters := make([]CharacterModel, 0)
    e.ForEach("tr", func(_ int, el * colly.HTMLElement) {
      fmt.Println(el)
      el.ForEach("td", func(12, ed * colly.HTMLElement) {

      })
    })
	})
	c.Visit(basestats_URL)
}

func main() {
	os.Create(db_name)
	db, err := sql.Open("sqlite3", db_name)

	defer db.Close()

	if err != nil {
		fmt.Println("Failed to Print Value")
	}

	c := colly.NewCollector()
  scrape_base_stats(db, c)
	// db.Exec(schema)
}
