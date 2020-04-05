package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/coapcloud/go-coap"
	"github.com/coapcloud/veetoo/router"
)

type apiHandler struct {
	router *router.Router
}

func (h apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "got %v from %v", r, r.RemoteAddr)

	type request struct {
		Verb string `json:"verb"`
		Path string `json:"path"`
		Name string `json:"name"`
		Lang string `json:"lang"`
		Src  string `json:"src"`
	}

	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "%s", err)
		return
	}
	r.Body.Close()

	coapVerb, err := translateVerbToCoapCode(req.Verb)
	if err != nil {
		if err != nil {
			log.Println(err)
			fmt.Fprintf(w, "%s", err)
			return
		}
	}

	switch req.Lang {
	case "go":
		node := h.router.HotRegisterRoute(coapVerb, req.Path, req.Name)
		h.router.RegisterFunc(node, req.Name, req.Src)
		fmt.Fprintf(w, "%s %s -> %s registered", req.Verb, req.Path, req.Name)
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ERROR: unsupported language: %s", req.Lang)
		return
	}
}

func translateVerbToCoapCode(httpVerb string) (coap.COAPCode, error) {
	switch httpVerb {
	case http.MethodGet:
		return coap.GET, nil
	case http.MethodPut:
		return coap.PUT, nil
	case http.MethodPost:
		return coap.POST, nil
	case http.MethodDelete:
		return coap.DELETE, nil
	default:
		return 0, fmt.Errorf("invalid coap verb: %s", httpVerb)
	}
}

func Start() {
	address := ":8091"

	h := apiHandler{}

	http.Handle("/api", h)

	log.Printf("API Server running on %v\n", address)
	err := http.ListenAndServe(":8091", nil)
	if err != nil {
		log.Fatal(err)
	}
}
