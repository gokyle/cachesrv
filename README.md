cachesrv - serve working directory
-------------------------------
cachesrv is a web server that serves the current working directory using a
file cache to speed up commonly-accessed files.

Dependencies
------------
`cachesrv` is based on a [file
cache](http://gokyle.github.com/filecache) I wrote in Go.

Compatibility
-------------
`cachesrv` has been tested on the following operating systems:
* OpenBSD (5.2)
* OS X (10.8)


Installation
------------
`go install` will install the binary to the `${GOROOT}/bin` directory.


Usage
-----
```
cachesrv [options] [dir]

Valid options:
        -c certfile     specify the TLS certificate
        -d duration     dump cache stats duration; by default, this is turned
                        off. Must be parsable with time.ParseDuration.
        -e seconds      seconds to expire cache items after; 0 to never expire
                        due to time in cache
        -g seconds      seconds to delay between checks for expired items in
                        the cache.
        -k key          specify the TLS key
        -n items        maximum number of files to store in cache
        -p              the port to listen on
        -r              the server should chroot to the target directory
        -u user         the server should drop privileges to this user
        -v              display the version and exit
```

Specifying a key and certificate will cause the server to listen in TLS
mode. 

Use `^C` to halt the server.


Why require sudo for chroot?
----------------------------
From [chroot(2)](http://www.openbsd.org/cgi-bin/man.cgi?query=chroot&apropos=0&sektion=2&manpath=OpenBSD+Current&arch=i386&format=ascii):

     This call is restricted to the super-user.

The server uses root privileges to chroot to the target directory, then
immediately drops privileges.


Known bugs / caveats
--------------------
setgrp isn't implemented here, as no good solution exists in Go.


History
-------
This is version 3.0.0 of the `srvwd` file server. The [original
version](http://tyrfingr.is/projects/srvwd/)
(the 1.x series) was written in C; a subsequent [rewrite in
Go](http://gokyle.github.com/srvwd) (the 2.x series) used only the standard
library to serve files. This version, the 3.x series, uses a [file
cache](http://gokyle.github.com/filecache) to speed up commonly accessed files.
