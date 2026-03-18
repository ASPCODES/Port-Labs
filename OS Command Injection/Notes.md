## What is OS command injection?

OS command injection is also known as shell injection. It allows an attacker to execute operating system (OS) commands on the server that is running an application, and typically fully compromise the application and its data. Often, an attacker can leverage an OS command injection vulnerability to compromise other parts of the hosting infrastructure, and exploit trust relationships to pivot the attack to other systems within the organization.


## How it works

Most web applications need to interact with the underlying OS to perform tasks like sending emails, processing files, or checking network status. The vulnerability arises when the application fails to sanitize the input before including it in a command.



The "Joining" Technique
Attackers use command separators to chain their own malicious instructions onto the intended command. Common separators include:

& or && (AND)

| or || (OR/Pipe)

; (Semicolon - Unix only)

newline (0x0a or \n)



## Example Scenario: Network Diagnostic Tool

Imagine a web tool that lets users ping an IP address. The backend code might look like this:
system("ping -c 4 " + user_input);

Normal usage: The user enters 8.8.8.8. The server runs ping -c 4 8.8.8.8.

The Attack: The user enters 8.8.8.8; cat /etc/passwd.

The Result: The server runs the ping, then immediately executes cat /etc/passwd, exposing sensitive system user data to the attacker.



## Useful commands:


Purpose of command	             Linux	              Windows


Name of current user             whoami	              whoami

Operating system                 uname -a             ver

Network configuration            ifconfig             ipconfig /all

Network connections              netstat -an          netstat -an

Running processes                ps -ef               tasklist



## Ways of injecting OS commands:

The following command separators work on both Windows and Unix-based systems:

&

&&

|

||