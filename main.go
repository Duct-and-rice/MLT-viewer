package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/pkg/browser"
	"golang.org/x/text/encoding/japanese"
	cli "gopkg.in/urfave/cli.v1"
)

type yhHandler struct {
	html string
}

func (h *yhHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, h.html)
}

func serve(html string, port int) error {
	addr := fmt.Sprintf(":%d", port)
	handler := &yhHandler{html}
	l, err := net.Listen("tcp", "localhost"+addr)
	if err != nil {
		return err
	}
	fmt.Printf("I'm listening on http://localhost:%d\n", port)
	browser.OpenURL("http://localhost" + addr)
	http.Serve(l, handler)
	return nil
}

func main() {
	app := cli.NewApp()

	app.Name = "YaruoHTML"
	app.Usage = "AST,MLT -> HTML"
	app.Version = "0.0.1"

	app.Action = func(context *cli.Context) error {
		args := context.Args()
		inputFileName := args.First()
		var filetype string
		if context.Bool("ast") {
			filetype = "ast"
		} else {
			if strings.HasSuffix(inputFileName, ".ast") {
				filetype = "ast"
			} else {
				filetype = "mlt"
			}
		}

		headerHTML := `
		<html>
			<head>
			<style>
			@font-face { 
				font-family: 'Stmr'; 
				src: url("https://cdn.rawgit.com/Duct-and-rice/yaruo-blog/1317c979/fonts/Saitamaar.woff2") 
						format("woff2"), 
						url("https://cdn.rawgit.com/Duct-and-rice/yaruo-blog/1317c979/fonts/Saitamaar.woff") 
						format("woff"), 
						url("https://cdn.rawgit.com/Duct-and-rice/yaruo-blog/1317c979/fonts/Saitamaar.ttf") 
						format("truetype"), 
						url("https://cdn.rawgit.com/Duct-and-rice/yaruo-blog/1317c979/fonts/Saitamaar.eot"), 
						url("https://cdn.rawgit.com/Duct-and-rice/yaruo-blog/1317c979/fonts/Saitamaar.eot?#iefix") 
						format("embedded-opentype"); 
						font-weight: normal; 
						font-display: auto; 
					}
			body {
				font-family: Stmr;
			}
			</style>
		</head>
		<body>`

		var builder strings.Builder
		var writer io.Writer
		var reader io.Reader
		var scanner *bufio.Scanner

		if context.Bool("server") {
			writer = &builder
		} else if context.IsSet("output") {
			file, err := os.OpenFile(context.String("output"), os.O_WRONLY, 0644)
			if err != nil {
				return err
			}
			writer = bufio.NewWriter(file)
		} else {
			writer = os.Stdout
		}

		if inputFileName != "" {
			file, err := os.OpenFile(inputFileName, os.O_RDONLY, 0644)
			if err != nil {
				return err
			}
			reader = file
		} else {
			reader = os.Stdin
		}

		if context.Bool("utf") {
			scanner = bufio.NewScanner(reader)
		} else {
			scanner = bufio.NewScanner(japanese.ShiftJIS.NewDecoder().Reader(reader))
		}

		fmt.Fprint(writer, headerHTML)

		for scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return err
			}
			line := scanner.Text()
			if filetype == "mlt" && strings.HasPrefix(line, "[SPLIT]") ||
				filetype == "ast" && strings.HasPrefix(line, "[AA][") && strings.HasSuffix(line, "]") {
				fmt.Fprintln(writer, "<hr>")
			} else {
				fmt.Fprintln(writer, strings.Replace(strings.Replace(line, "<", "&gt;", -1), ">", "&lt;", -1)+"<br />")
			}
		}

		fmt.Fprint(writer, "</body></html>")

		if context.Bool("server") {
			html := builder.String()
			return serve(html, context.Int("port"))
		}
		return nil
	}
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "ast, A",
			Usage: "File is an ast",
		},
		cli.StringFlag{
			Name:  "output, o",
			Usage: "File to output",
		},
		cli.BoolFlag{
			Name:  "server, s",
			Usage: "Open http server",
		},
		cli.BoolFlag{
			Name:  "utf, u",
			Usage: "When open utf8 encoded file, flag this",
		},
		cli.IntFlag{
			Name:  "port, p",
			Usage: "Port number",
			Value: 8080,
		},
	}

	app.Run(os.Args)
}
