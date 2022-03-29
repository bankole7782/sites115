package sites115s

import (
  "strings"
  "math"
  "time"
)

func ToLower(s string) string {
  return strings.ToLower(s)
}

func ToUpper(s string) string {
  return strings.ToUpper(s)
}

func Modulo(num, num2 int) int {
  tmp :=  math.Mod(float64(num), float64(num2))
  return int(tmp)
}

func Plus(num, num2 int) int {
  return num + num2
}

func Minus(num, num2 int) int {
  return num - num2
}

func ToLongDate(s string) string {
  t, err := time.Parse("2006-01-02", s)
  if err != nil {
    return ""
  }

  return t.Format("02 January, 2006")
}
