import signal
import grpc
import trader_pb2_grpc
import trader_pb2
from collections import defaultdict
import threading

# This is a dictionary to store the notifications
inbox = defaultdict(list)
subscription_thread = None
stop_event = threading.Event()

def list_offers(stub)-> dict:
    def pprint_offer(offer):
        print(f"Offer ID: {offer.symbol}, Current price: {offer.currentPrice}, "
              f"Offer currency: {trader_pb2.Currency.Name(offer.baseCurrency)}, instrument type: {trader_pb2.InstrumentType.Name(offer.type)}")
        print("-----------------------------------------------------")
        
    offers = stub.GetOffers(trader_pb2.Filter()).offers
    for offer in offers:
        pprint_offer(offer)
    return [offer.symbol for offer in offers]


def read_notifications():
    if not inbox:
        print("No notifications to read.")
        return
    for key, values in inbox.items():
        print(f"Notifications for {key}:")
        for value in values:
            print(value)
        print()
    print("-----------------------------------------------------")

def subscribe(stub, subscribe_to: list[str]):
    global stop_event
    stop_event.clear()  # Reset the stop event

    def subscription_worker():
        try:
            results = stub.Subscribe(trader_pb2.Subscription(symbols=subscribe_to))
            for newOffer in results:
                if stop_event.is_set():
                    break
                for offer in newOffer.offers:
                    inbox[offer.symbol].append(offer.currentPrice)
                    inbox[offer.symbol] = inbox[offer.symbol][-5:]  # Keep only the last 5 notifications

        except grpc.RpcError as e:
            if e.code() == grpc.StatusCode.CANCELLED:
                print("Subscription cancelled")
            else:
                print(f"Error: {e.code()} - {e.details()}")

    # Start the subscription worker in a separate thread
    thread = threading.Thread(target=subscription_worker, daemon=True)
    thread.start()
    return thread


client_dialog = """
Instructions:
1. 'list' to list offers
2. 'subscribe' to subscribe to offers
3. 'read' to read notifications
4. 'exit' to exit
> """
def run():
    global subscription_thread, stop_event

    available_offers_symbols = []
    with grpc.insecure_channel('localhost:50051') as channel:
        stub = trader_pb2_grpc.TraderStub(channel)
        while True:
            instr = input(client_dialog).lower().strip()
            if instr.startswith('list'):
                # List offers
                available_offers_symbols = list_offers(stub)
                print("Available offers:", available_offers_symbols)

            elif instr == 'read':
                # Read notifications
                read_notifications()

            elif instr.startswith('sub'):
                symbols = list(map(str.upper, instr.split()[1:]))
                if not symbols:
                    print("Please provide at least one symbol to subscribe to.")
                    continue
                subscribe_to = []
                for symbol in symbols:
                    if symbol in available_offers_symbols:
                            subscribe_to.append(symbol)
                    else:
                        print(f"Symbol {symbol} not found in available offers.")
                if not subscribe_to:
                    print("No valid symbols to subscribe to.")
                    continue

                # Stop the current subscription thread if it exists
                if subscription_thread:
                    print("Stopping current subscription...")
                    stop_event.set()
                    subscription_thread.join()

                # Start a new subscription thread
                print("Subscribing to offers:", subscribe_to)
                subscription_thread = subscribe(stub, subscribe_to)
            elif instr in ('exit', 'quit', 'q'):
                if subscription_thread:
                    stop_event.set()  # Stop the current subscription thread
                    subscription_thread.join()
                break
            else:
                print("Invalid command. Please try again.")
                continue



if __name__ == '__main__':
    run()
