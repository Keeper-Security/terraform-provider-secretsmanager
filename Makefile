TEST?=$$(go list ./...)

$(eval AVERSION=$(shell git describe --tags --abbrev=0))
$(eval VERSION=$(patsubst v%,%,$(AVERSION)))
$(eval VVERSION=v$(VERSION))

build:
	go build

testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

install:
	go build -o ~/.terraform.d/plugins/terraform-provider-keeper

build-all:
	$(eval GOOS=darwin)  $(eval GOARCH=amd64) GOOS=${GOOS} GOARCH=${GOARCH} go build -o terraform-provider-keeper_${VVERSION}_${GOOS}_${GOARCH}
	$(eval GOOS=darwin)  $(eval GOARCH=arm64) GOOS=${GOOS} GOARCH=${GOARCH} go build -o terraform-provider-keeper_${VVERSION}_${GOOS}_${GOARCH}
	$(eval GOOS=freebsd) $(eval GOARCH=386)   GOOS=${GOOS} GOARCH=${GOARCH} go build -o terraform-provider-keeper_${VVERSION}_${GOOS}_${GOARCH}
	$(eval GOOS=freebsd) $(eval GOARCH=amd64) GOOS=${GOOS} GOARCH=${GOARCH} go build -o terraform-provider-keeper_${VVERSION}_${GOOS}_${GOARCH}
	$(eval GOOS=freebsd) $(eval GOARCH=arm)   GOOS=${GOOS} GOARCH=${GOARCH} go build -o terraform-provider-keeper_${VVERSION}_${GOOS}_${GOARCH}
	$(eval GOOS=linux)   $(eval GOARCH=386)   GOOS=${GOOS} GOARCH=${GOARCH} go build -o terraform-provider-keeper_${VVERSION}_${GOOS}_${GOARCH}
	$(eval GOOS=linux)   $(eval GOARCH=amd64) GOOS=${GOOS} GOARCH=${GOARCH} go build -o terraform-provider-keeper_${VVERSION}_${GOOS}_${GOARCH}
	$(eval GOOS=linux)   $(eval GOARCH=arm)   GOOS=${GOOS} GOARCH=${GOARCH} go build -o terraform-provider-keeper_${VVERSION}_${GOOS}_${GOARCH}
	$(eval GOOS=linux)   $(eval GOARCH=arm64) GOOS=${GOOS} GOARCH=${GOARCH} go build -o terraform-provider-keeper_${VVERSION}_${GOOS}_${GOARCH}
	$(eval GOOS=windows) $(eval GOARCH=386)   GOOS=${GOOS} GOARCH=${GOARCH} go build -o terraform-provider-keeper_${VVERSION}_${GOOS}_${GOARCH}.exe
	$(eval GOOS=windows) $(eval GOARCH=amd64) GOOS=${GOOS} GOARCH=${GOARCH} go build -o terraform-provider-keeper_${VVERSION}_${GOOS}_${GOARCH}.exe

release-all:
	$(eval GOOS=darwin)  $(eval GOARCH=amd64) GOOS=${GOOS} GOARCH=${GOARCH} go build -o terraform-provider-keeper_${VVERSION}; zip -m terraform-provider-keeper_${VERSION}_${GOOS}_${GOARCH}.zip terraform-provider-keeper_${VVERSION}
	$(eval GOOS=darwin)  $(eval GOARCH=arm64) GOOS=${GOOS} GOARCH=${GOARCH} go build -o terraform-provider-keeper_${VVERSION}; zip -m terraform-provider-keeper_${VERSION}_${GOOS}_${GOARCH}.zip terraform-provider-keeper_${VVERSION}
	$(eval GOOS=freebsd) $(eval GOARCH=386)   GOOS=${GOOS} GOARCH=${GOARCH} go build -o terraform-provider-keeper_${VVERSION}; zip -m terraform-provider-keeper_${VERSION}_${GOOS}_${GOARCH}.zip terraform-provider-keeper_${VVERSION}
	$(eval GOOS=freebsd) $(eval GOARCH=amd64) GOOS=${GOOS} GOARCH=${GOARCH} go build -o terraform-provider-keeper_${VVERSION}; zip -m terraform-provider-keeper_${VERSION}_${GOOS}_${GOARCH}.zip terraform-provider-keeper_${VVERSION}
	$(eval GOOS=freebsd) $(eval GOARCH=arm)   GOOS=${GOOS} GOARCH=${GOARCH} go build -o terraform-provider-keeper_${VVERSION}; zip -m terraform-provider-keeper_${VERSION}_${GOOS}_${GOARCH}.zip terraform-provider-keeper_${VVERSION}
	$(eval GOOS=linux)   $(eval GOARCH=386)   GOOS=${GOOS} GOARCH=${GOARCH} go build -o terraform-provider-keeper_${VVERSION}; zip -m terraform-provider-keeper_${VERSION}_${GOOS}_${GOARCH}.zip terraform-provider-keeper_${VVERSION}
	$(eval GOOS=linux)   $(eval GOARCH=amd64) GOOS=${GOOS} GOARCH=${GOARCH} go build -o terraform-provider-keeper_${VVERSION}; zip -m terraform-provider-keeper_${VERSION}_${GOOS}_${GOARCH}.zip terraform-provider-keeper_${VVERSION}
	$(eval GOOS=linux)   $(eval GOARCH=arm)   GOOS=${GOOS} GOARCH=${GOARCH} go build -o terraform-provider-keeper_${VVERSION}; zip -m terraform-provider-keeper_${VERSION}_${GOOS}_${GOARCH}.zip terraform-provider-keeper_${VVERSION}
	$(eval GOOS=linux)   $(eval GOARCH=arm64) GOOS=${GOOS} GOARCH=${GOARCH} go build -o terraform-provider-keeper_${VVERSION}; zip -m terraform-provider-keeper_${VERSION}_${GOOS}_${GOARCH}.zip terraform-provider-keeper_${VVERSION}
	$(eval GOOS=windows) $(eval GOARCH=386)   GOOS=${GOOS} GOARCH=${GOARCH} go build -o terraform-provider-keeper_${VVERSION}.exe; zip -m terraform-provider-keeper_${VERSION}_${GOOS}_${GOARCH}.zip terraform-provider-keeper_${VVERSION}.exe
	$(eval GOOS=windows) $(eval GOARCH=amd64) GOOS=${GOOS} GOARCH=${GOARCH} go build -o terraform-provider-keeper_${VVERSION}.exe; zip -m terraform-provider-keeper_${VERSION}_${GOOS}_${GOARCH}.zip terraform-provider-keeper_${VVERSION}.exe
	shasum -a 256 terraform-provider-keeper_${VERSION}_*.zip > terraform-provider-keeper_${VERSION}_SHA256SUMS
#	detached signature terraform-provider-keeper_keeper_${VERSION}_SHA256SUMS.{asc|sig}
#	gpg -ab terraform-provider-keeper_keeper_${VERSION}_SHA256SUMS
#	gpg -sb terraform-provider-keeper_keeper_${VERSION}_SHA256SUMS
