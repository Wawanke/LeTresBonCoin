#!/bin/sh
while true
do
	cd ..
	cp -r shared/tmpl ./tmpl
	echo "copyDone"
	sleep 30

	mkdir tmpl ./shared
	cp -r tpml ./shared/tmpl
	cd tmpl
	echo "sleep end"
	
done

