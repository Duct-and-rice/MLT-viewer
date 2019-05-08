package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()

	app.Name = "YaruoHTML"
	app.Usage = "AST,MLT -> HTML"
	app.Version = "0.0.1"

	app.Action = func(context *cli.Context) error {
		var filetype string
		if context.Bool("ast") {
			filetype = "ast"
		} else {
			filetype = "mlt"
		}

		fmt.Print("<html><head>")
		fmt.Print(`<style>
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
		</style>`)
		fmt.Print("</head><body>")

		stdin := bufio.NewScanner(os.Stdin)
		for stdin.Scan() {
			if err := stdin.Err(); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			line := stdin.Text()
			if filetype == "mlt" && strings.HasPrefix(line, "[SPLIT]") ||
				filetype == "ast" && strings.HasPrefix(line, "[AA][") && strings.HasSuffix(line, "]") {
				fmt.Println("<hr>")
			} else {
				fmt.Println(strings.Replace(strings.Replace(line, "<", "&gt;", -1), ">", "&lt;", -1) + "<br />")
			}
		}

		fmt.Print("</body></html>")
		return nil
	}
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "ast, A",
			Usage: "File is an ast",
		},
	}

	app.Run(os.Args)
}
