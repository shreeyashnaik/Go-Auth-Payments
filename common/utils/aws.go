package utils

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/spf13/viper"
)

var (
	awsSession *session.Session
	sesSession *ses.SES
)

func ConnectAWS() {
	var err error = nil
	awsSession, err = session.NewSession(
		&aws.Config{
			Region: aws.String(viper.GetString("AWS_REGION")),
			Credentials: credentials.NewStaticCredentials(
				viper.GetString("AWS_ACCESS_KEY_ID"),
				viper.GetString("AWS_SECRET_KEY"),
				"",
			),
		},
	)

	if err != nil {
		log.Fatalln("Cannot connect to AWS: ", err)
	}

	log.Println("Successfully connected to AWS: ", awsSession.Config.Endpoint)

	// Initialize SES session
	sesSession = ses.New(awsSession)
}

func SendOTPEmail(recipient string, otp int) {
	subject := "Go Auth OTP Verification"
	htmlBody := fmt.Sprintf("<b>%v</b> is your OTP. Your OTP expires in <em>5 minutes</em>", otp)

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{&recipient},
		},
		Message: &ses.Message{
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(subject),
			},
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(htmlBody),
				},
			},
		},
		Source: aws.String(viper.GetString("AWS_SES_SENDER")),
	}

	// Attempt to send the email.
	result, err := sesSession.SendEmail(input)

	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
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
}
