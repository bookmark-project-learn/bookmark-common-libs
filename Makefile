.PHONY: test test-nocache

# Coverage exclude files and directories out of report
COVERAGE_EXCLUDE=mocks|main.go|test|config.go|infrastructure/**
COVERAGE_THRESHOLD ?= 80

# Process check run make test whether cache or no-cache
CACHECMD := 
OPTION ?= cache
ifeq ($(OPTION),nocache)
    CACHECMD := go clean -testcache    
endif


test: 
	@mkdir -p ./test
	# run clean test cache
	$(CACHECMD)

	go test ./... -coverprofile=./test/coverage_tmp -covermode=atomic -coverpkg=./... -p 1	
	grep -vE "$(COVERAGE_EXCLUDE)" ./test/coverage_tmp > ./test/coverage_out
	go tool cover -html=./test/coverage_out -o ./test/coverage.html
	@total=$$(go tool cover -func=./test/coverage_out | grep total: | awk '{print $$3}' | sed 's/%//'); \
	if [ $$(echo "$$total < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
		echo "❌ Coverage ($$total%) is below threshold ($(COVERAGE_THRESHOLD)%)"; \
		exit 1; \
	else \
		echo "✅ Coverage ($$total%) meets threshold ($(COVERAGE_THRESHOLD)%)"; \
	fi

test-nocache: test OPTION=nocache

