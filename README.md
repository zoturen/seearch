# Search Engine thats crawls the web

This is very incomplete and uses a hardcoded start point for the crawl. It also has a hardcoded crawldepth. 

It uses TF-IDF to rank the documents.

The program keeps everything in-memory. You should probably not crawl the entire web.

## How to run it

1. Change the url and the depth of your choice inside main.go

2. Run the program, eg. `go run main.go`

3. Open up the index.html in your browser 

4. Search
