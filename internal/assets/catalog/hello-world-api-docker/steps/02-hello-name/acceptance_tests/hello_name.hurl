GET http://api:8080/hello?name=John

HTTP 200
[Asserts]
`Hello John!`

#---

GET http://api:8080/hello?name=Alice

HTTP 200
[Asserts]
`Hello Alice!`

##---

GET http://api:8080/hello?name=Bob

HTTP 200
[Asserts]
`Hello Bob!`
