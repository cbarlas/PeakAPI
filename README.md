# PeakAPI

A RESTful web API for consuming Event data and keeping track of the response times of the requests.
Written in GO Programming Language.

### Prerequisites
In this app, I used MariaDB as external storage for data in my local machine. It is a new RDBMS and considered as
an upgraded version of MySQL, and is also developed by the original developers of MySQL.
```
GO Environment
```
```
MariaDB as external data storage
```
Here are the links for installing guides:
* [GO](https://golang.org/doc/install)
* [MariaDB](https://downloads.mariadb.org/mariadb/repositories/)

### Installing

Assuming you have installed GO Environment, to download the project, open Terminal and type:
```bash
go get github.com/cbarlas/RESTful-Event-API
```
After that you have downloaded the project, go to the project directory and:
```bash
go install
```

### Configurations

Before running the app, some configurations might be necessary.
First, head to the `db/database.go` file

```golang
...
// Constant variables for creating and connecting to database instance
const (
	DB_DRIVER = "mysql"
	DB_USERNAME = "root"
	DB_PW = "root"
	CONNECTION = DB_USERNAME + ":" + DB_PW + "@tcp(127.0.0.1:3306)/"
	DB_NAME = "api_db"
	RESPONSE_TABLE = "responses"
)
...
```
In here, the default `DB_USERNAME` and `DB_PW` constant variables are for MariaDB database connection. Please change them
as your database username and password.

## Example Usage
Now take a look at `apis/apis.go` file.
```golang
// All of the predefined api keys are stored here in map
var API_KEYS map[string]bool = map[string]bool{
	"apiKey1": true,
	"apiKey2": true,
	"apiKey3": true,
	"apiKey4": true,
	"apiKey5": true,
	"apiKey6": true,
	"apiKey7": true,
	"apiKey8": true,
	"apiKey9": true,
}
```
Here are some example strings for representing vaild api key values for the app. They can be added, removed or changed from here.
The app runs on localhost and listens to `:8081` port.
Head to the project directory on Terminal and run the app using:
```bash
$GOPATH/bin/PeakAPI
```
To send Event data with POST method, use the `/Events` endpoint,
to get Event data with GET method, use `/Events/{ApiKey}` endpoint.

Now that the app runs in the background, open another Terminal to send some http requests to the app:
```bash
curl -X POST -d '{"apikey":"apiKey1", "userID":80241, "timestamp":1798669390}' http://127.0.0.1:8081/Events -v
```
Here, we sent a simple Event data using the POST method of http in JSON format. Make sure to set the apiKey field to valid value which are all defined in `apis/apis.go`.

To get Event datas that have the same Api Key values, we can again use `curl` to send a GET request:
```bash
curl -X GET http://127.0.0.1:8081/Events/apiKey1 -v
```
Or we can simply send a request on a web browser:
`http://localhost:8081/Events/apiKey1`

The response data is also represented as JSON list.

After we have stored some Event data, we can visualize a Histogram using the response times of the requests.
To see the visualization, go `http://localhost:8081/Histogram`.

## Authors

* **Çağatay Barlas** - [cbarlas](https://github.com/cbarlas)

## Libraries Used

* [Gorilla Webtoolkit](https://github.com/gorilla/mux)
* [Go-chart](https://github.com/wcharczuk/go-chart)
* [Go-sql-driver](https://github.com/go-sql-driver/mysql)
