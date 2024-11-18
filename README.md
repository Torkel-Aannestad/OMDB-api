# MovieMaze

## Database and Model design

- migrations are handled with goose from the sql/schema directory
- mapping between sql schema, sql queries and Go types are done with sqlc. See sqlc.yaml config file
- sqlc is configured for autogenerating json tags for structs. By default uses the database column name as json value. Caseing can be configued. json value can be overwritten in sqlc.yaml if needed.
- to start the handlers returns the generated types directly, however, set up a mapping to a seperate application type if needed. FOr example in case you don't want include some columns.
- With SQLC is not convenient to do build dynamic SLQ queries with go code such as sort column and sort direction.

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
