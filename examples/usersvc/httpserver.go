package usersvc

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func NewHTTPServer(svc Service) http.Handler {
	r := chi.NewRouter()
	r.Method(http.MethodGet, "/users/{name}", handleGetUser(svc))
	r.Method(http.MethodGet, "/users", handleListUsers(svc))
	r.Method(http.MethodPost, "/users", handleCreateUser(svc))
	r.Method(http.MethodPatch, "/users/{name}", handleUpdateUser(svc))
	r.Method(http.MethodDelete, "/users/{name}", handleDeleteUser(svc))
	return r
}

func handleGetUser(svc Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")

		user, err := svc.GetUser(r.Context(), name)
		if err != nil {
			errorToJSON(w, r, err)
			return
		}

		render.JSON(w, r, user)
	})
}

func handleListUsers(svc Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		users, err := svc.ListUsers(r.Context())
		if err != nil {
			errorToJSON(w, r, err)
			return
		}
		if users == nil {
			users = []*User{}
		}
		render.JSON(w, r, map[string]interface{}{"users": users})
	})
}

func handleCreateUser(svc Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := new(User)
		if err := render.DecodeJSON(r.Body, user); err != nil {
			errorToJSON(w, r, err)
			return
		}

		if err := svc.CreateUser(r.Context(), user); err != nil {
			errorToJSON(w, r, err)
			return
		}

		render.NoContent(w, r)
	})
}

func handleUpdateUser(svc Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")

		user := new(User)
		if err := render.DecodeJSON(r.Body, user); err != nil {
			errorToJSON(w, r, err)
			return
		}

		if err := svc.UpdateUser(r.Context(), name, user); err != nil {
			errorToJSON(w, r, err)
			return
		}

		render.NoContent(w, r)
	})
}

func handleDeleteUser(svc Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")

		if err := svc.DeleteUser(r.Context(), name); err != nil {
			errorToJSON(w, r, err)
			return
		}

		render.NoContent(w, r)
	})
}

func errorToJSON(w http.ResponseWriter, r *http.Request, err error) {
	render.Status(r, http.StatusBadRequest)
	render.JSON(w, r, map[string]string{"error": err.Error()})
}
