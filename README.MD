# Notifier

## Considerations

To meet the requirements of the test I have made the following decisions:

* Using the [concurrency pipelines pattern](https://go.dev/blog/pipelines), which in this case is composed
by two stages connected by channels, which will allow concurrent processing (reading and sending) of STDIN messages.
* Using a buffered channel in order to continue reading(not block the channel) the messages from the STDIN in case of the time interval for sending the messages is very high.
* Creating an independent channel to store the errors of each stage and based on the information from this [Google blog](https://go.dev/blog/pipelines), merge the two channels, 
this allows it to monitor the result of the pipeline execution and control if I want to terminate the process (canceling the context) when one or more errors are found. 
* The number of errors allowed can be configured using the environment variable `NOTIFIER_MAXALLOWEDERRORS`.
* Sharing the context with the `gracefullyShutdown` function to listen for the `SIGINT` signal and gracefully terminate the process by canceling the context.

Other aspects of the design of the solution:

* Use of dependency inversion creating an interface for the HTTP client and passing it to the processor together with the file io.Reader interface for reading the STDIN.
* The logs are categorized by level, it is possible to select the logging level of the service with the LogLevel environment variable.
* It is possible to configure different aspects of the application through environment variables.
* Unit tests are included.
* Including Make as interface to execute commands to compile, test, check(linter) the program.
* Docker configuration.

## Environment Variables

* NOTIFIER_INTERVAL: Allows to set a default interval for sending notifications.
* NOTIFIER_MAXALLOWEDERRORS:The application will terminate when it has reached this number of errors
* NOTIFIER_LOGLEVEL:Log level
* NOTIFIER_BUFFERSIZE:Buffer size for the message channel.
* NOTIFIER_HTTPREQUESTTIMEOUT:Configuring the timeout time for the http client.
* NOTIFIER_MAXIDLECONNS: controls the maximum number of idle(keep-alive) connections across all hosts.
* NOTIFIER_MAXIDLECONNSPERHOST: restricts how many connections there are which are not actively serving requests.

## Execute the program

* Locally

```
make compile
./bin/notifier -u http://localhost:8080/notify < messages.txt -i 10ms
```

* Docker

```
make make compile
docker build -t notifier .
docker run -i notifier -u http://localhost:8080/notify < messages.txt -i 200ms
```

## Run tests

```
make test
```

## Run checks
```
make check
```


