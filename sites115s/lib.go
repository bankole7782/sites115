package sites115s

import (
  "os"
  "encoding/json"
  "strings"
  "fmt"
  "github.com/kljensen/snowball"

)


var ALLOWED_CHARS = "abcdefghijklmnopqrstuvwxyz0123456789"
var STOP_WORDS []string


type PaginatorStruct struct {
  Page int
  PaginationCount int
  Pages []map[string]string
  TotalPages int
  TotalPagesArr []int
}


func init() {
  stopWordsList := make([]string, 0)
  err := json.Unmarshal(stopwordsBytes, &stopWordsList)
  if err != nil {
    panic(err)
  }

  STOP_WORDS = stopWordsList
}


func CleanWord(word string) string {
  word = strings.ToLower(word)

  allowedCharsList := strings.Split(ALLOWED_CHARS, "")

  if strings.HasSuffix(word, "'s") {
    word = word[: len(word) - len("'s")]
  }

  newWord := ""
  for _, ch := range strings.Split(word, "") {
    if FindIn(allowedCharsList, ch) != -1 {
      newWord += ch
    }
  }

  var toWriteWord string
  stemmed, err := snowball.Stem(newWord, "english", true)
  if err != nil {
    toWriteWord = newWord
    fmt.Println(err)
  }
  toWriteWord = stemmed

  return toWriteWord
}

func DoesPathExists(p string) bool {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return false
	}
	return true
}



func FindIn(container []string, elem string) int {
	for i, o := range container {
		if o == elem {
			return i
		}
	}
	return -1
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
