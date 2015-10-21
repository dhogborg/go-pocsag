package main

import (
	"io"
	"os"

	"github.com/codegangsta/cli"
	"github.com/fatih/color"

	"bitbucket.org/dhogborg/go-pocsag/internal/datatypes"
	"bitbucket.org/dhogborg/go-pocsag/internal/pocsag"
	"bitbucket.org/dhogborg/go-pocsag/internal/utils"
)

var (
	config *Config
)

var (
	green = color.New(color.FgGreen)
	red   = color.New(color.FgRed)
	blue  = color.New(color.FgBlue)
)

type Config struct {
	input       string
	baud        int
	debug       bool
	messagetype pocsag.MessageType
}

func main() {

	app := cli.NewApp()
	app.Name = "go-pocsag"
	app.Usage = "Parse audiostream for POCSAG messages"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "input,i",
			Value: "",
			Usage: "wav file or data dump with signed 16 bit ints",
		},
		cli.IntFlag{
			Name:  "baud,b",
			Value: 0,
			Usage: "Baud 600/1200/2400. Default auto",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Debug mode",
		},
		cli.StringFlag{
			Name:  "type,t",
			Value: "auto",
			Usage: "Force message type: alpha, bcd, auto",
		},
	}

	app.Action = func(c *cli.Context) {
		config = &Config{
			input:       c.String("input"),
			baud:        c.Int("baud"),
			debug:       c.Bool("debug"),
			messagetype: pocsag.MessageType(c.String("type")),
		}

		utils.SetDebug(config.debug)
		pocsag.SetDebug(config.debug)

		Run()

	}

	app.Run(os.Args)
}

func Run() {

	var source io.Reader

	if config.input == "-" {
		source = os.Stdin
	} else { // file reading
		source = pocsag.ReadWav(config.input)
	}

	if source == nil {
		println("invalid input")
		os.Exit(0)
	}

	reader := pocsag.NewStreamReader(source, config.baud)

	bitstream := make(chan []datatypes.Bit, 1)
	go reader.StartScan(bitstream)

	for {
		bits := <-bitstream
		pocsag.ParsePOCSAG(bits, config.messagetype)
	}
}
