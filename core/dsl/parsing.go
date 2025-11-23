package dsl

import (
	"strings"

	"github.com/Ether-Security/leviathan/utils"
)

func ExtractNmapIP(src, dest string) {
	var result []string
	for _, line := range utils.ReadFile(src) {
		split := strings.Split(line, " ")
		if split[1] != "Nmap" {
			result = append(result, split[1])
		}
	}
	utils.WriteFile(dest, strings.Join(result, "\n"))
}
