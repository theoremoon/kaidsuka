.PHONY: run
run: build
	cd build ;\
		./main

.PHONY: build
build: main.cpp
	mkdir -p build

	cd build ;\
		cmake .. ;\
		make

.PHONY: clean
clean:
	rm -rf build

