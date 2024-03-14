#!/bin/bash
cd ./internel/front
#npm install
#npm run build

#cp  ./dist ../../internel/serve/api/internal/handler/h5
file_list=$(find ./dist -type f)
#echo $file_list


for loop in $file_list
do
    if [[ "$(echo $loop | grep "css")" != "" ]]; then
      cp  $loop ../../internel/serve/api/internal/handler/h5/dist/assets/index.css
    elif [[ "$(echo $loop | grep "js")" != "" ]]; then
      cp  $loop ../../internel/serve/api/internal/handler/h5/dist/assets/index.js
    elif [[ "$(echo $loop | grep "html")" != "" ]]; then
      cp  $loop ../../internel/serve/api/internal/handler/h5/dist/index.html
#    elif [[ "$(echo $loop | grep "vite.svg")" != "" ]]; then
#      cp  $loop ../../internel/serve/api/internal/handler/h5/dist/vite.svg
#    elif [[ "$(echo $loop | grep "favicon.ico")" != "" ]]; then
#      cp  $loop ../../internel/serve/api/internal/handler/h5/dist/favicon.ico
    else
        echo $loop
    fi

done
