package routes

import (
	"go_API_server/controllers"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

//	func SetupRouter(redisClient *redis.Client) *mux.Router {
//		router := mux.NewRouter()
func SetupRouter(db *gorm.DB) *mux.Router {
	router := mux.NewRouter()

	// router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {

	// 	controllers.LoginHandler(redisClient, w, r)
	// }).Methods("POST")

	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		controllers.LoginHandler(db, w, r)
	}).Methods("POST")

	router.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		controllers.CreateItemHandler(db, w, r)
	}).Methods("POST")

	router.HandleFunc("/items/{id}", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetItemHandler(db, w, r)
	}).Methods("GET")

	router.HandleFunc("/items/{id}", func(w http.ResponseWriter, r *http.Request) {
		controllers.UpdateItemHandler(db, w, r)
	}).Methods("PUT")

	router.HandleFunc("/items/{id}", func(w http.ResponseWriter, r *http.Request) {
		controllers.DeleteItemHandler(db, w, r)
	}).Methods("DELETE")

	return router
}
