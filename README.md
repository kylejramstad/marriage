# Download, Build, and Run on Raspberry Pi
Follow these instructions to install the necessary software to build and run the code on your own Raspberry Pi.

## Install software
### Ensure you have git installed on your Raspberry Pi
```
sudo apt-get install git
```

### Install Go on your Raspberry Pi
```
sudo apt-get install golang
```

### Install Dependancies
```
sudo go get golang.org/x/net/publicsuffix
```
## Download Project
```
sudo git clone https://github.com/kylejramstad/marriage.git
```
### Clean up project files
```
sudo mv ./marriage/* ./ && sudo rm -r marriage
```

## Connect with IFTTT
### Make a webhooks event
Go to [IFTTT](https://ifttt.com/create) and create an event.
1. Click "This" and search for "webhooks"
1. Click "Receive a web request"
1. Give your webhooks event a name like "marriage"
1. Click "That" and search for "email"
1. Click "Send me an email"
	1. You can send an email to your phone number (send a text message) by following the instructions at [https://www.androidpolice.com/2018/07/28/get-around-ifttts-cap-sending-sms-messages/](https://www.androidpolice.com/2018/07/28/get-around-ifttts-cap-sending-sms-messages/)
1. Change the Subject to something like "Cupid Project" change the body to 
	1. ```There is a Project Cupid appointment available for the following date(s): {{Value1}}```
1. Click "Create Action"

#### Find your IFTTT webhooks key
Go to the [IFTTT Webhooks page](https://ifttt.com/maker_webhooks) and click "Documentation".
Here you will find your key.

Open main.go with your favorite text editor. Example: ```nano main.go```
Change the values of event and IFTTTkey to match your webhooks event name and your webhooks key.
```
const (
	unavailable        = "unavailable"
	event              = ""
	IFTTTkey           = ""
)
```
Save the file and return to the command line.

## Build and Run Project
### Build
```
sudo go build main.go
```

### Test IFTTT Notifications
```
./main -test
```
This is a test message.
I hope this makes it to its destination.
If not, then check to make sure your event name and IFTTTkey is correct.
Also check to make sure you set up your webhooks applet correctly.

### Set up Crontab to Run Program Automatically
```crontab -e```
Add this line to the end of the file
```0 * * * * ~/main```
This will run the program once every hour.
If you want to change this amount go to [Crontab Guru](https://crontab.guru/)

## View Logs to Ensure Proper Running
```
cat marriage.log
```
This will print out a log of all the times the program ran and it's results. This can be used to verify the programs operation.

