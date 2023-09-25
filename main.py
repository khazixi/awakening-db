from bs4 import BeautifulSoup
import requests
import sqlite3
from pprint import pprint

basestats_URL = "https://serenesforest.net/awakening/characters/base-stats/main-story/"

basestats_page = requests.get(basestats_URL)

soup = BeautifulSoup(basestats_page.content, "html.parser")
con = sqlite3.connect("awakening.db")
cur = con.cursor()

def sql():
    cur.execute(
        """
            CREATE TABLE IF NOT EXISTS
            basestats(
                name TEXT,
                class TEXT,
                level INTEGER,
                hp INTEGER,
                str INTEGER,
                mag INTEGER,
                skl INTEGER,
                spd INTEGER,
                lck INTEGER,
                def INTEGER,
                res INTEGER,
                mov INTEGER
            )
        """
    )

def bases():
    for element in soup.find_all('tr'):
        b = [x.get_text() for x in element.find_all('td', limit=12)]
        if (b and len(b) == 12):
            data = cur.execute(
                """
                    INSERT INTO basestats VALUES(
                        ?, ?, ?,
                        ?, ?, ?,
                        ?, ?, ?,
                        ?, ?, ?
                    )
                """,
                b
            )
            print(data.description)
            for d in data:
                print(d)
        else:
            print('Fucked Up')
    con.commit()

sql()
bases()
