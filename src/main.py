import re
import sqlite3
import requests

from http.server import HTTPServer, BaseHTTPRequestHandler
from bs4 import BeautifulSoup

# TODO: Add colorful prompts to make it a proper CLI?

basestats_URL = "https://serenesforest.net/awakening/characters/base-stats/main-story/"
basegrowths_URL = "https://serenesforest.net/awakening/characters/growth-rates/base/"
class_URL = "https://serenesforest.net/awakening/characters/class-sets/"
class_base_URL = "https://serenesforest.net/awakening/classes/base-stats/"
skills_URL = "https://serenesforest.net/awakening/miscellaneous/skills/"
character_assets_URL = "https://serenesforest.net/awakening/characters/maximum-stats/modifiers/"


def schema(cur: sqlite3.Cursor, con: sqlite3.Connection):
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
                res TEXT,
                affinity INT,
                growth INT
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

            CREATE TABLE IF NOT EXISTS
            skills(
                skill TEXT,
                effect TEXT,
                activation TEXT,
                class TEXT,
                level INTEGER
            );

            CREATE TABLE IF NOT EXISTS
            character_assets(
                name TEXT,
                str INTEGER,
                mag INTEGER,
                skl INTEGER,
                spd INTEGER,
                lck INTEGER,
                def INTEGER,
                res INTEGER
            );
        """
    )


def base_stats(cur: sqlite3.Cursor, con: sqlite3.Connection):
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
    con.commit()


def split_affinity(c: list[str]) -> tuple[list[str], list[str]]:
    asset = []
    flaw = []
    p1 = re.compile(r"\+\/-(\d{1,2})")
    p2 = re.compile(r"\+?(\d{1,2})\/(-\d{1,2})")
    for string in c:
        m1 = p1.match(string)
        m2 = p2.match(string)
        if (m1 and m1.lastindex == 1):
            v = m1.group(1)
            asset.append(v)
            flaw.append(v)
        elif (m2 and m2.lastindex == 2):
            a, f = m2.groups()
            asset.append(a)
            flaw.append(f)
        else:
            asset.append(string)
            flaw.append(string)
    return (asset, flaw)


# NOTE: Types kind of don't matter because of the way SQLite works
def base_growths(cur: sqlite3.Cursor, con: sqlite3.Connection):
    basegrowths_page = requests.get(basegrowths_URL)
    soup = BeautifulSoup(basegrowths_page.content, "html.parser")
    b = soup.find(string="Magic").find_parent("table").find_all("tr")
    for element in b:
        c = [x.get_text() for x in element.find_all("td")]
        if (len(c) == 9):
            d, e = split_affinity(c)
            cur.execute(
                """
                    INSERT INTO asset VALUES(
                        ?, ?, ?,
                        ?, ?, ?,
                        ?, ?, ?,
                        1, 1
                    )
                """,
                d
            )
            cur.execute(
                """
                    INSERT INTO asset VALUES(
                        ?, ?, ?,
                        ?, ?, ?,
                        ?, ?, ?,
                        0, 1
                    )
                """,
                e
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


def class_sets(cur: sqlite3.Cursor, con: sqlite3.Connection):
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


def class_base(cur: sqlite3.Cursor, con: sqlite3.Connection):
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


def class_skills(cur: sqlite3.Cursor, con: sqlite3.Connection):
    page = requests.get(skills_URL)
    soup = BeautifulSoup(page.content, "html.parser")
    classes = soup.find("table").find_all("tr")
    for element in classes:
        a = [x.get_text() for x in element.find_all("td")]
        if (len(a) == 5):
            cur.execute(
                """
                    INSERT INTO skills
                    VALUES(
                        ?, ?, ?, ?, ?
                    )
                """,
                a
            )
    con.commit()


# NOTE: Types kind of don't matter because of the way SQLite works
def character_assets(cur: sqlite3.Cursor, con: sqlite3.Connection):
    page = requests.get(character_assets_URL)
    soup = BeautifulSoup(page.content, "html.parser")
    classes = soup.find(string="Magic").find_parent("table").find_all("tr")
    for element in classes:
        c = [x.get_text().replace(" ", "") for x in element.find_all("td")]
        if (len(c) == 8):
            d, e = split_affinity(c)
            cur.execute(
                """
                    INSERT INTO asset VALUES(
                        ?, 0, ?,
                        ?, ?, ?,
                        ?, ?, ?,
                        1, 0
                    )
                """,
                d
            )
            cur.execute(
                """
                    INSERT INTO asset VALUES(
                        ?, 0, ?,
                        ?, ?, ?,
                        ?, ?, ?,
                        0, 0
                    )
                """,
                e
            )

    parent_assets = soup.find(string="Chrom").find_parent(
        "table").find_all("tr")
    for element in parent_assets:
        a = [x.get_text() for x in element.find_all("td")]
        if (len(a) == 8):
            cur.execute(
                """
                    INSERT INTO character_assets
                    VALUES(
                        ?, ?, ?, ?,
                        ?, ?, ?, ?
                    )
                """,
                a
            )
    con.commit()


html = """
<!DOCTYPE html>
<html>
  <head>
    <title> Awakening DB </title>
  </head>
  <body>
    <h1>Awakening DB</h1>
    <form method="POST">
      <label for="basestats"> Base Stats </label>
      <input type="checkbox" name="basestats" id="basestats">
      <br>
      <label for="basegrowths"> Base Growths </label>
      <input type="checkbox" name="basegrowths" id="basegrowths">
      <br>
      <label for="classsets"> Class Sets </label>
      <input type="checkbox" name="classsets" id="classsets">
      <br>
      <label for="classbase"> Class Base </label>
      <input type="checkbox" name="classbase" id="classbase">
      <br>
      <label for="classskills"> Class Skills </label>
      <input type="checkbox" name="classskills" id="classskills">
      <br>
      <label for="characterassets"> Character Assets </label>
      <input type="checkbox" name="characterassets" id="characterassets">
      <br>
      <button type="submit"> Submit </button>
    </form>
  </body>
</html>
"""


def activateMethods(b: list[str]):
    con = sqlite3.connect("awakening.db")
    cur = con.cursor()

    schema(cur, con)
    for line in b:
        if 'basestats' in line and 'on' in line:
            base_stats(cur, con)
        elif 'basegrowths' in line and 'on' in line:
            base_growths(cur, con)
        elif 'classsets' in line and 'on' in line:
            class_sets(cur, con)
        elif 'classbase' in line and 'on' in line:
            class_base(cur, con)
        elif 'classskills' in line and 'on' in line:
            class_skills(cur, con)
        elif 'characterassets' in line and 'on' in line:
            character_assets(cur, con)


class Server(BaseHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.send_header('Content-type', 'text/html')
        self.end_headers()

        self.wfile.write(bytes(html, "UTF-8"))
        return

    def do_POST(self):
        content_length = int(self.headers['Content-Length'])
        post_data_bytes = self.rfile.read(content_length)
        formdata: list[str] = str(post_data_bytes, 'UTF-8').split('&')
        try:
            activateMethods(formdata)
            self.send_response(201)
            self.send_header()
            self.end_headers()
            return
        except sqlite3.Error:
            self.send_response(500)
            self.send_header()
            self.end_headers()
            return


if __name__ == "__main__":
    httpd = HTTPServer(('localhost', 8080), Server)
    print("Go to http://localhost:8080")
    print("Press Ctrl+C to close out the server")
    try:
        httpd.serve_forever()
    except KeyboardInterrupt:
        pass
    httpd.server_close()