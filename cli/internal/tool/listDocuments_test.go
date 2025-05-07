package tool_test

import (
	"database/sql"
	"hyaline/internal/sqlite"
	"hyaline/internal/tool"
	"testing"

	_ "modernc.org/sqlite"
)

func TestListDocuments(t *testing.T) {
	// Create an in memory database and populate it
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	sqlite.CreateSchema(db)
	sqlite.InsertDocument(sqlite.Document{
		ID:              "doc1",
		DocumentationID: "app",
		SystemID:        "system1",
		Type:            "md",
		RawData:         "doc1 raw data",
		ExtractedData:   "doc1 extracted data",
	}, db)
	sqlite.InsertDocument(sqlite.Document{
		ID:              "doc2",
		DocumentationID: "app",
		SystemID:        "system1",
		Type:            "md",
		RawData:         "doc2 raw data",
		ExtractedData:   "doc2 extracted data",
	}, db)
	sqlite.InsertDocument(sqlite.Document{
		ID:              "doc3",
		DocumentationID: "app",
		SystemID:        "system2",
		Type:            "md",
		RawData:         "doc3 raw data",
		ExtractedData:   "doc3 extracted data",
	}, db)

	listDocuments := tool.ListDocuments("system1", db)

	expected := `<documents>
  <document name="app/doc1">
  <document name="app/doc2">
</documents>
`

	tests := []struct {
		stop   bool
		result string
		err    error
	}{
		{false, expected, nil},
	}

	for _, tc := range tests {
		stop, result, err := listDocuments.Callback("")
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
