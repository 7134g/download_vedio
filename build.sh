#!/bin/bash
for loop in "build" "build/etc"
do
    if [ -d "$loop" ]; then
        echo "Folder already exists" "$loop"
    else
        mkdir "$loop"
        echo "Folder created successfully"
    fi
done


systemName=$(uname -a)
if [[ "$(echo $systemName | grep "Darwin")" != "" ]]
then
    go build -o ./build/serve ./internel/serve/api/task_serve.go
    go build -o ./build/dv ./internel/serve/api/tool/main.go
elif [[ "$(echo $systemName | grep "Linux")" != "" ]]
then
    go build -o ./build/serve ./internel/serve/api/task_serve.go
    go build -o ./build/dv ./internel/serve/api/tool/main.go
else
    go build -o ./build/serve.exe ./internel/serve/api/task_serve.go
    go build -o ./build/dv.exe ./internel/serve/api/tool/main.go
fi

cp ./internel/serve/api/etc/task_serve.yaml ./build/etc/task_serve.yaml




