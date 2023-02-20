package druid

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_ "github.com/proullon/druid/driver"
)

var expectedJSON string = `[["__time","added","channel","user"],["LONG","LONG","STRING","STRING"],["2015-09-12T00:46:58.771Z",36,"#en.wikipedia","GELongstreet"],["2015-09-12T00:47:00.496Z",17,"#ca.wikipedia","PereBot"],["2015-09-12T00:47:08.770Z",18,"#vi.wikipedia","Cheers!-bot"],["2015-09-12T00:47:11.862Z",18,"#vi.wikipedia","ThitxongkhoiAWB"],["2015-09-12T00:47:13.987Z",18,"#vi.wikipedia","ThitxongkhoiAWB"],["2015-09-12T00:47:19.591Z",345,"#en.wikipedia","New Media Theorist"],["2015-09-12T00:47:21.578Z",121,"#en.wikipedia","WP 1.0 bot"],["2015-09-12T00:47:25.821Z",18,"#vi.wikipedia","ThitxongkhoiAWB"],["2015-09-12T00:47:33.004Z",18,"#vi.wikipedia","ThitxongkhoiAWB"],["2015-09-12T00:47:35.776Z",36,"#vi.wikipedia","TuHan-Bot"]]`

func TestQuery(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, expectedJSON)
	}))
	defer ts.Close()

	db, err := sql.Open("druid", ts.URL)
	if err != nil {
		t.Fatalf("cannot open connection : %s\n", err)
	}

	rows, err := db.Query(`SELECT __time, added, channel, user FROM  "wikipedia" LIMIT 10`)
	if err != nil {
		t.Fatalf("query: %s", err)
	}
	defer rows.Close()

	var __time time.Time
	var added int
	var channel, user string
	for rows.Next() {
		err = rows.Scan(&__time, &added, &channel, &user)
		if err != nil {
			t.Fatalf("scan: %s", err)
		}
		if __time.IsZero() {
			t.Errorf("time.Time: is zero")
		}
		if added == 0 {
			t.Errorf("int: is zero")
		}
		if channel == "" {
			t.Errorf("string: is empty")
		}
	}
}
