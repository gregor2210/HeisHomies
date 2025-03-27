Elevator driver for Go
======================

See [`main.go`](main.go) for usage example. The code is runnable with just `go run main.go `

---

Add these lines to your `go.mod` file:
```
require Driver-go v0.0.0
replace Driver-go => ./Driver-go
```
Where `./Driver-go` is the relative path to this folder, after you have downloaded it.




This program is designed to control multiple elevators across multiple computers in a distributed system. Each elevator operates independently but communicates with the others using a peer-to-peer (P2P) architecture over TCP. The system ensures coordinated elevator behavior by exchanging status updates and requests between all connected nodes.  

### Running the Program  
To start an elevator instance, use the following command:  

` go run main.go -id X  `

Where `X` is the unique identifier integer assigned to each elevator. The ID must be an integer between `0` and `NumElevators-1`, where `NumElevators` represents the total number of elevators in the system. Ensure that no two elevators share the same ID, as this may lead to conflicts in communication and system behavior.  

Before running the program, choose what run mode you want to use.

- Setting `UseIPs = false` will make the program usable on only one computer using multiple simulated elevator servers (`simElevatorServers`).

- Setting `UseIPs = true` will make the program attempt to connect to multiple computers. In this case, you must ensure that the `IPs` list is correctly set up. The list should contain the same number of IP addresses as `NumElevators`. The computer assigned ID `0` should have its IP address at index `0` in `IPs`, ID `1` at index `1`, and so on.

- Find the ip using `hostname -I`

You can change all the settings above in `configConnectivity.go`.

### Other Configs 
In`config_connectivity.go`: 
1) You can change the `TimeOut` value if you want the program to be more or less aggressive in handling connection TimeOuts and packet loss scenarios.

In `configFsm.go`:
1) You can change the `_timerPollRate`, `_motorErrorDuration`, `_obstrErrorDuration`, `NumFloors`

In `main.go`:
1) you can change the local port `PortServerID0` that acceses your elevatorserver. If you are running with simulators the port of the different simulators whould be `PortServerID0 + ID`



### The Code
The code is written using three packages along with one `main.go` file.
- The package `fsm` contains the general logic for managing a single elevator.
- The package `elevio` handles communication with the elevator server.
- The package `connectivity` manages communication and logic required for multiple elevators to work together.






