package main

import (
	"fmt"
	"os"
)

type Manifest struct {
	file     *os.File
	sstables []string
	counter  uint
}

func GetOrCreateManifest() (*Manifest, error) {
	file, err := GetOrCreateManifestFile()
	if err != nil {
		return nil, err
	}

	sstables, err := parseManifest(file)
	if err != nil {
		return nil, err
	}

	return &Manifest{
		file:     file,
		sstables: sstables,
		counter:  uint(len(sstables)), // should actually check the counter
	}, nil
}

func (m *Manifest) Flush(mt *MemTable) error {
	filename := fmt.Sprintf("sst-%d.json", m.counter+1)
	err := mt.Flush(filename)
	if err != nil {
		return err
	}
	m.counter++
	m.sstables = append(m.sstables, filename)
	m.file.WriteString(filename + "\n")
	return nil
}

// GetOrCreateManifestFile returns the MANIFEST file if it exists, otherwise creates it
func GetOrCreateManifestFile() (*os.File, error) {
	file, err := os.OpenFile("MANIFEST", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// parseManifest parses the MANIFEST file and returns a list of SSTable names.
// the MANIFEST file holds a table name per line. Each table name will follow
// the format `sst-[counter].json`.
func parseManifest(f *os.File) ([]string, error) {
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	fileSize := stat.Size()
	if fileSize == 0 {
		return []string{}, nil
	}

	data := make([]byte, fileSize)
	_, err = f.ReadAt(data, 0)
	if err != nil {
		return nil, err
	}

	var sstables []string
	var line []byte

	for _, b := range data {
		if b == '\n' {
			if len(line) > 0 {
				sstables = append(sstables, string(line))
				line = nil
			}
		} else {
			line = append(line, b)
		}
	}

	if len(line) > 0 {
		sstables = append(sstables, string(line))
	}

	return sstables, nil
}
