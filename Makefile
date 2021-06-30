all: listener vsockcat

listener:
	go build -v \
    	-trimpath \
    	-buildmode=pie \
    	-mod=readonly \
    	-modcacherw \
    	-ldflags="-s -w" \
    	./cmd/listener

vsockcat:
	go build -v \
    	-trimpath \
    	-buildmode=pie \
    	-mod=readonly \
    	-modcacherw \
    	-ldflags="-s -w" \
    	./cmd/vsockcat

clean:
	go clean
