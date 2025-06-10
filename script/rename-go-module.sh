#!/usr/bin/env sh

export CUR="template" 
export NEW="github.com/BC-Technology/my-project" 
go mod edit -module ${NEW}
find . -type f -name '*.go' -exec perl -pi -e 's/$ENV{CUR}/$ENV{NEW}/g' {} \;