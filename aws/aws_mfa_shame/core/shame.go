package core

import (
	"encoding/base64"
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"log"
	"net/http"
	"strings"
	"text/template"
	"path/filepath"
	"sync"
)

func Web_serv() {

	type templateHandler struct {
		once sync.Once
		filename string
		temp1 	*template.Template
	}

	func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
		t.once.Do(func() {
			t.temp1 = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
		})
		t.temp1.Execute(w, nil)
	}

	//port := flag.String("p", "8100", "port to serve on")
	//directory := flag.String("d", "./static", "the directory of static file to how")
	//flag.Parse()

	//http.Handle("/", http.FileServer(http.Dir(*directory)))
	http.Handle("/", &templateHandler{filename: "shame.html"})

	if err := http.ListenAndServe(":8100", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

//	log.Printf("Serving %v on HTTP port: %v\n", *directory, *port)
//	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func Shame(AwsEnv string, AwsReg string) {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: aws.String(AwsReg)},
		Profile: AwsEnv,
	}))
	svc := iam.New(sess)
	input := &iam.GenerateCredentialReportInput{}

	result, err := svc.GenerateCredentialReport(input)
	if err != nil {
		fmt.Println(err.Error())

	}
	fmt.Println(result)

	get_input := &iam.GetCredentialReportInput{}

	get_result, get_err := svc.GetCredentialReport(get_input)

	if get_err != nil {
		fmt.Println(get_err.Error())
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(get_result.Content))
	decoded, err := base64.StdEncoding.DecodeString(encoded)

	s := strings.Split(string(decoded), "\n")

	for _, each := range s {
		account := csv.NewReader(strings.NewReader(each))
		record, err := account.Read()

		if err != nil {
			log.Fatal(err)
		}

		if record[3] == "true" && record[7] == "false" {
			fmt.Println("MFA NOT ACTIVE FOR: ")
			fmt.Println(record[0])
		} else if record[3] == "true" && record[7] == "true" {
			fmt.Println("MFA IS ACTIVE FOR: ")
			fmt.Println(record[0])
		} else {
			continue
		}

	}
	Web_serv()
	return

}
