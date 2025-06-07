-----

## Run Z1_Consumer
```bash
javac -cp ".:amqp-client-5.11.0.jar:slf4j-api-1.7.29.jar:slf4j-simple-1.6.2.jar" Z1_Consumer.java && java -cp ".:amqp-client-5.11.0.jar:slf4j-api-1.7.29.jar:slf4j-simple-1.6.2.jar" Z1_Consumer
```

## Run Z1_Producer
```bash
javac -cp ".:amqp-client-5.11.0.jar:slf4j-api-1.7.29.jar:slf4j-simple-1.6.2.jar" Z1_Producer.java && java -cp ".:amqp-client-5.11.0.jar:slf4j-api-1.7.29.jar:slf4j-simple-1.6.2.jar" Z1_Producer
```

## Run Z2_Producer
```bash
javac -cp ".:amqp-client-5.11.0.jar:slf4j-api-1.7.29.jar:slf4j-simple-1.6.2.jar" Z2_Producer.java && java -cp ".:amqp-client-5.11.0.jar:slf4j-api-1.7.29.jar:slf4j-simple-1.6.2.jar" Z2_Producer
```

## Run Z2_Consumer
```bash
javac -cp ".:amqp-client-5.11.0.jar:slf4j-api-1.7.29.jar:slf4j-simple-1.6.2.jar" Z2_Consumer.java && java -cp ".:amqp-client-5.11.0.jar:slf4j-api-1.7.29.jar:slf4j-simple-1.6.2.jar" Z2_Consumer
```