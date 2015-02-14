# wikiparser
Convert wikitext to HTML (using parsoid)

#### There are two components:
 - A frontend (wikipedia.go) which parses wikipedia dump and extracts articles in wikitext format. These articles in wikitext format are sent over to the backend (server.js). HTML documents received as response from this service are simply dumped to stdout (TODO: store in some KV-store).
 - A backend service (server.js) which uses Wikimedia's parsoid parser to convert wikitext to HTML format.


#### Installation and Usage
 - In project source root, do: `npm install`
    - This installs parsoid parser nodejs module
 - Run backend service: `node server`
 - Parse wikipedia dump: `go run wikipedia.go enwiki-latest-pages-articles.xml.bz2`
   - Where enwiki...bz2 is wikipedia dump as downloaded from `http://meta.wikimedia.org/wiki/Data_dump_torrents#enwiki`
   - Dump reader can directly read compressed dump
