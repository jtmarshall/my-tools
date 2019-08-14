# Monitor

Monitors and crawls all sites located in the 'status' table of the MonitorDB.
Built with Go and React.

## Getting Started / Installing

Clone repo and run:
```
go run application.go
```

### Prerequisites

Things you need to have installed
```
Go(lang)
```

## Deployment

After aws cli configuration, use:
```
eb deploy -timeout=60
```

Alternatively:

Zip all files in parent directory -> "Archive.zip" -> upload to Elastic Beanstalk

## Built With

* [React](https://reactjs.org/) - Library used for UI
* [Go](https://golang.org/) - Server-side

## Authors

* **Jordan Marshall**

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
