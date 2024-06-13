#! /bin/bash
while true
do
    cd ..
	cp db/db.db ./shared
    echo "copy Done"
	sleep 60
    cd ./db
    echo " Sleep reussy"
done
