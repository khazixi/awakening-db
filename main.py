from bs4 import BeautifulSoup, BeautifulStoneSoup
import requests
from pprint import pprint

basestats_URL = "https://serenesforest.net/awakening/characters/base-stats/main-story/"

basestats_page = requests.get(basestats_URL)

soup = BeautifulSoup(basestats_page.content, "html.parser")

pprint(soup)
