package main

import (
        "crypto/tls"
        "crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type ContactServer struct {
	verbose bool
	path	string
	db      *gorm.DB
	server  *http.Server
}

func NewContactServer(verbose bool, dsn, path string) (*ContactServer, error) {
	s := &ContactServer{verbose: verbose, path: path}
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Migrate the schema
	db.AutoMigrate(&Contact{})

	s.db = db
	return s, nil
}

func (s *ContactServer) Serve(addr, certPath, keyPath, caPath string) error {
	mux := http.NewServeMux()
	mux.HandleFunc(s.path, func(w http.ResponseWriter, r *http.Request) {
		s.handleContacts(w, r)
	})

	s.server = &http.Server{Addr: addr, Handler: mux}
	if certPath != "" && keyPath != "" {
		fmt.Printf("Creating TLS config from cert %q, key %q, ca %q\n", certPath, keyPath, caPath)
        	cfg, err := newTLSConfig(certPath, keyPath, caPath)
		if err != nil {
			panic(err)
		}
        	s.server.TLSConfig = cfg
	}
	
	if s.server.TLSConfig == nil {
		fmt.Printf("Serving HTTP on %s\n", addr);
		return s.server.ListenAndServe()
	} else {
		fmt.Printf("Serving HTTPS on %s\n", addr);
		return s.server.ListenAndServeTLS("","")
	}

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

func (s *ContactServer) idFromPath(urlPath string) (int, error) {
	tail := strings.TrimPrefix(strings.TrimPrefix(urlPath, s.path), "/")
	id, err := strconv.Atoi(tail)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *ContactServer) handleContactsGet(w http.ResponseWriter, r *http.Request) {
	payload := make(map[string]interface{})

	if r.URL.Path == s.path {
		var contacts []Contact
		if err := s.db.Find(&contacts).Error; err != nil {
			s.writeError(w, err)
			return
		}
		payload["contacts"] = &contacts
		s.writePayload(w, payload)
		return
	}

	id, err := s.idFromPath(r.URL.Path)
	if err != nil {
		s.writeError(w, fmt.Errorf("Could not get ID from %q", r.URL.Path))
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
	id, err := s.idFromPath(r.URL.Path)
	if err != nil {
		s.writeError(w, fmt.Errorf("Could not get ID from %q", r.URL.Path))
		return
	}
	if err:= s.db.Delete(id).Error; err != nil {
		s.writeError(w, err)
	}
	return
}

func newTLSConfig(certPath, keyPath, caPath string) (*tls.Config, error) {
        cert, err := tls.LoadX509KeyPair(certPath, keyPath)
        if err != nil {
                return nil, fmt.Errorf("Could not load TLS cert: %s", err)
        }

        roots, err := loadRoots(caPath)
        if err != nil {
                return nil, err
        }
	
        return &tls.Config{Certificates: []tls.Certificate{cert}, RootCAs: roots}, nil
}

func loadRoots(caPath string) (*x509.CertPool, error) {
        if caPath == "" {
                return nil, nil
        }

        roots := x509.NewCertPool()
        pem, err := ioutil.ReadFile(caPath)
        if err != nil {
                return nil, fmt.Errorf("Error reading %s: %s", caPath, err)
        }
        ok := roots.AppendCertsFromPEM(pem)
        if !ok {
                return nil, fmt.Errorf("Could not read root certs: %s", err)
        }
        return roots, nil
}
