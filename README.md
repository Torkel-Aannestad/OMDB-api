# MovieMaze

This project is a personal learning project. MovieMaze is a movie and series database presented as a JSON-api and later a React-based user interface. The data is downloaded from [OMDB](https://www.omdb.org) which is a free community driven database for film media. The goal of the project has been to create a full featured ideomatic Go JSON API, which uses few dependencies and abstractions.

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

<br/>
<h3>1. User Signup</h3>

- Base url: moviemaze.torkelaannestad.com
- Endpoint: POST /v1/users
- Body: name, email, password

```shell
  BODY='{"name": "Jake Perolta","email": "jake.perolta@example.com", "password": "yourSecurePassword"}'
  curl -d "$BODY" moviemaze.torkelaannestad.com/v1/users
```

<br/>
<h3>2. User Activation</h3>
An email is sent to your email with activation token. Send the following request to active your account.

- Endpoint: POST /v1/users/activate
- Body: token

```shell
  BODY='{"token": token-from-email}'
  curl -d "$BODY" moviemaze.torkelaannestad.com/v1/users/activate
```

<br/>
<h3>3. Get Auth Token</h3>

- Endpoint: POST /v1/auth/authentication
- Body: email, password

```shell
  BODY='{"email": "youEmail@example.com", "password": "pa55word"}'
  curl -d "$BODY" moviemaze.torkelaannestad.com/v1/auth/authentication
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
[Healthcheck](#Healthcheck)
[Movies](#Movies)
[People](#People)
[Casts](#Casts)
[Jobs](#Jobs)
[Categories](#Categories)
[Movie Keywords](#Movie-Keywords)
[Movie Categories](#Movie-Categories)
[Movie Links](#Movie-Links)
[People Links](#People-Links)
[Trailers](#Trailers)
[Images](#Images)
[Users](#Users)
[Authentication](#Authentication)

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
  - page: defaults to 1
  - page_size: number of record for each page
  - sort: defaults to "id". Use "-" for descending order. Valid values are: "id", "name", "date", "runtime", "-id", "-name", "-date", "-runtime".
- Permission: movies:read

```shell
 curl -H "Authorization: Bearer yourTokenHere" "https://moviemaze.torkelaannestad.com/v1/movies?name=mad%20max&page=1&page_size=5&sort=id"
```

Response:

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

##### POST /v1/movies

- Description: Create a new movie.
- Permission: movies:write

##### GET /v1/movies/:id

- Description: Retrieve a specific movie by ID.
- Permission: movies:read

##### PATCH /v1/movies/:id

- Description: Update a specific movie by ID.
- Permission: movies:write

##### DELETE /v1/movies/:id

- Description: Delete a specific movie by ID.
- Permission: movies:write

#### People

##### GET /v1/people

- Description: Retrieve a list of people.
- Permission: people:read

##### POST /v1/people

- Description: Create a new person.
- Permission: people:write

##### GET /v1/people/:id

- Description: Retrieve a specific person by ID.
- Permission: people: read

##### PATCH /v1/people/:id

- Description: Update a specific person by ID.
- Permission: people:write

##### DELETE /v1/people/:id

- Description: Delete a specific person by ID.
- Permission: people:write

#### Casts

##### POST /v1/casts

- Description: Create a new cast entry.
- Permission: casts:write

##### GET /v1/casts/by-movie-id/:id

- Description: Retrieve casts associated with a specific movie by movie ID.
- Permission: casts:read

##### GET /v1/casts/by-person-id/:id

- Description: Retrieve casts associated with a specific person by person ID.
- Permission: casts:read

##### PATCH /v1/casts/:id

- Description: Update a specific cast entry by ID.
- Permission: casts:write

##### DELETE /v1/casts/:id

- Description: Delete a specific cast entry by ID.
- Permission: casts:write

#### Jobs

##### POST /v1/jobs

- Description: Create a new job.
- Permission: jobs:write

##### GET /v1/jobs/:id

- Description: Retrieve a specific job by ID.
- Permission: jobs:read

##### PATCH /v1/jobs/:id

- Description: Update a specific job by ID.
- Permission: jobs:write

##### DELETE /v1/jobs/:id

- Description: Delete a specific job by ID.
- Permission: jobs:write

#### Categories

##### POST /v1/categories

- Description: Create a new category.
- Permission: categories:write

##### GET /v1/categories/:id

- Description: Retrieve a specific category by ID.
- Permission: categories:read

##### PATCH /v1/categories/:id

- Description: Update a specific category by ID.
- Permission: categories:write

##### DELETE /v1/categories/:id

- Description: Delete a specific category by ID.
- Permission: categories:write

#### Movie Keywords

##### POST /v1/movie-keywords

- Description: Add keywords to a movie.
- Permission: category-items:write

##### GET /v1/movie-keywords/:id

- Description: Retrieve keywords associated with a movie by movie ID.
- Permission: category-items:read

##### DELETE /v1/movie-keywords

- Description: Delete keywords from a movie.
- Permission: category-items:write

#### Movie Categories

##### POST /v1/movie-categories

- Description: Add categories to a movie.
- Permission: category-items:write

##### GET /v1/movie-categories/:id

- Description: Retrieve categories associated with a movie by movie ID.
- Permission: category-items:read

##### DELETE /v1/movie-categories

- Description: Delete categories from a movie.
- Permission: category-items:write

#### Movie Links

##### POST /v1/movie-links

- Description: Create a movie link.
- Permission: movie-links:write

##### GET /v1/movie-links/:id

- Description: Retrieve links associated with a movie by movie ID.
- Permission: movie-links:read

##### DELETE /v1/movie-links/:id

- Description: Delete a movie link by ID.
- Permission: movie-links:write

#### People Links

##### POST /v1/people-links

- Description: Create a people link.
- Permission: people-links:write

##### GET /v1/people-links/:id

- Description: Retrieve links associated with a person by person ID.
- Permission: people-links:read

##### DELETE /v1/people-links/:id

- Description: Delete a people link by ID.
- Permission: people-links:write

#### Trailers

##### POST /v1/trailers

- Description: Add a trailer to a movie.
- Permission: trailers:write

##### GET /v1/trailers/:id

- Description: Retrieve trailers associated with a movie by movie ID.
- Permission: trailers:read

##### DELETE /v1/trailers/:id

- Description: Delete a trailer by ID.
- Permission: trailers:write

#### Images

##### POST /v1/images

- Description: Upload an image.
- Permission: images:write

##### GET /v1/images/:id

- Description: Retrieve a specific image by ID.
- Permission: images:read

##### GET /v1/images

- Description: Retrieve images by object ID.
- Permission: images:read

##### PATCH /v1/images/:id

- Description: Update a specific image by ID.
- Permission: images:write

##### DELETE /v1/images/:id

- Description: Delete a specific image by ID.
- Permission: images:write

#### Users

##### POST /v1/users

- Description: Register a new user.
  Authentication: None

##### PUT /v1/users/activated

- Description: Activate a registered user.
- Authentication: None

#### Authentication

##### POST /v1/auth/authentication

- Description: Authenticate a user and obtain a token.
- Authentication: None

### Error Handling

400 Bad Request: If any of the input validations fail.
404 Not Found: If a movie with the specified ID does not exist.
409 Conflict: If there's a concurrency edit conflict (version mismatch).
405 Method Not Allowed: If the method is not allowed on the specified route.

## Roadmap

- Improved testing
- More advance auth features like MFA and additional rate limiting on auth endpoints.
