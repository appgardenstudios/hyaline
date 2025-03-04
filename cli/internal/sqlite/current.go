package sqlite

import "database/sql"

func CreateCurrentSchema(db *sql.DB) (err error) {
	_, err = db.Exec(`
create table system(id);
create table code(id, system_id, path);
create table file(id, code_id, system_id, relative_path, raw_data);
create table documentation(id, system_id, type, path);
create table document(id, documentation_id, system_id, relative_path, format, raw_data, extracted_text);
create table section(id, document_id, documentation_id, system_id, parent_section_id, section_order, title, format, raw_data, extracted_text);
`)

	return
}
