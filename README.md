# logiQ
Blockchain Based Solution for Logistics

 Implementation of distributed blockchain to store transactions. Every Node connects to all other nodes. 

# Prerequisite
  1. Golang 
  2. Internet if you want to see web interface to download jQuery and bootstraps cdn.
  
# Feature
  No external Golang Library is used.
  
# Usage
  A. Clone Repositoy
  B. Build using command `go build main.go`
  
#### One same Machine
   This is when you have just one machine and need to open multiple Terminals/ Shells
    1. `./main.exe` for first Node. It will tell which ip:port is is listening on. Copy this ip:port address
    2. Open other terminals and start with `./main.exe -peer ip:port` where ip:port is what you copied
    3. Repeat step 2 as many number of time which ip:port of any Node
 
#### On Differnt Machine
   This is when you are running it in different nodes in a connected network on different Machines
   1. `./main.exe -local=false`
   2. On other machine `./main.exe -local=false -peer ip:port` where ip:port is print of console of any other running Nodes.
   
#### Open Web Interface
  Hit `http://ip:port/` of any Node
  Add Transaction Data and Fetch Block Chain
  
  
  
