# csv2json

This is a little command line utility to convert csv data into json data.
It uses the column headers as the field names on the resulting json objects.

## Install

Installing is using golang tooling.

```bash
go install github.com/Reisender/csv2json
```

## Usage

You can look at the available flags with the `--help` flag.

```bash
csv2json --help
```

This utility expects the csv data to be piped in to it through STDIN
while the json data is sent to STOUT.

```bash
cat ./data.csv | csv2json > ./data.json
```

### schema validation

You can validate the schema of the json with a [cuelang](https://cuelang.org/) file.

```bash
cat ./users.csv | csv2json --schema ./users.cue > ./users.json
```
