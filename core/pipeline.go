package core

type Pipeline struct {
	Importer Importer
	Exporter Exporter
	Filters  []Filter
}
