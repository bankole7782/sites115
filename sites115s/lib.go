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
