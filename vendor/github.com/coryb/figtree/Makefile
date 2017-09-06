GENERATOR_SRC = \
	rawoption.go \
	$(NULL)

GENERATED_SRC = $(GENERATOR_SRC:%.go=gen-%.go)

test: $(GENERATED_SRC)
	go get -t -v
	go get github.com/kr/pretty
	go get gopkg.in/alecthomas/kingpin.v2
	go test

gen-%.go: %.go
	# use github.com/cheekybits/genny after https://github.com/cheekybits/genny/pull/42 is merged
	go get github.com/coryb/genny
	go generate

.PHONY: test
