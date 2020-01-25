package core

import (
	"encoding/base64"
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"log"
	"net/http"
	"strings"
)

//set IAM fields
var User = 0
var Arn = 1
var User_Creation_Time = 2
var Password_Enabled = 3
var Password_Last_Used = 4
var Password_Last_Changed = 5
var Password_Next_Rotation = 6
var Mfa_Active = 7
var Access_Key_1_Active = 8
var Access_Key_1_Last_Rotated = 9
var Access_Key_1_Last_Used_Date = 10
var Access_Key_1_Last_Used_Region = 11
var Access_Key_1_Last_Used_Service = 12
var Access_Key_2_Active = 13
var Access_Key_2_Last_Rotated = 14
var Access_Key_2_Last_Used_Date = 15
var Access_Key_2_Last_Used_Region = 16
var Access_Key_2_Last_Used_Service = 17
var Cert_1_Active = 18
var Cert_1_Last_Rotated = 19
var Cert_2_Active = 20
var Cert_2_Last_Rotated = 21

func Web_serv() {
	port := flag.String("p", "8100", "port to serve on")
	directory := flag.String("d", "./static", "the directory of static file to host")
	flag.Parse()

	http.Handle("/", http.FileServer(http.Dir(*directory)))

	log.Printf("Serving %s on HTTP port: %s\n", *directory, *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func CIS_Benchmark(AwsEnv string, AwsReg string) {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: aws.String(AwsReg)},
		Profile: AwsEnv,
	}))

	svc := iam.New(sess)

	input := &iam.GenerateCredentialReportInput{}

	result, err := svc.GenerateCredentialReport(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeLimitExceededException:
				fmt.Println(iam.ErrCodeLimitExceededException, aerr.Error())
			case iam.ErrCodeServiceFailureException:
				fmt.Println(iam.ErrCodeServiceFailureException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
	}

	fmt.Println(result)

	get_input := &iam.GetCredentialReportInput{}
	//	get_output := &iam.GetCredentialReportOutput{
	//		ReportFormat: aws.String("text"),
	//	}

	get_result, get_err := svc.GetCredentialReport(get_input)

	if get_err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeCredentialReportNotPresentException:
				fmt.Println(iam.ErrCodeCredentialReportNotPresentException, aerr.Error())
			case iam.ErrCodeCredentialReportExpiredException:
				fmt.Println(iam.ErrCodeCredentialReportExpiredException, aerr.Error())
			case iam.ErrCodeCredentialReportNotReadyException:
				fmt.Println(iam.ErrCodeCredentialReportNotReadyException, aerr.Error())
			case iam.ErrCodeServiceFailureException:
				fmt.Println(iam.ErrCodeServiceFailureException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(get_err.Error())
		}
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(get_result.Content))
	//fmt.Println("Encoded")
	//fmt.Println(encoded)

	decoded, err := base64.StdEncoding.DecodeString(encoded)

	//fmt.Println("Decoded")
	if err != nil {
		fmt.Println(err)
	}

	s := strings.Split(string(decoded), "\n")

	for _, each := range s {
		//1.1 Avoid the use of the "root" account
		if strings.Contains(each, "<root_account>") {
			root_account := csv.NewReader(strings.NewReader(each))
			record, err := root_account.Read()
			if err != nil {
				log.Fatal(err)
			}
			if record[Access_Key_1_Last_Used_Date] != "N/A" && record[Access_Key_2_Last_Used_Date] != "N/A" {
				fmt.Println("Disable root keys")
			} else {
				fmt.Println("Root good")
			}
		}
		//fmt.Println(index, each)
	}
	//fmt.Println(string(decoded))
	//fmt.Println(get_result.Content)

	//	req, resp := svc.GetCredentialReportRequest(get_input)

	//	send_err := req.Send()

	//	if send_err == nil {
	//		fmt.Println("print resp")
	//		fmt.Println(resp)
	//	}

	Web_serv()

	return

	//os.Exit(0)
}
