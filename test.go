package main

// import (
// 	"fmt"
// 	"fmt"

// 	"github.com/vmihailenco/msgpack"
// )

// type TStruct struct {
// 	M    map[string]string
// 	Data []byte
// }

// func main() {
// 	tp := new(TStruct)
// 	tp.M = map[string]string{
// 		"hello": "hi",
// 	}
// 	tp.Data = []byte("hello world")

// 	b, err := msgpack.Marshal(tp)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("%s", b)

// 	tp2 := new(TStruct)
// 	err = msgpack.Unmarshal(b, tp2)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("%v", tp2)
// }

// func (t *TStruct) String() string {
// 	return fmt.Sprintf("%v %s", t.M, t.Data)
// }
