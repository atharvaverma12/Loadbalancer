package main

import(
	"net/http/httputil"
	"net/url"
	"net/http"
	"fmt"
	"os"
)

type Server interface {
	Address() string
	IsAlive() boot
	Serve(rw http.ResponseWriter, r* http.Request)
}

type simpleServer  struct {
	addr string
	proxy *httputil.ReverseProxy
}

func newSimpleServer(add string) *simpleServer{
	serverUrl,err := url.Parse(addr)
	handleErr(err)
	return &simpleServer{
		addr : addr,
		proxy : httputil.NewSingleHostReverseProxy(serverUrl),
	}
}

type LoadBalancer struct {
	port string
	roundRobinCount int 
	servers []Server
}

func NewLoadBalancer(port string, servers[] Server) *LoadBalancer {
	return &LoadBalancer{
		port : port,
		roundRobinCount : 0,
		servers : servers,
	}
}

func handleErr(err error){
	if err!= nil {
		fmt.Printf("error : %v\n", err)
		os.Exit(1)
	}
}

func Address(s* simpleServer) string {return s.addr}

func (s* simpleServer) IsAlive() bool {return true}

func (s* simpleServer) Serve(rw http.ResponseWriter, req * http.Request){
	s.proxy.ServeHTTP(rq,req)
}

func (lb *LoadBalancer) getNextAvailableServer() Server{
	server := lb.servers[lb.roundRobinCount%len(lb.servers)]
	for !server.IsAlive() {
		lb.roundRobinCount++
		server = lb.servers[lb.roundRobinCount%len(lb.servers)]
	}
	lb.roundRobinCount++
	return server
}

func (lb * LoadBalancer) serverProxy(rw http.ResponseWriter, r* http.Request){

	//get the next server
	targetServer := lb.getNextAvailableServer()
	fmt.Printf("forwarding request to address %q\n",targetServer.Address())
	targetServer.Serve(rq,req)
}

func main(){

	//creates a slice of multiple servers
	servers := []Server{
		newSimpleServer("https://www.instagram.com"),
		newSimpleServer("http://www.bing.com"),
		newSimpleServer("http://www.duckduckgo.com"),
	}
	lb := NewLoadBalancer("8000", servers)

	handleRedirect := func(rw http.ResponseWriter, req *http.Request){
		lb.serverProxy(rw,req);
	}
	http.HanldeFunc("/",handleRedirect)

	fmt.Printf("serving requests at 'localhost:%s'\n", lb.port)
	http.ListenAndServe(":"+lb.port,nil)

}