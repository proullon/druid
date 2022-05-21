# druid

## Driver usage

Example:
```go
import (
	"database/sql"
	"time"

	_ "github.com/proullon/druid/driver"
)
	
func main() {
  db, err := sql.Open("druid", "https://druid.domain.com")
	if err != nil {
		panic("cannot open connection : %s\n", err)
	}

	rows, err := db.Query(`SELECT __time, added, channel, user FROM  "wikipedia" LIMIT 10`)
	if err != nil {
		panic("query: %s", err)
	}
	defer rows.Close()

	var __time time.Time
	var added int
	var channel, user string
	for rows.Next() {
		err = rows.Scan(&__time, &added, &channel, &user)
		if err != nil {
			panic("scan: %s", err)
		}
	}
}
```

## CLI testing

Running test.sql against your server after [Quickstart](https://druid.apache.org/docs/latest/tutorials/index.html):
```sh
DRUID_DSN="http://druid.company.com/" druidql < test.sql
```

will grant you this output:
```
druidsql>
__time  |  added   |  channel  |  cityName  |  comment  |  countryIsoCode  |  countryName  |  deleted  |  delta   |  isAnonymous  |  isMinor  |  isNew   |  isRobot  |  isUnpatrolled  |  metroCode  |  namespace  |  page    |  regionIsoCode  |  regionName  |  user
------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
2015-09-12T00:46:58.771Z  |  36      |  #en.wikipedia  |          |  added project  |          |          |  0       |  36      |  false   |  false   |  false   |  false   |  false   |          |  Talk    |  Talk:Oswald Tilghman  |          |          |  GELongstreet
2015-09-12T00:47:00.496Z  |  17      |  #ca.wikipedia  |          |  Robot inserta {{Commonscat}} que enllaça amb [[commons:category:Rallicula]]  |          |          |  0       |  17      |  false   |  true    |  false   |  true    |  false   |          |  Main    |  Rallicula  |          |          |  PereBot
2015-09-12T00:47:05.474Z  |  0       |  #en.wikipedia  |  Auburn  |  /* Status of peremptory norms under international law */ fixed spelling of 'Wimbledon'  |  AU      |  Australia  |  0       |  0       |  true    |  false   |  false   |  false   |  false   |          |  Main    |  Peremptory norm  |  NSW     |  New South Wales  |  60.225.66.142
2015-09-12T00:47:08.77Z  |  18      |  #vi.wikipedia  |          |  fix Lỗi CS1: ngày tháng  |          |          |  0       |  18      |  false   |  true    |  false   |  true    |  false   |          |  Main    |  Apamea abruzzorum  |          |          |  Cheers!-bot
2015-09-12T00:47:11.862Z  |  18      |  #vi.wikipedia  |          |  clean up using [[Project:AWB|AWB]]  |          |          |  0       |  18      |  false   |  false   |  false   |  true    |  false   |          |  Main    |  Atractus flammigerus  |          |          |  ThitxongkhoiAWB
2015-09-12T00:47:13.987Z  |  18      |  #vi.wikipedia  |          |  clean up using [[Project:AWB|AWB]]  |          |          |  0       |  18      |  false   |  false   |  false   |  true    |  false   |          |  Main    |  Agama mossambica  |          |          |  ThitxongkhoiAWB
2015-09-12T00:47:17.009Z  |  0       |  #ca.wikipedia  |          |  /* Imperi Austrohongarès */  |          |          |  20      |  -20     |  false   |  false   |  false   |  false   |  false   |          |  Main    |  Campanya dels Balcans (1914-1918)  |          |          |  Jaumellecha
2015-09-12T00:47:19.591Z  |  345     |  #en.wikipedia  |          |  adding comment on notability and possible COI  |          |          |  0       |  345     |  false   |  false   |  true    |  false   |  true    |          |  Talk    |  Talk:Dani Ploeger  |          |          |  New Media Theorist
2015-09-12T00:47:21.578Z  |  121     |  #en.wikipedia  |          |  Copying assessment table to wiki  |          |          |  0       |  121     |  false   |  false   |  false   |  true    |  false   |          |  User    |  User:WP 1.0 bot/Tables/Project/Pubs  |          |          |  WP 1.0 bot
2015-09-12T00:47:25.821Z  |  18      |  #vi.wikipedia  |          |  clean up using [[Project:AWB|AWB]]  |          |          |  0       |  18      |  false   |  false   |  false   |  true    |  false   |          |  Main    |  Agama persimilis  |          |          |  ThitxongkhoiAWB
druidsql>
HourTime  |  LinesDeleted
-------------------------
2015-09-12T00:00:00Z  |  1761
2015-09-12T01:00:00Z  |  16208
2015-09-12T02:00:00Z  |  14543
2015-09-12T03:00:00Z  |  13101
2015-09-12T04:00:00Z  |  12040
2015-09-12T05:00:00Z  |  6399
2015-09-12T06:00:00Z  |  9036
2015-09-12T07:00:00Z  |  11409
2015-09-12T08:00:00Z  |  11616
2015-09-12T09:00:00Z  |  17509
druidsql> exit
```
