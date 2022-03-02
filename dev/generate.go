package main

import (
  "os"
  "html/template"
  "github.com/otiai10/copy"
  "strings"
  "fmt"
  "errors"
  "path/filepath"
  "bufio"
  "github.com/bankole7782/sites115/sites115s"
  "github.com/russross/blackfriday"
)


func generate(siteName string) {
  rootPath, _ := GetRootPath()

  copy.Copy(filepath.Join(rootPath, siteName, "site.zconf"), filepath.Join(rootPath, siteName, "out", "site.zconf"))
  // copy statics
  copy.Copy(filepath.Join(rootPath, siteName, "static"), filepath.Join(rootPath, siteName, "out", "static"))
  os.WriteFile(filepath.Join(rootPath, siteName, "static", "index.html"),
    []byte("index.html not generated for this page."), 0777)

  dir := filepath.Join(rootPath, siteName, "stuffs")

  // render pages
  err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
    if err != nil {
      return err
    }

    if info.IsDir() {
      if ! sites115s.DoesPathExists(filepath.Join(path, "index.html")) {
        return errors.New(fmt.Sprintf("The directory '%s' does not have an index.html. It is compulsory", path))
      }
      if ! sites115s.DoesPathExists(filepath.Join(path, "toc.txt")) && path != dir {
        return errors.New(fmt.Sprintf("The directory '%s' does not have an toc.txt . It is compulsory", path))
      }
    }

    if ! info.IsDir() {
      if ! strings.HasSuffix(dir, "/") {
        dir += "/"
      }

      if strings.HasSuffix(path, ".txt") {
        return nil
      }

      raw, err := os.ReadFile(path)
      if err != nil {
        return err
      }
      if ! strings.HasPrefix(string(raw), "---") {
        return errors.New(fmt.Sprintf("The file '%s' does not have a data part", path))
      }

      pathToWrite := strings.Replace(path, dir, "", 1)

      if strings.HasSuffix(path, ".html") {
        return RenderHTMLToFile(string(raw), pathToWrite, siteName)
      } else if strings.HasSuffix(path, ".md") {
        return RenderMDToFile(string(raw), pathToWrite, siteName)
      }

      // outStr := fmt.Sprintf("%s,,,%d\n", pathToWrite, info.Size())
      // indexFile.WriteString(outStr)
    }
    return nil
  })

  if err != nil {
    panic(err)
  }
}


func GetPartsOfMarkup(s string) (string, string) {
  parts := strings.Split(s, "\n")
  var endOfDataIndex int
  for index, part := range parts {
    if index == 0 {
      continue
    }
    if part == "---" {
      endOfDataIndex = index
      break
    }
  }

  dataPart := strings.Join(parts[1: endOfDataIndex], "\n")
  markupPart := strings.Join(parts[endOfDataIndex+1: ], "\n")
  return dataPart, markupPart
}


func ParsePageVariables(s string) map[string]string {
  parts := strings.Split(s, "\n")
  var colonIndex int
  ret := make(map[string]string)
  for _, line := range parts {
    for i, ch := range line {
      if fmt.Sprintf("%c", ch) == ":" {
        colonIndex = i
        break
      }
    }

    if colonIndex == 0 {
      continue
    }

    itemName := strings.ToLower(strings.TrimSpace(line[0: colonIndex]))
    itemValue := strings.TrimSpace(line[colonIndex + 1 : ])
    ret[itemName] = itemValue
  }

  return ret
}


func RenderHTMLToFile(s, path, siteName string) error {
  rootPath, _ := GetRootPath()
  dataPart, markupPart := GetPartsOfMarkup(s)
  pageVariables := ParsePageVariables(dataPart)

  tmpPath := filepath.Join(rootPath, siteName, "tmp")
  os.MkdirAll(tmpPath, 0777)
  tmpMarkupPath := filepath.Join(tmpPath, "m" + UntestedRandomString(10) + ".html")

  os.WriteFile(tmpMarkupPath, []byte(markupPart), 0777)

  tmpl, err := template.ParseFiles(filepath.Join(rootPath, siteName, "templates", pageVariables["template"]), tmpMarkupPath)
  if err != nil {
    return err
  }

  type Context struct {
    Page map[string]string
    TOC []map[string]string
  }

  tocObj := make([]map[string]string, 0)
  if strings.HasSuffix(path, "index.html") && path != "index.html" {
    tocFile := strings.ReplaceAll(path, "index.html", "toc.txt")
    tocPath := filepath.Join(rootPath, siteName, "stuffs", tocFile)
    rawTocFile, err := os.ReadFile(tocPath)
    if err != nil {
      return err
    }

    lines := strings.Split(string(rawTocFile), "\n")
    for _, line := range lines {
      line = strings.TrimSpace(line)
      if line == "" {
        continue
      }
      tocLinePath := strings.ReplaceAll(tocPath, "toc.txt", line)
      rawTocLinePath, err := os.ReadFile(tocLinePath)
      if err != nil {
        return err
      }

      dataPart, _ := GetPartsOfMarkup(string(rawTocLinePath))
      pageVariables := ParsePageVariables(dataPart)
      tocURL := strings.ReplaceAll(tocLinePath, filepath.Join(rootPath, siteName, "stuffs"), "")
      tocURL = strings.ReplaceAll(tocURL, ".md", "")
      pageVariables["url"] = tocURL

      tocObj = append(tocObj, pageVariables)
    }

  } else {
    tocObj = nil
  }

  baseDir := filepath.Dir(filepath.Join(rootPath, siteName, "out", path))
  os.MkdirAll(baseDir, 0777)
  outPathHandle, err := os.Create(filepath.Join(rootPath, siteName, "out", path))
  if err != nil {
    fmt.Println("Got to error Handler")
    return err
  }
  defer outPathHandle.Close()
  writer := bufio.NewWriter(outPathHandle)

  tmpl.Execute(writer, Context{pageVariables, tocObj})
  writer.Flush()

  return nil
}


func RenderMDToFile(s, path, siteName string) error {
  rootPath, _ := GetRootPath()
  dataPart, markupPart := GetPartsOfMarkup(s)
  pageVariables := ParsePageVariables(dataPart)


  tmpl, err := template.ParseFiles(filepath.Join(rootPath, siteName, "templates", pageVariables["template"]),
    filepath.Join(rootPath, siteName, "templates", pageVariables["md_template"]))
  if err != nil {
    return err
  }

  html := string(blackfriday.MarkdownCommon([]byte(markupPart)))

  type Context struct {
    Page map[string]string
    HTML template.HTML
  }

  baseDir := filepath.Dir(filepath.Join(rootPath, siteName, "out", path))
  os.MkdirAll(baseDir, 0777)

  newPath := strings.ReplaceAll(path, ".md", ".html")
  outPathHandle, err := os.Create(filepath.Join(rootPath, siteName, "out", newPath))
  if err != nil {
    return err
  }
  defer outPathHandle.Close()
  writer := bufio.NewWriter(outPathHandle)

  tmpl.Execute(writer, Context{pageVariables, template.HTML(html)})
  writer.Flush()
  return nil
}
