# MIDDiff
### Know when your summer assignment has changed
*This tool may not work due to the database system switching from HTTP Authentication to SAML*

## Prerequisites
Install Golang environment if you do not have it installed already<br/>
```
git clone https://go.googlesource.com/go
cd go
git checkout go1.4.2
cd src
./all.bash
export $PATH=$PATH:$HOME/go/bin
```

## Clone the code
In the directory you want the middiff folder, run `git clone https://github.com/anson/middiff`. 

## Build it
Build MIDDiff with `go build middiff.go mail.go`.

The binary will be named _middiff_.

## Setup credentials and URLs
Create credentials file with the format of *username \n password \n*...
Create URL file with the format of *human description\n url path \n*...

## Automate with a Cron job
Edit the user crontab file `crontab -e`.  
Add the lines:

```
MAILTO="user@usna.edu"
0 0 * * * ./middiff/BINARYLOCATION -credentials=./CREDENTIALSLOCATION -urllist=./URLLISTLOCATION
```
