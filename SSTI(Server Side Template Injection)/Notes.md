## Server Side Template Injection (SSTI)


## What is SSTI?

When web applications use template engines (like Jinja2, Twig, Freemarker), and user input is directly injected into the template without sanitization — that's SSTI.
The template engine's job is to render dynamic content, but if an attacker injects their code using template syntax, the server executes that code.



## Simple Example — Let's Understand

Normal flow:

Template: "Hello {{ name }}"
Input: "Ali"
Output: "Hello Ali"


SSTI attack:

https://site.com/profile?name={{7*7}}
https://site.com/search?q={{config}}




## Form Fields / Input Boxes

1. Name, email, message, search bar

2. Feedback forms, contact forms

3. Registration forms




## HTTP Headers

User-Agent, Referer, X-Forwarded-For

When the site displays these headers on any page.




## SSTI Detection Methodology

Step 1: Inject into every input field → {{7*7}}, ${7*7}, #{7*7}
Step 2: Check the response — if math is calculated → it's SSTI
Step 3: Identify the template engine (from error messages or behavior)
Step 4: Try engine-specific payloads
Step 5: Escalate to RCE (if within scope)




## Detection Cheatsheet:

{{7*7}} → Jinja2/Twig
${7*7} → Freemarker/Groovy
#{7*7} → Ruby ERB
*{7*7} → Spring (Java)
<%= 7\*7 %> → ERB (Ruby)
