Simple Image Server
===================

This is an implementation of the Image HTTP server. 

Build
=====

```bash
$ make build
```

Run
===

Session #1:
```bash
$ ./sis -pretty-json
```

Session #2:
```bash
$ curl -s http://localhost:8000/generate/png/100/100 > /dev/null
$ curl -s http://localhost:8000/generate/png/200/20 > /dev/null
$ curl -s http://localhost:8000/stats
{
  "num_images": 2,
  "average_width_px": 150,
  "average_height_px": 60
}
```

Help
====

```bash
$ ./sis --help
Usage of ./sis:
  -addr string
    	IP address to lister (default "127.0.0.1")
  -port string
    	Http port to lister (default "8000")
  -pretty-json
    	Return pretty json in stats
```

Tests
=====

```bash
$ make test
```
