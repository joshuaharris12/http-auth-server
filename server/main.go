package main 


import (
	"net/http" 
	"io/ioutil"
	"golang.org/x/time/rate"
	"fmt"
)

func main() {
	port := 8080
	http.HandleFunc("/authenticated", requiredAuthRequestHandler)
	fmt.Printf("Server Running on port: %d", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func requiredAuthRequestHandler(w http.ResponseWriter, r *http.Request) {
	limiter := rate.NewLimiter(100, 30)
	if !limiter.Allow() {
		w.WriteHeader(http.StatusTooManyRequests)
		return
	}

	username, password, ok := r.BasicAuth()
	if !isAuthenticated(username, password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !ok {
		w.Header().Add("WWW-Authenticate", `"Provide username and password"`)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message": "No basic auth present"}`))
		return
	}

	if r.Method == "GET" {
		w.Header().Add("Content-Type", "text/html")
		w.Write([]byte("<div>Hello, World</div>"))
		return
	}

	if r.Method == "POST" {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Header().Add("Content-Type", "text/html")
		w.Write([]byte(body))
		return
	}

}

func isAuthenticated(username string, password string) bool {
	return username == "SOME_USER_123" && password == "SOME_PASSWORD_123" 
}
