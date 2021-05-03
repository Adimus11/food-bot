from concurrent import futures
import time
import math
import logging

import grpc

import communication_pb2
import communication_pb2_grpc

import spacy

import re

import en_core_web_sm


class NLPService(communication_pb2_grpc.IngridientsServiceServicer):
    """Provides methods that implement functionality of route guide server."""

    def __init__(self):
        self.model = en_core_web_sm.load()
        self.measurements = re.compile(r'(bowl|bulb|cube|clove|cup|drop|ounce|oz|pinch|pound|teaspoon|tablespoon)s?')

    def GetIngridients(self, request, context):

        ingredients = []
        for token in self.model(request.text):
            if (token.pos_ in ['NOUN', 'PROPN']) and (not self.measurements.match(token.text)):
                ingredients.append(str(token))

        return communication_pb2.Ingridients(ingridients=ingredients)

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    communication_pb2_grpc.add_IngridientsServiceServicer_to_server(
        NLPService(), server)
    port = 9090
    server.add_insecure_port(f'[::]:{port}')
    server.start()
    print(f"Server started at: {port}")
    server.wait_for_termination()


if __name__ == '__main__':
    logging.basicConfig()
    serve()