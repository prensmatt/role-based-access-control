package main
import(
	"html/template"
	"path/filepath"
	"net/http"

	"rbac.mheybatova.net/internal/models"
)

type templateData struct{
	User *models.User
	Users []models.User
	Form any
	Flash string
	IsAuthenticated bool
}

func(app *application) newTemplateData(r *http.Request) *templateData{
	return &templateData{
		Flash: app.sessionManager.PopString(r.Context(),"flash"),
		IsAuthenticated: app.isAuthenticated(r),
	}
}

func newTemplateCache()(map[string]*template.Template, error){
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil{
		return nil, err
	}

	for _,page := range pages{
		name := filepath.Base(page)

		files := []string{
			"./ui/html/base.tmpl",
			"./ui/html/partials/nav.tmpl",
			page,
		}
		ts,err := template.ParseFiles(files...)
		if err != nil{
			return nil,err
		}
		cache[name] = ts
	}
	return cache, nil
}