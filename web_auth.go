package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// AllowedUser - Allowed User
type AllowedUser struct {
	Name     string `json:"Name"`
	Email    string `json:"Email"`
	Password string `json:"Password"` // This is fucking aweful switch to Hash
}

// AuthWrapHACK - Aweful Auth Wrapper - Replace with Github OAuth
type AuthWrapHACK struct {
	requireLogin bool
	allowed      map[string]AllowedUser
	logPage      []byte
	router       *mux.Router

	banIP map[string]int
}

// CheckUserPass - Check User & Pass
func (aw *AuthWrapHACK) CheckUserPass(r *http.Request) *AllowedUser {
	usr, pass, ok := r.BasicAuth()
	if !ok {
		return nil
	}

	uObj, ok := aw.allowed[usr]
	if !ok {
		return nil
	}

	if uObj.Password != pass {
		return nil
	}

	return &uObj
}

func (aw *AuthWrapHACK) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// time.Sleep(time.Second * 5)

	if aw.requireLogin {
		usr := aw.CheckUserPass(r)

		if usr == nil {
			banCount, ok := aw.banIP[r.RemoteAddr]
			if !ok {
				banCount = 0
			}
			banCount++
			aw.banIP[r.RemoteAddr] = banCount

			if banCount > 10 {
				fmt.Println("Banned -- ", r.RemoteAddr)
				http.Error(w, "Coffee", 304)
				return
			}

			w.Header().Set("WWW-Authenticate", `Basic realm="Access to Test"`)
			http.Error(w, "Teapot", 401)
			return
		}

		r.SetBasicAuth(usr.Name, usr.Email)
	}

	aw.router.ServeHTTP(w, r)
}
