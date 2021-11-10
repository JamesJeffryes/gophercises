package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	var (
		match   string
		outFmt  string
		recurse bool
		flatten bool
		force   bool
	)
	flag.StringVar(&match, "match", "", "A string to match in files")
	flag.StringVar(&outFmt, "outfmt", "", "Supply the output format. "+
		"Use $oldname, $oldext, $match, $count, $total to replace with those values. Defaults to $match$count$oldext")
	flag.BoolVar(&recurse, "recurse", false,
		"If true, the contents of subdirectories will be renamed as well")
	flag.BoolVar(&flatten, "flatten", false,
		"If true, matching files from subfolders will be copied to the root")
	flag.BoolVar(&force, "force", false,
		"If true, overwrite existing files")

	flag.Parse()
	if len(flag.Args()) != 1 {
		log.Fatalln("Please supply exactly 1 path to files/folders to rename")
	}
	rootPath, err := filepath.Abs(flag.Args()[0])
	if _, err := os.Stat(rootPath); err != nil {
		log.Fatalln("Unable to walk supplied path.")
	}
	if len(match) == 0 {
		fmt.Printf("You have not supplied an input format flag. All files in the folder will be renamed.\n")
	}
	if len(outFmt) == 0 {
		outFmt = "$match$count$oldext"
		fmt.Printf("You have not supplied an output format flag. "+
			"Matching files in the folder will be renamed '%s'.\n", outFmt)
	}
	var total int
	err = filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if d != nil && !d.IsDir() && strings.Contains(filepath.Base(path), match) {
			total++
		} else if !recurse && path != rootPath {
			return filepath.SkipDir
		}
		return err
	})
	handleError(err)
	fmt.Printf("%d files will be renamed.\n", total)
	fmt.Printf("Type 'y' to continue or any other key to quit.")
	consoleReader := bufio.NewReader(os.Stdin)
	response, _ := consoleReader.ReadString('\n')
	if strings.TrimSpace(response) != "y" {
		os.Exit(0)
	}
	var count int
	err = filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if d != nil && !d.IsDir() {
			if !strings.Contains(filepath.Base(path), match) {
				return nil
			}
			count++
			newDir := filepath.Dir(path)
			if flatten {
				newDir = rootPath
			}
			newName := strings.ReplaceAll(outFmt, "$oldname", filepath.Base(path))
			newName = strings.ReplaceAll(newName, "$oldext", filepath.Ext(path))
			newName = strings.ReplaceAll(newName, "$match", match)
			newName = strings.ReplaceAll(newName, "$count", strconv.Itoa(count))
			newName = strings.ReplaceAll(newName, "$total", strconv.Itoa(total))
			newPath := filepath.Join(newDir, newName)
			if _, err = os.Stat(rootPath); err != nil && !force {
				log.Printf("Skipping renameing %s as file with target name already exists", path)
				return nil
			}
			return os.Rename(path, newPath)
		} else if !recurse && path != rootPath {
			return filepath.SkipDir
		}
		return nil
	})
	handleError(err)
	fmt.Println("Rename complete!")
}

func handleError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
