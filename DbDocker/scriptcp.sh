#!/bin/sh
while true
do
	cd ..
	cp shared/db.db ./db
	echo "copyDone"
	sleep 30
	cp db/db.db ./shared
	cd db
	echo "sleep end"
	
done

