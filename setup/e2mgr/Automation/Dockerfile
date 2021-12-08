FROM python:3.6

RUN python3 -m pip install robotframework \
    && pip install --upgrade RESTinstance \
    && pip install docker \
    && pip install -U robotframework-requests\
    && apt-get update 

WORKDIR /opt
COPY ./Scripts /opt/Scripts
COPY ./Tests /opt/Tests
COPY ./run_tests.sh /opt/run_tests.sh

ENV DOCKER_HOST_IP  "127.0.0.1"

CMD [ "/opt/run_tests.sh" ]
