# multithreaded du
Proof of concept multithreaded `du -s` command written in go.

# Install
`make`

# Usage: 
`bin/mdu DIRECTORY_ROOT`

# Performance
Perhaps surprisingly, this is not an IO-bound operation. Here's a non-scientific comparison with the bsd du command on a directory with ~450,000 files (run on a quadcore Macbook Pro):

```
$ time du -s .
137308480   .

real    0m24.239s
user    0m0.377s
sys 0m20.132s
```

```
$ time multithread_du/bin/mdu .
137308480   .

real    0m6.008s
user    0m2.559s
sys 0m27.026s
```
