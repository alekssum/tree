package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}

	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out *os.File, path string, printFiles bool) error {

	var tree func(string, int, bool) (string, error)

	tree = func(path string, lvl int, isLast bool) (string, error) {
		file, err := os.Open(path)
		if err != nil {
			return "", err
		}

		list, err := file.Readdir(0)
		if err != nil {
			return "", err
		}

		if !printFiles {

		}
		resultList := list
		sort.SliceStable(resultList, func(i, j int) bool {
			if resultList[i].IsDir() && !resultList[j].IsDir() {
				return true
			}
			return resultList[i].Name() > resultList[j].Name()
		})
		resultString := ""
		for i, v := range resultList {
			if !v.IsDir() && !printFiles {
				continue
			}

			prefix := ""
			//If parent is the last item in it's dirrectory
			if isLast {
				prefix = strings.Repeat("	", lvl)
			} else {
				prefix = strings.Repeat("│	", lvl)
			}

			switch {
			// The last item in current dirrectory
			case i == len(list)-1:
				prefix = prefix + "└───"
			default:
				prefix = prefix + "├───"
			}

			resultString += fmt.Sprintf("%s\n", prefix+v.Name())
			if v.IsDir() {

				res, err := tree(path+"/"+v.Name(), lvl+1, i == len(list)-1)
				if err != nil {
					return "", err
				}
				resultString += res
			}

		}
		return resultString, nil
	}

	res, err := tree(path, 0, false)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(out, res)
	return nil
}
