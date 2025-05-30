# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc
import warnings

import trader_pb2 as trader__pb2

GRPC_GENERATED_VERSION = '1.71.0'
GRPC_VERSION = grpc.__version__
_version_not_supported = False

try:
    from grpc._utilities import first_version_is_lower
    _version_not_supported = first_version_is_lower(GRPC_VERSION, GRPC_GENERATED_VERSION)
except ImportError:
    _version_not_supported = True

if _version_not_supported:
    raise RuntimeError(
        f'The grpc package installed is at version {GRPC_VERSION},'
        + f' but the generated code in trader_pb2_grpc.py depends on'
        + f' grpcio>={GRPC_GENERATED_VERSION}.'
        + f' Please upgrade your grpc module to grpcio>={GRPC_GENERATED_VERSION}'
        + f' or downgrade your generated code using grpcio-tools<={GRPC_VERSION}.'
    )


class TraderStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.GetOffers = channel.unary_unary(
                '/trader.Trader/GetOffers',
                request_serializer=trader__pb2.Filter.SerializeToString,
                response_deserializer=trader__pb2.OfferList.FromString,
                _registered_method=True)
        self.Subscribe = channel.unary_stream(
                '/trader.Trader/Subscribe',
                request_serializer=trader__pb2.Subscription.SerializeToString,
                response_deserializer=trader__pb2.OfferList.FromString,
                _registered_method=True)


class TraderServicer(object):
    """Missing associated documentation comment in .proto file."""

    def GetOffers(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def Subscribe(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_TraderServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'GetOffers': grpc.unary_unary_rpc_method_handler(
                    servicer.GetOffers,
                    request_deserializer=trader__pb2.Filter.FromString,
                    response_serializer=trader__pb2.OfferList.SerializeToString,
            ),
            'Subscribe': grpc.unary_stream_rpc_method_handler(
                    servicer.Subscribe,
                    request_deserializer=trader__pb2.Subscription.FromString,
                    response_serializer=trader__pb2.OfferList.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'trader.Trader', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))
    server.add_registered_method_handlers('trader.Trader', rpc_method_handlers)


 # This class is part of an EXPERIMENTAL API.
class Trader(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def GetOffers(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/trader.Trader/GetOffers',
            trader__pb2.Filter.SerializeToString,
            trader__pb2.OfferList.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def Subscribe(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_stream(
            request,
            target,
            '/trader.Trader/Subscribe',
            trader__pb2.Subscription.SerializeToString,
            trader__pb2.OfferList.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)
