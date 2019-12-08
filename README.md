# testoy
Yet another test automation library.

Imagine you have a yaml seen as below. It allows you to define your request and responses also provides request chaining.
```
- case: 
    title: it should create something
    id: case1
    request:
        timeout: 3s
        method: POST
        url: https://jsonplaceholder.typicode.com/posts
        headers:
            - content-type: application/json
        body: '{"userId": 1, "id": 31, "title":"test title", "body"}'
    response:
        statusCode: 201
        body: '{"id": 102}'
        headers:
            - content-type: application/json
    then: 
        caseId: case2
        args:
            - id: body.id
- case: 
    title: it should create something
    id: case2
    request:
        timeout: 10s
        method: POST
        url: https://jsonplaceholder.typicode.com/posts/{{id}}
        header:
            - content-type: application/json
        body: '{"name": "test"}'
    response:
        statusCode: 201
        body: '{"id": 1}'
        header:
            - content-type: application/json
```

### TODO:

- Provide test file name as run argument
- Request chaining
- Paralel test execution
- User friendly test results