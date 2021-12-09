# findresources

## Usage

Write filters you want to find to config file.
`.findres.yml`

```yaml
tagfilters:
  - key: Project
    val: myproject
  - key: Stage
    val: staging
resourcetypefilters:
  - rds:db
    #  - ec2:instance
  - elasticloadbalancing:loadbalancer
```


```console
$ findresource

arn:aws:rds:ap-northeast-1:012345678912:db:mydb1
arn:aws:rds:ap-northeast-1:012345678912:db:mydb1-readreplica
arn:aws:elasticloadbalancing:ap-northeast-1:012345678912:loadbalancer/app/my-lb0/abcdefg012345678
arn:aws:elasticloadbalancing:ap-northeast-1:012345678912:loadbalancer/app/my-lb1/abcdffg012345678
arn:aws:elasticloadbalancing:ap-northeast-1:012345678912:loadbalancer/app/my-lb2/abcdffg012345678
```
