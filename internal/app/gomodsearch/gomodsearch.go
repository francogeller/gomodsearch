package gomodsearch

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/mod/semver"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type Mod struct {
	FullName           string
	Name               string
	Version            string
	IsDirectDependency bool
	Parents            []*Mod
}

func Run(path string, mods ...string) error {

	// Print search params
	for _, mod := range mods {

		// Load data
		out, err := load(&path)
		if err != nil {
			return err
		}

		// Init map
		modMap := make(map[string]*Mod)
		verMap := make(map[string][]*Mod)

		// Decode data
		for _, row := range bytes.Split(out[:len(out)-1], []byte("\n")) {
			lineSplit := bytes.Split(row, []byte(" "))
			parent := string(lineSplit[0])
			dep, depName, depVersion := decodeMod(&lineSplit[1])
			isDirectDependency := isDirectDependency(&parent)
			if existingMod, ok := modMap[dep]; ok {
				existingMod.Parents = append(existingMod.Parents, modMap[parent])
			} else {
				newMod := &Mod{
					FullName:           dep,
					Name:               depName,
					Version:            depVersion,
					IsDirectDependency: isDirectDependency,
					Parents:            []*Mod{},
				}
				modMap[dep] = newMod
				if !isDirectDependency {
					newMod.Parents = append(newMod.Parents, modMap[parent])
				}
				if _, ok := verMap[depName]; !ok {
					verMap[depName] = []*Mod{newMod}
				} else {
					verMap[depName] = append(verMap[depName], newMod)
				}
			}
		}

		_, modVersion := splitMod(&mod)

		if modVersion != "" {
			// Search mod
			search(&mod, &modMap)
		} else {
			if mods, ok := verMap[mod]; ok {
				// Order mods
				sort.Slice(mods, func(i, j int) bool {
					return semver.Compare(mods[i].Version, mods[j].Version) < 0
				})
				// Search for each mod
				for _, m := range mods {
					search(&(m.FullName), &modMap)
				}
			} else {
				fmt.Fprintf(os.Stdout, "%s\n", mod)
				fmt.Fprintf(os.Stdout, "└── No direct dependency found\n\n")
			}
		}
	}

	return nil
}

func load(path *string) ([]byte, error) {
	// Check go
	_, err := exec.LookPath("go")
	if err != nil {
		return nil, err
	}

	// Init
	cmd := exec.Command("go", "mod", "graph")
	cmd.Dir = *path
	cmd.Stdin = strings.NewReader("")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Run command
	err = cmd.Run()
	if err != nil {
		if stderr.Len() == 0 {
			return nil, err
		}
		return nil, errors.New(stderr.String())
	}
	return stdout.Bytes(), nil
}

func isDirectDependency(mod *string) bool {
	return !strings.Contains(*mod, "@")
}

func decodeMod(mod *[]byte) (raw string, name string, version string) {
	raw = string(*mod)
	name, version = splitMod(&raw)
	return
}

func splitMod(mod *string) (name string, version string) {
	modSplit := strings.Split(*mod, "@")
	if len(modSplit) == 2 {
		name = modSplit[0]
		version = modSplit[1]
	}
	return
}

func findDirectDependencies(mod *Mod, directDependencies *[]*Mod, processedDependencies *map[string]interface{}) {
	if mod.IsDirectDependency {
		*directDependencies = append(*directDependencies, mod)
	}
	(*processedDependencies)[fmt.Sprintf("%p", mod)] = nil
	for _, parent := range mod.Parents {
		if _, ok := (*processedDependencies)[fmt.Sprintf("%p", parent)]; !ok {
			findDirectDependencies(parent, directDependencies, processedDependencies)
		}
	}
	return
}

func search(mod *string, modMap *map[string]*Mod) {
	modch, ok := (*modMap)[*mod]
	fmt.Fprintf(os.Stdout, "%s\n", *mod)
	if ok {
		var directDependencies []*Mod
		processedDependencies := make(map[string]interface{})
		findDirectDependencies(modch, &directDependencies, &processedDependencies)
		for _, dp := range directDependencies {
			fmt.Fprintf(os.Stdout, "├── %s\n", dp.FullName)
		}
		fmt.Fprintf(os.Stdout, "└── %d direct dependencies\n\n", len(directDependencies))
	} else {
		fmt.Fprintf(os.Stdout, "└── No direct dependency found\n\n")
	}
}
