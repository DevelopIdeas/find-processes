package main

import (
  "strings"
  "regexp"
  "fmt"
  "os/exec"
  "github.com/alexflint/go-arg"
)

var args struct {
	Port int `help:"Processes binding to given Port"`
	Search string `help:"Search for running processes / scripts containing string"`
}

func portSearch(port int) {
  var lsofRe = regexp.MustCompile(`(?m)^([a-zA-Z]{1,}.*?)\s{1,}(\d{1,})\s{1,}([a-zA-Z]{1,}.*?)\s{1,}\s{1,}`)
  cmd := exec.Command("lsof", "-i", ":"+fmt.Sprint(port))
  stdout, err := cmd.CombinedOutput()
  if err != nil {
    fmt.Println(err.Error())
    return
  }
  stdoutStr := string(stdout)

  res := lsofRe.FindAllStringSubmatch(stdoutStr, -1)
  for i, _ := range res {
    // _ := res[i][0] // full match
    // _ := res[i][1] // command
    // _ := res[i][3] // user
    pid := res[i][2] // pid

    procSearch(pid)
  }
}

func procSearch(search string) {
  grep := exec.Command("grep", search)
  ps := exec.Command("ps", "-eo", "pid,args")

  // Get ps's stdout and attach it to grep's stdin.
  pipe, _ := ps.StdoutPipe()
  defer pipe.Close()

  grep.Stdin = pipe

  // Run ps first.
  ps.Start()

  // Run and get the output of grep.
  res, _ := grep.Output()

  psOutStr := string(res)
  // fmt.Println(psOutStr, "ps out")

  var psRe = regexp.MustCompile(`(?m)^^(\d{1,})\s{1,}(.*)$`)
  psReRes := psRe.FindAllStringSubmatch(psOutStr, -1)

  for ii, _ := range psReRes {
    cmd2 := psReRes[ii][2] // cmd
    if !strings.HasPrefix(cmd2, "grep") {
      pid2 := psReRes[ii][1] // pid
      fmt.Printf("PID=%s CMD=%s\n", pid2, cmd2)
      // fmt.Println(cmd2, "cmd2")
    }
  }
}

func main() {
  arg.MustParse(&args)

  if args.Port > 0 {
    fmt.Printf("Find Port \"%d\"\n", args.Port)
    portSearch(args.Port)
  } else if args.Search != "" {
    fmt.Printf("Find \"%s\"\n", args.Search)
    procSearch(args.Search)
  }
}