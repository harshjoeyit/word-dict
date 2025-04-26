package dict

// This file contains the code to build a new dictionary
// from the words.dat and index.dat files.

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// BuildNewDict creates a new dict.data file using the
// words.dat and index.dat files
func BuildNewDict() error {
	// Open words.dat file for reading
	wordsFile, err := os.OpenFile(wordsFilename, os.O_RDONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening words.dat: %v", err)
	}
	defer wordsFile.Close()

	// Build index

	// Read line by line
	scanner := bufio.NewScanner(wordsFile)

	// Offset in words.dat file
	var currOffset int64 = 0

	var indexEntries []IndexEntry

	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.Split(line, ",")

		word := parts[0]
		definition := parts[1]

		// Create an index entry
		idxe := IndexEntry{
			Word:    word,
			Offset:  currOffset,
			DefSize: int16(len(definition)),
		}

		indexEntries = append(indexEntries, idxe)

		// Update offset (current position + length of line + newline character)
		currOffset += int64(len(line)) + 1 // +1 for newline character
	}

	// flush the index to index.dat file

	err = flushIndex(indexEntries)
	if err != nil {
		return fmt.Errorf("error flushing index: %v", err)
	}

	err = mergeFiles()
	if err != nil {
		return fmt.Errorf("error merging files: %v", err)
	}

	return nil
}

// flushIndex serializes the index entries and flushes to index.dat file
func flushIndex(indexEntries []IndexEntry) error {
	// Clean up the old index file
	os.Remove(indexFilename)

	indexFile, err := os.OpenFile(indexFilename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer indexFile.Close()

	// Calculate the size of the index
	totalIndexSize := calcIndexSize(indexEntries)

	log.Println("Total index size:", totalIndexSize)

	var separator string = ":"

	var buf bytes.Buffer

	for _, idxe := range indexEntries {
		// In order to create dict file we have prepend indexFile to wordsFile
		// hence the offset for each word in index file would be shifted by indexSize bytes
		// so we need to update the offset of each word before flushing the index
		idxe.Offset += int64(totalIndexSize)

		// Write the word
		buf.WriteString(idxe.Word)
		buf.WriteString(separator)

		// Write the offset
		// buf.Write([]byte(fmt.Sprintf("%d", idxe.Offset))
		binary.Write(&buf, binary.BigEndian, idxe.Offset)
		buf.WriteString(separator)

		// Write the definition size
		// buf.Write([]byte(fmt.Sprintf("%d", idxe.DefSize)))
		binary.Write(&buf, binary.BigEndian, idxe.DefSize)
		buf.WriteString("\n")
	}

	// Write constant sized index header to the beginning of the file
	indexHeaderBuf := make([]byte, 8)
	binary.BigEndian.PutUint64(indexHeaderBuf, uint64(totalIndexSize))

	_, err = indexFile.WriteAt(indexHeaderBuf, 0)
	if err != nil {
		return err
	}

	// Write serialized index entries to the file
	_, err = indexFile.WriteAt(buf.Bytes(), 8)
	if err != nil {
		return err
	}

	return nil
}

// calcIndexSize calculates the number of bytes needed to store the index
// This includes the constant size header (8 bytes) + total size of all index entries
func calcIndexSize(indexEntries []IndexEntry) uint64 {
	var indexSize uint64 = 8 // 8 bytes for the header

	for _, idxe := range indexEntries {
		// len(word) bytes for word and 1 byte for separator
		indexSize += uint64(len(idxe.Word)) + 1
		// 8 bytes for offset and 1 byte for separator
		indexSize += 8 + 1
		// 2 bytes for definition size and 1 byte for newline
		indexSize += 2 + 1
	}

	return indexSize
}

// mergeFiles merges the words.dat and index.dat files into a single file - dict.dat
func mergeFiles() error {
	// Clean up the old dict file
	os.Remove(dictFilename)

	// Open words.dat and index.dat files for reading

	wordsFile, err := os.Open(wordsFilename)
	if err != nil {
		log.Fatalf("Error opening words.dat: %v", err)
	}
	defer wordsFile.Close()

	indexFile, err := os.Open(indexFilename)
	if err != nil {
		log.Fatalf("Error opening index.dat: %v", err)
	}
	defer indexFile.Close()

	// Open the merged file for writing
	dictFile, err := os.OpenFile(dictFilename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	// Write the index file followed by word file to the dict file

	// Buffered Copying: io.Copy copies data from the source
	// (indexFile and wordFile) to the destination (dictFile) in chunks
	// (typically 32KB by default, depending on the implementation).
	// This means it does not read the entire file into memory at once.

	_, err = io.Copy(dictFile, indexFile)
	if err != nil {
		return fmt.Errorf("error copying index file to dict: %v", err)
	}

	_, err = io.Copy(dictFile, wordsFile)
	if err != nil {
		return fmt.Errorf("error copying words file to dict: %v", err)
	}

	return nil
}
