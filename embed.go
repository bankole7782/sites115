package sites115

import (
	_ "embed"
)

//go:embed stopwords_en.txt
var StopWordsBytes []byte
