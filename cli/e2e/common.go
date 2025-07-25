package e2e

import (
	"bytes"
	"database/sql"
	"flag"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	_ "modernc.org/sqlite"
)

var update = flag.Bool("update", false, "Update golden files")

// runBinary will run the ../hyaline-e2e using _this_ directory as the working directory (i.e. cli/e2e)
func runBinary(args []string, t *testing.T) ([]byte, error) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("could not get current dir: %v", err)
	}

	binaryPath := filepath.Join(dir, "../hyaline-e2e")
	workingDir := dir
	t.Log("binaryPath", binaryPath)
	t.Log("workingDir", workingDir)

	cmd := exec.Command(binaryPath, args...)
	cmd.Env = append(os.Environ(), "GOCOVERDIR=../.coverdata/e2e")
	cmd.Dir = workingDir
	return cmd.CombinedOutput()
}

func updateGolden(goldenPath string, outputPath string, t *testing.T) {
	srcFile, err := os.Open(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	defer srcFile.Close()

	destFile, err := os.Create(goldenPath)
	if err != nil {
		t.Fatal(err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		t.Fatal(err)
	}

	err = destFile.Sync()
	if err != nil {
		t.Fatal(err)
	}
}

func compareDBs(path1 string, path2 string, t *testing.T) {
	// Get abs path for both dbs
	absPath1, err := filepath.Abs(path1)
	t.Log("absPath1", absPath1)
	if err != nil {
		t.Fatal(err)
	}
	absPath2, err := filepath.Abs(path2)
	t.Log("absPath2", absPath2)
	if err != nil {
		t.Fatal(err)
	}

	// Open DBs
	db1, err := sql.Open("sqlite", absPath1)
	if err != nil {
		t.Fatal(err)
	}
	defer db1.Close()
	db2, err := sql.Open("sqlite", absPath2)
	if err != nil {
		t.Fatal(err)
	}
	defer db2.Close()

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

func compareFiles(path1 string, path2 string, t *testing.T) {
	// Get abs path for both files
	absPath1, err := filepath.Abs(path1)
	t.Log("absPath1", absPath1)
	if err != nil {
		t.Fatal(err)
	}
	absPath2, err := filepath.Abs(path2)
	t.Log("absPath2", absPath2)
	if err != nil {
		t.Fatal(err)
	}

	// Read content of both files
	path1bytes, err := os.ReadFile(absPath1)
	if err != nil {
		t.Fatal(err)
	}
	path2bytes, err := os.ReadFile(absPath2)
	if err != nil {
		t.Fatal(err)
	}

	// Compare
	if !bytes.Equal(path1bytes, path2bytes) {
		t.Fatal("path1 and path2 do not have the same contents")
	}
}
