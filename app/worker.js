console.log("SW Startup!");

// Install Service Worker
self.addEventListener('install', function(event) {
  console.log('installed!');
});

// Service Worker Active
self.addEventListener('activate', function(event) {
  console.log('activated!');
});

self.addEventListener('push', function(event) {
  console.log('data:', event.data.text());
  var msg = event.data.text();
  clients.matchAll().then(clients => {
    clients.forEach(client => {
      send_message_to_client(client, msg).then(m =>
        console.log("SW Received Message: " + m));
    })
  })
});

function send_message_to_client(client, msg) {
  return new Promise(function(resolve, reject) {
    var msg_chan = new MessageChannel();

    msg_chan.port1.onmessage = function(event) {
      if (event.data.error) {
        reject(event.data.error);
      } else {
        resolve(event.data);
      }
    };

    client.postMessage("SW Says: '" + msg + "'", [msg_chan.port2]);
  });
}
