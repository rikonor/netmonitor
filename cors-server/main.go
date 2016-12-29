package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	neo4j "github.com/johnnadratowski/golang-neo4j-bolt-driver"
)

func main() {
	// setup neo4j
	driver := neo4j.NewDriver()
	conn, err := driver.OpenNeo("bolt://neo4j:1234@localhost:7687")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// setup persistor
	p := &neo4jPersistor{conn: conn}
	// p := &logPersistor{}

	// setup collector
	c := &Collector{port: "8080", persistor: p}
	c.Start()
}

type Collector struct {
	port      string
	persistor Persistor
}

func (c *Collector) HandleRequest(w http.ResponseWriter, r *http.Request) {
	var br BrowserRequest
	if err := json.NewDecoder(r.Body).Decode(&br); err != nil {
		log.Fatal(err)
	}
	r.Body.Close()

	if err := c.persistor.Persist(&br); err != nil {
		log.Fatal(err)
	}
}

func (c *Collector) Start() error {
	r := mux.NewRouter()
	r.HandleFunc("/", c.HandleRequest)

	return http.ListenAndServe(":"+c.port, handlers.CORS()(r))
}

// BrowserRequest is a request made by the browser
type BrowserRequest struct {
	SiteURL string `json:"siteUrl"`
	// Method is GET, POST, etc
	Method string `json:"method"`
	// Type can be (xmlhttprequest|script|image|stylesheet|main_frame|ping|sub_frame)
	Type string `json:"type"`
	URL  string `json:"url"`
}

type Persistor interface {
	Persist(*BrowserRequest) error
}

type logPersistor struct{}

func (p *logPersistor) Persist(br *BrowserRequest) error {
	siteURL, err := url.Parse(br.SiteURL)
	if err != nil {
		log.Fatal(err)
	}
	a := siteURL.Host

	reqURL, err := url.Parse(br.URL)
	if err != nil {
		log.Fatal(err)
	}
	b := reqURL.Host

	fmt.Printf("%s -> %s\n", a, b)
	return nil
}

type neo4jPersistor struct {
	conn neo4j.Conn
	mu   sync.Mutex
}

func (p *neo4jPersistor) Persist(br *BrowserRequest) error {
	// locking is necesasry because this neo4j lib is shit
	p.mu.Lock()
	defer p.mu.Unlock()

	// Get relationship data: (siteHost) -> (reqHost)
	siteURL, err := url.Parse(br.SiteURL)
	if err != nil {
		log.Fatal(err)
	}
	a := siteURL.Host

	reqURL, err := url.Parse(br.URL)
	if err != nil {
		log.Fatal(err)
	}
	b := reqURL.Host

	// Create nodes on neo4j
	query := `
		MERGE (a:Site {host: {hostA}})
		MERGE (b:Site {host: {hostB}})
		MERGE (a)-[:MADE_REQUEST_TO]->(b)
	`
	params := map[string]interface{}{"hostA": a, "hostB": b}

	_, err = p.conn.ExecNeo(query, params)
	return err
}
