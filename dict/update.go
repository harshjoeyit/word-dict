package dict

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

/*
This file contains the code to update the dictionary based on the changelog
Requirements:
1. The changelog file should be present in root directory
Constraints on changelog:
1. Is of the same format as words.dat
2. Is sorted in ascending order of word
3. Contains updated definition of only existing version of dict.dat. No new words
*/

const (
	chglogFilename = "changelog.dat"
)

func UpdateDict() error {
	// Create a temp directory and words.dat file inside if does not exist
	tempDir, err := os.MkdirTemp(".", "tmp-dict-update-*")
	if err != nil {
		return fmt.Errorf("error creating temp directory: %v", err)
	}
	// Defer the removal of the temp directory
	defer func() {
		err := os.RemoveAll(tempDir)
		if err != nil {
			log.Printf("error removing temp directory: %v", err)
		}
	}()

	// Create/Open the words.dat file for writing
	newWordsFile, err := os.OpenFile(filepath.Join(tempDir, wordsFilename), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("error creating/opening words.dat: %v", err)
	}
	defer newWordsFile.Close()

	// Open the changelog file for reading
	chglogFile, err := os.OpenFile(chglogFilename, os.O_RDONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening changelog.dat: %v", err)
	}
	defer chglogFile.Close()

	// Open dict.dat file for reading
	dictFile, err := os.OpenFile(dictFilename, os.O_RDONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening dict file: %v", err)
	}
	defer dictFile.Close()

	// Get the first word offset in the dict file
	err = seekToFirstWordOffsetInDict(dictFile)
	if err != nil {
		return fmt.Errorf("error seeking to first word offset in dict: %v", err)
	}

	err = mergeSortedFiles(newWordsFile, chglogFile, dictFile)
	if err != nil {
		return fmt.Errorf("error merging files: %v", err)
	}

	// Archive the existing words, index and dict file
	err = archiveFiles()
	if err != nil {
		return fmt.Errorf("error archiving files: %v", err)
	}

	// Move the new words.dat file from temp to the current directory
	err = os.Rename(filepath.Join(tempDir, wordsFilename), wordsFilename)
	if err != nil {
		return fmt.Errorf("error moving words.dat to current directory: %v", err)
	}

	// Rebuild the dictionary
	err = BuildNewDict()
	if err != nil {
		return fmt.Errorf("error rebuilding dictionary: %v", err)
	}

	return nil
}

func seekToFirstWordOffsetInDict(dictFile *os.File) error {
	// Read the index size from the first 8 bytes
	var indexSize int64
	// Reading 8 bytes move the offset forward by 8 bytes
	err := binary.Read(dictFile, binary.BigEndian, &indexSize)
	if err != nil {
		return fmt.Errorf("error reading index size from dict file: %v", err)
	}

	// Seek to the index size offset since first word comes after the index
	_, err = dictFile.Seek(indexSize, 0)
	if err != nil {
		return err
	}

	log.Println("Seeked to first word offset in dict file:", indexSize)

	return nil
}

// mergeSortedFiles writes the merged the dict.dat and changelog.dat files into
// a single file - ./<temp-folder>/words.dat
// This file can then be used to build the new dictionary
func mergeSortedFiles(newWordsFile, chglogFile, dictFile *os.File) error {
	// Read files line by line and write to the new words file

	dictFileScanner := bufio.NewScanner(dictFile)
	chglogFileScanner := bufio.NewScanner(chglogFile)

	chglogLine := ""
	chglogEOL := false

	for dictFileScanner.Scan() {
		dictLine := dictFileScanner.Text()

		if !chglogEOL && chglogLine == "" {
			// Read the next line from the changelog file
			if chglogFileScanner.Scan() {
				chglogLine = chglogFileScanner.Text()
			} else {
				chglogEOL = true
			}
		}

		if chglogEOL {
			// Write the dict line to the new words file as it is
			_, err := newWordsFile.WriteString(dictLine + "\n")
			if err != nil {
				return fmt.Errorf("error writing dict line to new words file: %v", err)
			}

			continue
		}

		// changelog file is not EOL
		// compare words from both files
		dictWord := strings.Split(dictLine, ",")[0]
		chglogWord := strings.Split(chglogLine, ",")[0]

		if dictWord == chglogWord {
			log.Println("Updating word:", dictWord)

			// Write the changelog line to the new words file
			_, err := newWordsFile.WriteString(chglogLine + "\n")
			if err != nil {
				return fmt.Errorf("error writing changelog line to new words file: %v", err)
			}

			// Reset the changelog line
			chglogLine = ""

		} else if dictWord < chglogWord {
			// Write the dict line to the new words file as it is
			_, err := newWordsFile.WriteString(dictLine + "\n")
			if err != nil {
				return fmt.Errorf("error writing dict line to new words file: %v", err)
			}

		} else {
			// This means that the changelog word is not in the dict file
			// this should not happen as per the constraints
			return fmt.Errorf("error: changelog word %s not found in dict file", chglogWord)
		}
	}

	log.Println("Merged files successfully")

	return nil
}

// archiveFiles moves the old words.dat, index.dat, dict.dat and changelog.dat files to an archive directory
// with the current timestamp - YYYYMMDDHHMMSS
func archiveFiles() error {
	// Create a new acrchive directory
	dir := filepath.Join("archive", time.Now().Format("20060102150405"))

	err := os.MkdirAll(dir, 0755) // Create the archive directory if it does not exist
	if err != nil {
		return fmt.Errorf("error creating archive directory: %v", err)
	}

	// Move the old words.dat, index.dat and dict.dat files to the archive directory

	err = os.Rename(wordsFilename, filepath.Join(dir, wordsFilename))
	if err != nil {
		return fmt.Errorf("error moving words.dat to archive: %v", err)
	}

	err = os.Rename(indexFilename, filepath.Join(dir, indexFilename))
	if err != nil {
		return fmt.Errorf("error moving words.dat to archive: %v", err)
	}

	err = os.Rename(dictFilename, filepath.Join(dir, dictFilename))
	if err != nil {
		return fmt.Errorf("error moving dict.dat to archive: %v", err)
	}

	// Move the changelog file to the archive directory

	err = os.Rename(chglogFilename, filepath.Join(dir, chglogFilename))
	if err != nil {
		return fmt.Errorf("error moving changelog.dat to archive: %v", err)
	}

	log.Println("Archived old files successfully")

	return nil
}
