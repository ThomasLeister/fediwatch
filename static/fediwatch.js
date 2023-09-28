// Protobuf objects
var Connection;

// Globals
var WebSettings = {};
var globe;

// Constants 
const ARC_REL_LEN = 0.4; // relative to whole arc
const FLIGHT_TIME = 1000;
const NUM_RINGS = 1;
const RINGS_MAX_R = 5; // deg
const RING_PROPAGATION_SPEED = 5; // deg/sec
const OPACITY = 0.8;


// Hide loading screen
function hideLoading() {
    console.log("Hiding loading screen");
    document.getElementById("loading").style.display = "none";
}

// Synchronous function loads settings via the /settings API endpoint
function loadSettings() {
    const request = new XMLHttpRequest();
    request.open("GET", "/settings", false); // `false` makes the request synchronous
    request.send(null);

    if (request.status === 200) {
        WebSettings = JSON.parse(request.responseText);
        console.log("Loaded settings");
    } else {
        console.err("Failed to load settings via /settings API!");
    }
}

function setupWebSocket() {
    let socket;

    const connect = () => {
        socket = (location.hostname == "localhost") ? new WebSocket("ws://localhost:8011/ws") : new WebSocket(WebSettings.websocketUrl);
        socket.binaryType = 'arraybuffer';

        // Event-Handler für die Verbindungsherstellung
        socket.addEventListener("open", (event) => {
            console.log("Connected to Websocket");
        });

        // Event-Handler für den Empfang von WebSocket-Nachrichten
        socket.addEventListener("message", (event) => {
            const array = new Uint8Array(event.data)

            // Decode bytesArray to Protobuf Connection Message
            const connMsg = Connection.decode(array);

            // Decode Connection message to Connection object
            const conn = Connection.toObject(connMsg, {
                enums: String,
                floats: Number,
                defaults: true		// Include attribute in object even if it's set to the default value
            });

            if (conn.dir == "IN") {
                emitArcIn({ lat: conn.lat, lng: conn.lng});
            } else if (conn.dir == "OUT"){
                emitArcOut({ lat: conn.lat, lng: conn.lng});
            }
        });

        // Event-Handler für das Schließen der WebSocket-Verbindung
        socket.addEventListener("close", (event) => {
            if (event.wasClean) {
                console.log("WebSocket connection was closed");
            } else {
                console.error("WebSocket connection closed unexpectedly");
                setTimeout(connect, 1000);
            }
        });

        // Event-Handler für Fehler in der WebSocket-Verbindung
        socket.addEventListener("error", (event) => {
        console.error("WebSocket error:", event);
        });
    };

    connect();
}

    
function emitArcIn({ lat: startLat, lng: startLng }) {
    console.log("emitArc:", WebSettings.homeLocation.lat);

    const { lat: endLat, lng: endLng } = { lat: WebSettings.homeLocation.lat, lng: WebSettings.homeLocation.long };

    // add and remove arc after 1 cycle
    const id = Math.random(9999);
    const arc = { startLat, startLng, endLat, endLng, id };
    globe.arcsData([...globe.arcsData(), arc]);
    setTimeout(() => globe.arcsData(globe.arcsData().filter(d => d.id !== arc.id)), FLIGHT_TIME * 2);
}

function emitArcOut({ lat: endLat, lng: endLng }) {
    const { lat: startLat, lng: startLng } = { lat: WebSettings.homeLocation.lat, lng: WebSettings.homeLocation.long };

    // add and remove arc after 1 cycle
    const id = Math.random(9999);
    const arc = { startLat, startLng, endLat, endLng, id };
    globe.arcsData([...globe.arcsData(), arc]);
    setTimeout(() => globe.arcsData(globe.arcsData().filter(d => d.id !== arc.id)), FLIGHT_TIME * 2);
}

function initGlobe() {
    globe = Globe({animateIn: true})
        .globeImageUrl('/earth.webp')
        .arcColor(d => [`rgba(0, 255, 255, ${OPACITY})`, `rgba(0, 255, 0, ${OPACITY})`])
        .arcDashLength(ARC_REL_LEN)
        .arcDashGap(2)
        .arcDashInitialGap(1)
        .arcDashAnimateTime(FLIGHT_TIME)
        .arcsTransitionDuration(0)
        .arcAltitudeAutoScale(0.3)
        .arcStroke(0.2)
        .ringColor(() => t => `rgba(52, 195, 235,${1-t})`)
        .ringMaxRadius(RINGS_MAX_R)
        .ringPropagationSpeed(RING_PROPAGATION_SPEED)
        .ringRepeatPeriod(FLIGHT_TIME * ARC_REL_LEN / NUM_RINGS)
        .pointOfView({ lat: WebSettings.homeLocation.lat, lng: WebSettings.homeLocation.long, altitude: 1.7}, 2500)
        .onGlobeReady(hideLoading)
        (document.getElementById('globeViz'));
}




// Load settings
loadSettings();

// Init globe
initGlobe();

// Load protobuf and set up websocket connection for connection updates
protobuf.load("/fediwatch.proto", function (err, root) {
    if (err) throw err;
    console.log("Loaded fediwatch.proto");
    Connection = root.lookupType("fediwatch.Connection");

    setupWebSocket();
});