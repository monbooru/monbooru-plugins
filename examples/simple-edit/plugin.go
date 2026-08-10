// Common mechanisms you can copy for any monbooru plugin: pairing, the
// credentials a pairing leaves behind, calling the REST API, and checking that
// an inbound request really is monbooru. 

package main

import (
	"bytes"
	"cmp"
	"crypto/rand"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"
)

// button is one entry monbooru renders. The list is part of what the
// user approves at pairing, so changing it means pairing again.
type button struct {
	Slot  string `json:"slot"`  // detail-actions or batch-bar
	Label string `json:"label"` // what it reads
	Mode  string `json:"mode"`  // relay (monbooru posts the scope) or open (a page of the plugin)
	Path  string `json:"path"`  // where on the plugin the click lands
	Media string `json:"media"` // image, archive, animated - blank for every medium
}

var (
	monbooruURL string // where monbooru answers
	selfURL     string // where monbooru calls us back
	stateDir    string // where the credentials are kept
)

// credentials are the two halves of a pairing: the API token monbooru issued
// us, and the secret it presents on everything it sends back.
type credentials struct {
	Token string `json:"token"`
	Peer  string `json:"peer"`
}

// Pairing lands while the server is already answering, so both halves travel in one swap.
var paired atomic.Pointer[credentials]

func creds() credentials {
	if c := paired.Load(); c != nil {
		return *c
	}
	return credentials{}
}

// run serves mux, offers to pair, and returns when monbooru stops the plugin.
func run(mux *http.ServeMux) {
	addr := flag.String("addr", "127.0.0.1:7301", "address to serve on")
	monbooru := flag.String("monbooru", cmp.Or(os.Getenv("MONBOORU_URL"), "http://127.0.0.1:8080"), "monbooru base url")
	advertised := flag.String("url", "", "address monbooru should call us at (default http://<addr>)")
	state := flag.String("state", ".", "directory holding the pairing credentials")
	flag.Parse()

	log.SetFlags(0)
	monbooruURL = strings.TrimRight(*monbooru, "/")
	selfURL = cmp.Or(*advertised, "http://"+*addr)
	stateDir = *state

	// monbooru probes this before the user can approve, and every 30 s
	// after to decide whether to render our buttons. 
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		sendJSON(w, http.StatusOK, map[string]string{"version": version})
	})
	// The teardown channel: monbooru calls it when the user removes the
	// pairing, so one click ends the pairing on both sides.
	mux.HandleFunc("POST /api/v1/pair/remove", fromMonbooru(unpair))

	srv := &http.Server{Addr: *addr, Handler: mux}
	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()
	log.Printf("%s %s listening on %s", appName, version, *addr)

	// Pairing waits on the user, so it runs beside the server rather than before it.
	go func() {
		if err := pair(); err != nil {
			log.Fatal(err)
		}
	}()

	<-stopped()
	log.Print("shutting down")
	_ = srv.Close()
}

// stopped closes its channel when monbooru's launcher closes our stdin.
func stopped() <-chan struct{} {
	stop := make(chan struct{})
	if info, err := os.Stdin.Stat(); err == nil && info.Mode()&os.ModeNamedPipe != 0 {
		go func() {
			_, _ = io.Copy(io.Discard, os.Stdin)
			close(stop)
		}()
	}
	return stop
}

// ---- pairing -----------------------------------------------------------

// pairing guards the offer loop. Startup makes one offer and a teardown makes another
var pairing atomic.Bool

// pair fills in the credentials: the stored ones when the plugin has run here
// before and monbooru still accepts them, else a fresh offer waiting on
// approval under Settings > Plugins. An offer nobody answers ages out, so it
// is made again.
func pair() error {
	if !pairing.CompareAndSwap(false, true) {
		return nil
	}
	defer pairing.Store(false)
	if c := storedCredentials(); c.Token != "" && c.Peer != "" {
		paired.Store(&c)
		if tokenAccepted() {
			return nil
		}
		// The pairing went while we were not running: removed by hand in
		// monbooru.toml, or a monbooru that was reinstalled. 
		log.Print("monbooru no longer accepts the stored credentials; offering again")
		forgetCredentials()
	}
	for {
		peer := rand.Text() // the secret monbooru will present back to us
		token, err := offer(peer)
		if err != nil {
			return err
		}
		if token == "" {
			log.Print("the pairing offer expired, offering again")
			continue
		}
		c := credentials{Token: token, Peer: peer}
		if err := storeCredentials(c); err != nil {
			return err
		}
		paired.Store(&c)
		log.Print("paired")
		return nil
	}
}

// offer makes one pairing offer and polls until monbooru answers it. An empty
// token and no error means the offer aged out before anyone approved it.
func offer(peer string) (string, error) {
	var offered struct {
		RequestID string `json:"request_id"`
	}

	if err := postJSON(api("/pair/request"), map[string]any{
		"app":              appName,
		"url":              selfURL,
		"requested_scopes": []string{"read", "write"},
		"peer_token":       peer,
		"version":          version,
		"buttons":          buttons,
	}, &offered); err != nil {
		return "", err
	}
	log.Print("waiting for approval in monbooru: Settings > Plugins")

	for {
		reply, code, err := call("GET", api("/pair/status?id=%s", offered.RequestID), "", nil)
		if code == http.StatusNotFound {
			return "", nil
		}
		if err != nil {
			return "", err
		}
		var answer struct {
			Token  string `json:"token"`
			Status string `json:"status"`
		}
		if err := json.Unmarshal(reply, &answer); err != nil {
			return "", err
		}
		switch {
		case answer.Token != "":
			return answer.Token, nil
		case answer.Status == "denied":
			return "", errors.New("the user denied the pairing")
		}
		time.Sleep(2 * time.Second)
	}
}

// tokenAccepted asks monbooru whether the token we hold is still one of its own. 
func tokenAccepted() bool {
	_, code, _ := call("GET", api("/galleries"), "", nil)
	return code != http.StatusUnauthorized && code != http.StatusServiceUnavailable
}

func credentialsPath() string { return filepath.Join(stateDir, "credentials.json") }

func storedCredentials() credentials {
	var c credentials
	if raw, err := os.ReadFile(credentialsPath()); err == nil {
		_ = json.Unmarshal(raw, &c)
	}
	return c
}

func storeCredentials(c credentials) error {
	raw, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(credentialsPath(), raw, 0o600)
}

// forgetCredentials drops both halves, in memory and on disk.
func forgetCredentials() {
	paired.Store(nil)
	if err := os.Remove(credentialsPath()); err != nil && !os.IsNotExist(err) {
		log.Printf("dropping %s: %v", credentialsPath(), err)
	}
}

// unpair answers monbooru's teardown. The user removed the pairing there,
// which revoked the token we hold, so the copy on our disk authenticates to
// nothing and the only thing left to do is stop keeping it. 
func unpair(w http.ResponseWriter, _ *http.Request) {
	forgetCredentials()
	log.Print("the pairing was removed in monbooru; offering again")
	sendJSON(w, http.StatusOK, map[string]string{"status": "removed"})
	go func() {
		if err := pair(); err != nil {0
			log.Printf("offering again after the teardown: %v", err)
		}
	}()
}

// ---- talking to monbooru -----------------------------------------------

var client = &http.Client{Timeout: 60 * time.Second}

// api builds a REST API address: api("/images/%d/file", id).
func api(path string, args ...any) string {
	return monbooruURL + "/api/v1" + fmt.Sprintf(path, args...)
}

// call runs one exchange, carrying the token we were issued once pairing has
// landed (the pairing endpoints take none). It answers with the reply body,
// the status code, and an error for anything monbooru refused.
func call(method, url, contentType string, body io.Reader) ([]byte, int, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, 0, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	if token := creds().Token; token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err == nil && resp.StatusCode >= 400 {
		err = fmt.Errorf("%s %s: %s: %s", method, url, resp.Status, strings.TrimSpace(string(raw)))
	}
	return raw, resp.StatusCode, err
}

// postJSON posts a JSON body and decodes the reply into out, which may be nil.
func postJSON(url string, body, out any) error {
	raw, err := json.Marshal(body)
	if err != nil {
		return err
	}
	reply, _, err := call("POST", url, "application/json", bytes.NewReader(raw))
	if err != nil || out == nil {
		return err
	}
	return json.Unmarshal(reply, out)
}

// ---- answering monbooru ------------------------------------------------

// fromMonbooru refuses anything not carrying the secret we minted at pairing.
// A relay click and a request for one of our own pages both present it, so
// every route but /health is wrapped in this.
func fromMonbooru(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		presented, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
		peer := creds().Peer
		if !ok || peer == "" || subtle.ConstantTimeCompare([]byte(presented), []byte(peer)) != 1 {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		h(w, r)
	}
}

func sendJSON(w http.ResponseWriter, code int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(body)
}
