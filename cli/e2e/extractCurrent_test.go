package e2e

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

var binaryName = "hyaline-coverage"

var binaryPath = ""

func TestMain(m *testing.M) {
	err := os.Chdir("..")
	if err != nil {
		fmt.Printf("could not change dir: %v", err)
		os.Exit(1)
	}

	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("could not get current dir: %v", err)
	}

	binaryPath = filepath.Join(dir, binaryName)

	os.Exit(m.Run())
}

func runBinary(args []string) ([]byte, error) {
	cmd := exec.Command(binaryPath, args...)
	cmd.Env = append(os.Environ(), "GOCOVERDIR=.coverdata")
	return cmd.CombinedOutput()
}

func TestExtractCurrent(t *testing.T) {
	output := fmt.Sprintf("./e2e/_output/extract-current-%d.db", time.Now().UnixMilli())
	args := []string{
		"extract", "current",
		"--config", "./e2e/_input/extract-current/config.yml",
		"--system", "my-app",
		"--output", output,
	}

	stdOutStdErr, err := runBinary(args)
	if err != nil {
		t.Log(string(stdOutStdErr))
		t.Fatal(err)
	}

	// Ensure golden file and output are the same
	expectedOutput, err := filepath.Abs("./e2e/_golden/extract-current.sqlite")
	if err != nil {
		t.Fatal(err)
	}
	actualOutput, err := filepath.Abs(output)
	if err != nil {
		t.Fatal(err)
	}
	compareDBs(expectedOutput, actualOutput, t)
}

func compareDBs(path1 string, path2 string, t *testing.T) {
	// Open DBs
	db1, err := sql.Open("sqlite", path1)
	if err != nil {
		t.Fatal(err)
	}
	db2, err := sql.Open("sqlite", path2)
	if err != nil {
		t.Fatal(err)
	}

	// Query and compare tables
	db1Tables := getTables(db1, t)
	t.Log(db1Tables)
	db2Tables := getTables(db2, t)
	t.Log(db2Tables)
	if !reflect.DeepEqual(db1Tables, db2Tables) {
		t.Fatal("db1 and db2 do not have the same tables")
	}

	// Compare the contents of each table
	for _, table := range db1Tables {
		db1Rows := getRows(table, db1, t)
		db2Rows := getRows(table, db2, t)
		if !reflect.DeepEqual(db1Rows, db2Rows) {
			t.Log(db1Rows)
			t.Log(db2Rows)
			t.Fatal("db1 and db2 do not have the same rows for table " + table)
		}
	}

}

func getTables(db *sql.DB, t *testing.T) []string {
	dbRows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' ORDER BY name")
	if err != nil {
		t.Fatal(err)
	}
	defer dbRows.Close()

	var rows []string
	for dbRows.Next() {
		var row string
		if err := dbRows.Scan(&row); err != nil {
			t.Fatal(err)
		}
		rows = append(rows, row)
	}
	if err = dbRows.Err(); err != nil {
		t.Fatal(err)
	}

	return rows
}

func getColumns(table string, db *sql.DB, t *testing.T) []string {
	dbRows, err := db.Query("SELECT name FROM pragma_table_info('" + table + "') ORDER BY name;")
	if err != nil {
		t.Fatal(err)
	}
	defer dbRows.Close()

	var rows []string
	for dbRows.Next() {
		var row string
		if err := dbRows.Scan(&row); err != nil {
			t.Fatal(err)
		}
		rows = append(rows, row)
	}
	if err = dbRows.Err(); err != nil {
		t.Fatal(err)
	}

	return rows
}

func getRows(table string, db *sql.DB, t *testing.T) [][]interface{} {
	columns := getColumns(table, db, t)
	numColumns := len(columns)

	// Order by the columns so we can make this as deterministic as possible
	dbRows, err := db.Query("SELECT * FROM " + table + " ORDER BY " + strings.Join(columns, ","))
	if err != nil {
		t.Fatal(err)
	}
	defer dbRows.Close()

	var rows [][]interface{}
	for dbRows.Next() {
		row := make([]interface{}, numColumns)
		rowPtrs := make([]interface{}, numColumns)
		for i := range numColumns {
			rowPtrs[i] = &row[i]
		}
		if err := dbRows.Scan(rowPtrs...); err != nil {
			t.Fatal(err)
		}
		rows = append(rows, row)
	}
	if err = dbRows.Err(); err != nil {
		t.Fatal(err)
	}

	return rows
}
