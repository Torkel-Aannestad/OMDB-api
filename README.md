# MovieMaze

This project is a personal learning project. MovieMaze is a movie and series database presented as a JSON-api. The data is downloaded from [OMDB](https://www.moviemaze.org) which is a free community driven database for film media. The goal of the project has been to create a full featured ideomatic Go JSON API, which uses few dependencies and abstractions.

Stack:

- HTTPRouter for routing
- PostgreSQL database
- Models built with raw SQL and Go standard libary
- Auth and user management built from cratch
- Mailer implemented with Go-mail and Mailtrap SMTP relay
- Makefile for automations
- Hosting from a Ubuntu VM with Caddy reverse proxy and rsync for file transfer and shell scripts

Sentral to the design of the application is Alex Edwards' books [Let's Go](https://lets-go.alexedwards.net/) and [Let's Go Further](https://lets-go-further.alexedwards.net/). You can read more about the technology stack and design desitions taken in the project below.

<strong>Created by Torkel Aannestad</strong>

- [torkelaannestad.com](torkelaannestad.com)
- [Github](github.com/Torkel-Aannestad/)

Jump to the API-documentation:
[API documentation](#api-documentaion)

## Quickstart

The API is open for anyone to use. Sign up and activate you user account and you are ready to use the API. The database will occationally be reset, but please don't perform to many destructive operations.

You'll get access with 3 steps:

1. User signup
2. Confirm your email with token
3. Authenticate with email and password to get a request token

### 1. User Signup

- Base url: moviemaze.torkelaannestad.com
- Endpoint: POST /v1/users
- Body: name, email, password

```shell
  BODY='{"name": "Jake Perolta","email": "jake.perolta@example.com", "password": "yourSecurePassword"}'
  curl -d "$BODY" moviemaze.torkelaannestad.com/v1/users
```

### 2. User Activation

An email is sent to your email with activation token. Send the following request to active your account.

- Endpoint: PUT /v1/users/activate
- Body: token

```shell
  BODY='{"token": token-from-email}'
  curl -X PUT -d "$BODY" moviemaze.torkelaannestad.com/v1/users/activate
```

Reponse: a user object with updated "activated" value.

### 3. Get Auth Token

- Endpoint: POST /v1/auth/authentication
- Body: email, password

```shell
  BODY='{"email": "yourEmail@example.com", "password": "pa55word"}'
  curl -d "$BODY" moviemaze.torkelaannestad.com/v1/auth/authentication
```

Response:

```JSON
{"authentication_token":{"token":"JZ5B4SABDN7PBKUFKROWAQM7DU","expiry":"2024-12-13T09:16:07.9290901Z"}}
```

<br/>
<strong>Request data</strong>

- Endpoint: GET /v1/movies/:id

```shell
  curl -H "Authorization: Bearer G5TU7Y46GRENMNUDZP2T75QGNE" moviemaze.torkelaannestad.com/v1/movies/35819
```

<br/>

## Database and Model design

- OMDB is imported from CSV files. An import SQL script is created to set up the data model in a good starting state. See /sql/data-import/run.sql for all the set up steps.
- Added to Makefile to transfer csv data and import to DB in production.
- After initial data import migrations are handled with goose from the sql/schema directory.
- sqlc is configured for autogenerating json tags for Go structs. The generated types are not used directly but copied and modified. This way we get better control over the context.Context instance and error handling. We also get full control when needing to build dynamic queries.
- PostgreSQL configured with citext plugin for user email column to make string case insensitive.
- Full text search features in PostgreSQL is configured to enabled a good search experience with for examample movies or people resources.

## Mailer

- MailTrap for sending transational emails.
- go-mail for handling SMTP.
- Sending email with background Go routine
- Email templates are found in assets/templates

## Middleware

- Authenticate middleware ensures that we retrieve the user from the database or know that the user is anonymous. The user is added to the request context your later use.
- The protectedRoute middleware bounces the user if she is not activated, has the right permission or is anonymous.

## MISC

- IP based rate limiting with x/time/rate package
- Getting user's IP with Realip package
  - github.com/tomasen/realip
- Error triage is implemented in readJSON() helper function to ensure that potensial errors from reading json body is caught and error messages regarding what the issue is can be sendt to the user. A standardized set of response messages are found in cmd/api/errors.go to ensure that only known formulations will reach the end user, and thus hiding for example error messages bubbling from PostgreSQL.

## Roadmap

- Improved testing
- More advance auth features like MFA and additional rate limiting on auth endpoints.

## API Documentaion

Base URL and example endpoint:

```shell
 moviemaze.torkelaannestad.com/v1/healthcheck
```

- To get access to the API please see the [quickstart section](#Quickstart) or the [users](#Users) and [auth](#Authentication) resources.
- Error handling and expected status codes are found in [error handling](#Error-Handling).
- Permissions. The api is implementet with permission based authorization. Upon signup your user will be granted both read and write access to most resources. Please behave nicely.
- Optimistic concurrency control is applied to any records that can be updated thought the version field. This way multile simultanious requests to update a will fail with status code 409 conflict.

### Resources

Overview

- [Healthcheck](#Healthcheck)
- [Movies](#Movies)
- [People](#People)
- [Casts](#Casts)
- [Jobs](#Jobs)
- [Categories](#Categories)
- [Movie Keywords](#Movie-Keywords)
- [Movie Categories](#Movie-Categories)
- [Movie Links](#Movie-Links)
- [People Links](#People-Links)
- [Trailers](#Trailers)
- [Images](#Images)
- [Users](#Users)
- [Authentication](#Authentication)

#### Healthcheck

```shell
 GET /v1/healthcheck
```

- Description: Check the health status of the API.
- Authentication: None

#### Movies

The movies table include both movies, series and episodes. Series uses the parent_id and series_id fields to form a hierarchy between top level series, seasons and episodes.

- Movie Response example:

```JSON
{
  "id": 1,
  "parent_id": null,
  "series_id": null,
  "name": "Inception",
  "date": "2010-07-16T00:00:00Z",
  "kind": "movie",
  "runtime": 148,
  "budget": 160000000,
  "revenue": 829895144,
  "homepage": "https://www.inceptionmovie.com",
  "vote_average": 8.8,
  "vote_count": 210000,
  "abstract": "A skilled thief is given a chance at redemption if he can successfully perform inception.",
  "version": 1
}
```

##### GET /v1/movies

- Description: Retrieve a list of movies. The endpoint allows full text search thought query parameters.
- Query parameters:
  - name: Full text search
  - kind: movie, series, season, episode, movieseries (enum)
  - page: default 1
  - page_size: number of record for each page
  - sort: default "id". Use "-" for descending order. Valid values are: "id", "name", "date", "runtime", "-id", "-name", "-date", "-runtime".
- Permission: movies:read

```shell
 curl -H "Authorization: Bearer JZ5B4SABDN7PBKUFJROWAQM7DU" "https://moviemaze.torkelaannestad.com/v1/movies?name=dark%20knight&kind=movie&page=1&page_size=2&sort=id"
```

Example response:

```JSON
{
  "metadata": {
    "current_page": 1,
    "page_size": 2,
    "first_page": 1,
    "last_page": 3,
    "total_records": 5
  },
  "movies": [
    {
      "id": 155,
      "parent_id": 158459,
      "series_id": null,
      "name": "The Dark Knight",
      "date": "2008-07-18T00:00:00Z",
      "kind": "movie",
      "runtime": 150,
      "budget": 185000000,
      "revenue": 970000000,
      "homepage": "",
      "vote_average": 6.9402985573,
      "vote_count": 67,
      "abstract": "",
      "version": 1
    },
    {
      "id": 32937,
      "parent_id": 158459,
      "series_id": null,
      "name": "The Dark Knight Rises",
      "date": "2012-07-20T00:00:00Z",
      "kind": "movie",
      "runtime": 164,
      "budget": 250000000,
      "revenue": 576798000,
      "homepage": "http://www.thedarkknightrises.com/",
      "vote_average": 6.125,
      "vote_count": 8,
      "abstract": "",
      "version": 1
    }
  ]
}
```

##### POST /v1/movies

- Description: Create a new movie.
- Body: send a movie object as the body of the request. id and version should not be included. Parent_id and series_id are nullable int64 types and can be omited.
- Permission: movies:write

```shell
 BODY='{"name":"Go programming is awesome","date":"2024-12-02T00:00:00Z", "kind":"movie", "runtime":108,"budget":0,"revenue":0,"homepage":"", "vote_average": 5.4, "votes_count": 23, "abstract": ""}'
 curl -d "$BODY" -H "Authorization: Bearer 2FEPZTRMM6WDQXEX7SLO47RJFE" localhost:4000/v1/movies
```

Response

```JSON
{
  "movie": {
    "id": 272776,
    "parent_id": {
      "Int64": 0,
      "Valid": false
    },
    "series_id": {
      "Int64": 0,
      "Valid": false
    },
    "name": "Go programming is awesome",
    "date": "2024-12-02T00:00:00Z",
    "kind": "movie",
    "runtime": 108,
    "budget": 0,
    "revenue": 0,
    "homepage": "",
    "vote_average": 5.4,
    "vote_count": 23,
    "abstract": "",
    "version": 1
  }
}

```

##### GET /v1/movies/:id

- Description: Retrieve a specific movie by ID.
- Query Parameter: id is a movie_id
- Permission: movies:read

```shell
 curl -H "Authorization: Bearer 2FEPZTRMM6WDQXEX7SLO47RJFE" localhost:4000/v1/movies/35819
```

Response: same as create movie.

##### PATCH /v1/movies/:id

- Description: Update a specific movie by ID.
- Query Parameter: id is a movie_id
- Body: fields that you want to update.
- Permission: movies:write

```shell
 BODY='{ "vote_average": 5.6, "votes_count": 25}'
 curl -X PATCH -d "$BODY" -H "Authorization: Bearer 2FEPZTRMM6WDQXEX7SLO47RJFE" http://localhost:4000/v1/movies/272775
```

##### DELETE /v1/movies/:id

- Description: Delete a specific movie by ID.
- Query Parameter: id is a movie_id
- Permission: movies:write

```shell
 curl -X DELETE -H "Authorization: Bearer 2FEPZTRMM6WDQXEX7SLO47RJFE" http://localhost:4000/v1/movies/272775
```

#### People

##### GET /v1/people

- Description: Retrieve a list of people. The endpoint supports full-text search using query parameters.
- Query Parameters:
  - name: Full text search by name.
  - page: Page number for pagination, default is 1.
  - page_size: Number of records for each page.
  - sort: Default is "id". Use "-" for descending order. Valid values are: id, name, birthday, -id, -name, -birthday.
- Permission: people:read
- Gender: 0=male, 1=female, 2=non-binary, 99=not spesified

```shell
 curl -H "Authorization: Bearer YOUR_TOKEN" "http://localhost:4000/v1/people?name=john&page=1&page_size=5&sort=id"
```

##### POST /v1/people

- Description: Create a new person.
- Body: Send a person object as the body of the request. id and version should not be included.
- Permission: people:write

```shell
 BODY='{"name":"John Doe","birthday":"1944-05-14T00:00:00Z","gender":"0","aliases":["foo, bar"]}'
 curl -d "$BODY" -H "Authorization: Bearer 2FEPZTRMM6WDQXEX7SLO47RJFE" http://localhost:4000/v1/people
```

##### GET /v1/people/:id

- Description: Retrieve a specific person by ID.
- Query Parameter: person id
- Permission: people: read

```shell
 curl -H "Authorization: Bearer 2FEPZTRMM6WDQXEX7SLO47RJFE" localhost:4000/v1/people/311418
```

##### PATCH /v1/people/:id

- Description: Update a specific person by ID.
- Path Parameter: id is the person_id.
- Body: Fields that you want to update.
- Permission: people:write

```shell
 BODY='{"name":"John Doeski","gender":"99"'
 curl -X PPATCH -d "$BODY" -H "Authorization: Bearer 2FEPZTRMM6WDQXEX7SLO47RJFE" http://localhost:4000/v1/people
```

Response: the updated user object.

##### DELETE /v1/people/:id

- Description: Delete a specific person by ID.
- Query Parameter: person id.
- Permission: people:write

```shell
 curl -X DELETE -H "Authorization: Bearer 2FEPZTRMM6WDQXEX7SLO47RJFE" localhost:4000/v1/people/311418
```

#### Casts

##### POST /v1/casts

- Description: Create a new cast entry.
- Body: see field below. Position field is the position that the cast will appear in the list. Meaning position 1 should be example be lead actor.
- Permission: casts:write

```shell
 BODY='{"movie_id":35819,"person_id":287,"job_id":15,"role":"Very cool role","position":1}'
curl -d "$BODY" -H "Authorization: Bearer 2FEPZTRMM6WDQXEX7SLO47RJFE" http://localhost:4000/v1/casts
```

##### GET /v1/casts/by-movie-id/:id

- Description: Retrieve casts associated with a specific movie by movie ID.
- Query Parameter: movie id.
- Permission: casts:read

```shell
 curl -H "Authorization: Bearer 2FEPZTRMM6WDQXEX7SLO47RJFE" "localhost:4000/v1/casts/by-movie-id/35819"
```

Reponse: list of casts

##### GET /v1/casts/by-person-id/:id

- Description: Retrieve casts associated with a specific person by person ID.
- Query Parameter: person id.
- Permission: casts:read

```shell
 curl -H "Authorization: Bearer 2FEPZTRMM6WDQXEX7SLO47RJFE" "localhost:4000/v1/casts/by-person-id/2524"
```

Reponse: list of casts

##### PATCH /v1/casts/:id

- Description: Update a specific cast entry by ID.
- Query Parameter: cast id.
- Body: Fields that you want to update.
- Permission: casts:write

```shell
 BODY='{"role":"Luke Skywalker","position":2}'
 curl -X PATCH -d "$BODY" -H "Authorization: Bearer 2FEPZTRMM6WDQXEX7SLO47RJFE" http://localhost:4000/v1/casts/1260480
```

##### DELETE /v1/casts/:id

- Description: Delete a specific cast entry by ID.
- Query Parameter: cast id.
- Permission: casts:write

```shell
 curl -X DELETE -H "Authorization: Bearer 2FEPZTRMM6WDQXEX7SLO47RJFE" http://localhost:4000/v1/casts/1260479
```

#### Jobs

##### POST /v1/jobs

- Description: Create a new job.
  Body: name of the job
- Permission: jobs:write

```shell
 BODY='{"name":"Executive Producer"}'
 curl -d "$BODY" -H "Authorization: Bearer 2FEPZTRMM6WDQXEX7SLO47RJFE" http://localhost:4000/v1/jobs
```

##### GET /v1/jobs/:id

- Description: Retrieve a specific job by ID.
- Query Parameter: job id.
- Permission: jobs:read

```shell
 curl -H "Authorization: Bearer 2FEPZTRMM6WDQXEX7SLO47RJFE" "localhost:4000/v1/jobs/1050"
```

##### PATCH /v1/jobs/:id

- Description: Update a specific job by ID.
- Query Parameter: job id.
- Body: fields you want do update.
- Permission: jobs:write

```shell
 BODY='{"name":"Super Executive Producer"}'
 curl -X PATCH -d "$BODY" -H "Authorization: Bearer 2FEPZTRMM6WDQXEX7SLO47RJFE" http://localhost:4000/v1/jobs/1050
```

##### DELETE /v1/jobs/:id

- Description: Delete a specific job by ID.
- Query Parameter: job id.
- Permission: jobs:write

```shell
 curl -X DELETE -H "Authorization: Bearer 2FEPZTRMM6WDQXEX7SLO47RJFE" http://localhost:4000/v1/jobs/1050
```

#### Categories

##### POST /v1/categories

- Description: Create a new category.
- Body: name of the category, parent_id, root_id are nullable int64 types. Use parent and root id to create sub-categories.
- Permission: categories:write

```shell
 BODY='{"name":"Family Drama"}'
 BODY='{"name":"Family Drama", "parent_id":12,"root_id":1}'
 curl -d "$BODY" -H "Authorization: Bearer 2FEPZTRMM6WDQXEX7SLO47RJFE" http://localhost:4000/v1/categories
```

##### GET /v1/categories/:id

- Description: Retrieve a specific category by ID.
- Query Parameter: category id.
- Permission: categories:read

```shell
 curl -H "Authorization: Bearer 2FEPZTRMM6WDQXEX7SLO47RJFE" "localhost:4000/v1/categories/20705"
```

##### PATCH /v1/categories/:id

- Description: Update a specific category by ID.
- Query Parameter: category id.
- Body: fields you want do update.
- Permission: categories:write

```shell
 BODY='{"name":"Genre Adventure Subcategory", "parent_id":12,"root_id":1}'
 curl -X PATCH -d "$BODY" -H "Authorization: Bearer 2FEPZTRMM6WDQXEX7SLO47RJFE" "localhost:4000/v1/categories/20705"
```

##### DELETE /v1/categories/:id

- Description: Delete a specific category by ID.
- Query Parameter: category id.
- Permission: categories:write

```shell
 curl -X DELETE -H "Authorization: Bearer 2FEPZTRMM6WDQXEX7SLO47RJFE" http://localhost:4000/v1/categories/20705
```

#### Movie Keywords

##### POST /v1/movie-keywords

- Description: Add keywords to a movie.
- Body: movie_id and category_id
- Permission: category-items:write

```shell
 BODY='{"movie_id":35819,"category_id": 10}'
 curl -d "$BODY" -H "Authorization: Bearer 2FEPZTRMM6WDQXEX7SLO47RJFE" http://localhost:4000/v1/movie-keywords
```

##### GET /v1/movie-keywords/:id

- Description: Retrieve keywords associated with a movie by movie ID.
- Query Parameter: movie_keyword id.
- Permission: category-items:read

```shell
 curl -H "Authorization: Bearer 2FEPZTRMM6WDQXEX7SLO47RJFE" "localhost:4000/v1/movie-keywords/35819"
```

##### DELETE /v1/movie-keywords

- Description: Delete keywords from a movie.
- Body: movie_id and category_id
- Permission: category-items:write

```shell
 BODY='{"movie_id":35819,"category_id":10}'
 curl -X DELETE -d "$BODY" -H "Authorization: Bearer 2FEPZTRMM6WDQXEX7SLO47RJFE" http://localhost:4000/v1/movie-keywords
```

#### Movie Categories

##### POST /v1/movie-categories

- Description: Add categories to a movie.
- Body: movie_id and category_id
- Permission: category-items:write

```shell
 BODY='{"movie_id":35819,"category_id": 10}'
 curl -d "$BODY" -H "Authorization: Bearer 2FEPZTRMM6WDQXEX7SLO47RJFE" http://localhost:4000/v1/movie-categories
```

##### GET /v1/movie-categories/:id

- Description: Retrieve categories associated with a movie by movie ID.
- Query Parameter: movie_category id.
- Permission: category-items:read

```shell
 curl -H "Authorization: Bearer 2FEPZTRMM6WDQXEX7SLO47RJFE" "localhost:4000/v1/movie-categories/35819"
```

##### DELETE /v1/movie-categories

- Description: Delete categories from a movie.
- Body: movie_id and category_id
- Permission: category-items:write

```shell
 BODY='{"movie_id":35819,"category_id":10}'
 curl -X DELETE -d "$BODY" -H "Authorization: Bearer 2FEPZTRMM6WDQXEX7SLO47RJFE" http://localhost:4000/v1/movie-categories
```

#### Movie Links

##### POST /v1/movie-links

- Description: Create a movie link.
- Body:
  - source: query parameter or id for the source.
  - key: must be: "wikidata", "wikipedia", "imdbperson"
  - language: english is the only language used. When using xx as value english is assumed.
- Permission: movie-links:write

```shell
 BODY='{"source":"Fight_Club","key":"wikipedia","movie_id":550,"language":"xx"}'
 curl -d "$BODY" -H "Authorization: Bearer 2FEPZTRMM6WDQXEX7SLO47RJFE" http://localhost:4000/v1/movie-links
```

##### GET /v1/movie-links/:id

- Description: Retrieve links associated with a movie by movie ID.
- Query Parameter: movie id.
- Permission: movie-links:read

```shell
 curl -H "Authorization: Bearer 2FEPZTRMM6WDQXEX7SLO47RJFE" "localhost:4000/v1/movie-links/550"
```

##### DELETE /v1/movie-links/:id

- Description: Delete a movie link by ID.
- Query Parameter: movie-link id.
- Permission: movie-links:write

```shell
 curl -X DELETE -H "Authorization: Bearer 2FEPZTRMM6WDQXEX7SLO47RJFE" http://localhost:4000/v1/movie-links/348425
```

#### People Links

##### POST /v1/people-links

- Description: Create a people link.
- Body:
  - source: query parameter or id for the source.
  - key: must be: "wikidata", "wikipedia", "imdbperson"
  - language: english is the only language used. When using xx as value english is assumed.
- Permission: people-links:write

```shell
 BODY='{"source":"Brad_Pitt","key":"wikipedia","person_id":287,"language":"xx"}'
 curl -d "$BODY" -H "Authorization: Bearer 2FEPZTRMM6WDQXEX7SLO47RJFE" http://localhost:4000/v1/people-links
```

##### GET /v1/people-links/:id

- Description: Retrieve links associated with a person by person ID.
- Query Parameter: person link id.
- Permission: people-links:read

```shell
 curl -H "Authorization: Bearer 2FEPZTRMM6WDQXEX7SLO47RJFE" "localhost:4000/v1/people-links/287"
```

##### DELETE /v1/people-links/:id

- Description: Delete a people link by ID.
- Query Parameter: person-link id.
- Permission: people-links:write

```shell
 curl -X DELETE -H "Authorization: Bearer 2FEPZTRMM6WDQXEX7SLO47RJFE" http://localhost:4000/v1/people-links/213048
```

#### Trailers

##### POST /v1/trailers

- Description: Add a trailer to a movie.
- Body:
  - key: id of the yotube or vimeo video
  - movie_id
  - language
  - source: "youtube" or "vimeo"
- Permission: trailers:write

```shell
 BODY='{"key":"youtube-id-from-url","movie_id":35819,"language":"en","source":"youtube"}'
 curl -d "$BODY" -H "Authorization: Bearer 2FEPZTRMM6WDQXEX7SLO47RJFE" http://localhost:4000/v1/trailers
```

##### GET /v1/trailers/:id

- Description: Retrieve trailers associated with a movie by movie ID.
- Query Parameter: trailer id.
- Permission: trailers:read

```shell
 curl -H "Authorization: Bearer 2FEPZTRMM6WDQXEX7SLO47RJFE" "localhost:4000/v1/trailers/35819"
```

##### DELETE /v1/trailers/:id

- Description: Delete a trailer by ID.
- Query Parameter: trailer id.
- Permission: trailers:write

```shell
 curl -X DELETE -H "Authorization: Bearer 2FEPZTRMM6WDQXEX7SLO47RJFE" http://localhost:4000/v1/trailers/10922
```

#### Images

##### POST /v1/images

- Description: Upload an image.
- Body:
  - object_id: id of person, movie, job, category
  - object_type: "Movie", "Person", "Job", "Category"
- Permission: images:write

```shell
 BODY='{"object_id":272775, "object_type": "Movie"}'
 curl -d "$BODY" -H "Authorization: Bearer 2FEPZTRMM6WDQXEX7SLO47RJFE" http://localhost:4000/v1/images
```

##### GET /v1/images/:id

- Description: Retrieve a specific image by ID.
- Query Parameter: image id.
- Permission: images:read

```shell
 curl -H "Authorization: Bearer 2FEPZTRMM6WDQXEX7SLO47RJFE" "localhost:4000/v1/images/60360"
```

##### GET /v1/images

- Description: Retrieve images by object ID.
- Query Parameter:
  - object_id: id of person, movie, job, category
  - object_type: "Movie", "Person", "Job", "Category"
- Permission: images:read

```shell
 curl -H "Authorization: Bearer 2FEPZTRMM6WDQXEX7SLO47RJFE" "localhost:4000/v1/images?object_id=272775&object_type=Movie"
```

##### PATCH /v1/images/:id

- Description: Update a specific image by ID.
- Body:
  - object_id: id of person, movie, job, category
  - object_type: "Movie", "Person", "Job", "Category"
- Permission: images:write

```shell
 BODY='{"object_id":272775, "object_type": "Movie"}'
 curl -X PATCH -d "$BODY" -H "Authorization: Bearer 2FEPZTRMM6WDQXEX7SLO47RJFE" http://localhost:4000/v1/images/60360
```

##### DELETE /v1/images/:id

- Description: Delete a specific image by ID.
- Query Parameter: image id.
- Permission: images:write

```shell
 curl -X DELETE -H "Authorization: Bearer 2FEPZTRMM6WDQXEX7SLO47RJFE" http://localhost:4000/v1/images/60360
```

#### Users

##### POST /v1/users

- Description: Register a new user.
- Authentication: None

```shell
  BODY='{"name": "Jake Perolta","email": "jake.perolta@example.com", "password": "yourSecurePassword"}'
  curl -d "$BODY" moviemaze.torkelaannestad.com/v1/users
```

##### PUT /v1/users/activated

- Description: Activate a registered user.
- Body: token from email
- Authentication: None

```shell
  BODY='{"token": token-from-email}'
  curl -X PUT -d "$BODY" moviemaze.torkelaannestad.com/v1/users/activate
```

##### POST /v1/users/resend-activation-token

- Description: Resend activate token to users email.
- Body: email
- Authentication: None

```shell
  BODY='{"email": yourEmail@example.com}'
  curl -d "$BODY" moviemaze.torkelaannestad.com/v1/users/resend-activation-token
```

Response: json message confirmation.

#### Authentication

Authentication endpoint are protected with an additional rate limiter with a slow refill rate. This is to protect agains password sprays or other brute force methods agains the auth resources.

##### POST /v1/auth/authentication

- Description: Authenticate a user and obtain a token.
- Body: email, password
- Authentication: None

```shell
  BODY='{"email": "yourEmail@example.com", "password": "pa55word"}'
  curl -d "$BODY" moviemaze.torkelaannestad.com/v1/auth/authentication
```

Response: auth token

##### POST /v1/auth/change-password

- Description: Change user password and obtain a token. Auth token is regenerated and old sessions are invalidated.
- Body: current_password and new_password
- Authenticated

```shell
  BODY='{"current_password": "pa55word", "password": "NewStrongerpa55word"}'
  curl -d "$BODY" localhost:4000/v1/auth/change-passord
```

Response: auth token

##### POST v1/auth/password-reset-verify-email

- Description: Reset user password. Send user email in body and receive verification code by email.
- Body: email

```shell
 BODY='{"email": "yourEmail@example.com"}'
 curl -d "$BODY" localhost:4000/v1/auth/password-reset-verify-email
```

Response: json message

##### POST v1/auth/password-reset

- Description: Reset user password. Send email verification token from email in body along with new password.
- Body: token

```shell
  BODY='{"token": "ZAAKSYTB2CV2RU2OOQ2JA5K35Y", "new_password": "NewStrongerpa55word"}'
  curl -d "$BODY" localhost:4000/v1/auth/reset-password
```

Response: new authentication token

##### POST v1/auth/revoke

- Description: Revokes all sessions for user.
- Authentication: True

```shell
  curl -X POST -H "Authorization: Bearer W2LMNM6XED32GE3WW5GRFR7PMU" localhost:4000/v1/auth/revoke
```

Response: json message confirmation

### Error Handling

- 400 Bad Request: General response when body or query params are invalid.
- 422 Unprocessable Entity: Failed validations response. Fields are errors are specified.
- 404 Not Found: If a movie with the specified ID does not exist.
- 409 Conflict: If there's a concurrency edit conflict (version mismatch).
- 405 Method Not Allowed: If the method is not allowed on the specified route.
- 429 Too Many Request: Failed due to rate limiting.
- 401 Unauthorized: any issue with auth token.
- 403 Forbidden: When trying to access resources without respective permission or account not active.
- 500 Server error: General server error response
