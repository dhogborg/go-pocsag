package main

import (
	"io"
	"os"

	"github.com/urfave/cli"
	"github.com/fatih/color"

	"github.com/dhogborg/go-pocsag/internal/datatypes"
	"github.com/dhogborg/go-pocsag/internal/pocsag"
	"github.com/dhogborg/go-pocsag/internal/utils"
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
	output      string
	baud        int
	debug       bool
	messagetype pocsag.MessageType
	verbosity   int
}

func main() {

	app := cli.NewApp()
	app.Name = "go-pocsag"
	app.Usage = "Parse audiostream for POCSAG messages"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "input,i",
			Value: "",
			Usage: "wav file with signed 16 bit ints, - for sttdin",
		},
		cli.StringFlag{
			Name:  "output,o",
			Value: "",
			Usage: "Output decoded messages to a folder",
		},
		cli.IntFlag{
			Name:  "verbosity",
			Value: 0,
			Usage: "Verbosity level, 0 lowest level",
		},
		cli.IntFlag{
			Name:  "baud,b",
			Value: 0,
			Usage: "Baud 600/1200/2400. Default auto",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Output debug information",
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
			output:      c.String("output"),
			baud:        c.Int("baud"),
			debug:       c.Bool("debug"),
			verbosity:   c.Int("verbosity"),
			messagetype: pocsag.MessageType(c.String("type")),
		}

		utils.SetDebug(config.debug, config.verbosity)
		pocsag.SetDebug(config.debug, config.verbosity)

		Run()

	}

	app.Run(os.Args)
}

func Run() {

	var source io.Reader

	if config.input == "-" || config.input == "" {
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
		messages := pocsag.ParsePOCSAG(bits, config.messagetype)

		for _, m := range messages {
			m.Print(config.messagetype)

			if config.output != "" {
				m.Write(config.output, config.messagetype)
			}
		}

	}
}
