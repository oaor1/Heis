#chmod +x ./install.sh

#!/bin/bash
export GOPATH=$(pwd)
#mkdir $GOPATH
go install driver

#go help gopath