package helper

type DataTables struct {
	Data            []map[string]any `json:"data"`
	RecordsTotal    int              `json:"recordsTotal"`
	RecordsFiltered int              `json:"recordsFiltered"`
}

func BuildDatatTables(data []map[string]any, filtered int, total int) (dataTables DataTables) {
	if filtered == 0 {
		dataTables.Data = make([]map[string]any, 0)
	} else {
		dataTables.Data = data
	}

	dataTables.RecordsFiltered = filtered
	dataTables.RecordsTotal = total

	return dataTables
}
