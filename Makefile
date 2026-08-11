# Dock-Diet Makefile

APP_NAME = dock-diet

# Dependencies clean karna
tidy:
	go mod tidy

# Tool ko build karna
build: tidy
	go build -o $(APP_NAME)

# Saare tests run karna
test:
	go test ./... -v

# Compiled file aur kachra clean karna
clean:
	rm -f $(APP_NAME)
	rm -f *.optimized

# Direct build kar ke run karna (Test file ke liye)
run-scan: build
	./$(APP_NAME) scan Dockerfile