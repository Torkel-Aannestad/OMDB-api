package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/movies", app.protectedRoutePermission("movies:read", app.listMoviesHandler))
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.protectedRoutePermission("movies:write", app.createMovieHandler))
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.protectedRoutePermission("movies:read", app.getMovieHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.protectedRoutePermission("movies:write", app.updateMovieHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.protectedRoutePermission("movies:write", app.deleteMovieHandler))

	router.HandlerFunc(http.MethodGet, "/v1/people", app.protectedRoutePermission("people:read", app.listPeopleHandler))
	router.HandlerFunc(http.MethodPost, "/v1/people", app.protectedRoutePermission("people:write", app.createPeopleHandler))
	router.HandlerFunc(http.MethodGet, "/v1/people/:id", app.protectedRoutePermission("people:read", app.getPeopleHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/people/:id", app.protectedRoutePermission("people:write", app.updatePeopleHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/people/:id", app.protectedRoutePermission("people:write", app.deletePeopleHandler))

	router.HandlerFunc(http.MethodPost, "/v1/casts", app.protectedRoutePermission("casts:write", app.createCastHandler))
	router.HandlerFunc(http.MethodGet, "/v1/casts/by-movie-id/:id", app.protectedRoutePermission("casts:read", app.getCastsByMovieIdHandler))
	router.HandlerFunc(http.MethodGet, "/v1/casts/by-person-id/:id", app.protectedRoutePermission("casts:read", app.getCastsByPersonIdHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/casts/:id", app.protectedRoutePermission("casts:write", app.updateCastHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/casts/:id", app.protectedRoutePermission("casts:write", app.deleteCastHandler))

	router.HandlerFunc(http.MethodPost, "/v1/jobs", app.protectedRoutePermission("jobs:write", app.createJobHandler))
	router.HandlerFunc(http.MethodGet, "/v1/jobs/:id", app.protectedRoutePermission("jobs:read", app.getJobHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/jobs/:id", app.protectedRoutePermission("jobs:write", app.updateJobHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/jobs/:id", app.protectedRoutePermission("jobs:write", app.deleteJobHandler))

	router.HandlerFunc(http.MethodPost, "/v1/categories", app.protectedRoutePermission("categories:write", app.createCategoryHandler))
	router.HandlerFunc(http.MethodGet, "/v1/categories/:id", app.protectedRoutePermission("categories:read", app.getCategoryHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/categories/:id", app.protectedRoutePermission("categories:write", app.updateCategoryHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/categories/:id", app.protectedRoutePermission("categories:write", app.deleteCategoryHandler))

	router.HandlerFunc(http.MethodPost, "/v1/movie-keywords", app.protectedRoutePermission("category-items:write", app.createMovieKeywordsHandler))
	router.HandlerFunc(http.MethodGet, "/v1/movie-keywords/:id", app.protectedRoutePermission("category-items:read", app.getMovieKeywordsHandler)) //expects movieId
	router.HandlerFunc(http.MethodDelete, "/v1/movie-keywords", app.protectedRoutePermission("category-items:write", app.deleteMovieKeywordHandler))
	router.HandlerFunc(http.MethodPost, "/v1/movie-categories", app.protectedRoutePermission("category-items:write", app.createMovieCategoriesHandler))
	router.HandlerFunc(http.MethodGet, "/v1/movie-categories/:id", app.protectedRoutePermission("category-items:read", app.getMovieCategoriesHandler)) //expects movieId
	router.HandlerFunc(http.MethodDelete, "/v1/movie-categories", app.protectedRoutePermission("category-items:write", app.deleteMovieCategoryHandler))

	router.HandlerFunc(http.MethodPost, "/v1/movie-links", app.protectedRoutePermission("movie-links:write", app.createMovieLinkHandler))
	router.HandlerFunc(http.MethodGet, "/v1/movie-links/:id", app.protectedRoutePermission("movie-links:read", app.getMovieLinksHandler))       //expects movieId
	router.HandlerFunc(http.MethodDelete, "/v1/movie-links/:id", app.protectedRoutePermission("movie-links:write", app.deleteMovieLinkHandler)) //expects id from movie_links

	router.HandlerFunc(http.MethodPost, "/v1/people-links", app.protectedRoutePermission("people-links:write", app.createPeopleLinkHandler))
	router.HandlerFunc(http.MethodGet, "/v1/people-links/:id", app.protectedRoutePermission("people-links:read", app.getPeopleLinksHandler))       //expects personId
	router.HandlerFunc(http.MethodDelete, "/v1/people-links/:id", app.protectedRoutePermission("people-links:write", app.deletePeopleLinkHandler)) //expects id from people_links

	router.HandlerFunc(http.MethodPost, "/v1/trailers", app.protectedRoutePermission("trailers:write", app.createTrailerHandler))
	router.HandlerFunc(http.MethodGet, "/v1/trailers/:id", app.protectedRoutePermission("trailers:read", app.getTrailersHandler)) //expects movieId
	router.HandlerFunc(http.MethodDelete, "/v1/trailers/:id", app.protectedRoutePermission("trailers:write", app.deleteTrailerHandler))

	router.HandlerFunc(http.MethodPost, "/v1/images", app.protectedRoutePermission("images:write", app.createImageHandler))
	router.HandlerFunc(http.MethodGet, "/v1/images/:id", app.protectedRoutePermission("images:read", app.getImageHandler))
	router.HandlerFunc(http.MethodGet, "/v1/images", app.protectedRoutePermission("images:read", app.getImagesObjektIdHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/images/:id", app.protectedRoutePermission("images:write", app.updateImageHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/images/:id", app.protectedRoutePermission("images:write", app.deleteImageHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activate", app.activateUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/auth/authentication", app.authenticateUserHandler) //open
	router.HandlerFunc(http.MethodPost, "/v1/auth/change-password", app.protectedRoute(app.changePasswordHandler))
	router.HandlerFunc(http.MethodGet, "/v1/auth/password-reset", app.protectedRoute(app.resetPasswordHandler))

	// router.Handler(http.MethodGet, "/metrics", expvar.Handler())

	//Admin swap permission
	router.HandlerFunc(http.MethodPost, "/v1/users/permissions/:id", app.protectedRoutePermission("admin:write", app.addUserPermissionsHandler))

	router.HandlerFunc(http.MethodGet, "/", app.getDocs)
	return app.panicRecovery(app.rateLimit(app.authenticate(router)))
}
