package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"slices"
	"strings"
)

func executeCommand(fullCommand string) string {

	stringFields := strings.Fields(fullCommand)
	command := stringFields[0]
	args := stringFields[1:]

	var getCommandResult *exec.Cmd
	getCommandResult = exec.Command(command, args...)

	var cmdOutput bytes.Buffer
	getCommandResult.Stdout = &cmdOutput
	err := getCommandResult.Run()

	if err != nil {
		fmt.Errorf("error executing command %s %v", command, args)
	}

	fmt.Printf("Command %s executed with success \n", fullCommand)

	return cmdOutput.String()

}

func transformStringToList(str string, sep string) []string {
	imageList := strings.Split(strings.TrimSpace(str), sep)

	for idx, image := range imageList {
		if !strings.Contains(image, ":") {
			imageList[idx] += ":latest"
		}
	}

	return imageList
}

func removeUnusedImages() {
	allImagesStr := executeCommand("docker images --format {{.Repository}}:{{.Tag}}")

	allImagesList := transformStringToList(allImagesStr, "\n")

	runningImagesStr := executeCommand("docker ps --format {{.Image}}")
	runningImagesList := transformStringToList(runningImagesStr, "\n")

	for _, imageName := range allImagesList {
		if !(slices.Contains(runningImagesList, imageName)) {
			executeCommand("docker rmi " + imageName)
		}
	}
}

func main() {
	removeUnusedImages()
}
