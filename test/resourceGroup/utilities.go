package main

import (
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"time"
)

func RandomName(prefix string, postfix string) string {
	rand.Seed(time.Now().UnixNano())
	expectedName := fmt.Sprintf("test-%s-%d-%s", prefix, rand.Intn(99999), postfix)
	return expectedName
}

func CopyFile(original string, copy string) {
	log.Println("Copying file " + original + " to " + copy)
	cpCmd := exec.Command("cp", "-rf", original, copy)
	err := cpCmd.Run()
	if err != nil {
		log.Println(err.Error())
		log.Panic(err)
	}
}
