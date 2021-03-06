EXEC = kbdashboard
COMP = $(EXEC).bash-completion

VER = `grep "const VERSION" cmd_version.go  | cut -d "=" -f 2 | cut -d '"' -f 2`
TAR = $(EXEC)-$(VER).tar.gz

BUILD_TIME = `date +%Y-%m-%d_%H:%M:%S`
GIT_COMMIT=`git log --pretty=format:"%h" -1`
GIT_BRANCH=`git rev-parse --abbrev-ref HEAD`

# Add build-time-string into the executable file.
X_ARGS += -X main.BUILD_TIME=$(BUILD_TIME)
X_ARGS += -X main.GIT_COMMIT="$(GIT_COMMIT)@$(GIT_BRANCH)"
X_ARGS += -X main.COMP_FILENAME=$(COMP)

BIN = $(DESTDIR)/usr/bin
COMP_DIR = $(DESTDIR)/etc/bash_completion.d

all:bin $(COMP)

bin:
	@echo "Build Version: $(VER)"
	@go build -ldflags "$(X_ARGS)" -o $(EXEC)

$(COMP): $(EXEC)
	@./$(EXEC) completion

install:$(EXEC) $(COMP)
	install -d $(BIN) $(COMP_DIR)
	install $(EXEC) $(BIN)
	install $(COMP) $(COMP_DIR)/$(EXEC)

clean:
	@rm -rfv $(EXEC) $(COMP)

archive:
	@echo "archive to $(TAR)"
	@git archive master --prefix="$(EXEC)-$(VER)/" --format tar.gz -o $(TAR)

test:
	@go test
	@go test -v -race ./...
