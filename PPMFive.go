package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	"strconv"
	"io/ioutil"
	"github.com/hefju/PPMFive/models"
	"time"
)

// error response contains everything we need to use http.Error
type handlerError struct {
	Error   error
	Message string
	Code    int
}

// book model
type book struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Id     int    `json:"id"`
	Done   bool   `json:"done"`
}
// a custom type that we can use for handling errors and formatting responses
type handler func(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError)

// attach the standard ServeHTTP method to our handler so the http library can call it
func (fn handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// here we could do some prep work before calling the handler if we wanted to

	// call the actual handler
	response, err := fn(w, r)

	// check for errors
	if err != nil {
		log.Printf("ERROR: %v\n", err.Error)
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Message), err.Code)
		return
	}
	if response == nil {
		log.Printf("ERROR: response from method is nil\n")
		http.Error(w, "Internal server error. Check the logs.", http.StatusInternalServerError)
		return
	}

	// turn the response into JSON
	bytes, e := json.Marshal(response)
	if e != nil {
		http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
		return
	}

	// send the response and log
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
	log.Printf("%s %s %s %d", r.RemoteAddr, r.Method, r.URL, 200)
}
// list of all of the books
var books = make([]book, 0)

func main() {
	// command line flags
	port := flag.Int("port", 8000, "port to serve on")
	dir := flag.String("directory", "static/web/", "directory of web files")
	flag.Parse()

	// handle all requests by serving a file of the same name
	fs := http.Dir(*dir)
	fileHandler := http.FileServer(fs)

	// setup routes
	router := mux.NewRouter()
	router.Handle("/", http.RedirectHandler("/st/", 302))
	router.Handle("/books", handler(listBooks)).Methods("GET")
	router.Handle("/books", handler(addBook)).Methods("POST")
	router.Handle("/books/{id}", handler(getBook)).Methods("GET")
	router.Handle("/books/{id}", handler(updateBook)).Methods("POST")
	router.Handle("/books/{id}", handler(removeBook)).Methods("DELETE")
	router.PathPrefix("/st/").Handler(http.StripPrefix("/st/", fileHandler))//
//	router.Handle("/static/",fileHandler)
	http.Handle("/", router)

	log.Printf("Running on port %d\n", *port)

	addr := fmt.Sprintf("127.0.0.1:%d", *port)
	err := http.ListenAndServe(addr, nil)
	fmt.Println(err.Error())
}

func  addBook(w http.ResponseWriter,r *http.Request)(interface{},*handlerError){
	payload,e:=parseTaskRequest(r)
	if e!=nil{
		return nil,e
	}
	_,err:=models.AddTask(&payload)

	if err!=nil{
		log.Fatal(err)
	}
	return payload, nil
}
func  listBooks(w http.ResponseWriter, r *http.Request)(interface{}, *handlerError)  {
   list,_:=	models.GetTaskList(time.Now())
	return list, nil
}

func  getBook(w http.ResponseWriter, r *http.Request)(interface{}, *handlerError)  {
	parma:=mux.Vars(r)["id"]
	id,e:=strconv.Atoi(parma)
	if e!=nil{//handlerError
		return nil,&handlerError{e,"convert to int faile",http.StatusBadRequest}
	}
	b,err:=models.GetTaskByID(int64(id)) //getBookById(id)
	if err!=nil{
		return nil, &handlerError{e,"not find",http.StatusBadRequest}
	}
	return b,nil
}

func removeBook(w http.ResponseWriter,r *http.Request)(interface {},*handlerError) {
	log.Println("removeBook")
	parm := mux.Vars(r)["id"]
	log.Println("removeBook: parm=",parm)
	id, e := strconv.Atoi(parm)
	if e != nil {
		return nil, &handlerError{e, "id should be an int", http.StatusBadRequest}
	}
	_,err:= models.DeleteTask(int64(id))
	if err !=nil {
		return nil, &handlerError{nil, "no entity", http.StatusBadRequest}
	}
	return make(map[string]string), nil
}
func updateBook(w http.ResponseWriter, r *http.Request)(interface {},*handlerError){
	payload,e:=parseTaskRequest(r)
	if e!=nil{
		return nil,e
	}
	_,err:=models.UpdateTask(&payload)
	if err!=nil{
		log.Println("updateBook:",err)
	}
	return make(map[string]string), nil
}

func parseTaskRequest(r *http.Request)(models.TaskItem,*handlerError){
	data,e:=ioutil.ReadAll(r.Body)
	if e!=nil{
		return models.TaskItem{}, &handlerError{e,"could not read request",http.StatusBadRequest}
	}

	var payload models.TaskItem
	e=json.Unmarshal(data,&payload)
	if e!=nil{
		return models.TaskItem{}, &handlerError{e, "could not parse JSON", http.StatusBadRequest}
	}
	return payload,nil
}
