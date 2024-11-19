- Refactor uten sqlc. Implementer selv, men bruk sqlc til å generere for deg. Da kan du copy pasta og endre litt. Mindre sannsynelighet for feil i mapping mellom feltene.
  - vi får bedre håndtering av ctx og error
  - full kontrol der vi må bruk Sprintf til å bygge opp størring.
  - vi kan legge valideringslogikk sammen med types.
  - vi kan legge til response types der vi trenger der. eks userResponse som ikke skal inneholde alle feltene.
- Legge til reactivateTokenHandler
- Styling av eposter
- Tokens med kortere kode? annen base ecoding
- Add openID Connect https://github.com/coreos/go-oidc
- mux.HandleFunc("POST /api/login", cfg.HandlerAuthLogin)
  mux.HandleFunc("POST /api/refresh", cfg.HandlerAuthRefresh)
  mux.HandleFunc("POST /api/revoke", cfg.HandlerAuthRevoke)
