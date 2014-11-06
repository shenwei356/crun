package main

import (
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/shenwei356/util/exec"
)

var (
	logErr  *log.Logger
	logOut  *log.Logger
	jobs    []*Job
	threads int
)

const (
	sequential   = 1
	parallelized = 2
	usage        = `
crun - Run commands in a chain and concurrently (V2014.09.29)
    
    Golang edition, by Wei Shen <shenwei356\@gmail.com> 
    http://github.com/shenwei356/crun

USAGE:
    
    crun OPTIONS...

OPTIONS:

    -c STRING    Add a concurrent command
    -s STRING    Add a sequential command
    -n INT       Maximum concurrency (threads number) [CPUs number]
    -h           Help message
    
NOTE:

    The order of options decides the work flow! The two cases below are different:
    
    crun -s ls -s date
    crun -s date -s ls

EXAMPLE:

    crun -n 4 -s job1 -c job2 -c job3 -c job4 -s job5 -s job6

    The work flow is: 

        1. job1 must be executed fist. 
        2. job2,3,4 are independent, so they could be executed parallelly.
        3. job5 must wait for job2,3,4 being done.
        4. job6 could not start before job5 done.
     
    See the workflow graph below, long arrow means long excution time.
        
                  |----> job2     |  
        job1 ---> |--> job3       | -------> job5 --> job6
                  |--------> job4 |

`
)

type Job struct {
	Type int
	Cmds []string
}

func init() {
	logErr = log.New(os.Stderr, "", 0)
	logOut = log.New(os.Stdout, "", 0)

	threads = runtime.NumCPU()
	jobs = make([]*Job, 0)

	typehelp := "\nType \"crun -h\" for help\n"
	if len(os.Args) < 2 {
		logErr.Print("no arguments given" + typehelp)
		os.Exit(1)
	}

	i, nArgs := 1, len(os.Args)
	flag, value := "", ""
	for i <= nArgs-1 {
		flag = os.Args[i]
		i++
		if flag == "-h" || flag == "-help" || flag == "--help" {
			logErr.Print(usage)
			os.Exit(0)
		} else if !strings.HasPrefix(flag, "-") {
			logErr.Printf("invalid option: %s%s", flag, typehelp)
			os.Exit(1)
		}

		if i > nArgs-1 {
			logErr.Printf("no value for option: %s%s", flag, typehelp)
			os.Exit(1)
		}

		value = os.Args[i]
		i++

		if flag == "-n" {
			n, err := strconv.Atoi(value)
			if err != nil || n == 0 {
				logErr.Printf("illegale value for option -n: %s, positive integer needed", flag, typehelp)
				os.Exit(1)
			}
			threads = n
		} else if flag == "-c" {
			if len(jobs) > 0 {
				lastjob := jobs[len(jobs)-1]
				if lastjob.Type == parallelized {
					lastjob.Cmds = append(lastjob.Cmds, value)
				} else {
					jobs = append(jobs, &Job{parallelized, []string{value}})
				}
			} else {
				jobs = append(jobs, &Job{parallelized, []string{value}})
			}
		} else if flag == "-s" {
			jobs = append(jobs, &Job{sequential, []string{value}})
		} else {
			logErr.Printf("invalid option: %s%s", flag, typehelp)
			os.Exit(1)
		}
	}

	if len(jobs) == 0 {
		logErr.Print("no commands added" + typehelp)
		os.Exit(1)
	}
}

func main() {
	runtime.GOMAXPROCS(threads)

	for _, job := range jobs {
		if job.Type == parallelized {
			var wg sync.WaitGroup
			tokens := make(chan int, threads)
			for _, cmd := range job.Cmds {
				tokens <- 1
				wg.Add(1)
				go func(cmd string) {
					defer func() {
						wg.Done()
						<-tokens
					}()
					run(cmd)
				}(cmd)
			}
			wg.Wait()
		} else {
			run(job.Cmds[0])
		}
	}
}

func run(name string) {
	var err error
	cmd, err := exec.Command(name)
	if err != nil {
		logErr.Printf("job: %s: %s", name, err)
		os.Exit(1)
	}

	stderr, err := cmd.StderrChannel()
	if err != nil {
		logErr.Printf("fail to get stderr channel of: %s", name)
		os.Exit(1)
	}

	stdout, err := cmd.StdoutChannel()
	if err != nil {
		logErr.Printf("fail to get stdout channel of: %s", name)
		os.Exit(1)
	}

	go func() {
		var str string
		var ok bool = false
		var errEnd bool = false
		var outEnd bool = false
		for {
			select {
			case str, ok = <-stdout:
				if !ok {
					outEnd = true
					if errEnd {
						return
					}
				} else {
					logOut.Print(str)
				}
			case str, ok = <-stderr:
				if !ok {
					errEnd = true
					if outEnd {
						return
					}
				} else {
					logErr.Print(str)
				}
			}
		}
	}()

	err = cmd.Start()
	if err != nil {
		logErr.Printf("fail to start job: %s", name)
		os.Exit(1)
	}

	err = cmd.Wait()
	if err != nil {
		logErr.Printf("fail to wait job: %s", name)
		os.Exit(1)
	}
}
