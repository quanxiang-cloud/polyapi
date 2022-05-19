// to setup a service that check and set cookie

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// access url:
// http://127.0.0.1:9000

func writeJSON(w http.ResponseWriter, d interface{}) {
	b, err := json.Marshal(d)
	if err == nil {
		w.Write(b)
	}
}

func handleSession(w http.ResponseWriter, r *http.Request) {
	cookies := r.Cookies()
	cs := ""
	if len(cookies) > 0 {
		for _, v := range cookies {
			fmt.Printf("%#v\n", v)
			cs += v.String()
			cs += ";"
		}
	} else {
		for i := 1; i <= 2; i++ {
			c := &http.Cookie{
				Name: fmt.Sprintf("testcookie%d", i),
				Raw:  fmt.Sprintf("testraw%d", i),
				//Secure: true,
				//Domain: "expect-domain",
				Value: fmt.Sprintf("testcookie-value%d-%s", i, time.Now().Format("2006-01-02T15:04:05")),
			}
			fmt.Println("addcookie", c)
			http.SetCookie(w, c)
		}
	}
	s := fmt.Sprintf("cookie-%d {%s}", len(cookies), cs)
	fmt.Println(s)
	w.Write([]byte(s))
}

type req struct {
	ID int `json:"id"`
}

func handleA(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var req = &req{}
	if err := dec.Decode(req); err != nil {
		return
	}

	d := map[string]interface{}{
		"a": "aaa",
		"b": "bbb",
	}
	if req.ID > 0 {
		d = map[string]interface{}{
			"a": "a++",
			"b": fmt.Sprintf("id=%d", req.ID),
		}
	}
	writeJSON(w, d)
}

func handleB(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var req = &req{}
	if err := dec.Decode(req); err != nil {
		return
	}

	d := map[string]interface{}{
		"c": "ccc",
		"d": "ddd",
	}
	if req.ID > 0 {
		d = map[string]interface{}{
			"c": fmt.Sprintf("id=%d", req.ID),
			"d": "d++",
		}
	}
	writeJSON(w, d)
}

func main() {
	http.HandleFunc("/api/v1/testSession", handleSession)
	http.HandleFunc("/api/v1/test1", handleA)
	http.HandleFunc("/api/v1/test2", handleB)
	err := http.ListenAndServe("127.0.0.1:9000", nil)
	if err != nil {
		fmt.Printf("ListenAndServe fail:%v\n", err)
		return
	}
}
