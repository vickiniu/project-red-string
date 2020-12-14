package main

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

var unmatchedNames = make(map[string]bool)

// filename of latest filing to update with
const latest = "CFB_20201116121431978.csv"

var headers = []string{
	"election", // election (skip)
	"office_cd",
	"recip_id",
	"can_class",
	"recipient_name",
	"committee",
	"filing",
	"schedule",
	"", // pageno
	"", // sequenceno
	"ref_no",
	"date",
	"", // refunddate
	"contributor_name",
	"c_code",
	"", // strno
	"", //strname
	"", //apartment
	"borough",
	"city",
	"state",
	"zip",
	"occupation",
	"employer_name",
	"", // empstrno
	"", // empstrname
	"", // empcity
	"", // empstate
	"amount",
	"", // matchamnt
	"", // prevamnt
	"", // pay_method
	"", // intermno
	"", // intermname
	"", // intstrno
	"", // instrnm
	"", // intaptno
	"", // intcity
	"", // intst
	"", // intzip
	"", // intempname
	"", // intempstno
	"", // intempstnm
	"", // intempcity
	"", // intempst
	"", // intoccupa
	"", // purposecd
	"", // exemptcd
	"", // adjtypecd
	"", // rr_ind
	"", // seg_ind
	"", // int_c_code
}

func main() {
	ctx := context.Background()
	dbURL := envString("DATABASE_URL", "postgres:///redstring?sslmode=disable")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error opening database connection: %v\n", err)
	}

	f, err := os.Open("csv/" + latest)
	if err != nil {
		log.Fatalf("error opening latest csv: %v", err)
	}
	r := csv.NewReader(f)
	// Skip the header
	_, err = r.Read()
	if err != nil {
		log.Fatalf("error reading header: %v", err)
	}
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error reading row from csv: %v", err)
		}
		err = handleRecord(ctx, db, record)
		if err != nil {
			log.Fatalf("error handling record: %v", err)
		}
	}
	// Write all unmatched names to a file
	var names []string
	for n := range unmatchedNames {
		names = append(names, n)
	}
	err = ioutil.WriteFile("unmatched_names.txt", []byte(strings.Join(names, "\n")), 0644)
	if err != nil {
		log.Println("unable to write file")
		log.Println(names)
	}
}

// cfbrecord defines a CFB record, fields corresponding to our
// DB columns.
type cfbrecord struct {
	election        string
	officeCD        string
	recipientID     string
	cfbRecipientID  string
	canClass        string
	recipientName   string
	committee       string
	filing          string
	schedule        string
	refNo           string
	date            time.Time
	contributorName string
	cCode           string
	borough         string
	city            string
	state           string
	zip             string
	occupation      string
	employerName    string
	amount          int
}

func handleRecord(ctx context.Context, db *sql.DB, record []string) error {
	c := cfbrecord{}
	for i, val := range record {
		if val == "" || headers[i] == "" {
			continue
		}
		switch headers[i] {
		case "election":
			c.election = val
		case "office_cd":
			c.officeCD = val
		case "recip_id":
			c.cfbRecipientID = val
		case "can_class":
			c.canClass = val
		case "recipient_name":
			c.recipientName = val
			var recipID string
			const q = `SELECT id FROM individuals WHERE cfb_name = $1`
			err := db.QueryRowContext(ctx, q, val).Scan(&recipID)
			if err != nil {
				// Try again, stripping the last space
				parts := strings.Split(val, " ")
				trimmedVal := strings.Join(parts[0:len(parts)-1], " ")
				const q = `SELECT id FROM individuals WHERE cfb_name = $1`
				err := db.QueryRowContext(ctx, q, trimmedVal).Scan(&recipID)
				if err != nil {
					unmatchedNames[val] = true
					fmt.Println(val)
					fmt.Println(trimmedVal)
				}
			}
			c.recipientID = recipID
		case "committee":
			c.committee = val
		case "filing":
			c.filing = val
		case "schedule":
			c.schedule = val
		case "ref_no":
			c.refNo = val
		case "date":
			d, err := time.Parse("1/2/2006", val)
			if err != nil {
				log.Fatalf("error parsing time: %v\n", err)
			}
			c.date = d
		case "contributor_name":
			c.contributorName = val
		case "c_code":
			c.cCode = val
		case "borough":
			c.borough = val
		case "city":
			c.city = val
		case "state":
			c.state = val
		case "zip":
			c.zip = val
		case "occupation":
			c.occupation = val
		case "employer_name":
			c.employerName = val
		case "amount":
			amt, err := strconv.ParseFloat(val, 64)
			if err != nil {
				log.Fatalf("error parsing amount: %v\n", err)
			}
			amtUnits := int(amt * 100)
			c.amount = amtUnits
		}
	}
	return upsertRecord(ctx, db, c)
}

func upsertRecord(ctx context.Context, db *sql.DB, c cfbrecord) error {
	if c.refNo == "" {
		return errors.New("record is missing a reference number")
	}
	var refno string
	const q = `SELECT id FROM contributions WHERE refno = $1`
	err := db.QueryRowContext(ctx, q, c.refNo).Scan(&refno)
	if err == sql.ErrNoRows {
		const insertQ = `
		INSERT INTO contributions (
			refno,
			amount,
			date,
			contributor_name,
			recipient_name,
			recipient_id,
			cfb_recipient_id,
			election,
			office_cd,
			can_class,
			committee,
			filing,
			schedule,
			c_code,
			borough,
			city,
			state,
			zip,
			occupation,
			employer_name
		) VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20
		)
	`
		_, err := db.ExecContext(
			ctx,
			insertQ,
			c.refNo,
			c.amount,
			c.date,
			c.contributorName,
			c.recipientName,
			c.recipientID,
			c.cfbRecipientID,
			c.election,
			c.officeCD,
			c.canClass,
			c.committee,
			c.filing,
			c.schedule,
			c.cCode,
			c.borough,
			c.city,
			c.state,
			c.zip,
			c.occupation,
			c.employerName,
		)
		if err != nil {
			return errors.Wrap(err, "inserting CFB record")
		}
	} else if err != nil {
		return errors.Wrap(err, "querying for record by refno")
	}
	return nil
}

// envString returns the value of the named environment variable.
// If name isn't in the environment os ir empty, it returns value.
func envString(name, value string) string {
	if s := os.Getenv(name); s != "" {
		value = s
	}
	return value
}
