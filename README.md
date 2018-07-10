# HTTP POST MESG service

Emits event on request and process POST calls

# Events

## onRequest

Event key: `onRequest`

Event emitted when server gets a HTTP POST request

| **Key** | **Type** | **Description** |
| --- | --- | --- |
| **body** | `String` | the body of the request |
| **date** | `String` | the date and time of the request |
| **id** | `String` | an unique id of the request |


# Tasks

## execute

Task key: `execute`

Execute HTTP POST call

### Inputs

| **Key** | **Type** | **Description** |
| --- | --- | --- |
| **URL** | `String` | the URL of the request |
| **body** | `String` | the body of the request |


### Outputs

##### error

Output key: `error`



| **Key** | **Type** | **Description** |
| --- | --- | --- |
| **message** | `String` | the error&#39;s message |

##### success

Output key: `success`



| **Key** | **Type** | **Description** |
| --- | --- | --- |
| **body** | `String` | the body of the response |
| **statusCode** | `Number` | the status code of the response |



# Testing

### with go test tool

`go test`

### with mesg test tool

Listen for event

`mesg-core service test --event-filter onRequest`
and execute to http://localhost:8080 POST Request

Test HTTP POST success

`mesg-core service test --task execute --data tests/service_post.json`

Test HTTP POST error

`mesg-core service test --task execute --data tests/service_should_fail.json`

# License

This library is under the MIT License