install:
	go build -o denv cmd/denv/main.go && sudo mv denv /usr/local/bin/denv && sudo chmod 777 /usr/local/bin/denv