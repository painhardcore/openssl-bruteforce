package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

type results struct {
	password string
	fileName string
}


func argParse() (string, string, int) {
	wordlistPath := flag.String("wordlist", "/usr/share/wordlists/rockyou.txt", "Wordlist to use.")
	encFile := flag.String("file", "", "File to decrypt. (Required)")
	concurr:= flag.Int("concurrency", runtime.NumCPU(), "Specify number of concurrent openssl executions")
	flag.Parse()


	if *encFile == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	return *wordlistPath, *encFile, *concurr
}


func printResults(info results) {
	fmt.Println(strings.Repeat("-", 50))
	fmt.Printf("Found password [ %s ] !!\n", info.password)
	fmt.Println(strings.Repeat("-", 50))

	data, _ := ioutil.ReadFile(info.fileName)
	fmt.Println(string(data))
	fmt.Println(strings.Repeat("-", 50))
}

func crack(encFile string, wordlistPath chan string, wg *sync.WaitGroup, found chan<- results, stop <-chan bool) {
	defer wg.Done()
	cmdFormat := "openssl ec -in %s -out %s -passin pass:%s"
	fileName := "result"
	// loop line by line
	for {
		select {
		case <-stop:
			return
		case word:=<-wordlistPath:
			cmd := fmt.Sprintf(cmdFormat, encFile, fileName, word)
			command := strings.Split(cmd, " ")
			_, err := exec.Command(command[0], command[1:]...).CombinedOutput()
			// if no errors and file is ascii => found correct pass
			if err == nil && isASCIITextFile(fileName) {
				found <- results{word, fileName}
				return
			}
		}
	}
}

func watcher(wg *sync.WaitGroup, watch chan<- bool) {
	defer close(watch)
	wg.Wait()
}

func removeJunkExcept(goodFile string) {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		fileName := f.Name()
		if strings.HasPrefix(fileName, "result") && fileName != goodFile {
			err := os.Remove(fileName)
			if err != nil {
				panic(err)
			}
		}
	}
}

func isASCIITextFile(filePath string) bool {

	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	buffer := make([]byte, 256)
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		panic(err)
	}
	text := string(buffer)
	for char := 0; char < 256; char++ {
		if text[char] >= 128 {
			return false
		}
	}
	return true
}

func main() {
	wordlist, encryptedFile, concurr := argParse()
	println("Bruteforcing Started")
	var info results
	alreadyFound := false
	found := make(chan results)
	stop := make(chan bool)
	watch := make(chan bool)
	pwd:= make(chan string)
	var wg sync.WaitGroup
	// loop throught the ciphers and start a routine
	go func() {
		inFile, err := os.Open(wordlist)
		defer inFile.Close()
		if err != nil {
			panic(err)
		}
		scanner := bufio.NewScanner(inFile)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			pwd<-scanner.Text()
		}
	}()
	for i := 0; i <= concurr; i++ {
		wg.Add(1)
		go crack(encryptedFile, pwd, &wg, found, stop)
	}
	go watcher(&wg, watch)
Waiting:
	for {
		select {
 		case <-watch:
			break Waiting
		case info = <-found:
			if !alreadyFound {
				alreadyFound = true
				close(stop)
			}
		}
	}
	if alreadyFound {
		fmt.Printf("Success!! Results in file [ %s ]\n", info.fileName)
		printResults(info)
	} else {
		fmt.Println("Couldnt find password in that file")
	}
	removeJunkExcept(info.fileName)
}
