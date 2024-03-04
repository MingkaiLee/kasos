package main

import "time"

func main() {
	print(time.Now().UnixMilli())
	time.Sleep(time.Second)
	print(time.Now().UnixMilli())
}
