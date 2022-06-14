# MySCP

A simple SCP client, that can help me out for copy files between machines, due to some policy restrictions(Kejserens nye Klæder凸).

## Install

```shell
go install github.com/elvuel/myscp@latest
```

## Usage

```shell
# remote to local
myscp -h host:2222 -u user -k ~/.ssh/id_rsa -l /home/user/some.txt -r /home/user/some.txt -ori rtl

# local to remote
myscp -h host:2222 -u user -k ~/.ssh/id_rsa -l /home/user/some.txt -r /home/user/some.txt -ori rtl
```
