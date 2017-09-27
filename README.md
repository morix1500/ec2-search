# ec2-search
ec2-search can search for EC2 instance, and it can cross accounts.

## Example

AWS Account: hoge  

|tagName|instance id|private ip|public ip|
|---|---|---|---|
|hoge-instance1|i-hogexxx|10.0.0.0|11.22.33.44|
|hoge-instance2|i-hogeyyy|10.0.0.1|22.33.44.55|

AWS Account: fuga

|tagName|instance id|private ip|public ip|
|---|---|---|---|
|fuga-instance1|i-fugaxxx|10.0.0.2|33.44.55.66|
|fuga-instance2|i-fugayyy|10.0.0.3|44.55.66.77|

AWS Account: puyo

|tagName|instance id|private ip|public ip|
|---|---|---|---|
|puyo-instance1|i-puyoxxx|10.0.0.4|55.66.77.88|
|puyo-instance2|i-puyoyyy|10.0.0.5|66.77.88.99|

```
$ cat ~/.aws/credentials
[hoge]
aws_access_key_id = XXXXXXXXXXX
aws_secret_access_key = XXXXXXXXXXX

[fuga]
aws_access_key_id = XXXXXXXXXXX
aws_secret_access_key = XXXXXXXXXXX

[piyo]
aws_access_key_id = XXXXXXXXXXX
aws_secret_access_key = XXXXXXXXXXX

$ ec2-search --name fuga-instance1
Name    InstanceID      InstanceType    PrivateIP       PublicIP        PublicDNSName   LaunchTime
fuga-instance1   i-fugaxxx      t2.micro        10.0.0.2     33.44.55.66   ec2-33-44-55-66.ap-northeast-1.compute.amazonaws.com  2017-08-18 07:04:04

$ ec2-search --name fuga-instance1 | ruler -t tsv
+----------------+------------+--------------+-------------+---------------+--------------------------------------------------------+---------------------+
| Name           | InstanceID | InstanceType | PrivateIP   | PublicIP      | PublicDNSName                                          | LaunchTime          |
+----------------+------------+--------------+-------------+---------------+--------------------------------------------------------+---------------------+
| fuga-instance1 | i-fugaxxx  | t2.micro     | 10.0.0.2    | 33.44.55.66   | ec2-33-44-55-66.ap-northeast-1.compute.amazonaws.com   | 2017-08-18 07:04:04 |
+----------------+------------+--------------+-------------+---------------+--------------------------------------------------------+---------------------+
```

## Usage
```
$ ec2-search --help
Usage of ec2-search:
  -config string
        Specify aws credential file path. (default "~/.aws/credentials")
  -eip string
        Specify instance public ip.
  -id string
        Specify instance id.
  -name string
        Specify instance name.
  -pip string
        Specify instance private ip.
  -region string
        Specify region. (default "ap-northeast-1")
  -v    Output version number.
```

### option: config
This option specify aws credentials file path.  
if not exists "~/.saws/credentials", follow AWS authentication rules.  
<http://docs.aws.amazon.com/ja_jp/sdk-for-go/v1/developer-guide/configuring-sdk.html>

### option: name
This option spcify EC2 Name Tag.  
We can use filter query.  
<http://docs.aws.amazon.com/cli/latest/reference/ec2/describe-instances.html>  
```
$ ec2-search --name fuga-* | ruler -t tsv
+----------------+------------+--------------+-----------+-------------+------------------------------------------------------+---------------------+
| Name           | InstanceID | InstanceType | PrivateIP | PublicIP    | PublicDNSName                                        | LaunchTime          |
+----------------+------------+--------------+-----------+-------------+------------------------------------------------------+---------------------+
| fuga-instance1 | i-fugaxxx  | t2.small     | 10.0.0.2  | 33.44.55.66 | ec2-33-44-55-66.ap-northeast-1.compute.amazonaws.com | 2017-08-18 07:04:04 |
| fuga-instance2 | i-fugaxxx  | t2.small     | 10.0.0.3  | 44.55.66.77 | ec2-44-55-66-77.ap-northeast-1.compute.amazonaws.com | 2017-08-18 08:01:01 |
+----------------+------------+--------------+-----------+-------------+------------------------------------------------------+---------------------+
```

## Installation
```
$ go get -u github.com/morix1500/ec2-search 
```

or

```
$ wget https://github.com/morix1500/ec2-search/releases/download/v0.1.0/ec2-search_linux_amd64 -O /usr/local/bin/ec2-search
$ chmod u+x /usr/local/bin/ec2-search
```

## License
Please see the [LICENSE](./LICENSE) file for details.  

## Author
Shota Omori(Morix)  
https://github.com/morix1500
