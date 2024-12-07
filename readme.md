### **Radius**
---
Modern companies host critical systems, APIs, and tools within secure private networks or isolated remote environments. These systems are often shielded behind firewalls, restricted networks, or lack public access for security reasons. While this enhances protection, it creates a major hurdle for developers or teams who need to interact with these systems seamlessly.  

Currently, developers rely on tools like Remote Desktop to access these environments. However, Remote Desktop tools are bandwidth-heavy and inefficient. They stream everything—screen visuals, keyboard inputs, and mouse movements—resulting in a laggy and frustrating experience, especially over slower connections. Simple tasks like testing an API, accessing a dashboard, or debugging a server can feel unnecessarily complex and slow.  

**Radius** solves this problem by acting as a secure bridge that connects developers to remote environments without exposing any ports or compromising security. Unlike other solutions, the private environment establishes an **outgoing connection** to Radius, so there is no need to open firewall rules or expose backend systems to the public internet. For developers, this means they can connect directly to internal systems—using their browsers or standard development tools—through Radius.  

Radius enables a smooth, fast, and lightweight experience. Instead of streaming a remote screen, it transfers only the necessary data, which significantly reduces bandwidth usage and eliminates lag. Developers can test APIs, access web applications, or interact with backend services in real time, as though they were on the same network.  

At its core, **Radius** combines security and simplicity. Internal systems remain completely hidden, while developers get the direct, seamless access they need to work efficiently. No more high-bandwidth streaming, no complex network setups—just fast, secure connections that enhance productivity and keep systems protected.

### **Architecture**  
---
The architecture of **Radius** is designed to provide secure and seamless access to remote environments while ensuring that no inbound connections or ports are exposed. It consists of two main components: the **Mid-Server** and the **Rad-Server**. The **Rad-Server** runs inside the private or remote environment where the company’s systems, APIs, or applications reside. Instead of waiting for external connections, the Rad-Server initiates a secure **outgoing connection** to the Mid-Server over WebSocket. This outgoing connection bypasses network restrictions and firewalls without requiring the backend systems to expose any ports to the public internet.  

The **Mid-Server**, which runs on a publicly accessible server or laptop, acts as a bridge between developers and the Rad-Server. It listens for incoming connections from two directions: developers connect via a **SOCKS5 proxy**, and the Rad-Server connects over WebSocket. When a developer sends a request to the Mid-Server, the Mid-Server securely relays the request to the Rad-Server over the pre-established WebSocket connection. The Rad-Server processes the request and sends the response back through the same WebSocket tunnel to the Mid-Server, which then forwards it to the developer.  

This approach enables developers to access private systems as though they were local, using familiar tools like browsers, Postman, or curl. The architecture ensures that backend systems remain hidden and secure because no direct public access is allowed. All communication flows through the Mid-Server, which acts as the controlled entry point while maintaining security by leveraging only outbound connections from the Rad-Server. This combination of security and simplicity provides a fast, lightweight, and efficient alternative to traditional remote desktop tools or high-bandwidth solutions.

### **Setup**  
---
To set up **Radius**, you need to deploy two components: the **Rad-Server** inside the private or remote environment and the **Mid-Server** on a publicly accessible machine. First, download the precompiled binaries for both components from the project’s GitHub repository.  

Start by running the **Rad-Server** in the private environment where your backend systems, APIs, or services are located. The Rad-Server will establish a secure outgoing connection to the Mid-Server. Run the Rad-Server binary with the following command:  
```bash  
./rad-server -mid ws://<mid-server-address>:8080/connect  
```  
Here, replace `<mid-server-address>` with the public address of the Mid-Server.

Next, run the **Mid-Server** binary on your laptop or any server that can be accessed publicly. The Mid-Server listens for incoming SOCKS5 connections on port `1080` and accepts WebSocket connections from the Rad-Server on port `8080`. Use the following command:  
```bash  
./mid-server  
```  

Once the Mid-Server is running, configure your browser or tools to use the Mid-Server as a **SOCKS5 proxy**. Set the proxy host to the address of the Mid-Server and port `1080`. For example, in your browser’s network settings, configure the proxy to:  
- **SOCKS Host**: `<mid-server-address>`  
- **Port**: `1080`  

At this point, all your browser or tool requests will flow through the Mid-Server. The Mid-Server will securely relay those requests to the Rad-Server over the WebSocket connection, allowing you to access APIs, test servers, or internal applications in the private environment directly and seamlessly.