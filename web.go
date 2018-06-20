package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/acme/autocert"
)

const (
	cookieID = "my-token-check"
	email    = "evilkimau@gmail.com"
)

var (
	webServ     *http.Server
	serverParam *Server
)

// Server - Server Details
type Server struct {
	Addr            string        `json:"Address"`
	UseSSL          bool          `json:"Use SSL"`
	UseAutoCert     bool          `json:"Use AutoCert"`
	ServerCert      string        `json:"Server Cert"`
	ServerKey       string        `json:"Server Key"`
	RequireLogin    bool          `json:"Use Login"`
	AuthUsers       []AllowedUser `json:"Allowed Users"`
	WhitelistDomain []string      `json:"Domain Whitelist"`
}

func launchWeb(configFile string) *mux.Router {
	// Loading Config
	configData, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(configData, &serverParam)
	if err != nil {
		panic(err)
	}

	// Launch Server
	log.Println("--- Launch Website ---")
	router := mux.NewRouter()

	router.PathPrefix("/{filepath:[0-9A-Za-z_/]+\\.[[:word:]]+}").Handler(
		http.FileServer(http.Dir("./static")))

	aw := AuthWrapHACK{
		requireLogin: serverParam.RequireLogin,
		allowed:      make(map[string]AllowedUser),
		router:       router,
		banIP:        make(map[string]int),
	}

	for _, u := range serverParam.AuthUsers {
		aw.allowed[u.Email] = u
	}

	webServ = &http.Server{
		Handler:      &aw,
		Addr:         serverParam.Addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Listening: ", serverParam.Addr)

	go func() {
		var err error
		if serverParam.UseSSL && serverParam.UseAutoCert {

			certManager := autocert.Manager{
				Prompt:     autocert.AcceptTOS,
				HostPolicy: autocert.HostWhitelist(serverParam.WhitelistDomain...),
				Cache:      autocert.DirCache("./.certs"),
				Email:      email,
			}

			go http.ListenAndServe(":80", certManager.HTTPHandler(nil))

			webServ.TLSConfig = &tls.Config{
				GetCertificate: certManager.GetCertificate,
			}
			err = webServ.ListenAndServeTLS("", "")

		} else if serverParam.UseSSL {
			cert := path.Clean(serverParam.ServerCert)
			key := path.Clean(serverParam.ServerKey)

			err = webServ.ListenAndServeTLS(cert, key)
		} else {
			err = webServ.ListenAndServe()
		}
		if err != nil {
			panic(err)
		}
		log.Println("------------- Server Closed --------------")
	}()

	return router
}

// ReadReqBody - Read the Body of a Request how the fuck is this not a default function
func ReadReqBody(r *http.Request) (string, error) {
	var bText string

	bTextBuffer := bytes.Buffer{}
	_, err := io.Copy(&bTextBuffer, r.Body)
	if err != nil {
		return "", err
	}

	bText = bTextBuffer.String()

	err = r.Body.Close()
	if err != nil {
		return "", err
	}

	return bText, nil
}

func corsMe(w http.ResponseWriter, req *http.Request) {
	corsAllow := ""
	for _, h := range serverParam.WhitelistDomain {
		if req.Host == h {
			corsAllow = "https://" + h
		}
	}

	if corsAllow == "" {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", corsAllow)
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X-Auth-Token, channel-id,client-id,user-id,x-extension-jwt")
}
