# openssl-bruteforce
Fixed to bruteforce with Elliptic Curve (EC)
## Usage
```
./openssl-brute -file <encrypted file>
```
Needs to be in base64 (will be fixed later).

### All options
```
$ ./openssl-brute 
  -concurrency int
    	Specify number of concurrent openssl executions. (default "number of cpu")
  -file string
    	File to decrypt. (Required)
  -print
    	Set to print all available ciphers and exit.
  -wordlist string
    	Wordlist to use. (default "/usr/share/wordlists/rockyou.txt")
```

## Build
```
cd openssl-bruteforce/
go build -o openssl-brute
```
