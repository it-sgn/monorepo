.PHONY: app
app:
	cd 	app/app/service && mkdir -p bin/ && \
	go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./... && \
	 ./bin/server -conf ./configs


.PHONY: user
user:
	cd 	app/user/service && mkdir -p bin/ && \
	go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./... && \
	  ./bin/server -conf ./configs

.PHONY: order
order:
	cd 	app/order/service && mkdir -p bin/ && \
	go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./... && \
	 ./bin/server -conf ./configs

.PHONY: spu
spu:
	cd 	app/spu/service && mkdir -p bin/ && \
	go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./... && \
	 ./bin/server -conf ./configs

.PHONY: interface
interface:
	cd 	app/mall/interface && mkdir -p bin/ && \
	go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./... && \
	 ./bin/server -conf ./configs

.PHONY: employers
employers:
	cd 	module/employers/service && mkdir -p bin/ && \
	go build -buildvcs=false -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./... && \
	 ./bin/server -conf ./configs	 

.PHONY: department
department:
	cd 	module/department/service && mkdir -p bin/ && \
	go build -buildvcs=false -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./... && \
	 ./bin/server -conf ./configs	


.PHONY: biometric
biometric:
	cd 	module/biometric/service && mkdir -p bin/ && \
	go build -buildvcs=false -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./... && \
	 ./bin/server -conf ./configs		 