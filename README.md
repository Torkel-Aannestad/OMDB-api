# MovieMaze

## API Documentaion

### Base URL

All API endpoints are prefixed by /v1.

### Endpoints

#### Healthcheck

GET /v1/healthcheck

- Description: Check the health status of the API.
- Authentication: None

#### Movies

- Movie Response example:

```{
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
    },
```

##### GET /v1/movies

- Description: Retrieve a list of movies.
- Permission: movies:read

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

## Database and Model design

- migrations are handled with goose from the sql/schema directory
- mapping between sql schema, sql queries and Go types are done with sqlc. See sqlc.yaml config file
- sqlc is configured for autogenerating json tags for structs. By default uses the database column name as json value. Caseing can be configued. json value can be overwritten in sqlc.yaml if needed.
- to start the handlers returns the generated types directly, however, set up a mapping to a seperate application type if needed. FOr example in case you don't want include some columns.
- With SQLC is not convenient to do build dynamic SLQ queries with go code such as sort column and sort direction.

### Refactored

Refactor uten sqlc. Implementer selv, men bruk sqlc til å generere for deg. Da kan du copy pasta og endre litt. Mindre sannsynelighet for feil i mapping mellom feltene.

- vi får bedre håndtering av ctx og error
- full kontrol der vi må bruk Sprintf til å bygge opp størring.
- vi kan legge valideringslogikk sammen med types.
- vi kan legge til response types der vi trenger der. eks userResponse som ikke skal inneholde alle feltene.
- vi kan nå generere kode med sqlc, men kopiere den over får full kontrol.

### OMDB

- Added to Makefile to transfer csv data and import to DB

## Mailer

- Mailtrap
- MailTrap for sending transational emails. Easy free setup with no requiment for adding a domain.
- go-mail for handling SMTP. https://pkg.go.dev/github.com/go-mail/mail
- Sending email with background Go routine

## Misc

- IP based rate limiting with x/time/rate package
- Getting user's IP with Realip package
  - github.com/tomasen/realip
- users email. Use postgresql plugin citext to make string case insensitive. This way we don't need to worry case.
- error triage

```

```
