#  psychic-octo-pancake

psychic-octo-pancake is a producer-consumer application created as a homework

# Usage

1. Run `go run . --help` to see available commands
2. Run `docker-compose up` to run all dependencies such as RabbitMQ
3. Run `go run . server` to run a server. Application config will be created automatically in the running filder.
4. Run commands: 
   1. `add "some-key" "some-value"`: add item to memory store
   2. `get "some-key"`: write "some-value" into a data-log file. 
   3. `get-all`: write all values into a data-log file
   4. `remove "some-key"`: delete "some-key" from the memory store

## Configuration
The application is well-configured out of the box. 

By default it creates config file "config.toml" in a working directory and the file is ready to use with default configuration. 

By default, data-file is stored in a working directory as well, but it can be changed in config.toml

- `amqp_queue = 'command_requests'`: client/server queue name
- `amqp_url = 'amqp://guest:guest@localhost:5672/'`: amqp 0.9 (RabbitMQ) url
- `data-file = 'data.log'`: path to a data-log file

