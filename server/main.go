package main 


import (
	"net/http" 
	"io/ioutil"
	"fmt"
	"golang.org/x/time/rate"
)

func main() {
	maxRate := 100 
	maxBurst := 30
	limiter := rate.NewLimiter(maxRate, maxBurst)
	http.HandleFunc("/authenticated", requiredAuthRequestHandler)

	


	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			queryParams := r.URL.Query()
			if len(queryParams) != 0 {
				for param := range queryParams {
					w.Write([]byte(param))
				}
			}

 			w.Header().Add("Content-Type", "text/html")
			w.Write([]byte("Hello, World"))
			
		}

		if r.Method == "POST" {
			defer r.Body.Close()
			body, _ := ioutil.ReadAll(r.Body)
			w.Header().Add("Content-Type", "text/html")
			w.Write([]byte(body))
		}
	})

	http.HandleFunc("/authenticated", func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			w.Header().Add("WWW-Authenticate", `Basic realm="Give username and password"`)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"message": "No basic auth present"}`))
			return
		}
		fmt.Println(username)
		fmt.Println(password)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Hello, %s", username)))
	})

	http.HandleFunc("/limited", func(w http.ResponseWriter, r *http.Request) {
		limiter := rate.NewLimiter(100, 30)

		if limiter.Allow() {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("That's good!!!!"))
		}
		w.WriteHeader(http.StatusTooManyRequests)
	})

	http.ListenAndServe(":8080", nil)
}

func requiredAuthRequestHandler(w http.ResponseWriter, r *http.Request) {
	
	username, password, ok := r.BasicAuth()

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
		w.Header().Add("Content-Type", "text/html")
		w.Write([]byte(body))
		return
	}

}

