Implements a simple read/write file server in go.

Exposes the file system directory `FILESERVE_ROOT` at the root URL /.

Setup is via environment variables:

- `FILESERVE_ROOT` - this is the folder in the file system that will be exposed for read/write access, along with all of its subdirectories (required)
- `FILESERVE_IP` - the IP address to expose the server on (optional)
- `FILESERVE_PORT` - the port to expose the server on (defaults to 8080)

To read a file, simply GET the filename. For example:

```
curl http://localhost:8080/test.txt
```

To write a file, PUT or POST to the path you want to write the file, with the data in the request body. So for example:

```
curl -i -X PUT http://localhost:8080/newfile.txt --data-binary "@mydata.txt"
```

If you request a directory, as opposed to a specific filename, a HTML-formatted directory listing will be delivered. This behavior is an implementation side effect, and may be removed at any time.

Note there is no authentication or authorization being done here - anyone with network access can read and write files. This is intended for use in a secure private network environment.