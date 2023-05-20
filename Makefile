build:
	mkdir -p bin 
	go build -C $(d) -o ../bin/$(d)
