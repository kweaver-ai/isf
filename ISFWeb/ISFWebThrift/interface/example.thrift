namespace py example

struct ExampleRequest {
    1: required string message
    2: optional i32 count = 1
}

struct ExampleResponse {
    1: required string result
    2: optional i32 status_code
}

service ExampleService {
    ExampleResponse process(1: required ExampleRequest request)
    string echo(1: required string message)
} 