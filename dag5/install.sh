#chmod a+x ./install.sh
#chmod a+x $(pwd)install.sh
#scp -r project sindrevh@129.241.187.148:/Documents/project

#!/bin/bash
export GOPATH=$(pwd)
#mkdir $GOPATH
go install driver
echo "det kjorte"

#go help gopath
