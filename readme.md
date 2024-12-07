### **Radius**
-----
Modern companies host critical systems, APIs, and tools within secure private networks or isolated remote environments. These systems are often shielded behind firewalls, restricted networks, or lack public access for security reasons. While this enhances protection, it creates a major hurdle for developers or teams who need to interact with these systems seamlessly.  

Currently, developers rely on tools like Remote Desktop to access these environments. However, Remote Desktop tools are bandwidth-heavy and inefficient. They stream everything—screen visuals, keyboard inputs, and mouse movements—resulting in a laggy and frustrating experience, especially over slower connections. Simple tasks like testing an API, accessing a dashboard, or debugging a server can feel unnecessarily complex and slow.  

**Radius** solves this problem by acting as a secure bridge that connects developers to remote environments without exposing any ports or compromising security. Unlike other solutions, the private environment establishes an **outgoing connection** to Radius, so there is no need to open firewall rules or expose backend systems to the public internet. For developers, this means they can connect directly to internal systems—using their browsers or standard development tools—through Radius.  

Radius enables a smooth, fast, and lightweight experience. Instead of streaming a remote screen, it transfers only the necessary data, which significantly reduces bandwidth usage and eliminates lag. Developers can test APIs, access web applications, or interact with backend services in real time, as though they were on the same network.  

At its core, **Radius** combines security and simplicity. Internal systems remain completely hidden, while developers get the direct, seamless access they need to work efficiently. No more high-bandwidth streaming, no complex network setups—just fast, secure connections that enhance productivity and keep systems protected.

### **Architecture**  
-----
The architecture of **Radius** is designed to provide secure and seamless access to remote environments while ensuring that no inbound connections or ports are exposed. It consists of two main components: the **Mid-Server** and the **Rad-Server**. The **Rad-Server** runs inside the private or remote environment where the company’s systems, APIs, or applications reside. Instead of waiting for external connections, the Rad-Server initiates a secure **outgoing connection** to the Mid-Server over WebSocket. This outgoing connection bypasses network restrictions and firewalls without requiring the backend systems to expose any ports to the public internet.  

The **Mid-Server**, which runs on a publicly accessible server or laptop, acts as a bridge between developers and the Rad-Server. It listens for incoming connections from two directions: developers connect via a **SOCKS5 proxy**, and the Rad-Server connects over WebSocket. When a developer sends a request to the Mid-Server, the Mid-Server securely relays the request to the Rad-Server over the pre-established WebSocket connection. The Rad-Server processes the request and sends the response back through the same WebSocket tunnel to the Mid-Server, which then forwards it to the developer.  

This approach enables developers to access private systems as though they were local, using familiar tools like browsers, Postman, or curl. The architecture ensures that backend systems remain hidden and secure because no direct public access is allowed. All communication flows through the Mid-Server, which acts as the controlled entry point while maintaining security by leveraging only outbound connections from the Rad-Server. This combination of security and simplicity provides a fast, lightweight, and efficient alternative to traditional remote desktop tools or high-bandwidth solutions.

### **Setup**  
-----

To set up **Radius**, you need to run two components: the **Mid-Server** and the **Rad-Server**. The Mid-Server acts as a public bridge, while the Rad-Server connects securely to it from a private environment. Both binaries can be downloaded from the project’s GitHub repository.  

**Step 1: Run the Mid-Server**  
Start by running the **Mid-Server** on your laptop or any machine that can be accessed publicly. The Mid-Server listens for incoming **SOCKS5 proxy connections** on port `1080` and accepts **WebSocket connections** from the Rad-Server on port `8080`. Use the following command to start the Mid-Server:  
```bash  
./mid-server  
```  
Once running, the Mid-Server will wait for the Rad-Server to connect.

**Step 2: Run the Rad-Server**  
Next, run the **Rad-Server** in the private or remote environment where your backend systems, APIs, or internal tools are hosted. The Rad-Server will establish a secure **outgoing WebSocket connection** to the Mid-Server. Run the Rad-Server with the following command:  
```bash  
./rad-server -mid ws://<mid-server-address>:8080/connect  
```  
Replace `<mid-server-address>` with the publicly accessible address of the Mid-Server.

**Step 3: Configure Your Browser or Tools**  
To use the Mid-Server, configure your browser or tools (like Postman or curl) to connect via the **SOCKS5 proxy** exposed by the Mid-Server. Set the proxy settings as follows:  
- **SOCKS Host**: `<mid-server-address>`  
- **Port**: `1080`  

For example, you can test access using the following `curl` command:  
```bash  
curl --socks5-hostname <mid-server-address>:1080 http://your-internal-api  
```  

At this point, all requests from your browser or tools will flow through the Mid-Server, which securely relays them to the Rad-Server via the WebSocket connection. The Rad-Server will process the requests and forward responses back through the same connection, enabling seamless access to private systems without exposing any ports.