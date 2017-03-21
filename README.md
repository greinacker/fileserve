# Fileserve

Implements a read/write file server in go, with optional signature verification for writes. This server is intended for use on secure private networks.

Exposes the file system directory `FILESERVE_ROOT` at the root URL /.

### Setup

Setup is via environment variables:

- `FILESERVE_ROOT` - this is the folder in the file system that will be exposed for read/write access, along with all of its subdirectories (required)
- `FILESERVE_IP` - the IP address to expose the server on (optional)
- `FILESERVE_PORT` - the port to expose the server on (defaults to 8080)
- `SIGN_SECRET` - secret key used to verify write requests are authentic. If present, all requests must be signed appropriately. If this variable is not present, all writes will succeed.

### Basic usage

To start the server:

```
$ SIGN_SECRET="1234567890" FILESERVE_ROOT=`pwd` fileserve
```

To read a file, simply GET the filename. For example:

```
curl http://localhost:8080/test.txt
```

To write a file, PUT or POST to the path you want to write the file, with the data in the request body. So for example:

```
curl -i -X PUT http://localhost:8080/newfile.txt --data-binary "@mydata.txt"
```

### Signing requests

The goal of signed requests is to verify that the client has the secret to sign the request with. The other parameters are present to prevent casual replay-style attacks if a particular request is intercepted on the private network. A given signature, for a given file, will remain valid for 3 minutes, during which time the exact same content could be re-saved to the same file path.

If `SIGN_SECRET` is specified, all requests must be signed, with a HTTP request header that looks like

```
FileserveSignature: 1490135003:436e7a178376123f71fe813d479e0d11c545f488
```

If the signature is invalid, the server will return 403 Forbidden.

The signature is made up of two parts:

```
timestamp:signature
```

The "timestamp" part is simply the number of seconds from the Unix epoch. The signature is calculated as follows:

1. Concatenate the following strings together:
   - the timestamp used above, converted to a string
   - the signing secret, as specified by `SIGN_SECRET`
   - the path for the request (e.g. /filename.txt or /A/B/file.txt)
   - the contents to be written to the file
2. Calculate the SHA1 hash of the concatenated string
3. Write out the hash as a hex string, with all letters in lower case.

Here is an example, in Ruby, of making a properly signed request:

```
def save_to_fileserve(filepath, str)
  sign_token = "1234567890"

  ts = Time.now.to_i
  sig = Digest::SHA1.hexdigest(ts.to_s + sign_token + filepath + str)
  signature = "#{ts}:#{sig}"

  uri_str = "http://localhost:8080" + filepath
  response = nil
  url = URI.parse(uri_str)
  http = Net::HTTP.new(url.host, url.port)
  req = Net::HTTP::Put.new(url.request_uri)
  req["FileserveSignature"] = signature
  req.body = str
  response = http.start {|http| http.request(req) }

  puts response.code
end
```

And an example request in curl:

```
curl -i \
  -H "FileserveSignature: 1490134367:bd51661eff323ff4315fe91f3088f700c34dbd1e" \
  -X PUT http://localhost:8080/newfile.txt \
  --data-binary "@README.md"
```

### Other

If you request a directory, as opposed to a specific filename, a HTML-formatted directory listing will be delivered. This behavior is an implementation side effect, and may be removed at any time.

Note there is no authentication or explicit authorization being done here - anyone with network access can read and write files. This is intended for use in a secure private network environment. The signature features are here to remove the casual attack surface, in the event the network interface is inadvertently (and temporarily) exposed, but should not be considered hardened.
