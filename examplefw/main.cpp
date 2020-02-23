#include "WiFi.h"
#include "coap_client.h"

//instance for coapclient
coapClient coap;

//WiFi connection info
const char* ssid = "";
const char* password = "";

//ip address and default port of coap server in which your interested in
// IPAddress ip(134,102,218,18);// coap.me
IPAddress ip(192,168,183,63);
int port =5683;

// coap client response callback
// void callback_response(coapPacket &packet, IPAddress ip, int port);

// coap client response callback
void callback_response(coapPacket &packet, IPAddress ip, int port) {
Serial.println("callback");
    char p[packet.payloadlen + 1];
    memcpy(p, packet.payload, packet.payloadlen);
    p[packet.payloadlen] = NULL;

    //response from coap server
 if(packet.type==3 && packet.code==0){
      Serial.println("ping ok");
    }

    Serial.print("packet: ");
    Serial.println(p);
}

void setup() {
   
    Serial.begin(115200);

    WiFi.begin(ssid, password);
    Serial.println(" ");

    // Connection info to WiFi network
    Serial.println();
    Serial.println();
    Serial.print("Connecting to ");
    Serial.println(ssid);
    WiFi.begin(ssid, password);
    while (WiFi.status() != WL_CONNECTED) {
    //delay(500);
    yield();
    Serial.print(".");
    }
    Serial.println("");
    Serial.println("WiFi connected");
    // Print the IP address of client
    Serial.println(WiFi.localIP());

    // client response callback.
    // this endpoint is single callback.
    coap.response(callback_response);

    Serial.println("callback registered!");

    // start coap client
    coap.start();
}

void loop() {
    Serial.print("hello");
    bool state;
 
    // Requests

    // get request
    int msgid = coap.get(ip,port,"/");

    char* payload = "+2";
    msgid =coap.post(ip,port,"/",payload,strlen(payload));


    payload = "+2";
    msgid =coap.post(ip,port,"/",payload,strlen(payload));


    payload = "-2";
    msgid =coap.post(ip,port,"/",payload,strlen(payload));


    msgid =coap.delet(ip,port,"/");

    coap.loop();
    delay(1000);
}