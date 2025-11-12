install-deps:
	pip install -r requirements.txt

install-pip:
	python3 -m pip install --upgrade pip

clean:
	rm -rf venv

venv: clean
	python3 -m venv venv
	. venv/bin/activate
	pip3 install -r requirements.txt

run:
	. venv/bin/activate && python3 ./main.py