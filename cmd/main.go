package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

// Code Snippets & Ideas from:
// https://flaviocopes.com/go-list-files/
// https://stackoverflow.com/a/20877193
// https://topic.alibabacloud.com/a/how-to-use-golang-to-get-the-creationmodification-time-of-files-on-linux_1_16_30132202.html
// https://golang.org/pkg/os/#FileInfo

func visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		*files = append(*files, path)
		return nil
	}
}

func main() {
	var files []string
	var pathErr error
	var root, userConsent string

	fpath := flag.String("p", "", "Path inside which to arrange the files into a YYYY/mm/dd folder structure. Defaults to the current directory. PLEASE NOTE: This program does not recurse into sub-directories")

	flag.Parse()

	if *fpath == "" {
		*fpath, pathErr = os.Getwd()
		if pathErr != nil {
			fmt.Printf("Error getting current working directory. Details: %s", pathErr.Error())
		} else {
			root = *fpath
		}
	}

	fmt.Printf("This is your LAST CHANCE to stop this. Enter Y/N(Y=Yes, continue. N=No, stop and exit)? ")
	fmt.Scanf("%s", &userConsent)
	if userConsent != "Y" {
		fmt.Printf("You did not give consent to continue. Exiting...")
		os.Exit(-1)
	}

	err := filepath.Walk(root, visit(&files))
	if err != nil {
		panic(err)
	}
	for fileKey, file := range files {
		fi, err := os.Stat(file)
		if err != nil {
			fmt.Printf("Error stating file %s. Details: %s", file, err.Error())
		}
		modifiedTime := fi.ModTime()
		fiSys := fi.Sys().(*syscall.Stat_t)
		atime := time.Unix(fiSys.Atimespec.Unix())
		ctime := time.Unix(fiSys.Ctimespec.Unix())
		mtime := time.Unix(fiSys.Mtimespec.Unix())
		btime := time.Unix(fiSys.Birthtimespec.Unix())
		//fmt.Printf("%d. [%s] %s\n", fileKey, modifiedTime.String(), file)
		fmt.Printf("%d. [%s] %s\nLast Access Time: %s\nLast Modified Time:%s\nCreated Time: %s\nBirth Time: %s\n", fileKey, modifiedTime.String(), file, atime.String(), mtime.String(), ctime.String(), btime.String())

		btime_yyyy := btime.Year()
		btime_mm := int(btime.Month())
		btime_dd := btime.Day()

		dirString := root + "/" + strconv.Itoa(btime_yyyy) + "/" + strconv.Itoa(btime_mm) + "/" + strconv.Itoa(btime_dd)

		organizedFilePath := dirString + "/" + path.Base(file)

		if _, err := os.Stat(dirString); os.IsNotExist(err) {
			fmt.Printf("Directory %s Nonexistent. Creating...\n", dirString)
			mkdirErr := os.MkdirAll(dirString, 0755)
			if mkdirErr != nil {
				fmt.Printf("Error creating directory %s. Details: %s\n", dirString, mkdirErr.Error())
			}
		}

		if ! fi.IsDir() {
			moveFileErr := os.Rename(file, organizedFilePath)
			if err != nil {
				fmt.Printf("Error moving file from %s to %s. Details: %s\n", file, organizedFilePath, moveFileErr.Error())
			} else {
				fmt.Printf("Successfully moved file from %s to %s.\n", file, organizedFilePath)
			}
		} else {
			fmt.Printf("Not doing anything as %s is a directory\n", fi.Name())
		}

		fmt.Printf("----------------------\n")
	}
}