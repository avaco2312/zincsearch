package main

import (
	"bytes"
	"fmt"
	"io"
	"net/mail"
	"os"
	"prueba/zinc"
	"strings"
	"sync"
	"time"
)

const (
	emailsDir         = "F:/enron_mail_20110402"
	cConcEmailProcess = 10
	cConcZincProcess  = 10
	cConcEmails       = 500
	cConcZincData     = 500
	batchEmail        = 50
	batchEmailMaxSize = 600000
)

var concEmail = make(chan struct{}, cConcEmailProcess)
var concZinc = make(chan struct{}, cConcZincProcess)
var emails = make(chan string, cConcEmails)
var zincData = make(chan string, cConcZincData)
var wge, wgd sync.WaitGroup
var inserted, rejected, formaterr, emailtobig int
var mus sync.Mutex

func main() {
	_ = zinc.DeleteIndex()
	err := zinc.CreateIndex()
	if err != nil {
		panic(err)
	}
	ti := time.Now()
	wgd.Add(1)
	// Read available emails, form a package of email and index it with ZincSearch
	go func() {
		defer wgd.Done()
		b := strings.Builder{}
		i := 0
		for d := range zincData {
			if b.Len()+len(d) > batchEmailMaxSize || i == batchEmail {
				wgd.Add(1)
				go zincCreate(b.String(), i)
				i = 0
				b.Reset()
			}
			b.WriteString(d)
			b.WriteByte(10)
			i++
		}
		if b.Len() != 0 {
			wgd.Add(1)
			go zincCreate(b.String(), i)
		}
	}()
	// Read available email files, read file and unpackage each email, send it to further process
	go func() {
		for e := range emails {
			go processEmail(e)
		}
	}()
	// Traverse emails directories, send email file names to further process
	findDirs(emailsDir)
	close(emails)
	wge.Wait() // waits for or email file names found and processed
	close(zincData)
	wgd.Wait() // waits for or emails packaged and processed
	fmt.Println("Duracion: ", time.Since(ti))
	fmt.Println("Inserted: ", inserted)
	fmt.Println("Error formato: ", formaterr)
	fmt.Println("Rechazados batch: ", rejected)
	fmt.Println("Email to big: ", emailtobig)
}

// Traverse directories recursively. Send email file names for further process
func findDirs(iniDir string) {
	dir, err := os.Open(iniDir)
	if err != nil {
		panic(err)
	}
	files, err := dir.ReadDir(-1)
	if err != nil {
		panic(err)
	}
	dir.Close()
	for _, f := range files {
		if f.IsDir() {
			findDirs(iniDir + "/" + f.Name())
		} else {
			fi, err := f.Info()
			if err != nil {
				panic(err)
			}
			if fi.Size() > batchEmailMaxSize {
				emailtobig++
				continue
			}
			wge.Add(1)
			emails <- iniDir + "/" + f.Name()
		}
	}
}

// Read one email file. Unpack email and send for further process
func processEmail(email string) {
	defer wge.Done()
	concEmail <- struct{}{}
	fileContent, err := os.ReadFile(email)
	if err != nil {
		panic(err.Error())
	}
	r := bytes.NewReader(fileContent)
	m, err := mail.ReadMessage(r)
	if err != nil {
		mus.Lock()
		formaterr++
		mus.Unlock()
		<-concEmail
		return
	}
	body, err := io.ReadAll(m.Body)
	if err != nil {
		mus.Lock()
		formaterr++
		mus.Unlock()
		<-concEmail
		return
	}
	zincData <- fmt.Sprintf(`{"_id": "%s", "from": "%s", "to": "%s", "subject": "%s", "content": "%s"}`,
		email, clean(m.Header.Get("From")), clean(m.Header.Get("To")), clean(m.Header.Get("Subject")), clean(string(body)))
	<-concEmail
}

// Clean email special chars
func clean(s string) string {
	s1 := strings.ReplaceAll(s, `\`, `\\`)
	s2 := strings.ReplaceAll(s1, string(rune(10)), `\n`)
	return strings.ReplaceAll(s2, `"`, `\"`)
}

// Index one emails package
func zincCreate(s string, can int) {
	defer wgd.Done()
	concZinc <- struct{}{}
	count, err := zinc.CreateData(s)
	if err != nil {
		panic(err)
	}
	mus.Lock()
	inserted += count
	rejected += (can - count)
	mus.Unlock()
	<-concZinc
}
