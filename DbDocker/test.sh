#!/bin/sh
while true
do
	cd ..
	cp db/db.db ./shared
	echo "copyDone"
	sleep 20
	cd db
	echo "sleep end"
	
done