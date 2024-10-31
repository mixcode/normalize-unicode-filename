
## normalize-unicode-filename: utility to rename files in Unicode normalized form


Some Unicode characters can be represented by different combinations of code points. For example, the e-acute character 'é' can be represented either in a composed form, '\u00e9', or a decomposed form, 'e\u0301'. These forms are theoretically equivalent, but they may lead to differences in actual usage. For instance, macOS typically uses the NFD (decomposed) form for filenames, while Windows generally uses the NFC (composed) form. Due to this discrepancy, filenames can appear completely different across operating systems.

This program renames filenames to their normalized Unicode forms to account for these differences.


### Install

```
$ go install github.com/mixcode/normalize-unicode-filename@latest
```

### Usage

```
Usage: ./normalize-unicode-filename [option] filename [filename...]

  -b	shorthand for '-both'
  -both
    	print both original and changed filename
  -d	shorthand for '-dryrun'
  -dryrun
    	dry-run: do not change file name; print only
  -f string
    	shorthand for '-form' (default "NFC")
  -form string
    	Unicode normalization type. One of NFC, NFD, NFKC, NFKD,
    	or WIN, MAC (default "NFC")
  -q	quiet; do not print filenames
  -r	recurse subdirectories
```

### Examples


Change filenames in the current directory to Windows-friendly form.
```
$ ./normalize-unicode-filename -form=win *
```

Change filenames to macOS-friendly form, recursively renaming files in its subdirectories.
```
$ ./normalize-unicode-filename -form=mac -r *
```

Print possible filenames for NFKD form, without changing filenames.
```
$ ./normalize-unicode-filename -form=NFKD -r -dryrun -both *
```

### Memo

Please note that NFKC and NFKD may cause irreversible changes. Be careful to using them.


