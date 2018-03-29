SHELL = /bin/bash 

.PHONY: clean docs venv

clean:
	rm -rf docs/output
	rm -rf docs/contracts
	rm -rf docs/doxity/pages/docs

venv:
	python3 -m venv ~/venv
	source ~/venv/bin/activate
	pip3 install -r requirements.txt

docs: venv
	bash scripts/docs.sh