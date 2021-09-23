TEST?=$$(go list ./...)

VERSION=`git describe --always`
#VERSION=1.0.0

build-all:
	$(eval GOOS=darwin)  $(eval GOARCH=amd64) GOOS=${GOOS} GOARCH=${GOARCH} go build -o terraform-provider-keeper_v${VERSION}_${GOOS}_${GOARCH}
	$(eval GOOS=freebsd) $(eval GOARCH=386)   GOOS=${GOOS} GOARCH=${GOARCH} go build -o terraform-provider-keeper_v${VERSION}_${GOOS}_${GOARCH}
	$(eval GOOS=freebsd) $(eval GOARCH=amd64) GOOS=${GOOS} GOARCH=${GOARCH} go build -o terraform-provider-keeper_v${VERSION}_${GOOS}_${GOARCH}
	$(eval GOOS=freebsd) $(eval GOARCH=arm)   GOOS=${GOOS} GOARCH=${GOARCH} go build -o terraform-provider-keeper_v${VERSION}_${GOOS}_${GOARCH}
	$(eval GOOS=linux)   $(eval GOARCH=386)   GOOS=${GOOS} GOARCH=${GOARCH} go build -o terraform-provider-keeper_v${VERSION}_${GOOS}_${GOARCH}
	$(eval GOOS=linux)   $(eval GOARCH=amd64) GOOS=${GOOS} GOARCH=${GOARCH} go build -o terraform-provider-keeper_v${VERSION}_${GOOS}_${GOARCH}
	$(eval GOOS=linux)   $(eval GOARCH=arm)   GOOS=${GOOS} GOARCH=${GOARCH} go build -o terraform-provider-keeper_v${VERSION}_${GOOS}_${GOARCH}
	$(eval GOOS=linux)   $(eval GOARCH=arm64) GOOS=${GOOS} GOARCH=${GOARCH} go build -o terraform-provider-keeper_v${VERSION}_${GOOS}_${GOARCH}
	$(eval GOOS=windows) $(eval GOARCH=386)   GOOS=${GOOS} GOARCH=${GOARCH} go build -o terraform-provider-keeper_v${VERSION}_${GOOS}_${GOARCH}.exe
	$(eval GOOS=windows) $(eval GOARCH=amd64) GOOS=${GOOS} GOARCH=${GOARCH} go build -o terraform-provider-keeper_v${VERSION}_${GOOS}_${GOARCH}.exe

release-all:
	$(eval GOOS=darwin)  $(eval GOARCH=amd64) GOOS=${GOOS} GOARCH=${GOARCH} go build -o terraform-provider-keeper_v${VERSION}; zip -m terraform-provider-keeper_${VERSION}_${GOOS}_${GOARCH}.zip terraform-provider-keeper_v${VERSION}
	$(eval GOOS=freebsd) $(eval GOARCH=386)   GOOS=${GOOS} GOARCH=${GOARCH} go build -o terraform-provider-keeper_v${VERSION}; zip -m terraform-provider-keeper_${VERSION}_${GOOS}_${GOARCH}.zip terraform-provider-keeper_v${VERSION}
	$(eval GOOS=freebsd) $(eval GOARCH=amd64) GOOS=${GOOS} GOARCH=${GOARCH} go build -o terraform-provider-keeper_v${VERSION}; zip -m terraform-provider-keeper_${VERSION}_${GOOS}_${GOARCH}.zip terraform-provider-keeper_v${VERSION}
	$(eval GOOS=freebsd) $(eval GOARCH=arm)   GOOS=${GOOS} GOARCH=${GOARCH} go build -o terraform-provider-keeper_v${VERSION}; zip -m terraform-provider-keeper_${VERSION}_${GOOS}_${GOARCH}.zip terraform-provider-keeper_v${VERSION}
	$(eval GOOS=linux)   $(eval GOARCH=386)   GOOS=${GOOS} GOARCH=${GOARCH} go build -o terraform-provider-keeper_v${VERSION}; zip -m terraform-provider-keeper_${VERSION}_${GOOS}_${GOARCH}.zip terraform-provider-keeper_v${VERSION}
	$(eval GOOS=linux)   $(eval GOARCH=amd64) GOOS=${GOOS} GOARCH=${GOARCH} go build -o terraform-provider-keeper_v${VERSION}; zip -m terraform-provider-keeper_${VERSION}_${GOOS}_${GOARCH}.zip terraform-provider-keeper_v${VERSION}
	$(eval GOOS=linux)   $(eval GOARCH=arm)   GOOS=${GOOS} GOARCH=${GOARCH} go build -o terraform-provider-keeper_v${VERSION}; zip -m terraform-provider-keeper_${VERSION}_${GOOS}_${GOARCH}.zip terraform-provider-keeper_v${VERSION}
	$(eval GOOS=linux)   $(eval GOARCH=arm64) GOOS=${GOOS} GOARCH=${GOARCH} go build -o terraform-provider-keeper_v${VERSION}; zip -m terraform-provider-keeper_${VERSION}_${GOOS}_${GOARCH}.zip terraform-provider-keeper_v${VERSION}
	$(eval GOOS=windows) $(eval GOARCH=386)   GOOS=${GOOS} GOARCH=${GOARCH} go build -o terraform-provider-keeper_v${VERSION}.exe; zip -m terraform-provider-keeper_${VERSION}_${GOOS}_${GOARCH}.zip terraform-provider-keeper_v${VERSION}.exe
	$(eval GOOS=windows) $(eval GOARCH=amd64) GOOS=${GOOS} GOARCH=${GOARCH} go build -o terraform-provider-keeper_v${VERSION}.exe; zip -m terraform-provider-keeper_${VERSION}_${GOOS}_${GOARCH}.zip terraform-provider-keeper_v${VERSION}.exe
	shasum -a 256 terraform-provider-keeper_${VERSION}_*.zip > terraform-provider-keeper_${VERSION}_SHA256SUMS
#	detached signature terraform-provider-keeper_keeper_${VERSION}_SHA256SUMS.{asc|sig}
#	gpg -ab terraform-provider-keeper_keeper_${VERSION}_SHA256SUMS
#	gpg -sb terraform-provider-keeper_keeper_${VERSION}_SHA256SUMS

testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m
