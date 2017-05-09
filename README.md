[![License](http://img.shields.io/:license-mit-brightgreen.svg?style=flat-square)](http://yudppp.mit-license.org)


## goduplicheck
duplication line check tool

### install

```
$ go get github.com/yudppp/goduplicheck/cmd/goduplicheck
```

### help

```
$goduplicheck --help
NAME:
   goduplicheck - duplication check tool

USAGE:
   goduplicheck [options]

OPTIONS:
   --dir value, -d value           check directory
   --file value, -f value          check file
   --extension value, --ext value  filter extension
   --verbose, -v                   verbose
   --help                          To show help for the tool
```


### check examples

```
$ goduplicheck -f target.csv
$ goduplicheck -d target_dir
$ goduplicheck -f target_1.csv -f target_2.csv
$ goduplicheck -d target_dir1 -d target_dir2
$ goduplicheck -d target_dir -ext csv
$ goduplicheck -d target_dir -ext csv -ext tsv
```