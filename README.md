crun
====

Run commands in a chain and concurrently


DOWNLOAD
-------------

Two versions are provided:

- [**recommended**] crun is a Perl script, you may need to install a Perl modules by:

    cpan install Parallel::Runner

- [need more test] crun.go is written in Golang, no third-party packages are used. Please use [gobuild - Cross-Platform Go Project Compiler](http://gobuild.io/download/github.com/shenwei356/crun/go) to compile it. It's simple and fast!
    

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

    crun -s 'ls -lh' -s date
    
    
    crun -n 4 -s job1 -c job2 -c job3 -c job4 -s job5 -s job6

The work flow is:

            |-> job2 |  
    job1 -> |-> job3 | -> job5 -> job6
            |-> job4 |


Copyright
--------

Copyright (c) 2014, Wei Shen (shenwei356@gmail.com)


[MIT License](https://github.com/shenwei356/crun/blob/master/LICENSE)