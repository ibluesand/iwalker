package main
import "strings"

func main() {

    var s = "@baidu:你好"
    c := []rune(s)
    print(len(c))
    print(strings.Index(s,":"))
    print(string(c[7:]))
}
