import com.rabbitmq.client.BuiltinExchangeType;
import com.rabbitmq.client.Channel;
import com.rabbitmq.client.Connection;
import com.rabbitmq.client.ConnectionFactory;
import java.io.BufferedReader;
import java.io.InputStreamReader;

public class Z2_Producer {

    public static void main(String[] argv) throws Exception {

        // info
        System.out.println("Z2 PRODUCER");

        // connection & channel
        ConnectionFactory factory = new ConnectionFactory();
        factory.setHost("localhost");
        Connection connection = factory.newConnection();
        Channel channel = connection.createChannel();

        // exchange
        // String EXCHANGE_NAME = "exchange1";
        // channel.exchangeDeclare(EXCHANGE_NAME, BuiltinExchangeType.FANOUT);
        // String EXCHANGE_NAME = "exchange_direct";
        // channel.exchangeDeclare(EXCHANGE_NAME, BuiltinExchangeType.DIRECT);
        String EXCHANGE_NAME = "exchange_topic";
        channel.exchangeDeclare(EXCHANGE_NAME, BuiltinExchangeType.TOPIC);

        while (true) {

            // read msg
            BufferedReader br = new BufferedReader(new InputStreamReader(System.in));
            System.out.println("Enter message (<key>: <body>):");
            String message = br.readLine();
            String[] parts = message.split(":");
            String routingKey = parts[0].strip();
            String body = parts[1].strip();


            // break condition
            if ("exit".equals(message)) {
                break;
            }

            // publish
            // channel.basicPublish(EXCHANGE_NAME, "", null, message.getBytes("UTF-8"));
            channel.basicPublish(EXCHANGE_NAME, routingKey, null, body.getBytes("UTF-8"));
            System.out.println("Sent: " + body);
        }
    }
}
