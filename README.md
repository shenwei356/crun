crun
====

Run workflow. Generally all parts are sequential executed, and some parts could be done parallelly.


DOWNLOAD
-------------

Two editions are provided:

- [**recommended**] perl/crun is a Perl script. You may need to install a Perl module by:

    cpan install Parallel::Runner

- [need more test] go/crun.go is written in Golang, no third-party packages are used. Please use [gobuild - Cross-Platform Go Project Compiler](http://gobuild.io/download/github.com/shenwei356/crun/go) to compile it. It's simple and fast!
    

USAGE
-----
   
    crun OPTIONS...

    OPTIONS:

        -c STRING    Add a concurrent command
        -s STRING    Add a sequential command
        -n INT       Maximum concurrency (threads number) [4]
        -h           Help message
        
Note: The order of options decides the work flow! The two cases below are different:
    
    crun -s ls -s date
    crun -s date -s ls

EXAMPLE
-------    
    
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
     
You can also concurrently run commands from STDIN:

    cat jobs.list | crun -t 8 -s "echo start" -stdin -s "echo end" 


Copyright
--------

Copyright (c) 2014, Wei Shen (shenwei356@gmail.com)


[MIT License](https://github.com/shenwei356/crun/blob/master/LICENSE)