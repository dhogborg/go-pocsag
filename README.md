# go-pocsag

A parser for POCSAG pager protocol implemented in Go

## Usage
Read a recorded wav file `gopocsag -i path/to/file.wav`

Listen to stream from rtl_fm: `rtl_fm -f <freq> -E deemp | gopocsag`

Parse a raw datadump: `cat dump.bin | gopocsag`

## Options
* `--type` force message parsing type, one of `auto` `bcd` `alpha`
* `--debug` print debugging and extra information about transmission.
* `--verbosity` regulate the detail of debugging information

## Resource usage
Not much. About 0.2% of a i5 during normal operations. Just above 5 mb of RAM.