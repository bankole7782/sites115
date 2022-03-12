package main

import (
  "os"
  "github.com/gookit/color"
  "path/filepath"
  "fmt"
  "github.com/saenuma/zazabul"
  "github.com/radovskyb/watcher"
  "log"
  "time"
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
  		fmt.Println(`sites115 is a static site generator. It brings in search functionality as well
as generates a sitemap at '/static/sitemap.xml' of your site.

Directory Commands:
  pwd     Print working directory. This is the directory where the files needed by any command
          in this cli program must reside.

Main Commands:
  cs      Creates a newsite from a template. It expects the name of the site. The site would be created
          in the 'working directory'

  rs      Render Site. It expects the name of the site. It supports reload when the site changes.

  rso     Render Site Only. It expects the path to the site. This is necessary for building a docker image.


  			`)

	case "pwd":
		fmt.Println(rootPath)

  case "cs":
    if len(os.Args) != 3 {
      color.Red.Println("Expected three arguments. Please check the help")
      os.Exit(1)
    }

    var	siteConfig = `// pagination_count
// pagination_count is the number of links to put in a page for the purposes of pagination
pagination_count: 10

// port
// port is the number that the server would listen on.
port: 8080

// base_url
// base_url is the base URL used for generating sitemaps
base_url:

	`
		configFileName := "site.zconf"
    siteName := os.Args[2]
    // create a new site
    dirsToMake := []string{"templates", "stuffs", "static", "out" }
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
title:
meta:
keywords:
---
`

    gitignoreFile := `
out/
tmp/
`
    os.WriteFile(filepath.Join(rootPath, siteName, "templates", "base.html"), baseHtmlBytes, 0777)
    os.WriteFile(filepath.Join(rootPath, siteName, "stuffs", "index.html"), []byte(indexHtml), 0777)
    os.WriteFile(filepath.Join(rootPath, siteName, "stuffs", "404.html"), []byte(indexHtml), 0777)
    os.WriteFile(filepath.Join(rootPath, siteName, "static", jqueryName), jqueryBytes, 0777)
    os.WriteFile(filepath.Join(rootPath, siteName, ".gitignore"), []byte(gitignoreFile), 0777)
    os.WriteFile(filepath.Join(rootPath, siteName, "Dockerfile"), dockerfileBytes, 0777)

    fmt.Printf("Your site is created at '%s'.\n", filepath.Join(rootPath, siteName))

  case "rso":
    if len(os.Args) != 3 {
      color.Red.Println("Expected three arguments. Please check the help")
      os.Exit(1)
    }

    path := os.Args[2]

    os.RemoveAll(filepath.Join(path, "out"))
    os.MkdirAll(filepath.Join(path, "out"), 0777)

    render(path)
    renderIndexes(path)
    generateSitemap(path)
    os.RemoveAll(filepath.Join(path, "out", "tmp"))
    os.RemoveAll(filepath.Join(path, "tmp"))


  case "rs":
    if len(os.Args) != 3 {
      color.Red.Println("Expected three arguments. Please check the help")
      os.Exit(1)
    }

    siteName := os.Args[2]
    path := filepath.Join(rootPath, siteName)

    os.RemoveAll(filepath.Join(path, "out"))
    os.MkdirAll(filepath.Join(path, "out"), 0777)

    render(path)
    renderIndexes(path)
    generateSitemap(path)
    os.RemoveAll(filepath.Join(path, "out", "tmp"))

    fmt.Println("Started...")

    // watch for new files
    w := watcher.New()

    go func() {
      for {
        select {
        case event := <-w.Event:
          fmt.Println(event)
          os.RemoveAll(filepath.Join(path, "out"))
          os.MkdirAll(filepath.Join(path, "out"), 0777)

          render(path)
          renderIndexes(path)
          generateSitemap(path)
          os.RemoveAll(filepath.Join(path, "out", "tmp"))
          os.RemoveAll(filepath.Join(path, "tmp"))

        case err := <-w.Error:
          log.Fatalln(err)
        case <-w.Closed:
          return
        }
      }
    }()

    if err := w.AddRecursive(filepath.Join(rootPath, siteName, "static")); err != nil {
      log.Fatalln(err)
    }
    if err := w.AddRecursive(filepath.Join(rootPath, siteName, "stuffs")); err != nil {
      log.Fatalln(err)
    }
    if err := w.AddRecursive(filepath.Join(rootPath, siteName, "templates")); err != nil {
      log.Fatalln(err)
    }
    if err := w.Add(filepath.Join(rootPath, siteName, "site.zconf")); err != nil {
      log.Fatalln(err)
    }

    if err := w.Start(time.Millisecond * 100); err != nil {
      log.Fatalln(err)
    }


	default:
		color.Red.Println("Unexpected command. Run the cli with --help to find out the supported commands.")
		os.Exit(1)
	}

}
