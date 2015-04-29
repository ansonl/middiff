# MIDDiff
### Know when your summer assignment has changed

## Prerequesites
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
In the directory you want the middiff folder, run<code>git clone https://github.com/anson/middiff</code>. 

## Build it
Build MIDDiff with <code>go build middiff.go mail.go</code>.

The binary will be named <i>middiff</i>.

## Automate with a Cron job
Edit the user crontab file<code>crontab -e</code>.<br/>
Add the lines:<br/>
* <code>
MAILTO="user@usna.edu"</code><br/>
*  <code>
0 0 * * * ./middiff/BINARYLOCATION -credentials=./CREDENTIALSLOCATION
</code>
