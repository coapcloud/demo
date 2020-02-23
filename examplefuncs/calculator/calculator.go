package calculator

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/go-ocf/go-coap"
	"github.com/go-ocf/go-coap/codes"
)

func Run(port int) {
	dr := &CalculatorRouter{}

	mux := coap.NewServeMux()
	mux.Handle("", dr)

	fmt.Printf("handling calculator requests on %d\n", port)
	log.Fatal(coap.ListenAndServe("udp", fmt.Sprintf(":%d", port), mux))
}

type CalculatorRouter struct {
	Total int
}

func (r *CalculatorRouter) ServeCOAP(w coap.ResponseWriter, req *coap.Request) {
	log.Printf("Got message: from %v\n", req.Client.RemoteAddr())

	w.SetContentFormat(coap.AppJSON)

	var resp interface{}

	switch req.Msg.Code() {
	case codes.GET:
		resp = r.getHandler(w, req)
	case codes.POST, codes.PUT:
		resp = r.postHandler(w, req)
	case codes.DELETE:
		resp = r.deleteHandler(w, req)
	}

	b, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error encountered while trying to encode response: %v", err)
		return
	}

	_, err = w.Write(b)
	if err != nil {
		log.Printf("Error encountered while trying to send response: %v", err)
	}
}

func parseMessage(msg []byte) (int, error) {
	msgStr := string(msg)
	if msgStr[0] != '+' && msgStr[0] != '-' {
		return 0, fmt.Errorf("malformed delta: %s", string(msg))
	}

	// parse sign in delta
	var negative bool
	if msgStr[0] == '-' {
		negative = true
	}

	// parse number in delta
	deltaNum, err := strconv.Atoi(msgStr[1:])
	if err != nil {
		return 0, err
	}

	if negative {
		deltaNum = -deltaNum
	}

	return deltaNum, nil
}

func (r *CalculatorRouter) getHandler(w coap.ResponseWriter, req *coap.Request) interface{} {
	log.Printf("received GET: returning total of %d to the caller...\n", r.Total)

	return map[string]interface{}{
		"total": r.Total,
	}
}

func (r *CalculatorRouter) postHandler(w coap.ResponseWriter, req *coap.Request) interface{} {

	oldTotal := r.Total

	delta, err := parseMessage(req.Msg.Payload())
	if err != nil {
		w.SetCode(codes.BadRequest)

		return map[string]interface{}{
			"error": err,
		}
	}

	r.Total = oldTotal + delta

	log.Printf("received POST/PUT: changing total to %d\n", r.Total)

	return map[string]interface{}{
		"old_total": oldTotal,
		"delta":     delta,
		"total":     r.Total,
	}
}

func (r *CalculatorRouter) deleteHandler(w coap.ResponseWriter, req *coap.Request) interface{} {
	log.Println("received DELETE: clearing calculator total...")

	r.Total = 0

	return map[string]interface{}{
		"total": r.Total,
	}
}
