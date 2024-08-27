SHELL := /bin/bash

all: venv install

venv:
	python3 -m pip install virtualenv
	virtualenv .venv
	# source .venv/bin/activate

install:
	source .venv/bin/activate
	python3 -m pip install  -r requirements.txt

# writes requiremets list to file
requirements:
	pip freeze > requirements.txt

run-web:
	locust -f tests,loads --class-picker --host http://localhost:8080 --processes 16


