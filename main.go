package main

//See https://github.com/lspaccatrosi16/verace.js/#readme for detailed documentation

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var extensions = flag.String("extensions", "ts, js", "The extensions to filter files")
var help = flag.Bool("help", false, "Show help")

func main() {
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	var args []string

	if len(os.Args) > 1 {
		args = os.Args[1:]
	} else {
		args = make([]string, 0)
	}

	var dataToAdd []byte
	var err error

	if len(args) >= 1 && args[0] == "file" {
		dataToAdd, err = os.ReadFile("toAdd.txt")

		if err != nil {
			if os.IsNotExist(err) {
				panic("File not found: toAdd.txt")
			}
			panic(err)
		}
	} else {
		dataToAdd = []byte(strings.Join(args, " "))
	}
	extensionList := strings.Split(*extensions, ",")
	allFiles := searchDir(".", extensionList)

	commaSep := strings.Join(allFiles, ",\n")

	formatted := fmt.Sprintf("[%s]", commaSep)

	fmt.Println("Files found: ")
	fmt.Println(formatted)

	var result string

	fmt.Print("Confirm [Y/n]: ")
	fmt.Scan(&result)

	if strings.ToUpper(result) != "Y" {
		fmt.Println("Operation cancelled.")
		os.Exit(0)
	}

	for i := 0; i < len(allFiles); i++ {
		contents, err := os.ReadFile(allFiles[i])
		if err != nil {
			fmt.Println("Error")
			break
		}
		newContents := []byte(dataToAdd)
		newContents = append(newContents, contents...)
		os.WriteFile(allFiles[i], newContents, 0644)
	}

}

func searchDir(searchdir string, extensions []string) []string {
	var fileList = make([]string, 0)

	files, err := os.ReadDir(searchdir)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		name := filepath.Join(searchdir, f.Name())

		if f.IsDir() {
			fileList = append(fileList, searchDir(name, extensions)...)
		} else {

			matches := false

			for i := 0; i < len(extensions); i++ {
				if strings.HasSuffix(strings.TrimSpace(name), strings.TrimSpace(extensions[i])) {
					matches = true
				}
			}

			if matches {
				fileList = append(fileList, name)
			}

		}

	}

	return fileList
}
