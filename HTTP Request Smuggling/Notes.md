## HTTP REQUEST SMUGGLING

HTTP Request is when you send a specifically crafted HTTP that's interpreted differently by the Front-end
and Back-end servers. Front-end HTTP server thinks the requests end here, while at the same time the Back-end
server thinks the request begins here. Boom that is our DeSync, In this Gap an attacker smuggle and extra Data.



## How do HTTP request smuggling vulnerabilities arise?

Most HTTP request smuggling vulnerabilities arise because the HTTP/1 specification provides two different ways to specify where a request ends: 
The Content-Length header and the Transfer-Encoding header.

The Content-Length header is straightforward: 
It specifies the length of the message body in bytes.


1. CL.TE: the front-end server uses the Content-Length header and the back-end server uses the Transfer-Encoding header.

2. TE.CL: the front-end server uses the Transfer-Encoding header and the back-end server uses the Content-Length header.

3. TE.TE: the front-end and back-end servers both support the Transfer-Encoding header, but one of the servers can be induced not to process it by obfuscating the header in some way.



## Requirement of HTTP Request smuggling

 It needs a 2 system architecture, means the Request first go to Load Balancer and Reverse Proxy and the load balancer distribute the millions of request to different servers as to avoid System overflow. The reverse Proxy on the other hand is used to apply filter or checks on the Access control or can say it can used as a WAF or firewall also it is used to cashing the request.

 Proxy means Intermediatry thing between client and the server