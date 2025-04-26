package dict

import (
	"encoding/binary"
	"log"
	"os"
	"strings"
)

const (
	wordsFilename = "words.dat" // The file containing the words and their definitions
	indexFilename = "index.dat" // The file containing the index entries. It's a temporary file.
	dictFilename  = "dict.dat"  // The file containing the merged words and index. It's used by API to run queries.
)

type Dict struct {
	f     *os.File
	index map[string]IndexEntry
}

type IndexEntry struct {
	Word    string
	Offset  int64 // offset of the word in the file
	DefSize int16 // size of the definition
}

func NewDict() (*Dict, error) {
	// Check if the dictionary file exists
	if !Exists() {
		// build a new dictionary
		err := BuildNewDict()
		if err != nil {
			return nil, err
		}
	}

	// Open the dictionary file in read only mode
	f, err := os.OpenFile(dictFilename, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	// Read the index from the file
	index, err := readIndex(f)
	if err != nil {
		f.Close()
		return nil, err
	}

	d := &Dict{
		f:     f,
		index: index,
	}

	return d, nil
}

// Exists checks if the dictionary file - dict.dat exists
func Exists() bool {
	if _, err := os.Stat(dictFilename); os.IsNotExist(err) {
		return false
	}

	return true
}

// QueryWord queries the dictionary for a word and returns its definition
func (d *Dict) QueryWord(word string) (string, bool) {
	// Find the word in the index
	idxe, ok := d.index[word]
	if !ok {
		log.Printf("Error: '%s' word not found", word)
		return "", false
	}

	// Read the definition
	// <word>,<definition>
	def := make([]byte, idxe.DefSize)
	defOffset := idxe.Offset + int64(len(idxe.Word)+1) // +1 for the comma

	_, err := d.f.ReadAt(def, defOffset) // +1 for the comma
	if err != nil {
		log.Printf("error reading definition: %v", err)
		return "", false
	}

	return string(def), true
}

// Close closes the dictionary file
func (d *Dict) Close() {
	if d.f != nil {
		d.f.Close()
	}
}

// readIndex reads the serialized index entries from the dict file
// and returns a map of word to IndexEntry
func readIndex(f *os.File) (map[string]IndexEntry, error) {
	// Read the index size from the first 8 bytes
	var indexSize uint64
	err := binary.Read(f, binary.BigEndian, &indexSize)
	if err != nil {
		return nil, err
	}

	index := make(map[string]IndexEntry)

	// Create a buffer to required size to read all index entries at once
	buf := make([]byte, indexSize)
	_, err = f.Read(buf)
	if err != nil {
		return nil, err
	}

	// Since each index is separated by a newline, we can split the buffer into lines
	lines := strings.Split(string(buf), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, ":")
		if len(parts) != 3 {
			continue
		}

		word := parts[0]
		offset := int64(binary.BigEndian.Uint64([]byte(parts[1])))
		defSize := int16(binary.BigEndian.Uint16([]byte(parts[2])))

		index[word] = IndexEntry{
			Word:    word,
			Offset:  offset,
			DefSize: defSize,
		}
	}

	return index, nil
}
