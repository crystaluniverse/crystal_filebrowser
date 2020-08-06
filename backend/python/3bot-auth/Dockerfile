FROM ubuntu:latest

RUN apt-get update && apt-get upgrade -y
RUN apt-get install -y python3 python3-pip
RUN pip3 --version

COPY ./python/3bot-auth/ /app

RUN pip3 install -r /app/requirements.txt
EXPOSE 5001
CMD ["python3", "/app/app.py"]