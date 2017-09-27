package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"io"
	"os"
	"strings"
	"sync"
)

type CLI struct {
	outStream, errStream io.Writer
}

const (
	output_header string = "Name\tInstanceID\tInstanceType\tPrivateIP\tPublicIP\tPublicDNSName\tLaunchTime"
	app_name             = "ec2-search"
)

const (
	ExitCodeOK = iota
	ExitCodeErr
)

func (c *CLI) outputError(msg string) {
	fmt.Fprintln(c.errStream, app_name+" : "+msg)
}

func generateInstanceFilterByName(name string) *ec2.DescribeInstancesInput {
	return &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: []*string{aws.String(name)},
			},
		},
	}
}

func generateInstanceFilterByPrivateIP(pip string) *ec2.DescribeInstancesInput {
	return &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("private-ip-address"),
				Values: []*string{aws.String(pip)},
			},
		},
	}
}

func generateInstanceFilterByPublicIP(eip string) *ec2.DescribeInstancesInput {
	return &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("ip-address"),
				Values: []*string{aws.String(eip)},
			},
		},
	}
}

func generateInstanceFilterByInstanceId(id string) *ec2.DescribeInstancesInput {
	return &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("instance-id"),
				Values: []*string{aws.String(id)},
			},
		},
	}
}

func (c *CLI) getInstances(filter ec2.DescribeInstancesInput, region string, profile string) (ret []string) {
	var conf *aws.Config
	if profile != "" {
		cred := credentials.NewSharedCredentials("", profile)
		conf = &aws.Config{
			Credentials: cred,
			Region:      &region,
		}
	} else {
		conf = &aws.Config{
			Region: &region,
		}
	}
	sess, err := session.NewSession(conf)
	if err != nil {
		c.outputError("Invalid credential.[" + profile + "]")
		return nil
	}

	svc := ec2.New(sess)

	res, err := svc.DescribeInstances(&filter)
	if err != nil {
		c.outputError("Unauthorized EC2 DescribeInstances.[" + profile + "]")
		return nil
	}

	for _, r := range res.Reservations {
		for _, instance := range r.Instances {
			ret = append(ret, format(instance))
		}
	}
	return
}

func format(ins *ec2.Instance) string {
	var name string = ""
	for _, t := range ins.Tags {
		if *t.Key == "Name" {
			name = *t.Value
			break
		}
	}

	var eip string
	var public_dns string

	if ins.PublicIpAddress != nil {
		eip = *ins.PublicIpAddress
	}
	if ins.PublicDnsName != nil {
		public_dns = *ins.PublicDnsName
	}

	res := []string{
		name,
		*ins.InstanceId,
		*ins.InstanceType,
		*ins.PrivateIpAddress,
		eip,
		public_dns,
		(*ins.LaunchTime).Format("2006-01-02 15:04:05"),
	}
	return strings.Join(res[:], "\t")
}

func (c *CLI) Run(args []string) int {
	var name string
	var pip string
	var eip string
	var id string
	var region string
	var config string
	var version bool

	flags := flag.NewFlagSet("ec2-search", flag.ContinueOnError)
	flags.SetOutput(c.errStream)
	flags.StringVar(&name, "name", "", "Specify instance name.")
	flags.StringVar(&pip, "pip", "", "Specify instance private ip.")
	flags.StringVar(&eip, "eip", "", "Specify instance public ip.")
	flags.StringVar(&id, "id", "", "Specify instance id.")
	flags.StringVar(&region, "region", "ap-northeast-1", "Specify region.")
	flags.StringVar(&config, "config", "~/.aws/credentials", "Specify aws credential file path.")
	flags.BoolVar(&version, "v", false, "Output version number.")

	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeErr
	}

	if version {
		fmt.Fprintln(c.outStream, Version)
		return 0
	}

	var filter ec2.DescribeInstancesInput

	if name != "" {
		filter = *generateInstanceFilterByName(name)
	} else if pip != "" {
		filter = *generateInstanceFilterByPrivateIP(pip)
	} else if eip != "" {
		filter = *generateInstanceFilterByPublicIP(eip)
	} else if id != "" {
		filter = *generateInstanceFilterByInstanceId(id)
	} else {
		flags.PrintDefaults()
		return ExitCodeErr
	}
	profiles, err := load_file(config)
	if err != nil {
		c.outputError("Not found specify config file.")
		return ExitCodeErr
	}

	output := make([][]string, len(profiles))
	hit := false

	var wg sync.WaitGroup
	for i, profile := range profiles {
		wg.Add(1)
		go func(index int, profile string) {
			ret := c.getInstances(filter, region, profile)
			if ret != nil {
				hit = true
			}
			output[index] = ret
			defer wg.Done()
		}(i, profile)
	}
	wg.Wait()

	if hit {
		fmt.Fprintln(c.outStream, output_header)
		for i := 0; i < len(output); i++ {
			if len(output[i]) == 0 {
				continue
			}
			for j := 0; j < len(output[i]); j++ {
				fmt.Fprintln(c.outStream, output[i][j])
			}
		}
	} else {
		c.outputError("No hits.")
		return ExitCodeErr
	}

	return ExitCodeOK
}

func main() {
	cli := &CLI{outStream: os.Stdout, errStream: os.Stderr}
	os.Exit(cli.Run(os.Args))
}
