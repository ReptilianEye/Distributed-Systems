from /trader

run to generate golang code from proto file
```bash
protoc --go_out=../server/trader --go_opt=paths=source_relative --go-grpc_out=../server/trader --go-grpc_opt=paths=source_relative trader.proto
```


run to generate python code from proto file
```bash
python -m grpc_tools.protoc -I. --python_out=../client --grpc_python_out=../client trader.proto
```