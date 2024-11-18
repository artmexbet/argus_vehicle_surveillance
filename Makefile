PYTHON_BUILD_IMAGE_NAME=python-build
PYTHON_BUILD_IMAGE_PATH=deploy/python/Dockerfile

GO_BUILD_IMAGE_NAME=go-build
GO_BUILD_IMAGE_PATH=deploy/golang/build/Dockerfile

build:
	docker build -t $(GO_BUILD_IMAGE_NAME) -f $(GO_BUILD_IMAGE_PATH)  .

build-python:
	( \
		python -m venv YOLO/venv; \
		source YOLO/venv/Scripts/activate; \
		pip install torch torchvision torchaudio --index-url https://download.pytorch.org/whl/cu124; \
		pip install -r requirements.txt; \
	)

start:
	docker-compose up -d

stop:
	docker-compose down