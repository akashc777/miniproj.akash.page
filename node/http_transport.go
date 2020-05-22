package node

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type HTTPTransport struct {
	Address  string
	node     *Node
	listener net.Listener
}

func (t *HTTPTransport) Serve(node *Node) error {
	t.node = node

	httpListener, err := net.Listen("tcp", t.Address)
	if err != nil {
		log.Fatalf("FATAL: listen (%s) failed - %s", t.Address, err)
	}
	t.listener = httpListener

	go func() {
		log.Printf("[%s] starting HTTP server", t.node.ID)
		server := &http.Server{
			Handler: t,
		}
		err := server.Serve(httpListener)
		// theres no direct way to detect this error because it is not exposed
		if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
			log.Printf("ERROR: http.Serve() - %s", err)
		}
		time.Sleep(10 * time.Second)

		close(t.node.exitChan)
		log.Printf("[%s] exiting Serve()", t.node.ID)
	}()

	return nil
}

func (t *HTTPTransport) Close() error {
	return t.listener.Close()
}

func (t *HTTPTransport) String() string {
	return t.listener.Addr().String()
}

func (t *HTTPTransport) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//if req.Method != "POST" {
//		apiResponse(w, 405, nil)
//	}

	switch req.URL.Path {

	case "/frontApi":
		t.frontEnd(w, req)
	case "/ping":
		t.pingHandler(w, req)
	case "/request_vote":
		t.requestVoteHandler(w, req)
	case "/append_entries":
		t.appendEntriesHandler(w, req)
	case "/command":
		t.commandHandler(w, req)
	default:
		apiResponse(w, 404, nil)
	}
}

func apiResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	response, err := json.Marshal(data)
	if err != nil {
		response = []byte("INTERNAL_ERROR")
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(response)))
	w.WriteHeader(statusCode)
	w.Write(response)
}

func (t *HTTPTransport) pingHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Length", "2")
	io.WriteString(w, "OK")
}

func (t *HTTPTransport) requestVoteHandler(w http.ResponseWriter, req *http.Request) {
	var vr VoteRequest

	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		apiResponse(w, 500, nil)
		return
	}

	err = json.Unmarshal(data, &vr)
	if err != nil {
		apiResponse(w, 500, nil)
		return
	}

	resp, err := t.node.RequestVote(vr)
	if err != nil {
		apiResponse(w, 500, nil)
		return
	}

	apiResponse(w, 200, resp)
}

func (t *HTTPTransport) appendEntriesHandler(w http.ResponseWriter, req *http.Request) {
	var er EntryRequest

	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		apiResponse(w, 500, nil)
		return
	}

	err = json.Unmarshal(data, &er)
	if err != nil {
		apiResponse(w, 500, nil)
		return
	}

	resp, err := t.node.AppendEntries(er)
	if err != nil {
		apiResponse(w, 500, nil)
	}

	apiResponse(w, 200, resp)
}

func (t *HTTPTransport) frontEnd(w http.ResponseWriter, req *http.Request) {


	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	pwd, _ := os.Getwd()
	p := fmt.Sprint(pwd, "/logs/", t.node.Lastcommit, ".json")
	err, resp := Read(p)
	if err != nil {
		apiResponse(w, 500, nil)
	}

	response, err := json.Marshal(resp)
	if err != nil {
		apiResponse(w, 500, nil)
	}

	w.Write(response)
}

// TODO: split this out into peer transport and client transport
// move into client transport
func (t *HTTPTransport) commandHandler(w http.ResponseWriter, req *http.Request) {
	var cr CommandRequest

	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		apiResponse(w, 500, nil)
		return
	}

	err = json.Unmarshal(data, &cr)
	if err != nil {
		apiResponse(w, 500, nil)
		return
	}

	cr.ResponseChan = make(chan CommandResponse, 1)
	t.node.Command(cr)
	resp := <-cr.ResponseChan
	if !resp.Success {
		apiResponse(w, 500, nil)
	}

	apiResponse(w, 200, resp)
}

func (t *HTTPTransport) RequestVoteRPC(address string, voteRequest VoteRequest) (VoteResponse, error) {
	var vresp VoteResponse
	endpoint := fmt.Sprintf("http://%s/request_vote", address)
	log.Printf("[%s] RequestVoteRPC %+v to %s", t.node.ID, voteRequest, endpoint)
	data, err := apiRequest("POST", endpoint, voteRequest, 100*time.Millisecond)
	if err != nil {
		return VoteResponse{}, err
	}
	//term, _ := data.Get("term").Int64()
	//voteGranted, _ := data.Get("vote_granted").Bool()
	//vresp := VoteResponse{
	//	Term:        term,
	//	VoteGranted: voteGranted,
	//}
	err = json.Unmarshal(data, &vresp)
	return vresp, nil
}

func (t *HTTPTransport) AppendEntriesRPC(address string, entryRequest EntryRequest) (EntryResponse, error) {
	var er EntryResponse
	endpoint := fmt.Sprintf("http://%s/append_entries", address)
	log.Printf("[%s] AppendEntriesRPC %+v to %s", t.node.ID, entryRequest, endpoint)
	data, err := apiRequest("POST", endpoint, entryRequest, 3*time.Second)
	if err != nil {
		return EntryResponse{}, err
	}

	err = json.Unmarshal(data, &er)
	//fmt.Printf("\n\n%v\n\n", respdata)

	return er, nil
}
