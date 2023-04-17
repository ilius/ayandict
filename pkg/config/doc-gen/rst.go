package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

func codeValue(x any) string {
	return fmt.Sprintf("``%v``", x)
}

func jsonCodeValue(x any) string {
	switch typed := x.(type) {
	case time.Duration:
		return fmt.Sprintf("``%#v``", typed.String())
	case int, float64, string, bool:
		return fmt.Sprintf("``%#v``", x)
	}
	b, _ := json.Marshal(x)
	return fmt.Sprintf("``%s``", string(b))
}

// toml.NewEncoder can only encode maps or structs
func tomlObject(x any) string {
	buf := bytes.NewBuffer(nil)
	encoder := toml.NewEncoder(buf)
	err := encoder.Encode(x)
	if err != nil {
		log.Printf("failed to encode %#v", x)
		panic(err)
	}
	return buf.String()
}

func tableRowSep(width []int, c string) string {
	parts := make([]string, len(width))
	for i, w := range width {
		parts[i] = strings.Repeat(c, w)
	}
	return "+" + c + strings.Join(parts, c+"+"+c) + c + "+"
}

/*
func renderTable(header []string, rows [][]any) string {
	colN := len(header)
	width := make([]int, colN)
	for i:=0; i<colN
}

def renderTable(rows):
	"""
		rows[0] must be headers
	"""
	colN = len(rows[0])
	width = [
		max(
			max(len(line) for line in row[i].split("\n"))
			for row in rows
		)
		for i in range(colN)
	]
	rowSep = tableRowSep(width, "-")
	headerSep = tableRowSep(width, "=")

	lines = [rowSep]
	for rowI, row in enumerate(rows):
		newRows = []
		for colI, cell in enumerate(row):
			for lineI, line in enumerate(cell.split("\n")):
				if lineI >= len(newRows):
					newRows.append([
						" " * width[colI]
						for colI in range(colN)
					])
				newRows[lineI][colI] = line.ljust(width[colI], " ")
		for row in newRows:
			lines.append("| " + " | ".join(row) + " |")
		if rowI == 0:
			lines.append(headerSep)
		else:
			lines.append(rowSep)

	# widthsStr = ", ".join([str(w) for w in width])
	# header = f".. table:: my table\n\t:widths: {widthsStr}\n\n"
	# return header + "\n".join(["\t" + line for line in lines])

	return "\n".join(lines)
*/
