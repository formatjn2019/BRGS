package utils

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

func OsOperation(operation string, path ...string) {
	fmt.Println(operation, path)
	switch operation {
	//create dir
	case "md":
		os.MkdirAll(path[0], os.ModeDir)
		//create file
	case "cf":
		os.WriteFile(path[0], []byte("create file"), 0644)
		//delete
	case "rm":
		//update
		os.RemoveAll(path[0])
	case "up":
		origin, _ := ioutil.ReadFile(path[0])
		p := make([]byte, 100)
		rand.Read(p)
		origin = append(origin, p[rand.Intn(90):]...)
		os.WriteFile(path[0], origin, 0644)
	case "mv":
		os.Rename(path[0], path[1])
	}
	time.Sleep(1e9)
}
