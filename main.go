package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/vipinvkartha/lenslocked/controllers"
	"github.com/vipinvkartha/lenslocked/templates"
	"github.com/vipinvkartha/lenslocked/views"
)

func executeTemplate(w http.ResponseWriter, filepath string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tpl, err := template.ParseFiles(filepath)
	if err != nil {
		log.Printf("processing template: %v", err)
		http.Error(w, "There was an error processing the template.", http.StatusInternalServerError)
		return
	}
	err = tpl.Execute(w, nil)
	if err != nil {
		log.Printf("executing template: %v", err)
		http.Error(w, "There was an error executing the template.", http.StatusInternalServerError)
		return
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, filepath.Join("templates", "home.gohtml"))
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, filepath.Join("templates", "contact.gohtml"))
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, filepath.Join("templates", "faq.gohtml"))
}

type Router struct{}

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		homeHandler(w, r)
	case "/contact":
		contactHandler(w, r)
	case "/faq":
		faqHandler(w, r)
	default:
		http.Error(w, "Page not found", http.StatusNotFound)
	}
}

func pathHandler(w http.ResponseWriter, r *http.Request) {
	// if r.URL.Path == "/" {
	// 	homeHandler(w, r)
	// } else if r.URL.Path == "/contact" {
	// 	contactHandler(w, r)
	// }

	switch r.URL.Path {
	case "/":
		homeHandler(w, r)
	case "/contact":
		contactHandler(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func main() {
	// http.HandleFunc("/", homeHandler)
	// http.HandleFunc("/contact", contactHandler)
	// http.HandleFunc("/", pathHandler)
	// fmt.Println("Startting the server")
	// http.ListenAndServe(":3000", nil)
	// var router Router
	// fmt.Println("Starting the server on :3000...")
	// http.ListenAndServe(":3000", router)

	// var router http.HandlerFunc
	// router = pathHandler
	// fmt.Println("Starting the server on :3000...")
	// http.ListenAndServe(":3000", router)

	// fmt.Println("Starting the server on :3000...")
	// http.ListenAndServe(":3000", http.HandlerFunc(pathHandler))

	// var router Router
	// fmt.Println("Starting the server on :3000...")
	// http.Handle("/", router)
	// http.ListenAndServe(":3000", nil)

	// r := chi.NewRouter()

	// r.Get("/", homeHandler)
	// r.Get("/contact", contactHandler)
	// r.Get("/faq", faqHandler)
	// r.NotFound(func(w http.ResponseWriter, r *http.Request) {
	// http.Error(w, "Page not found", http.StatusNotFound)
	// })
	// fmt.Println("Starting the server on :3000...")
	// http.ListenAndServe(":3000", r)

	// r := chi.NewRouter()

	// tpl := views.Must(views.Parse(filepath.Join("templates", "home.gohtml")))
	// r.Get("/", controllers.StaticHandler(tpl))
	// r.Get("/contact", controllers.StaticHandler(
	// 	views.Must(views.Parse(filepath.Join("templates", "contact.gohtml")))))
	// r.Get("/faq", controllers.StaticHandler(
	// 	views.Must(views.Parse(filepath.Join("templates", "faq.gohtml")))))

	// r.NotFound(func(w http.ResponseWriter, r *http.Request) {
	// 	http.Error(w, "Page not found", http.StatusNotFound)
	// })
	// fmt.Println("Starting the server on :3000...")
	// http.ListenAndServe(":3000", r)

	r := chi.NewRouter()

	r.Get("/", controllers.StaticHandler(views.Must(
		views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))))
	r.Get("/contact", controllers.StaticHandler(views.Must(
		views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))))
	r.Get("/faq", controllers.FAQ(
		views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))))
	r.Get("/signup", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml"))))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", r)

}

// func pathHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprint(w, r.URL.Path)
// }

// func main() {
// 	http.HandleFunc("/", pathHandler)
// 	fmt.Println("Starting the server on :3000...")
// 	http.ListenAndServe(":3000", nil)
// }
