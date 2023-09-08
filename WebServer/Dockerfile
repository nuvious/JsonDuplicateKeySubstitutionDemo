FROM python:3.8-slim-buster

WORKDIR /python-docker

RUN pip3 install Flask==2.3.3 requests==2.31.0

COPY . .

CMD [ "python3", "-m" , "flask", "run", "--host=0.0.0.0"]