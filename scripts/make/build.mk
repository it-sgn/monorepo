.PHONY: run-all

run-all:
	$(MAKE) employers-run
	$(MAKE) department-run
	$(MAKE) biometric-run
.PHONY: build-all
build-all:
	$(MAKE) employers
	$(MAKE) department
	$(MAKE) biometric 
	$(MAKE) log_downloader 
	$(MAKE) attendance-raw

.PHONY: employers department biometric log_downloader attendance-raw

build: employers department biometric log_downloader attendance-raw

employers:
	cd module/employers/service && mkdir -p bin/ && \
	go build -buildvcs=false -ldflags="-X main.Version=1.0" -o ./bin/ ./...

department:
	cd module/department/service && mkdir -p bin/ && \
	go build -buildvcs=false -ldflags="-X main.Version=1.0" -o ./bin/ ./...

biometric:
	cd module/biometric/service && mkdir -p bin/ && \
	go build -buildvcs=false -ldflags="-X main.Version=1.0" -o ./bin/ ./...

log_downloader:
	cd module/log_downloader/service && mkdir -p bin/ && \
	go build -buildvcs=false -ldflags="-X main.Version=1.0" -o ./bin/ ./...

attendance-raw:
	cd module/attendance-raw/service && mkdir -p bin/ && \
	go build -buildvcs=false -ldflags="-X main.Version=1.0" -o ./bin/ ./...






employers-run:
	cd module/employers/service && mkdir -p bin/ && \
	go build -buildvcs=false -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./... && \
	./bin/server -conf ./configs

department-run:
	cd module/department/service && mkdir -p bin/ && \
	go build -buildvcs=false -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./... && \
	./bin/server -conf ./configs


biometric-run:
	cd module/biometric/service && mkdir -p bin/ && \
	go build -buildvcs=false -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./... && \
	./bin/server -conf ./configs


# .PHONY: app
# app:
# 	cd 	app/app/service && mkdir -p bin/ && \
# 	go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./... && \
# 	 ./bin/server -conf ./configs


# .PHONY: user
# user:
# 	cd 	app/user/service && mkdir -p bin/ && \
# 	go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./... && \
# 	  ./bin/server -conf ./configs

# .PHONY: order
# order:
# 	cd 	app/order/service && mkdir -p bin/ && \
# 	go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./... && \
# 	 ./bin/server -conf ./configs

# .PHONY: spu
# spu:
# 	cd 	app/spu/service && mkdir -p bin/ && \
# 	go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./... && \
# 	 ./bin/server -conf ./configs

# .PHONY: interface
# interface:
# 	cd 	app/mall/interface && mkdir -p bin/ && \
# 	go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./... && \
# 	 ./bin/server -conf ./configs

# .PHONY: employers
# employers:
# 	cd 	module/employers/service && mkdir -p bin/ && \
# 	go build -buildvcs=false -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./... && \
# 	 ./bin/server -conf ./configs	 

# .PHONY: department
# department:
# 	cd 	module/department/service && mkdir -p bin/ && \
# 	go build -buildvcs=false -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./... && \
# 	 ./bin/server -conf ./configs	


# .PHONY: biometric
# biometric:
# 	cd 	module/biometric/service && mkdir -p bin/ && \
# 	go build -buildvcs=false -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./... && \
# 	 ./bin/server -conf ./configs		 