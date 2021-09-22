GOOS=linux
GOARCH=amd64
GOARM=6


.PHONY: clean
clean:
ifneq ("$(wildcard cam-viewer)","")
	@rm media-player
endif
ifneq ("$(wildcard cam-viewer-arm)","")
	@rm media-player-arm
endif


.PHONY: build
build:
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o bin/cam-viewer . || (echo "build failed $$?"; exit 1)
	@echo 'Build suceeded... done'


.PHONY: buildarm
buildarm: GOARCH=arm GOARM=6
buildarm:
	@GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) go build -o bin/cam-viewer-arm . || (echo "build failed $$?"; exit 1)
	@echo 'Build suceeded... done'


.PHONY: push-to-samcam
push-to-samcam: buildarm
	cp bin/cam-viewer-arm .
	7z a cv.zip -y '-x!static/videos' static/ templates/ cam-viewer-arm
	rm cam-viewer-arm
	ssh samcam "sudo systemctl stop cam-viewer.service"
	scp cv.zip samcam:~/cam-viewer/
	ssh samcam "cd ~/cam-viewer/ && unzip -o cv.zip && rm cv.zip"
	ssh samcam "sudo systemctl start cam-viewer.service"
