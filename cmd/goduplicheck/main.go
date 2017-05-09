package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/urfave/cli"
)

type Match struct {
	FileID  int
	LineNum uint64
}

type CheckOption struct {
	Dirs  []string
	Files []string
	Exts  []string
}

var verbose bool

func main() {
	opt := CheckOption{}

	dirs := cli.StringSlice{}
	files := cli.StringSlice{}
	exts := cli.StringSlice{}
	cli.HelpFlag = cli.BoolFlag{
		Name:  "help",
		Usage: "To show help for the tool",
	}
	app := cli.NewApp()
	app.Name = "goduplicheck"
	app.Usage = "duplication check tool"
	app.Version = "0.0.1"
	app.HideVersion = true
	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:  "dir, d",
			Usage: "check directory",
			Value: &dirs,
		},
		cli.StringSliceFlag{
			Name:  "file, f",
			Usage: "check file",
			Value: &files,
		},
		cli.StringSliceFlag{
			Name:  "extension, ext",
			Usage: "filter extension",
			Value: &exts,
		},
		cli.BoolFlag{
			Name:        "verbose, v",
			Usage:       "verbose",
			Destination: &verbose,
		},
	}
	app.Action = func(c *cli.Context) error {
		for i, ext := range exts {
			if !strings.HasPrefix(ext, ".") {
				exts[i] = fmt.Sprintf(".%s", ext)
			}
		}
		opt.Dirs = dirs
		opt.Files = files
		opt.Exts = exts
		if len(opt.Dirs)+len(opt.Files) == 0 {
			cli.ShowAppHelp(c)
			return nil
		}
		logln("duplication check start")
		logf("search target dir: %v, file: %v, ext: %v\n", opt.Dirs, opt.Files, opt.Exts)
		return exec(opt)
	}

	app.Run(os.Args)
}

func exec(opt CheckOption) error {

	checklist := make([]string, 0, 16)

	// Add dir to check list
	if len(opt.Dirs) != 0 {
		for _, dir := range opt.Dirs {
			fileInfos, err := ioutil.ReadDir(dir)
			if err != nil {
				return err
			}
			for _, file := range fileInfos {
				fileName := file.Name()
				if file.IsDir() {
					logf("[skip] %s is directory.\n", fileName)
					continue
				}
				checklist = append(checklist, path.Join(dir, fileName))
			}
		}
	}

	// Add file to check list
	if len(opt.Files) != 0 {
		for _, fileName := range opt.Files {
			checklist = append(checklist, fileName)
		}
	}

	// filter
	// checked extensions
	if len(opt.Exts) != 0 {
		filteredList := make([]string, 0, len(checklist))
		for _, fileName := range checklist {
			check := false
			for _, ext := range opt.Exts {
				if path.Ext(fileName) == ext {
					check = true
					break
				}
			}
			if !check {
				logln("[skip] don't match extensions: ", fileName)
				continue
			}
			filteredList = append(filteredList, fileName)
		}
		checklist = filteredList
	}

	// check
	existMap := make(map[string][]Match, 1000)
	fileNameMap := make(map[int]string, len(checklist))
	totalLine := uint64(0)
	for fileID, fileName := range checklist {
		fileNameMap[fileID] = fileName
		logln("search: ", fileName)
		err := func(fileName string) error {
			fp, err := os.Open(fileName)
			if err != nil {
				return err
			}
			defer fp.Close()
			reader := bufio.NewReaderSize(fp, 4096)
			i := uint64(0)
			for {
				line, _, err := reader.ReadLine()
				if err == io.EOF {
					break
				} else if err != nil {
					return err
				}
				i++
				totalLine++
				match := Match{FileID: fileID, LineNum: i}
				matchs, exist := existMap[string(line)]
				if exist {
					existMap[string(line)] = append(matchs, match)
				} else {
					existMap[string(line)] = []Match{match}
				}
			}
			return nil
		}(fileName)
		if err != nil {
			return err
		}
	}
	duplicationCount := uint64(0)
	duplicationUniqueCount := uint64(0)
	for k, matchs := range existMap {
		if len(matchs) == 1 {
			continue
		}
		duplicationCount += uint64(len(matchs))
		duplicationUniqueCount++
		fmt.Println("\nduplicate line: ", k)
		for _, v := range matchs {
			fmt.Printf("\t %s:%v\n", fileNameMap[v.FileID], v.LineNum)
		}
	}
	fmt.Println("\ncheck finished")
	fmt.Printf("\ttotal line count: %v\n", totalLine)
	fmt.Printf("\tduplication line count: %v(%v)\n", duplicationCount, duplicationUniqueCount)
	return nil
}

func logln(a ...interface{}) {
	if verbose {
		fmt.Println(a...)
	}
}

func logf(format string, a ...interface{}) {
	if verbose {
		fmt.Printf(format, a...)
	}
}
