package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type RowText struct {
	Row  int
	Text string
}

func hasValidQuotes(input string) bool {
	count := strings.Count(input, "\"")

	return count%2 == 0
}

func lintCommand(cmd string, rowNumber int, fullCommand string, totalCommandRows int) {
	switch cmd {
	case "FROM":
		if strings.Contains(fullCommand, ":latest") {
			fmt.Printf("Line %d: The use of latest images is not recommended\n", rowNumber)
		}
	case "ADD":
		fmt.Printf("Line %d: Recommended the use of COPY instead of ADD\n", rowNumber)
	case "RUN":
		if totalCommandRows > 1 {
			fmt.Printf("Warning Line %d: Recommended to group the RUN command into a single line using `&&`, this avoid multiple layers optimizing the image\n", rowNumber)
		}
	case "USER":
		if strings.Contains(strings.ToLower(fullCommand), "root") {
			fmt.Printf("Warning Line %d: The user `root` is not recommended \n", rowNumber)
		}
	case "ENTRYPOINT":
		if !hasValidQuotes(fullCommand) {
			fmt.Printf("ERROR Line %d: Quotes not closed\n", rowNumber)
		}

		if !strings.Contains(strings.ToLower(fullCommand), "[") && !strings.Contains(strings.ToLower(fullCommand), "]") {
			fmt.Printf("ERROR Line %d: The `ENTRYPOINT` only allows `EXEC FORM`, example `[\"ping\", \"-c\", \"3\"]\n", rowNumber)
		} else if !strings.Contains(strings.ToLower(fullCommand), "[") || !strings.Contains(strings.ToLower(fullCommand), "]") {
			fmt.Printf("Warning Line %d: Square-brackets not closed\n", rowNumber)
		}
	}
}

func dockerfileLint() {
	file, err := os.Open("C:/Users/mathe/code/go_docker_codes/Dockerfile")

	if err != nil {
		fmt.Errorf("error: %w", err)
	}
	defer file.Close()

	commandsMap := mapDockerfileCommands(file)
	recommendedCommands := []string{"WORKDIR", "EXPOSE"}
	for _, command := range recommendedCommands {
		if _, ok := commandsMap[command]; !ok {
			fmt.Printf("Warning: It's recommended to set a %s in the Dockerfile\n", command)
		}
	}

	for cmd, details := range commandsMap {
		for _, detail := range details {
			lintCommand(cmd, detail.Row, detail.Text, len(details))
		}

	}

}

func mapDockerfileCommands(file *os.File) map[string][]RowText {
	var commandsMap = make(map[string][]RowText)
	scanner := bufio.NewScanner(file)
	rowLine := 1
	for scanner.Scan() {
		lineText := scanner.Text()
		if len(lineText) == 0 || strings.Contains(lineText, "#") {
			continue
		}

		commandString := strings.Split(strings.TrimSpace(lineText), " ")[0]
		commandsMap[commandString] = append(commandsMap[commandString], RowText{Row: rowLine, Text: lineText})

		rowLine++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	return commandsMap
}

func main() {

	dockerfileLint()
}
