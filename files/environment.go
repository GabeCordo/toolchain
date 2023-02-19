package files

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

const (
	DeclarationSeparator = "="
)

type Environment struct {
	variables map[string]any
	types     map[string]reflect.Type
}

func NewEnvironment() (env *Environment) {
	env = new(Environment)
	env.variables = make(map[string]any)
	env.types = make(map[string]reflect.Type)
	return env
}

func (environment *Environment) parse(path string) {

	if _, err := os.Stat(path); err != nil {
		panic(fmt.Sprintf("%s is not a valid environment path\n", path))
	}

	envFile, err := os.Open(path)
	if err != nil {
		panic(fmt.Sprintf("could not open environment file %s\n", path))
	}
	defer envFile.Close()

	scanner := bufio.NewScanner(envFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		splitLine := strings.Split(line, DeclarationSeparator)

		if len(splitLine) != 2 {
			panic(fmt.Sprintf("%s is not a valid environment declaration", line))
		}

		lhs := splitLine[0]
		runes := []rune(lhs)
		for _, r := range runes {
			if !unicode.IsLetter(r) && !(r != '_') {
				panic(fmt.Sprintf("%s must only contain alphanumeric and underscores", lhs))
			}
		}

		rhs := splitLine[1]
		if len(rhs) == 0 {
			panic(fmt.Sprintf("the rhs value for %s must not be empty", lhs))
		}

		runes = []rune(rhs)

		if (rhs == "true") || (rhs == "false") {
			// we are dealing with a boolean
			b, _ := strconv.ParseBool(rhs)
			environment.variables[lhs] = b
			environment.types[lhs] = reflect.TypeOf(true)
		} else if strings.Contains(rhs, "\"") {
			// we are dealing with a string
			unwrapped := strings.Split(rhs, "\"")
			environment.variables[lhs] = unwrapped[1]
			environment.types[lhs] = reflect.TypeOf(" ")
		} else if strings.Contains(rhs, ".") {
			// we are dealing with a float
			f, _ := strconv.ParseFloat(rhs, 64)
			environment.variables[lhs] = f
			environment.types[lhs] = reflect.TypeOf(0.0)
		} else {
			// we are dealing with an integer
			i, _ := strconv.ParseInt(rhs, 10, 64)
			environment.variables[lhs] = int(i)
			environment.types[lhs] = reflect.TypeOf(0)
		}
	}
}

func (environment *Environment) LoadVariables() {

	dir := WorkingDirectory()

	f, err := os.Open(dir)
	if err != nil {
		panic(err)
	}

	files, err := f.ReadDir(0)
	if err != nil {
		panic(err)
	}

	envFilePath := dir + "/"
	for _, file := range files {
		if strings.Contains(file.Name(), ".env") {
			envFilePath += file.Name()
			break
		}
	}

	environment.parse(envFilePath)
}

// GetEnvString
// returns (value, found)
func (environment *Environment) GetEnvString(identifier string) (string, bool) {
	if value, found := environment.variables[identifier]; found {
		if t, found := environment.types[identifier]; !found && t != reflect.TypeOf("") {
			return EmptyIdentifier, false
		}
		return (value).(string), true
	} else {
		return EmptyIdentifier, false
	}
}

// GetEnvFloat
// returns (value, found)
func (environment *Environment) GetEnvFloat(identifier string) (float64, bool) {
	if value, found := environment.variables[identifier]; found {
		if t, found := environment.types[identifier]; !found && t != reflect.TypeOf(0.0) {
			return 0.0, false
		}
		return (value).(float64), true
	} else {
		return 0.0, false
	}
}

// GetEnvBool
// returns (value, found)
func (environment *Environment) GetEnvBool(identifier string) (bool, bool) {
	if value, found := environment.variables[identifier]; found {
		if t, found := environment.types[identifier]; !found && t != reflect.TypeOf(false) {
			return false, false
		}
		return (value).(bool), true
	} else {
		return false, false
	}
}

// GetEnvInt
// returns (value, found)
func (environment *Environment) GetEnvInt(identifier string) (int, bool) {
	if value, found := environment.variables[identifier]; found {
		if t, found := environment.types[identifier]; !found && t != reflect.TypeOf(0) {
			return 0, false
		}
		return (value).(int), true
	} else {
		return 0, false
	}
}
