# Turtle

Network traffic monitoring tool that gather backend data into a database, currently supported backends are bandwhich and nethogs.

## Why

While programmes like nethogs and bandwhich are excellent for tracking network bandwidth, I wanted something that could save data persistently so that I could verify it later. As a result, turtle was developed. I can later check the logs by simply looking into a database.

## Install

If you already have go installed:

```bash
go install github.com/omarahm3/turtle@latest
```

Or from releases page

## Run

Since turtle is using nethogs & bandwhich which requires sudo access to collect network traffic data, turtle must be running with sudo so that it can spawn backend processes properly.

```bash
❯ sudo turtle -h
Log nethogs traffic per processes and applications

Usage:
  turtle [flags]
  turtle [command]

Available Commands:
  clear       Clear database
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  list        List collected information

Flags:
  -b, --bandwhich   use bandwhich as a underlayer backend (default true)
  -h, --help        help for turtle
  -n, --nethogs     use nethogs as a underlayer backend
```

Turtle will write all information on a DB under: `/var/log/.turtle.db`
