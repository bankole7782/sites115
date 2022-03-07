package main

import (
  "encoding/json"
  "path/filepath"
  "strings"
  "os"
  "fmt"
  "strconv"
  "github.com/microcosm-cc/bluemonday"
  "github.com/PuerkitoBio/goquery"
  "github.com/bankole7782/sites115/sites115s"
)

func renderIndexes(sitePath string) {

  indexesPath := filepath.Join(sitePath, "out", "_indexes")
  os.RemoveAll(indexesPath)
  os.MkdirAll(indexesPath, 0777)

  dir := filepath.Join(sitePath, "out")

  // get all pages
  allPages := make([]string, 0)
  err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
    if err != nil {
      return err
    }


    if ! info.IsDir() {
      if ! strings.HasSuffix(dir, "/") {
        dir += "/"
      }

      toRemove := filepath.Join(sitePath, "out")
      if strings.HasPrefix(path, filepath.Join(toRemove, "static")) {
        return nil
      }
      if strings.HasPrefix(path, filepath.Join(toRemove, "_indexes")) {
        return nil
      }
      if strings.HasPrefix(path, filepath.Join(toRemove, "_templates")) {
        return nil
      }
      if strings.HasPrefix(path, filepath.Join(toRemove, "_page_descs")) {
        return nil
      }
      pathToWrite := strings.Replace(path, dir, "", 1)
      allPages = append(allPages, pathToWrite)

    }
    return nil
  })

  if err != nil {
    panic(err)
  }

  allPagesMap := make(map[int]string)
  for i, page := range allPages {
    allPagesMap[i+1] = page
  }

  jsonBytes, err := json.Marshal(allPagesMap)
  if err != nil {
    panic(err)
  }
  os.WriteFile(filepath.Join(sitePath, "out", "allpages.json"), jsonBytes, 0777)

  // generate indexes
  err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
    if err != nil {
      return err
    }

    if ! info.IsDir() {
      if ! strings.HasSuffix(dir, "/") {
        dir += "/"
      }

      toRemove := filepath.Join(sitePath, "out")
      if strings.HasPrefix(path, filepath.Join(toRemove, "static")) {
        return nil
      }
      if strings.HasPrefix(path, filepath.Join(toRemove, "_indexes")) {
        return nil
      }
      if strings.HasPrefix(path, filepath.Join(toRemove, "_templates")) {
        return nil
      }
      if strings.HasPrefix(path, filepath.Join(toRemove, "_page_descs")) {
        return nil
      }
      if strings.HasSuffix(sitePath, "allpages.json") {
        return nil
      }

      pathToWrite := strings.Replace(path, dir, "", 1)
      index := 0
      for k, v := range allPagesMap {
        if v == pathToWrite {
          index = k
          break
        }
      }

      doIndex(path, sitePath, strconv.Itoa(index))
    }
    return nil
  })

  if err != nil {
    panic(err)
  }

  os.MkdirAll(filepath.Join(sitePath, "out", "_page_descs"), 0777)

  // get all dataPart's and save
  walkingDir := filepath.Join(sitePath, "stuffs")
  err = filepath.Walk(walkingDir, func(path string, info os.FileInfo, err error) error {
    if err != nil {
      return err
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

      dataPart, _ := sites115s.GetPartsOfMarkup(string(raw))

      pageVariables := sites115s.ParsePageVariables(dataPart)
      pathToWrite := strings.Replace(path, walkingDir, "", 1)
      pathToWrite = strings.ReplaceAll(pathToWrite, ".md", ".html")

      index := 0
      for k, v := range allPagesMap {
        if "/" + v == pathToWrite {
          index = k
          break
        }
      }

      jsonBytes, _ := json.Marshal(pageVariables)
      os.WriteFile(filepath.Join(sitePath, "out", "_page_descs", strconv.Itoa(index) + ".json"), jsonBytes, 0777)
    }
    return nil
  })

  if err != nil {
    panic(err)
  }

}


func doIndex(textPath, sitePath, index string) {
	raw, err := os.ReadFile(textPath)
	if err != nil {
		return
	}

  doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(raw)))
  if err != nil {
    panic(err)
  }
  html, err := doc.Find("body").Html()
  if err != nil {
    panic(err)
  }

	textStrippedOfHtml := bluemonday.StrictPolicy().Sanitize(html)
	words := strings.Fields(textStrippedOfHtml)

	wordCountMap := make(map[string]int64)
	for _, word := range words {
		// clean the word.
		word = sites115s.CleanWord(word)
		if word == "" {
			continue
		}
		if sites115s.FindIn(sites115s.STOP_WORDS, word) != -1 {
			continue
		}

		oldCount, ok := wordCountMap[word]
		if ! ok {
			wordCountMap[word] = 1
		} else {
			wordCountMap[word] = oldCount + 1
		}
	}

	for word, wordCount := range wordCountMap {
		dirToMake := filepath.Join(sitePath, "out", "_indexes", word)
		os.MkdirAll(dirToMake, 0777)
		err = os.WriteFile(filepath.Join(dirToMake, index), []byte(fmt.Sprintf("%d", wordCount)), 0777)
		if err != nil {
			fmt.Printf("word is : '%s'\n", word)
      fmt.Println(err)
      return
		}
	}

}
