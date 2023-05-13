VIRTUALENV_NAME=gcloud-run.$(shell pwd | rev | cut -d '/' -f 1 | rev)
IMAGE_URL=gcr.io/whatsapp-speech-to-text-bot/openai-transcribe:latest

.PHONY: install uninstall clean docker localhost

uninstall:
	@./bin/remove_pyenv.sh

install: uninstall
	@./bin/setup_pyenv.sh $(VIRTUALENV_NAME)

clean:
	@./bin/cleanup.sh

docker:
	docker build -f Dockerfile -t $(IMAGE_URL) .

localhost: docker
	@./bin/launch.sh $(IMAGE_URL)
