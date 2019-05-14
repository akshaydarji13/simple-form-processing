package main

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func sayHelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //Parse url parameters passed, then parse the response packet for the POST body (request body)
	// attention: If you do not call ParseForm method, the following data can not be obtained form
	fmt.Println(r.Form) // print information on server side.
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello Akshay!") // write data to response
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Method: ", r.Method)
	if r.Method == "GET" {
		currentTime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(currentTime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(w, token)
	} else {
		r.ParseForm()
		token := r.Form.Get("token")
		if token != "" {
			// check token validity
			if len(r.Form["username"][0]) == 0 || len(r.Form["password"][0]) == 0 {
				fmt.Fprintf(w, "Username or Password empty!")
			} else {
				fmt.Println("Username: ", template.HTMLEscapeString(r.Form.Get("username")))
				fmt.Println("Password: ", template.HTMLEscapeString(r.Form.Get("password")))
			}
		} else {
			// give error if no token
			fmt.Fprintf(w, "Token not received!")
		}
	}
}

func main() {
	http.HandleFunc("/", sayHelloName)
	http.HandleFunc("/login", login)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServeError: ", err)
	}
}
