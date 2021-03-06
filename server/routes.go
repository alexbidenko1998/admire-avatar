package main

import (
	"admire-avatar/controllers"
	"admire-avatar/middlewares"
	"admire-avatar/ws"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func initRoutes() http.Handler {
	go ws.PrintsPool.Start()
	go ws.NotificationsPool.Start()

	r := mux.NewRouter()
	s := r.PathPrefix("/api").Subrouter()

	s.HandleFunc("/user/sign-up", controllers.SignUp).Methods("POST")
	s.HandleFunc("/user/sign-in", controllers.SignIn).Methods("POST")
	s.HandleFunc("/user/logout", controllers.Logout).Methods("POST")
	s.HandleFunc("/user/refresh", controllers.Refresh).Methods("POST")
	s.HandleFunc("/user/password", middlewares.Auth(controllers.ChangePassword)).Methods("POST")
	s.HandleFunc("/user", middlewares.Auth(controllers.GetUserByToken)).Methods("GET")
	s.HandleFunc("/users", middlewares.Auth(controllers.SearchUsers)).Methods("GET")

	s.HandleFunc("/images", middlewares.Auth(controllers.SaveImage)).Methods("POST")
	s.HandleFunc("/images", middlewares.Auth(controllers.GenerateImage)).Methods("PUT")
	s.HandleFunc("/images/tags", middlewares.Auth(controllers.GetTags)).Methods("GET")
	s.HandleFunc("/images/folder/{id}", middlewares.Auth(controllers.GetFolderImages)).Methods("GET")
	s.HandleFunc("/images/folder/{id}/archive", middlewares.Auth(controllers.DownloadFolder)).Methods("GET")
	s.HandleFunc("/images/share", middlewares.Auth(controllers.ShareImage)).Methods("POST")
	s.HandleFunc("/images/{id}/folder/{folderId}", middlewares.Auth(controllers.ImageToFolder)).Methods("PUT")
	s.HandleFunc("/images/{id}", middlewares.Auth(controllers.RemoveImage)).Methods("DELETE")
	s.HandleFunc("/images/{id}", middlewares.Auth(controllers.CreateAvatar)).Methods("PUT")
	s.HandleFunc("/images/{id}/data", middlewares.Auth(controllers.GetImage)).Methods("GET")
	s.HandleFunc("/images/{id}", middlewares.Auth(controllers.GetImageFile)).Methods("GET")
	s.HandleFunc("/images/{offset}/{limit}", middlewares.Auth(controllers.GetPaginatedImages)).Methods("GET")

	s.HandleFunc("/prints", middlewares.Auth(controllers.GetPrints)).Methods("GET")
	s.HandleFunc("/prints", middlewares.Auth(controllers.GeneratePrints)).Methods("POST")
	s.HandleFunc("/prints/{id}", middlewares.Auth(controllers.PrintToAvatar)).Methods("PUT")
	s.HandleFunc("/prints", middlewares.Auth(controllers.Clear)).Methods("DELETE")
	s.HandleFunc("/prints/archive", middlewares.Auth(controllers.DownloadArchive)).Methods("GET")

	s.HandleFunc("/folders/public", middlewares.Auth(controllers.GetPublic)).Methods("GET")
	s.HandleFunc("/folders", middlewares.Auth(controllers.GetFolders)).Methods("GET")
	s.HandleFunc("/folders", middlewares.Auth(controllers.CreateFolder)).Methods("POST")
	s.HandleFunc("/folders/{id}", middlewares.Auth(controllers.DeleteFolder)).Methods("DELETE")
	s.HandleFunc("/folders/{id}", middlewares.Auth(controllers.UpdateFolder)).Methods("PUT")
	s.HandleFunc("/folders/{id}", middlewares.Auth(controllers.GetFolder)).Methods("GET")

	s.HandleFunc("/admire-avatar/{emailHash}", controllers.GetImageByEmail).Methods("GET")

	s.PathPrefix("/files/temporary/").Handler(http.StripPrefix("/api/files/temporary/", http.FileServer(http.Dir("files/temporary"))))
	s.PathPrefix("/files/images/").Handler(http.StripPrefix("/api/files/images/", http.FileServer(http.Dir("files/images"))))

	s.HandleFunc("/prints/channel", middlewares.Auth(func(w http.ResponseWriter, r *middlewares.AuthorizedRequest) {
		ws.ServeWs(ws.PrintsPool, w, r)
	}))
	s.HandleFunc("/user/channel", middlewares.Auth(func(w http.ResponseWriter, r *middlewares.AuthorizedRequest) {
		ws.ServeWs(ws.NotificationsPool, w, r)
	}))

	if os.Getenv("MODE") == "production" {
		r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := filepath.Join("dist", r.URL.Path)

			_, err := os.Stat(path)
			if os.IsNotExist(err) {
				t := template.New("index.html")
				parser, _ := t.ParseFiles("dist/index.html")

				ogImage := "https://picart.admire.social/android-chrome-512x512.png"
				if strings.HasPrefix(r.URL.Path, "/images/") {
					ogImage = "https://picart.admire.social/api" + r.URL.Path + ".png"
				}
				err = parser.Execute(w, map[string]string{
					"og_image": ogImage,
				})
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			} else if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if strings.HasSuffix(r.URL.Path, ".js") {
				w.Header().Set("Content-Type", "application/javascript")
			}
			http.FileServer(http.Dir("dist")).ServeHTTP(w, r)
		})
	}

	return r
}
