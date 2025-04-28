# Word Dictionary Prototype

This project demonstrates a simple word dictionary implementation without using a database and just a single file. 


## API

The `dict` package provides the following functions for interacting with the word dictionary:

*   **`NewDict() (*dict.Dict, error)`:** Creates and initializes a new dictionary. It opens `dict.dat` file and reads index into memory.

*   **`UpdateDict() error`:** Updates the dictionary based on changes specified in the `changelog.dat` file. This function performs the following steps:

    1.  Merges the `changelog.dat` file with the existing `dict.dat` to create a new `<temp-folder>/words.dat` file.
    2.  Archives the old dictionary files (words.dat, index.dat, dict.dat, and changelog.dat) to an archive directory.
    3.  Rebuilds the dictionary index using `<temp-folder>/words.dat` and creates a new `dict.dat` file.

    **Important:** The `changelog.dat` file must be in the same format as `words.dat` (i.e., `word,definition` on each line) and must be sorted in ascending order of words. The `changelog.dat` file should only contain updates to *existing* words in the dictionary; it should not contain new words.


*   **`(*Dict).QueryWord(word string) (string, bool)`:**  Using the index, API does pointed reades using offset to find definition of a word.

*   **`(*Dict).Close() error`:** Closes the dictionary file.

The `s3dict` package provides the following functions for interacting with a word dictionary stored in AWS S3:

*   **`New() (*s3dict.S3Dict, error)`:** Creates and initializes a new `S3Dict` object. It retrieves the dictionary file from S3, reads the index, and stores it in memory. This function requires the following environment variables to be set:

    *   `S3_BUCKET_NAME`, `AWS_REGION`, `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `DICT_KEY`: The key (path) of the dictionary file in the S3 bucket.

    Returns a pointer to an `S3Dict` object and an error if the dictionary cannot be created.

*   **`(*S3Dict).QueryWord(word string) (string, bool)`:** Queries the dictionary for a word and returns its definition. It first checks the in-memory index for the word. If found, it retrieves the definition from the S3 object using a byte range request. Returns the definition of the word (if found) and a boolean indicating whether the word was found.

## Workflow for building and querying the dictionary:

1.  **Create the words data file (`words.dat`):**
    *   This file is required and contains words and their meanings, with each entry formatted as `word,definition`.
    *   Example:

        ```
        a,first english alphabet
        abandon,to leave and never return to
        ability,power or skill to do something
        boast,talk with excessive pride and self-satisfaction about one's achievements or abilities.
        ...
        ```

2.  **Generate the index file (`index.dat`):**
    *   The program reads `words.dat` and creates an index that maps each word to its offset within the file.
    *   This index is then written to a separate file named `index.dat`.

3.  **Create the final dictionary file (`dict.dat`):**
    *   The program combines `index.dat` and `words.dat` into a single file named `dict.dat`.
    *   The `index.dat` content is prepended to the `words.dat` content, allowing for efficient word lookups using the index.

4. **Query words:**
    *   Create a NewDict() which loads the index in memory
    *   Use the Query() API to query the words 