#chmod a+x ./install.sh
#chmod a+x $(pwd)install.sh

#!/bin/bash
export GOPATH=$(pwd)
#mkdir $GOPATH
go install driver
echo "det kjorte"

#go help gopath