package driver

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/proullon/ramsql/engine/log"
)

// Stmt implements the Statement interface of sql/driver
type Stmt struct {
	endpoint  string
	userAgent string
	client    *http.Client

	query    string
	numInput int
}

// DruidQuery JSON format on /druid/v2/sql
type DruidQuery struct {
	Query        string `json:"query"`
	ResultFormat string `json:"resultFormat"`
	Header       bool   `json:"header"`
	TypesHeader  bool   `json:"typesHeader"`
}

// NewQuery initialise a marshal ready DruidQuery struct
func NewQuery(q string) ([]byte, error) {

	qs := DruidQuery{
		Query:        q,
		ResultFormat: "array",
		Header:       true,
		TypesHeader:  true,
	}
	data, err := json.Marshal(qs)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func countArguments(query string) int {
	for id := 1; id > 0; id++ {
		sep := fmt.Sprintf("$%d", id)
		if strings.Count(query, sep) == 0 {
			return id - 1
		}
	}

	return -1
}

func prepareStatement(endpoint, userAgent, query string, client *http.Client) *Stmt {

	// Parse number of arguments here
	// Should handle Postgres ($*) format
	numInput := countArguments(query)

	// Create statement
	stmt := &Stmt{
		endpoint:  endpoint,
		userAgent: userAgent,
		client:    client,
		query:     query,
		numInput:  numInput,
	}

	return stmt
}

// Close closes the statement.
//
// As of Go 1.1, a Stmt will not be closed if it's in use
// by any queries.
func (s *Stmt) Close() error {
	return nil
}

// NumInput returns the number of placeholder parameters.
//
// If NumInput returns >= 0, the sql package will sanity check
// argument counts from callers and return errors to the caller
// before the statement's Exec or Query methods are called.
//
// NumInput may also return -1, if the driver doesn't know
// its number of placeholders. In that case, the sql package
// will not sanity check Exec or Query argument counts.
func (s *Stmt) NumInput() int {
	return s.numInput
}

// Exec executes a query that doesn't return rows, such
// as an INSERT or UPDATE.
func (s *Stmt) Exec(args []driver.Value) (r driver.Result, err error) {
	return nil, fmt.Errorf("not implemented")
}

/*
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("fatalf error: %s", r)
			return
		}
	}()

	if s.query == "" {
		return nil, fmt.Errorf("empty statement")
	}

	var finalQuery string

	// replace $* by arguments in query string
	finalQuery = replaceArguments(s.query, args)
	log.Info("Exec <%s>\n", finalQuery)

	// Send query to server
	data, err := NewQuery(finalQuery)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", s.endpoint, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("cannot prepare query (%s)", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", s.userAgent)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cannot query server (%s)", err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot query server (%s)", err)
	}

	// Get answer from server
	lastInsertedID, rowsAffected, err := s.conn.conn.ReadResult()
	if err != nil {
		return nil, err
	}

	// Create a driver.Result
	return newResult(lastInsertedID, rowsAffected), nil
}
*/

// Query executes a query that may return rows, such as a
// SELECT.
func (s *Stmt) Query(args []driver.Value) (r driver.Rows, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("fatalf error: %s", r)
			return
		}
	}()

	if s.query == "" {
		return nil, fmt.Errorf("empty statement")
	}

	finalQuery := replaceArguments(s.query, args)
	log.Info("Query < %s >\n", finalQuery)

	// Send query to server
	data, err := NewQuery(finalQuery)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", s.endpoint, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("cannot prepare query (%s)", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", s.userAgent)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cannot query server (%s)", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("cannot query: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return r, err
	}

	return parseResponse(body)
}

// replace $* by arguments in query string
func replaceArguments(query string, args []driver.Value) string {

	holder := regexp.MustCompile(`[^\$]\$[0-9]+`)
	replacedQuery := ""

	if strings.Count(query, "?") == len(args) {
		return replaceArgumentsODBC(query, args)
	}

	allloc := holder.FindAllIndex([]byte(query), -1)
	queryB := []byte(query)
	for i, loc := range allloc {
		match := queryB[loc[0]+1 : loc[1]]

		index, err := strconv.Atoi(string(match[1:]))
		if err != nil {
			log.Warning("Matched %s as a placeholder but cannot get index: %s\n", match, err)
			return query
		}

		var v string
		if args[index-1] == nil {
			v = "null"
		} else {
			v = fmt.Sprintf("$$%v$$", args[index-1])
		}
		if i == 0 {
			replacedQuery = fmt.Sprintf("%s%s%s", replacedQuery, string(queryB[:loc[0]+1]), v)
		} else {
			replacedQuery = fmt.Sprintf("%s%s%s", replacedQuery, string(queryB[allloc[i-1][1]:loc[0]+1]), v)
		}
	}
	// add remaining query
	replacedQuery = fmt.Sprintf("%s%s", replacedQuery, string(queryB[allloc[len(allloc)-1][1]:]))

	return replacedQuery
}

func replaceArgumentsODBC(query string, args []driver.Value) string {
	var finalQuery string

	queryParts := strings.Split(query, "?")
	finalQuery = queryParts[0]
	for i := range args {
		arg := fmt.Sprintf("%v", args[i])
		_, ok := args[i].(string)
		if ok && !strings.HasSuffix(query, "'") {
			arg = "$$" + arg + "$$"
		}
		finalQuery += arg
		finalQuery += queryParts[i+1]
	}

	return finalQuery
}
