
OUT_PATH := ./bin
MAIN_PATH := ./cmd
CMDS := $(wildcard $(MAIN_PATH)/*)
OUT_FILES := $(patsubst $(MAIN_PATH)/%, $(OUT_PATH)/% , $(CMDS))

all: build

$(OUT_PATH)/%: $(MAIN_PATH)/%
	go build -o $(OUT_PATH)/ ./$^

build: $(OUT_FILES)

clean: $(OUT_PATH)
	rm -rf $(OUT_PATH)/*
	rmdir $(OUT_PATH)
