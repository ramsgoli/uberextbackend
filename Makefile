build:
	docker build -t ramsgoli/uberextbackend .

run_dev:
	docker run -it -p 8000:8000 ramsgoli/uberextbackend
