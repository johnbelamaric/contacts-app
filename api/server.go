package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type ContactServer struct {
	verbose bool
	db      *gorm.DB
	server  *http.ServeMux
}

func NewContactServer(verbose bool, dsn string) (*ContactServer, error) {
	s := &ContactServer{verbose: verbose}
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Migrate the schema
	db.AutoMigrate(&Contact{})

	s.db = db
	return s, nil
}

func (s *ContactServer) Serve(addr string) error {

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s.server = http.NewServeMux()
	s.server.HandleFunc("/contacts", func(w http.ResponseWriter, r *http.Request) {
		s.handleContacts(w, r)
	})

	return http.Serve(ln, s.server)
}

func (s *ContactServer) writeError(w http.ResponseWriter, err error) {
	s.writeErrorStatus(w, err, http.StatusBadRequest)
}

func (s *ContactServer) writeErrorStatus(w http.ResponseWriter, err error, status int) {
        j, _ := json.Marshal(err.Error())
	w.WriteHeader(status)
	io.WriteString(w, `{"error":`)
	w.Write(j)
	io.WriteString(w, `}\n`)
}

func (s *ContactServer) writePayload(w http.ResponseWriter, payload interface{}) {
        j, err := json.Marshal(payload)
	if err != nil {
		s.writeErrorStatus(w, err, http.StatusInternalServerError)
		return
	}
	w.Write(j)
}

func (s *ContactServer) handleContacts(w http.ResponseWriter, r *http.Request) {
	switch m := r.Method; m {
	case http.MethodGet:
		s.handleContactsGet(w,r)
	case http.MethodPost:
		s.handleContactsPost(w,r)
	case http.MethodPut:
		s.handleContactsPut(w,r)
	case http.MethodDelete:
		s.handleContactsDelete(w,r)
	default:
		s.writeError(w, fmt.Errorf("Unhandled method %q", m))
	}
}

func (s *ContactServer) handleContactsGet(w http.ResponseWriter, r *http.Request) {
	payload := make(map[string]interface{})

	if r.URL.Path == "/contacts" {
		var contacts []Contact
		if err := s.db.Find(&contacts).Error; err != nil {
			s.writeError(w, err)
			return
		}
		payload["contacts"] = &contacts
		s.writePayload(w, payload)
		return
	}

	tail := strings.TrimPrefix(r.URL.Path, "/contacts/")
	id, err := strconv.Atoi(tail)
	if err != nil {
		s.writeError(w, fmt.Errorf("Invalid ID %q", tail))
		return
	}
	
	var contact Contact
	if err := s.db.Find(&contact, id).Error; err != nil {
		s.writeError(w, err)
		return
	}
	payload["contacts"] = &[]Contact{contact}
	s.writePayload(w, payload)
}

func (s *ContactServer) handleContactsPost(w http.ResponseWriter, r *http.Request) {
	s.handleContactsPut(w,r)
}

func (s *ContactServer) handleContactsPut(w http.ResponseWriter, r *http.Request) {
	var c Contact
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.writeError(w, err)
		return
	}

	err = json.Unmarshal(b, &c)
	if err != nil {
		s.writeError(w, err)
		return
	}

	if err:= s.db.Create(&c).Error; err != nil {
		s.writeError(w, err)
	}
	return
}

func (s *ContactServer) handleContactsDelete(w http.ResponseWriter, r *http.Request) {
	s.writeError(w, fmt.Errorf("Not yet implemented"))
}
