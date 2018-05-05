package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	"github.com/inconshreveable/log15"
	"github.com/urfave/cli"
)

const (
	_satwp    = "b1674191a88ec5cdd733e4240a81803105dc412d6c6708d53ab94fc248f4f553"
	_genhash  = "000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"
	_genmroot = "4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b"
)

const version = "0.2"

var (
	gitCommit = ""
	printf    = fmt.Printf
	fprintf   = fmt.Fprintf
	newerr    = fmt.Errorf
	logger    = log15.New()
	app       = cli.NewApp()
)

func init() {
	app.Author = "aerth <aerth@riseup.net>"
	app.Copyright = "Copyright 2018  aerth (GPLv3 License) https://github.com/aerth/aevocab"
	app.UsageText = app.Name + " -h"
	app.Usage = ""
	app.HelpName = "base58 tool"
	app.Version = fmt.Sprint(version, gitCommit)
	app.Name = "aevocab"
	app.Flags = []cli.Flag{
		cli.Int64Flag{
			Name:  "r",
			Usage: "rounds/repeat (no hex)",
			Value: 1,
		},
		cli.BoolFlag{
			Name:  "d",
			Usage: "decoder mode",
		},
		cli.StringFlag{
			Name:  "input",
			Usage: "input path filename (default stdin)",
		},
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "verbose",
		},
	}
	app.Action = switchmenu // entry command
	logger.SetHandler(log15.StreamHandler(os.Stderr, log15.TerminalFormat()))
}

func main() {
	if err := app.Run(os.Args); err != nil {
		logger.Crit("fatal error", "err", err)
	}
}

func switchmenu(ctx *cli.Context) error {
	input := os.Stdin
	output := os.Stdout
	switch ctx.NArg() {
	case 0:
		return defaultAction(ctx, input, output)
	default:
		return ctx.App.Run([]string{os.Args[0], "-h"})
	}

	return defaultAction(ctx, input, output)
}

func defaultAction(ctx *cli.Context, input io.Reader, output io.Writer) error {
	scanner := bufio.NewScanner(input)
	line := ""

	// scan input lines
	for scanner.Scan() {
		line = scanner.Text()
		line = strings.TrimPrefix(line, "0x")
		if ctx.Bool("verbose") {
			println(line)
		}
		rounds := ctx.Int64("r")

		// single round
		if rounds < 2 {
			if !ctx.Bool("d") {
				fmt.Fprintln(output, base58.Encode([]byte(line)))
			} else {
				output.Write(base58.Decode(line))
				output.Write([]byte("\n"))
			}
			continue
		}

		// multiple rounds
		buf := &bytes.Buffer{}
		for i := int64(0); i < rounds; i++ {
			if ctx.Bool("verbose") {
				println("r", i)
			}
			switch ctx.Bool("d") {
			case false:
				if buf.Len() == 0 {
					fmt.Fprint(buf, base58.Encode([]byte(line)))
					continue
				}
				b := buf.Bytes()
				buf.Reset()
				fmt.Fprint(buf, base58.Encode(b))
			case true:
				if buf.Len() == 0 {
					buf.Write(base58.Decode(line))
					continue
				}
				s := buf.String()
				buf.Reset()
				buf.Write(base58.Decode(s))
			}
		}

		// print new line
		fmt.Fprintln(output, buf.String())
	}

	return nil
}
