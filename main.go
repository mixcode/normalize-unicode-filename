package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/text/unicode/norm" // unicode normalizer
)

// command line arguments
var (
	formName  string = "NFC"
	recurse          = false
	quiet            = false
	dryrun           = false
	printBoth        = false
)

// runtime variables
var (
	formCode  norm.Form
	fileCount = 0

	dirFixed = make(map[string]string)

	sep = string(filepath.Separator) // path separator in string
)

const (
	help_details = `Some Unicode characters can be represented by different combinations of code points. For example, the e-acute character 'é' can be represented either in a composed form, '\u00e9', or a decomposed form, 'e\u0301'. These forms are theoretically equivalent, but they may lead to differences in actual usage. For instance, macOS typically uses the NFD (decomposed) form for filenames, while Windows generally uses the NFC (composed) form. Due to this discrepancy, filenames can appear completely different across operating systems.

This program renames filenames to their normalized Unicode forms to account for these differences.

`

	help_examples = `Change filenames in the current directory to Windows-friendly form:
  $ %[1]s -form=win *

Change filenames to macOS-friendly form, recursively renaming files in its subdirectoreis:
  $ %[1]s -form=mac -r *

Print possible filenames for NFKD form, without changing filenames 
  $ %[1]s -form=NFKD -r -dryrun -both *

`
)

func normalize(s string) string {
	return formCode.String(s)
}

func process(originalName string) (err error) {
	var fInfo os.FileInfo
	fInfo, err = os.Stat(originalName)
	if err != nil {
		return
	}

	dir, fname := filepath.Split(originalName)

	actualName := originalName // the name of actual file based on dryrun flag
	newf := normalize(fname)
	newName := filepath.Join(dir, newf)

	if newf != fname { // name normalized
		fileCount++

		// for dry-run; get possibly renamed file path
		fixedDir := dirFixed[dir]
		if fixedDir == "" {
			fixedDir = dir
		}
		newName = filepath.Join(fixedDir, newf)

		// print the filePath
		if !quiet {
			if printBoth {
				fmt.Printf("%s\n  -> %s\n", originalName, newName)
			} else {
				fmt.Printf("%s\n", newName)
			}
		}

		// rename the file
		if !dryrun {
			err = os.Rename(originalName, newName)
			if err != nil {
				return
			}
			actualName = newName
		}
	}

	if fInfo.IsDir() {
		originalName = filepath.Join(originalName, "") + sep
		newName = filepath.Join(newName, "") + sep
		dirFixed[originalName] = newName
		if recurse {
			d, e := os.ReadDir(actualName)
			if e != nil {
				return e
			}
			for _, f := range d {
				subf := filepath.Join(actualName, f.Name())
				err = process(subf)
				if err != nil {
					return
				}
			}
			return nil
		}
	}

	return nil
}

func run() (err error) {

	formName = strings.ToUpper(formName)
	switch formName {
	case "NFC", "WIN": // Canonical equivalence, Composing
		formCode = norm.NFC
	case "NFD", "MAC": // Canonical equivalence, Decomposing
		formCode = norm.NFD

	case "NFKC": // Kompatibility equivalence, Composing
		formCode = norm.NFKC
	case "NFKD": // Kompatibility equivalence, Decomposing
		formCode = norm.NFKD
	default:
		return fmt.Errorf("invalid normalization form")
	}

	args := flag.Args()
	for _, pattern := range args {
		var l []string
		l, err = filepath.Glob(pattern)
		if err != nil {
			return
		}
		if l == nil {
			continue
		}

		for _, name := range l {

			err = process(name)
			if err != nil {
				return
			}
		}
	}

	return
}

func main() {
	var err error

	flag.StringVar(&formName, "form", formName, "Unicode normalization type. One of NFC, NFD, NFKC, NFKD,\nor WIN, MAC")
	flag.StringVar(&formName, "f", formName, "shorthand for '-form'")

	flag.BoolVar(&recurse, "r", recurse, "recurse subdirectories")

	flag.BoolVar(&quiet, "q", quiet, "quiet; do not print filenames")

	flag.BoolVar(&dryrun, "d", dryrun, "shorthand for '-dryrun'")
	flag.BoolVar(&dryrun, "dryrun", dryrun, "dry-run: do not change file name; print only")

	flag.BoolVar(&printBoth, "both", printBoth, "print both original and changed filename")
	flag.BoolVar(&printBoth, "b", printBoth, "shorthand for '-both'")

	flag.Usage = func() {
		o := flag.CommandLine.Output()
		execName := os.Args[0]
		fmt.Fprintln(o)
		fmt.Fprintf(o, "%s: Rename files in Unicode normalized form\n\n", execName)
		fmt.Fprintf(o, help_details)
		fmt.Fprintf(o, "Usage: %s [option] filename [filename...]\n\n", execName)
		flag.PrintDefaults()
		fmt.Fprintln(o)
		fmt.Fprintf(o, "Examples:\n")
		fmt.Fprintf(o, help_examples, execName)
		fmt.Fprintf(o, "Memo:\n")
		fmt.Fprintf(o, "Please note that NFKC and NFKD may cause irreversible changes. Be careful using them.")
		fmt.Fprintln(o)
	}

	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(0)
	}

	// run main
	err = run()

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func init() {
	// init default normalization form based on the OS
	switch runtime.GOOS {
	case "windows": // windows
		formName = "NFC"
	case "darwin": // macos / ios
		formName = "NFD"
	//case "linux":
	//	formName = "NFC"
	default:
		formName = "NFC"
	}
}
