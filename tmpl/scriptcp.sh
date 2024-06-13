#!/bin/sh
while true
do
	cd ..
	cp shared/tmpl/. ./db
	echo "copyDone"
	sleep 30
	cp tmpl/. ./shared/tmpl
	cd tmpl
	echo "sleep end"
	
done

