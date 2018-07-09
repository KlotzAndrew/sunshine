#!/bin/bash

while true
do
	curl 0.0.0.0:1323/basic
	curl 0.0.0.0:1323/dbwrite
	curl 0.0.0.0:1323/dbwrite_slow

	curl 0.0.0.0:1324/basic
	for i in {1..40}; do curl 0.0.0.0:1324/dbwrite; done
	curl 0.0.0.0:1324/dbwrite_slow
done
