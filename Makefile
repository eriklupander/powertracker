.PHONY build-all:
build-all:
	make -C functions/exporter build
	make -C functions/powerrecorder build
	make -C functions/statusrecorder build
	make -C functions/statusapi build