package main

import (
	"flag"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"sort"
)

type Language struct {
	name          string
	compileString string
	runString     string
}

type RunConfiguration struct {
	compileString string
	runString     string
}

type Result struct {
	language string
	runtime  unixNano
}

type unixNano int64

func (time unixNano) String() string {
	return fmt.Sprintf("%d.%d", time/1000000000, (time%1000000000)/1000000)
}

func getRunTime(start time.Time, end time.Time) int64 {

	return end.UnixNano() - start.UnixNano()

}

func runString(s string) string {
	command := strings.Split(s, " ")
	args := command[1:]
	cmd := exec.Command(command[0], args...)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}

	output := string(out)

	return output
}

func (language Language) compileProgram() {
	if language.compileString != "" {
		runString(language.compileString)
	}
}

func (language Language) runProgram(target int64, correctResult int64) int64 {
	var start = time.Now()
	var out = runString(language.runString + " " + strconv.FormatInt(target, 10))
	var output, _ = strconv.ParseInt(out, 10, 64)
	if output != correctResult {
		//return -1
	}
	var end = time.Now()
	var programTime = getRunTime(start, end)
	return programTime
}

var java = Language{"Java",
	"javac prime.java",
	"java prime"}

var python = Language{"Python", "", "python prime.py"}
var python2 = Language{"Python (pypy)", "", "pypy prime.py"}

var cplusplus = Language{"C++ (clang)", "clang++ -O3 -o plusplus prime.cpp", "./plusplus"}
var cplusplus2 = Language{"C++ (gcc)", "g++ -O3 -o plusplus2 prime.cpp", "./plusplus2"}

var d = Language{"D (ldc)", "ldc2 -O3 -of primeD prime.d", "./primeD"}
var rust = Language{"Rust", "rustc -C opt-level=2 -o primeRust prime.rs", "./primeRust"}
var node = Language{"Node", "", "node prime.js"}
var Go = Language{"Go", "go build -o primeGo prime.go", "./primeGo"}
var cSharp = Language{"C# (mcs)", "mcs prime.cs", "./prime.exe"}

type ByTime []Result
func (a ByTime) Len() int           { return len(a) }
func (a ByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTime) Less(i, j int) bool { return a[i].runtime < a[j].runtime }

func main() {
	var target, result int

	flag.IntVar(&target, "target", 100, "Tell program to run from 2 to target")
	flag.IntVar(&result, "result", 23, "What is the correct result")
	flag.Parse()

	languages := []Language{java, python, python2, cplusplus, cplusplus2, d, rust, node, Go, cSharp}

	times := make([]Result,0)

	for _, language := range languages {
		language.compileProgram()
		programTime := language.runProgram(int64(target), int64(result))

		times =  append(times,Result{language.name,unixNano(programTime)})
	}

	/*
	Works from Go 1.8
	sort.Slice(times,func(i, j unixNano) bool {
		return i < j
	})
	*/

	sort.Sort(ByTime(times))

	for _, result := range times {
		fmt.Printf("%s: %s\n", result.language, result.runtime)
	}
}
