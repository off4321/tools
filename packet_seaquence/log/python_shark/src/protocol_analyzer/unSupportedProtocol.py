class AnalyzeUnsupportedProtocol:
    def __init__(self, packet, highest_layer):
        self.packet = packet
        self.highest_layer = highest_layer

    def analyze(self):
        return {
            "protocol_name": self.highest_layer
        }
