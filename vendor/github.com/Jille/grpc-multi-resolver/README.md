# grpc-multi-resolver

Package multiresolver allows you to Dial to multiple hosts/IPs as a single ClientConn.

Make sure to import this package:

```go
import _ "github.com/Jille/grpc-multi-resolver"
```

and then you can use it with grpc.Dial():

```go
grpc.Dial("multi:///127.0.0.1:1234,dns://example.org:1234")
```

Note the triple slash at the beginning.

Note: The ServiceConfig and Attributes from the first target are used, the rest are ignored.
