package bddcommon

import "github.com/cucumber/godog"

// Rows turns a data table into one map per row, keyed by the header row:
//
//	| name | price |
//	| Mug  | 8.00  |    ->  []map[string]string{{"name": "Mug", "price": "8.00"}}
func Rows(table *godog.Table) []map[string]string {
	if table == nil || len(table.Rows) == 0 {
		return nil
	}
	header := table.Rows[0].Cells
	rows := make([]map[string]string, 0, len(table.Rows)-1)
	for _, r := range table.Rows[1:] {
		row := make(map[string]string, len(header))
		for i, cell := range r.Cells {
			row[header[i].Value] = cell.Value
		}
		rows = append(rows, row)
	}
	return rows
}
