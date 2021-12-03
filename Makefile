
OUT_PATH := ./bin
MAIN_PATH := ./cmd
PKGS_PAHT := ./pkg
CMDS := $(wildcard $(MAIN_PATH)/*)
SOURCES := $(join $(wildcard $(MAIN_PATH)/**/*.go), $(wildcard $(PKGS_PAHT)/**/*.go))
# SOURCES := $(shell find ./ -name '*.go')
OUT_FILES := $(patsubst $(MAIN_PATH)/%, $(OUT_PATH)/% , $(CMDS))

all: build


$(OUT_PATH)/%: $(MAIN_PATH)/%
	go build -o $(OUT_PATH)/ ./$^

build: clean $(OUT_FILES)

clean:
	if [ -d "$(OUT_PATH)" ]; then \
		rm -rf $(OUT_PATH)/* \
		rmdir $(OUT_PATH); \
	fi
