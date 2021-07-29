# HTTP Parallel Getter
A tool that receives a list of addresses as a space separated strings in the command line, and do a get request for all of them in parallel, and prints for each address the md5 hash of the response body in case there is no HTTP error, or print the error otherwise.


## Usage
./http [-parallel NUMBER] `<addresse1 address2>`

## Options
-parallel NUMBER: Optional, defines the parallelization limit for the tool, Default 10

## Examples
./http google.com http://google.com

./http http://www.office.com http://google.com http://www.office.com 

./http -parallel 3 office.com google.com facebook.com yahoo.com yandex.com twitter.com
