from bs4 import BeautifulSoup
import re
import requests
import sqlite3

basestats_URL = "https://serenesforest.net/awakening/characters/base-stats/main-story/"
basegrowths_URL = "https://serenesforest.net/awakening/characters/growth-rates/base/"
class_URL = "https://serenesforest.net/awakening/characters/class-sets/"
class_base_URL = "https://serenesforest.net/awakening/classes/base-stats/"

con = sqlite3.connect("awakening.db")
cur = con.cursor()


def schema():
    # TODO: Refactor into one string?
    cur.executescript(
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
            );

            CREATE TABLE IF NOT EXISTS
            basegrowths(
                name TEXT,
                hp INTEGER,
                str INTEGER,
                mag INTEGER,
                skl INTEGER,
                spd INTEGER,
                lck INTEGER,
                def INTEGER,
                res INTEGER
            );

            CREATE TABLE IF NOT EXISTS
            asset(
                asset TEXT,
                hp TEXT,
                str TEXT,
                mag TEXT,
                skl TEXT,
                spd TEXT,
                lck TEXT,
                def TEXT,
                res TEXT
            );

            CREATE TABLE IF NOT EXISTS
            classes(
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                name TEXT
            );

            CREATE TABLE IF NOT EXISTS
            classbase(
                class TEXT,
                hp INTEGER,
                str INTEGER,
                mag INTEGER,
                skl INTEGER,
                spd INTEGER,
                def INTEGER,
                res INTEGER,
                mov INTEGER,
                rank TEXT
            );
        """
    )


def base_stats():
    basestats_page = requests.get(basestats_URL)
    soup = BeautifulSoup(basestats_page.content, "html.parser")
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


def base_growths():
    basegrowths_page = requests.get(basegrowths_URL)
    soup = BeautifulSoup(basegrowths_page.content, "html.parser")
    b = soup.find(string="Asset/Flaw").find_parent("table").find_all("tr")
    for element in b:
        c = [x.get_text() for x in element.find_all("td")]
        if (len(c) == 9):
            cur.execute(
                """
                    INSERT INTO asset VALUES(
                        ?, ?, ?,
                        ?, ?, ?,
                        ?, ?, ?
                    )
                """,
                c
            )

    d = soup.find(string="Chrom").find_parent("table").find_all("tr")
    for element in d:
        c = [x.get_text() for x in element.find_all("td")]
        if (len(c) == 9):
            cur.execute(
                """
                    INSERT INTO basegrowths VALUES(
                        ?, ?, ?,
                        ?, ?, ?,
                        ?, ?, ?
                    )
                """,
                c
            )

    e = soup.find(string="Lucina").find_parent("table").find_all("tr")
    for element in e:
        c = [x.get_text() for x in element.find_all("td")]
        if (len(c) == 9):
            cur.execute(
                """
                    INSERT INTO basegrowths VALUES(
                        ?, ?, ?,
                        ?, ?, ?,
                        ?, ?, ?
                    )
                """,
                c
            )
    con.commit()


def class_sets():
    page = requests.get(class_URL)
    soup = BeautifulSoup(page.content, "html.parser")
    classes = soup.find(string="Regular classes").find_parent(
        "p").get_text().split(':')[1]
    a = [re.sub('\(.+\)', '', c).strip() for c in classes.split(',')]
    for element in a:
        # PERF: Refactor into executemany?
        cur.execute(
            """
                INSERT INTO classes(name) VALUES(?)
            """,
            [element]
        )
    con.commit()


def class_base():
    page = requests.get(class_base_URL)
    soup = BeautifulSoup(page.content, "html.parser")
    classes = soup.find_all("tr")
    for element in classes:
        a = [x.get_text() for x in element.find_all("td")]
        if (len(a) == 10):
            cur.execute(
                """
                    INSERT INTO classbase
                    VALUES(
                        ?, ?, ?, ?, ?,
                        ?, ?, ?, ?, ?
                    )
                """,
                a
            )

    con.commit()


schema()
# base_stats()
# base_growths()
# class_sets()
# class_base()
