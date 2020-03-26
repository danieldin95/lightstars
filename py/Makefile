.PHONY: release install publish

all:
	@echo "make release - prepares a release and publishes it"
	@echo "make dev - prepares a development environment"
	@echo "make install - install on local system"
	@echo "make publish - upload to pypi"

release: publish

dev:
	pip install -requirements.txt
	python setup.py develop

install:
	python setup.py install

publish:
	python setup.py sdist upload
