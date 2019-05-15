all:
	go build -buildmode=c-shared -o out_gdetail.so .

fast:
	go build out_gdetail.go

clean:
	rm -rf *.so *.h *~