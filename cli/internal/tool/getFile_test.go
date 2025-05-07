package tool_test

import (
	"database/sql"
	"hyaline/internal/sqlite"
	"hyaline/internal/tool"
	"testing"

	_ "modernc.org/sqlite"
)

func TestGetFile(t *testing.T) {
	// Create an in memory database and populate it
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	sqlite.CreateSchema(db)
	sqlite.InsertFile(sqlite.File{
		ID:       "file1",
		CodeID:   "app",
		SystemID: "system1",
		RawData:  "file1 contents",
	}, db)
	sqlite.InsertFile(sqlite.File{
		ID:       "file2",
		CodeID:   "app",
		SystemID: "system1",
		RawData:  "file2 contents",
	}, db)
	sqlite.InsertFile(sqlite.File{
		ID:       "file3",
		CodeID:   "app",
		SystemID: "system2",
		RawData:  "file3 contents",
	}, db)

	// Initialize our tool
	listFiles := tool.GetFile("system1", db)

	// Test(s)
	expected := `<file>
  <file_name>app/file1</file_name>
  <file_content>
file1 contents
  </file_content>
</file>
`

	tests := []struct {
		input  string
		stop   bool
		result string
		err    error
	}{
		{`{"name":"app/file1"}`, false, expected, nil},
		{`{"name":"app/file3"}`, false, "(File Not Found)", nil},
	}

	for _, tc := range tests {
		stop, result, err := listFiles.Callback(tc.input)
		if stop != tc.stop {
			t.Errorf("stop: got %t, expected %t", stop, tc.stop)
		}
		if result != tc.result {
			t.Errorf("result: got \n%s\n, expected \n%s\n", result, tc.result)
		}
		if err != tc.err {
			t.Errorf("err: got %v, expected %v", err, tc.err)
		}
	}
}
