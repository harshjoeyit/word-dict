# Word Dictionary Prototype

This project demonstrates a simple word dictionary implementation without using a database and just a single file. 


## API

The `dict` package provides the following functions for interacting with the word dictionary:

*   **`NewDict() (*dict.Dict, error)`:** Creates and initializes a new dictionary. It opens `dict.dat` file and reads index into memory.

*   **`(*Dict).QueryWord(word string) (string, bool)`:**  Using the index, API does pointed reades using offset to find definition of a word.

*   **`(*Dict).Close() error`:** Closes the dictionary file.

## Workflow for buildinga and querying the dictionary:

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