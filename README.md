# MIDDiff
### Know when your summer assignment has changed

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

## Automate with a Cron job
Edit the user crontab file<code>crontab -e</code>.<br/>
Add the lines:<br/>
* <code>
MAILTO="user@usna.edu"</code><br/>
*  <code>
0 0 * * * ./middiff/BINARYLOCATION -credentials=./CREDENTIALSLOCATION
</code>
