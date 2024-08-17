SHELL := /bin/bash

all: venv install

venv:
	python3 -m pip install virtualenv
	virtualenv .venv
	# source .venv/bin/activate

install:
	source .venv/bin/activate
	python3 -m pip install  -r requirements.txt

run-web:
	locust --host http://localhost:8080