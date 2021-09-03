package azure

import (
	"bufio"
	"bytes"
	"errors"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/adeturner/azureBilling/cloudauth"
	"github.com/adeturner/azureBilling/observability"
	"github.com/adeturner/azureBilling/utils"
)

// useful
// https://github.com/Azure-Samples/azure-cli-samples

type AzureCli struct {
	subscriptionId string
	spn            *cloudauth.AzureServicePrincipalType
}

func NewAzureCli() *AzureCli {
	cli := AzureCli{}
	cli.spn = cloudauth.NewAzureServicePrincipalType()
	err := cli.spn.LoadFromFile()

	if err != nil {
		observability.Error(err.Error())
	} else {
		observability.Info("Initialised azcli")
	}
	return &cli
}

func (cli *AzureCli) SetSubscription(subscriptionId string) {
	cli.subscriptionId = subscriptionId
}

func (cli *AzureCli) parseCmd(cmdString string) []string {

	//observability.Info(fmt.Sprintf("Parsing command: %s", cmdString))

	var haveOpenQuote, haveCloseQuote bool

	commandArray := strings.Fields(cmdString)
	var newArray []string
	j := 0

	for i, v := range commandArray {
		//observability.Info(fmt.Sprintf("Parsed output %v", commandArray))
		//observability.Info(fmt.Sprintf("newArray output %v", newArray))

		switch v {
		case "{SUBSCRIPTIONID}":
			commandArray[i] = cli.subscriptionId
		case "{TENANTID}":
			commandArray[i] = cli.spn.TenantId
		case "{CLIENTID}":
			commandArray[i] = cli.spn.ClientId
		case "{CLIENTSECRET}":
			commandArray[i] = "'" + cli.spn.ClientSecret + "'"
			commandArray[i] = cli.spn.ClientSecret
		default:
			if commandArray[i][0:1] == "\"" {
				haveOpenQuote = true
			}
			if commandArray[i][len(v)-1:] == "\"" {
				haveCloseQuote = true
			}

			if commandArray[i][0:1] == "{" {
				s := v[1 : len(v)-1]
				observability.Info("Get environment variable: " + s)
				commandArray[i] = os.Getenv(s)
			}
		}

		if !haveOpenQuote {
			newArray = append(newArray, commandArray[i])
			j = j + 1

		} else if haveOpenQuote && !haveCloseQuote {

			if len(newArray)-1 < j {
				newArray = append(newArray, commandArray[i][1:])
			} else {
				newArray[j] = newArray[j] + " " + commandArray[i][0:]
			}

		} else if haveOpenQuote && haveCloseQuote {
			newArray[j] = newArray[j] + " " + commandArray[i][:len(commandArray[i])-1]
			j = j + 1
			haveOpenQuote = false
			haveCloseQuote = false
		}

	}
	return newArray
}

func (cli *AzureCli) Login() (err error) {

	defer utils.Quiet()()
	_, err = cli.ExecCmd("az login --service-principal -u {CLIENTID} -p {CLIENTSECRET} --tenant {TENANTID} --allow-no-subscriptions")
	return err
}

func (cli *AzureCli) ExecCmd(cmd string) (output string, err error) {
	var out bytes.Buffer
	out, err = cli.doCmd(cmd)
	return out.String(), err
}

func (cli *AzureCli) doCmd(cmdString string) (out bytes.Buffer, err error) {

	parsedCmd := cli.parseCmd(cmdString)
	//observability.Info(fmt.Sprintf("parsedCmd output %v", parsedCmd))

	var cmd *exec.Cmd
	var stdout, stderr bytes.Buffer

	if parsedCmd[0] == "export" {
		s := strings.Split(parsedCmd[1], "=")
		err = os.Setenv(s[0], s[1])
		if err == nil {
			out.WriteString("os.setenv: " + s[0] + "=" + s[1])
		} else {
			out.WriteString(err.Error())
		}

	} else if parsedCmd[0] == "az" {

		cmd = exec.Command("az", parsedCmd[1:]...)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err = cmd.Run()

	} else {

		observability.Info("ignoring line beginning " + parsedCmd[0])

	}

	if err != nil {
		// observability.Error(err.Error() + "\n" + stderr.String())
		out = stderr
	} else {
		// observability.Info(fmt.Sprintf("Executed command '%s' with output: \n%s", parsedCmd[0], stdout.String()))
		out = stdout
	}

	return out, err
}

func (cli *AzureCli) ExecFileLineByLine(filename string) (err error) {

	var f *os.File
	var text string

	observability.Info("Opening file: " + filename)

	if err == nil {
		f, err = os.Open(filename)
		if err != nil {
			err = errors.New("Unable to access file %s" + filename)
		}
		defer f.Close()
	}

	if err == nil {
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			text = scanner.Text()
			if text != "" {
				observability.Info("processing line: " + text)
				cli.doCmd(text)
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

	}
	if err == nil {
		observability.Info("Completed file: " + filename)
	}
	return err
}
