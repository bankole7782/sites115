package main

import (
  "os"
  "github.com/gookit/color"
  "path/filepath"
  "fmt"
  "github.com/saenuma/zazabul"
)


func main() {
  rootPath, err := GetRootPath()
  if err != nil {
    panic(err)
    os.Exit(1)
  }

  if len(os.Args) < 2 {
		color.Red.Println("Expecting a command. Run with help subcommand to view help.")
		os.Exit(1)
	}


	switch os.Args[1] {
	case "--help", "help", "h":
  		fmt.Println(`sites115 devc is a terminal program used during development of a static site.

Directory Commands:
  pwd     Print working directory. This is the directory where the files needed by any command
          in this cli program must reside.

Main Commands:
  ns      Creates a newsite from a template. It expects the name of the site. The site would be created
          in the 'working directory'

  			`)

	case "pwd":
		fmt.Println(rootPath)

  case "ns":
    if len(os.Args) != 3 {
      color.Red.Println("Expected three arguments. Please check the help")
      os.Exit(1)
    }

    var	siteConfig = `// pagination_count
// pagination_count is the number of links to put in a page for the purposes of pagination
pagination_count: 10

	`
		configFileName := "site.zconf"
    siteName := os.Args[2]
    // create a new site
    dirsToMake := []string{"templates", "stuffs", "static", ".out" }
    for _, dir := range dirsToMake {
      os.MkdirAll(filepath.Join(rootPath, siteName, dir), 0777)
    }


		writePath := filepath.Join(rootPath, siteName, configFileName)
		conf, err := zazabul.ParseConfig(siteConfig)
    if err != nil {
    	panic(err)
    }
    err = conf.Write(writePath)
    if err != nil {
      panic(err)
    }

    indexHtml := `---
template: base.html
---
`
    os.WriteFile(filepath.Join(rootPath, siteName, "templates", "base.html"), baseHtmlBytes, 0777)
    os.WriteFile(filepath.Join(rootPath, siteName, "stuffs", "index.html"), []byte(indexHtml), 0777)
    os.WriteFile(filepath.Join(rootPath, siteName, "static", jqueryName), jqueryBytes, 0777)

	default:
		color.Red.Println("Unexpected command. Run the cli with --help to find out the supported commands.")
		os.Exit(1)
	}

}
