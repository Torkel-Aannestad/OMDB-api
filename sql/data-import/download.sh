#!/bin/bash

# https://www.omdb.org/content/Help:DataDownload

set -eu

BASE="https://www.omdb.org/data/"

FILES="
all_movies
all_series
all_seasons
all_episodes
all_movieseries

all_people
all_people_aliases
people_links

all_casts
job_names

movie_categories
movie_keywords
category_names
all_categories
trailers
movie_links

image_ids
image_licenses

all_movie_aliases_iso

all_votes
movie_languages
movie_countries
movie_details
movie_references
movie_abstracts_de
movie_abstracts_en
movie_abstracts_fr
movie_abstracts_es
movie_content_updates
"

for f in $FILES; do
  wget --no-verbose --mirror -P sql/data-import/ --no-host-directories "$BASE$f.csv.bz2"
done

bunzip2 --force sql/data-import/data/*.bz2
