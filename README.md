

To set up the blacklist monitor, follow these steps:

Clone the repository from https://github.com/mikelynn2/blacklistmonitor.git.

Compile the main.go file using the command:
```
go build -o blmon-exporter main.go
```

To use the exporter, run the command:
```
blmon-exporter -api=<API key> [-port=<port number | default :8080>] [-link=<API link | default http://localhost/blacklistmonitor/api.php >]
```

In step 1, you will need to clone the repository to your local machine. This can be done by using Git or by downloading the repository as a zip file and extracting it.

In step 2, you will compile the main.go file using the go build command. This will create an executable file named blmon-exporter.

Finally, in step 3, you can run the exporter by providing your API key with the -api flag. Optionally, you can also specify a port number and API link using the -port and -link flags, respectively. By default, the exporter will use port 8080 and the API link http://localhost/blacklistmonitor/api.php.
