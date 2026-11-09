package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		app.notFound(w)
	})

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath",http.StripPrefix("/static", fileServer))


	dynamic := alice.New(app.sessionManager.LoadAndSave)
	protected := dynamic.Append(app.requireAuthentication)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.Home))
	// router.Handler(http.MethodGet,"/rbac/view", dynamic.ThenFunc(app.RbacView)
	router.Handler(http.MethodGet, "/handle/username", protected.ThenFunc(app.HandleSearchByUsername))
	router.Handler(http.MethodPost, "/handle/username", protected.ThenFunc(app.HandleSearchByUsernamePost))

	router.Handler(http.MethodGet, "/handle/id", protected.ThenFunc(app.HandleSearchByUserID))
	router.Handler(http.MethodPost, "/handle/id", protected.ThenFunc(app.HandleSearchByUserIDPost))

	router.Handler(http.MethodGet,"/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost,"/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet,"/user/login",dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost,"/user/login",dynamic.ThenFunc(app.userLoginPost))

	router.Handler(http.MethodPost,"/user/logout",protected.ThenFunc(app.userLogoutPost))


	standard := alice.New(app.recoverPanic,app.logRequest,secureHeaders)

	return standard.Then(router)
}