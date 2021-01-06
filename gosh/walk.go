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
	// if the directory where the search will take place is not specified use the current one
	currentDir, _ := os.Getwd()

	// fmt.Println("")
	// if the previous char is not a blank space it means a word was being typed
	if !isEmpty(t.line) && !isABlankSpace(t.line[t.position - 1]) {
		wordToComplete = t.getWordToComplete()

		// check for slashes ("/")
		if strings.Contains(wordToComplete, "/") {
			if strings.HasPrefix(wordToComplete, "/") {
				currentDir = ""

				if strings.Contains(wordToComplete[1:], "/") {
					a := strings.Split(wordToComplete[1:], "/")
					// fmt.Println("Tab Completer -- a: ", a)
					for i := 0; i < len(a) - 1; i++ {
						currentDir = currentDir + "/" + a[i]
					}
					wordToComplete = a[len(a) - 1]
				} else {
					currentDir = "/"
					wordToComplete = wordToComplete[1:]
				}

			} else if strings.HasSuffix(wordToComplete, "/") {
				currentDir = currentDir + "/" + wordToComplete
				wordToComplete = ""

			} else {
				a := strings.Split(wordToComplete, "/")
				// fmt.Println("Tab Completer -- a: ", a)
				for i := 0; i < len(a) - 1; i++ {
					currentDir = currentDir + "/" + a[i]
				}
				wordToComplete = a[len(a) - 1]
			}
		}
	}

	// fmt.Println("Tab Completer -- dir: ", currentDir, " | word: ", wordToComplete)
	t.walk(currentDir, wordToComplete)
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
	// fmt.Println("Walk -- dir: ", directory, " | word: ", wordToComplete)
	files := getFilesFromDir(directory, wordToComplete)

	// If there is a lot of results show them on the screen
	// fmt.Println("Walk -- files:")
	if len(files) > 1 {
		fmt.Println("")
		for _, f := range files {
			if f.info.IsDir() {
				fmt.Print(f.name, "/ ")
			} else {
				fmt.Print(f.name, " ")
			}
		}
		fmt.Println("")
	}

	// If there is only one result, add it to the prompt line
	if len(files) == 1 {
		// if there is a word to complete, add to the line only the complement of the word, the prefix is already there
		if wordToComplete != "" {
			// position 0 contains the prefix of the word, that is already typed
			// position 1 contains the complement
			substrings := strings.SplitAfterN(files[0].name, wordToComplete, 2)
			// fmt.Println("substrings: ", substrings)
			// fmt.Println("substrings: ", len(substrings[1]))

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

		// fmt.Println("line: ", string(t.line))
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
		fmt.Println("err: ", err)
		return files
	}

	filesInfo, err := descriptor.Readdir(0)
    if err != nil {
		fmt.Println("err: ", err)
        return files
	}

    for _, fileInfo := range filesInfo {
		if prefix == "" {
			files = append(files, file{
				name: fileInfo.Name(),
				info: fileInfo,
			})
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
