package main

import (
    "fmt"
    "os"
	"strings"
)

type file struct {
	name string
	info os.FileInfo
}

func (t *terminal) tabCompleter() {
	wordToComplete := ""

	// if the directory where the search will take place is not specified, use the current one
	directory, _ := os.Getwd()

	// if the previous char is not a blank space it means a word was being typed
	if !isEmpty(t.line) && !isABlankSpace(t.line[t.position - 1]) {
		wordToComplete = t.getWordToComplete()

		// check for slashes ("/")
		// if there are slashes, means a directory was specified for the search to take place in
		if strings.Contains(wordToComplete, "/") {

			// slash as a prefix, means to search in the root directory
			if strings.HasPrefix(wordToComplete, "/") {

				if strings.Contains(wordToComplete[1:], "/") {
					substrings := strings.Split(wordToComplete[1:], "/")

					directory = ""
					for _, substring := range substrings[:len(substrings) - 1] {
						directory = directory + "/" + substring
					}

					wordToComplete = substrings[len(substrings) - 1]

				// set root as the directory where the search will take place
				// and remove this prefix slash of the word to complete
				} else {
					directory = "/"
					wordToComplete = wordToComplete[1:]
				}

			// slash as a suffix, means to search in a subdirectory of the current one
			// and no word to complete
			} else if strings.HasSuffix(wordToComplete, "/") {
				directory = directory + "/" + wordToComplete
				wordToComplete = ""

			} else {
				substrings := strings.Split(wordToComplete, "/")

				for _, substring := range substrings[:len(substrings) - 1] {
					directory = directory + "/" + substring
				}

				wordToComplete = substrings[len(substrings) - 1]
			}
		}
	}

	t.walk(directory, wordToComplete)
}

// Get all chars from the previous nearest blank space until the current cursor position
func (t *terminal) getWordToComplete() string {
	cursorPosition := t.position
	spacePosition := t.position

	// find how many chars until reach a blank space or the begining of the line
	for {
		if spacePosition - 1 == -1 || isABlankSpace(t.line[spacePosition - 1]) {
			break
		}
		spacePosition--
	}

	return string(t.line[spacePosition:cursorPosition])
}

func (t *terminal) walk(directory string, wordToComplete string) {
	files := getFilesFromDir(directory, wordToComplete)

	// If there is a lot of results show them on the screen
	if len(files) > 1 {
		fmt.Fprintln(os.Stdout)
		for _, file := range files {
			if file.info.IsDir() {
				fmt.Fprint(os.Stdout, file.name, "/ ")
			} else {
				fmt.Fprint(os.Stdout, file.name, " ")
			}
		}
		fmt.Fprintln(os.Stdout)
	}

	// If there is only one result, add it to the prompt line
	if len(files) == 1 {
		// if there is a word to complete, add to the line only the complement of the word, the prefix is already there
		if wordToComplete != "" {
			// position 0 contains the prefix of the word, that is already typed
			// position 1 contains the complement
			substrings := strings.SplitAfterN(files[0].name, wordToComplete, 2)

			// if the complement is empty, means the role word is already typed
			// just add a slash in case of a folder or a space in case of a file
			if len(substrings[1]) == 0 {
				if files[0].info.IsDir() {
					t.line = append(t.line, '/')
				} else {
					t.line = append(t.line, ' ')
				}
			} else {
				// add the complement to the line
				rest := []rune(substrings[1])
				t.line = append(t.line, rest...)
			}
		} else {
			// if there is no word to complete, add the role file name to the line
			rest := []rune(files[0].name)
			t.line = append(t.line, rest...)
		}

		t.position = len(t.line)
	}

	t.needRefresh = true
}

// Get all files from the specified directory
// If a prefix is passed filter the files by it, otherwise return all of them
func getFilesFromDir(directory string, prefix string) []file {
	var files []file

	descriptor, err := os.Open(directory)
	defer descriptor.Close()
	if err != nil {
		return files
	}

	filesInfo, err := descriptor.Readdir(0)
    if err != nil {
        return files
	}

    for _, fileInfo := range filesInfo {
		if prefix == "" {
			files = append(files, file{
				name: fileInfo.Name(),
				info: fileInfo,
			})

		// if there is a prefix filter the files by it
		} else {
			if strings.HasPrefix(fileInfo.Name(), prefix) {
				files = append(files, file{
					name: fileInfo.Name(),
					info: fileInfo,
				})
			}
		}
	}

	return files
}
