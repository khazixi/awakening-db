package main

import (
	"database/sql"
	_ "embed"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
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


func main() {
  os.Create(db_name)
  db, err := sql.Open("sqlite3",db_name)

  defer db.Close()

  if err != nil {
    fmt.Println("Failed to Print Value")
  }

  db.Exec(schema)
}
