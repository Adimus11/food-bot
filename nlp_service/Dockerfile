FROM python:3.7-slim-buster

WORKDIR /app

COPY requirements.txt requirements.txt
RUN pip3 install -r requirements.txt

COPY . .

EXPOSE 9090:9090

RUN python3 -m spacy download en_core_web_sm

CMD [ "python3", "server.py"]