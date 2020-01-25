package core

import (
	"bufio"
	//"encoding/json"
	"fmt"
	"log"
	//"reflect"
	//	"github.com/aws/aws-sdk-go/service/iam"
	//	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/manifoldco/promptui"
	"github.com/sethvargo/go-password/password"
	"os"
	//	"os/user"
	"strings"
)

type IamMenu struct {
	Name   string
	Number int
}

func Iam(AwsEnv string, AwsReg string) {

	iam_menu := []IamMenu{
		{Name: "ChangePassword", Number: 0},
		{Name: "CreateAccessKey", Number: 1},
		{Name: "CreateMFADevice", Number: 2},
		{Name: "DeactivateMFADevice", Number: 3},
		{Name: "DeleteAccessKey", Number: 4},
		{Name: "DeleteMFADevice", Number: 5},
		{Name: "EnableMFADevice", Number: 6},
		{Name: "ListAccessKeys", Number: 7},
		{Name: "ListMFADevices", Number: 8},
		{Name: "GetAccessKeyLastUsed", Number: 9},
		{Name: "UpdateAccessKey", Number: 10},
		{Name: "Main Menu", Number: 11},
		{Name: "Exit", Number: 12},
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F300 ({{.Number | red }}) {{ .Name | cyan }}",
		Inactive: " ({{ .Number | red }}) {{ .Name | cyan }}",
		Selected: "\U0001F300 {{ .Name | red | cyan }}",
		Details: `
--------- Option ----------
{{ "Name:" | faint }}   {{ .Name }}
{{ "Option:" | faint }} {{ .Number }}`,
	}

iam_Menu:
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: aws.String(AwsReg)},
		Profile: AwsEnv,
	}))

	svc := iam.New(sess)

	prompt := promptui.Select{
		Label:     "Select Option",
		Items:     iam_menu,
		Templates: templates,
		Size:      4,
	}

	i, _, err := prompt.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	switch iam_menu[i].Number {
	case 0:
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Enter Old password:")
		aws_old_pass, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		res, err := password.Generate(14, 4, 4, false, false)
		if err != nil {
			log.Fatal(err)
		}
		input := &iam.ChangePasswordInput{
			NewPassword: aws.String(strings.TrimSpace(res)),
			OldPassword: aws.String(strings.TrimSpace(aws_old_pass)),
		}
		result, err := svc.ChangePassword(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case iam.ErrCodeNoSuchEntityException:
					fmt.Println(iam.ErrCodeNoSuchEntityException, aerr.Error())
				case iam.ErrCodeInvalidUserTypeException:
					fmt.Println(iam.ErrCodeInvalidUserTypeException, aerr.Error())
				case iam.ErrCodeLimitExceededException:
					fmt.Println(iam.ErrCodeLimitExceededException, aerr.Error())
				case iam.ErrCodeEntityTemporarilyUnmodifiableException:
					fmt.Println(iam.ErrCodeEntityTemporarilyUnmodifiableException, aerr.Error())
				case iam.ErrCodePasswordPolicyViolationException:
					fmt.Println(iam.ErrCodePasswordPolicyViolationException, aerr.Error())
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
		fmt.Printf("New Password: %v\n", res)

		goto iam_Menu
	case 1:
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Which User would you like to create access key for: ")
		user_name, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}

		input := &iam.CreateAccessKeyInput{
			UserName: aws.String(strings.TrimSpace(user_name)),
		}

		result, err := svc.CreateAccessKey(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case iam.ErrCodeNoSuchEntityException:
					fmt.Println(iam.ErrCodeNoSuchEntityException, aerr.Error())
				case iam.ErrCodeLimitExceededException:
					fmt.Println(iam.ErrCodeLimitExceededException, aerr.Error())
				case iam.ErrCodeServiceFailureException:
					fmt.Println(iam.ErrCodeServiceFailureException, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
			return
		}

		fmt.Println(result)
		goto iam_Menu
	case 2:
		reader := bufio.NewReader(os.Stdin)
		//fmt.Println("Which User would you like to create MFA device for: ")
		fmt.Println("Enter name of new virtual device: ")
		user, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		input := &iam.CreateVirtualMFADeviceInput{
			VirtualMFADeviceName: aws.String(strings.TrimSpace(user)),
			//Path:                 aws.String("./QR_Code.png"),
		}
		result, err := svc.CreateVirtualMFADevice(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case iam.ErrCodeLimitExceededException:
					fmt.Println(iam.ErrCodeLimitExceededException, aerr.Error())
				case iam.ErrCodeEntityAlreadyExistsException:
					fmt.Println(iam.ErrCodeEntityAlreadyExistsException, aerr.Error())
				case iam.ErrCodeServiceFailureException:
					fmt.Println(iam.ErrCodeServiceFailureException, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
			return
		}
		fmt.Println("Enter QR code to google authenticator")

		//fmt.Printf("%+v\n", result.VirtualMFADevice.QRCodePNG)
		fmt.Println(result.VirtualMFADevice.QRCodePNG)
		n := len(result.VirtualMFADevice.QRCodePNG)
		blah := string(result.VirtualMFADevice.QRCodePNG[:n])
		fmt.Print("HERE IS BLAH")
		fmt.Print(blah)
		//res, _ := json.Marshal(result)
		//fmt.Println(string(res))

		//m := make(map[string]interface{})
		//json.Unmarshal(res, &m)
		//fmt.Println(m["VirtualMFADevice"])
		//fmt.Println(m)
		//fmt.Println(len(m))
		//		fmt.Println(m.VirtualMFADevice)
		//fmt.Println("QRCode: ", m["QRCodePNG"])
		//bytes := []byte(res)

		//var mfa []iam.CreateVirtualMFADeviceOutput
		//json.Unmarshal(bytes, &mfa)

		//for l := range mfa {
		//	fmt.Printf("%v", mfa[l])
		//}
		//fmt.Println(string(res.QRCodePNG))
		//		fmt.Println(string(res)) //works
		//fmt.Println(string(res))         //works
		//fmt.Printf("%d\n", res)          //prints raw bytes
		//fmt.Println(reflect.TypeOf(res)) //works

		//myString := string(res[:])
		//fmt.Println(myString)
		//fmt.Println(string(res{VirtualMFADevice})) //works
		//		for l := range res {
		//			fmt.Printf(" %v \n", string(res[l]))
		//		}

		//	fmt.Println("%v\n", res)
		//fmt.Printf("%v\n", res[0].QRCodePNG)
		//		fmt.Println(string(json.Marshal(result(1))))
		fmt.Println(result)
		goto iam_Menu

	case 3:
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("UserName of MFA Device to Deactivate: ")
		user, err := reader.ReadString('\n')
		fmt.Println("SerialNumber of MFA Device to Deactivate: ")
		serial, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		input := &iam.DeactivateMFADeviceInput{
			UserName:     aws.String(strings.TrimSpace(user)),
			SerialNumber: aws.String(strings.TrimSpace(serial)),
		}
		result, err := svc.DeactivateMFADevice(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case iam.ErrCodeNoSuchEntityException:
					fmt.Println(iam.ErrCodeNoSuchEntityException, aerr.Error())
				case iam.ErrCodeDeleteConflictException:
					fmt.Println(iam.ErrCodeDeleteConflictException, aerr.Error())
				case iam.ErrCodeLimitExceededException:
					fmt.Println(iam.ErrCodeLimitExceededException, aerr.Error())
				case iam.ErrCodeServiceFailureException:
					fmt.Println(iam.ErrCodeServiceFailureException, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
			return
		}
		fmt.Println("Deactivating device: %v for User: %v", serial, user)
		fmt.Println(result)
		goto iam_Menu

	case 4:
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Which Access Key would you like to delete: ")
		access_key, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		input := &iam.DeleteAccessKeyInput{
			AccessKeyId: aws.String(strings.TrimSpace(access_key)),
		}

		result, err := svc.DeleteAccessKey(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case iam.ErrCodeNoSuchEntityException:
					fmt.Println(iam.ErrCodeNoSuchEntityException, aerr.Error())
				case iam.ErrCodeLimitExceededException:
					fmt.Println(iam.ErrCodeLimitExceededException, aerr.Error())
				case iam.ErrCodeServiceFailureException:
					fmt.Println(iam.ErrCodeServiceFailureException, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
			goto iam_Menu
		}

		fmt.Println(result)
		goto iam_Menu
	case 5:
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Serial Number of MFA Device to Delete: ")
		serial, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		input := &iam.DeleteVirtualMFADeviceInput{
			SerialNumber: aws.String(strings.TrimSpace(serial)),
		}
		result, err := svc.DeleteVirtualMFADevice(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case iam.ErrCodeNoSuchEntityException:
					fmt.Println(iam.ErrCodeNoSuchEntityException, aerr.Error())
				case iam.ErrCodeDeleteConflictException:
					fmt.Println(iam.ErrCodeDeleteConflictException, aerr.Error())
				case iam.ErrCodeLimitExceededException:
					fmt.Println(iam.ErrCodeLimitExceededException, aerr.Error())
				case iam.ErrCodeServiceFailureException:
					fmt.Println(iam.ErrCodeServiceFailureException, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
			return
		}
		fmt.Println("Deleting %v device", serial)
		fmt.Println(result)
		goto iam_Menu

	case 6:
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("UserName to associate MFA Device with: ")
		user, err := reader.ReadString('\n')
		fmt.Println("MFA Device SerialNumber: ")
		serial, err := reader.ReadString('\n')
		fmt.Println("Authentication Code 1: ")
		code1, err := reader.ReadString('\n')
		fmt.Println("Authentication Code 2: ")
		code2, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		input := &iam.EnableMFADeviceInput{
			UserName:            aws.String(strings.TrimSpace(user)),
			SerialNumber:        aws.String(strings.TrimSpace(serial)),
			AuthenticationCode1: aws.String(strings.TrimSpace(code1)),
			AuthenticationCode2: aws.String(strings.TrimSpace(code2)),
		}
		result, err := svc.EnableMFADevice(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case iam.ErrCodeNoSuchEntityException:
					fmt.Println(iam.ErrCodeNoSuchEntityException, aerr.Error())
				case iam.ErrCodeDeleteConflictException:
					fmt.Println(iam.ErrCodeDeleteConflictException, aerr.Error())
				case iam.ErrCodeLimitExceededException:
					fmt.Println(iam.ErrCodeLimitExceededException, aerr.Error())
				case iam.ErrCodeServiceFailureException:
					fmt.Println(iam.ErrCodeServiceFailureException, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
			return
		}
		fmt.Println(result)
		goto iam_Menu
	case 7:
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("UserName of Access Keys: ")
		user_name, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		input := &iam.ListAccessKeysInput{
			UserName: aws.String(strings.TrimSpace(user_name)),
		}
		result, err := svc.ListAccessKeys(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case iam.ErrCodeNoSuchEntityException:
					fmt.Println(iam.ErrCodeNoSuchEntityException, aerr.Error())
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
		goto iam_Menu
	case 8:
		input := &iam.ListVirtualMFADevicesInput{}

		result, err := svc.ListVirtualMFADevices(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				fmt.Println(err.Error())
			}
			return
		}
		fmt.Println(result)
		goto iam_Menu
	case 9:
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("What key ID would you like to look up?")
		aws_key_id, err := reader.ReadString('\n')
		policy, err := svc.GetAccessKeyLastUsed(&iam.GetAccessKeyLastUsedInput{
			AccessKeyId: aws.String(strings.TrimSpace(aws_key_id)),
		})
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		fmt.Printf("Key Last Used: %v\n", policy)
		goto iam_Menu

	case 10:
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Which Access Key would you like to update: ")
		access_key, err := reader.ReadString('\n')
		fmt.Println("Activate or Deactivate this key (A/D): ")
		act_deact, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		if strings.TrimSpace(act_deact) == "A" {
			fmt.Println("Activating access key.....")
			input := &iam.UpdateAccessKeyInput{
				AccessKeyId: aws.String(strings.TrimSpace(access_key)),
				Status:      aws.String("Active"),
			}
			result, err := svc.UpdateAccessKey(input)
			if err != nil {
				if aerr, ok := err.(awserr.Error); ok {
					switch aerr.Code() {
					case iam.ErrCodeNoSuchEntityException:
						fmt.Println(iam.ErrCodeNoSuchEntityException, aerr.Error())
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
				fmt.Println(result)
				goto iam_Menu
			}
		} else if strings.TrimSpace(act_deact) == "D" {
			fmt.Println("Deactivating access key.....")
			input := &iam.UpdateAccessKeyInput{
				AccessKeyId: aws.String(strings.TrimSpace(access_key)),
				Status:      aws.String("Inactive"),
			}
			result, err := svc.UpdateAccessKey(input)
			if err != nil {
				if aerr, ok := err.(awserr.Error); ok {
					switch aerr.Code() {
					case iam.ErrCodeNoSuchEntityException:
						fmt.Println(iam.ErrCodeNoSuchEntityException, aerr.Error())
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
				fmt.Println(result)
				goto iam_Menu
			}
		} else {
			fmt.Println("ERROR: Invalid option\n")
			goto iam_Menu
		}
		goto iam_Menu
	case 11:
		return
	case 12:
		os.Exit(0)
	}

}
