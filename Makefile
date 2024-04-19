d:
	$(MAKE) d-build
	$(MAKE) d-run
d-build:
	docker build -t go-currency-rates .
d-run:
	 docker run -p 80:9999 go-currency-rates
d-sh:
	docker run -it go-currency-rates sh
d-log:
	docker logs -f $(shell docker ps -qf "ancestor=go-currency-rates")
