console.log("SW: startup");

// Install Service Worker
self.addEventListener('install', function(event) {
  console.log('SW: install');
});

// Service Worker Active
self.addEventListener('activate', function(event) {
  console.log('SW: activate');
  clients.claim();
  sendToAll("you have been claimed!");
});

self.addEventListener('push', function(event) {
  console.log('data:', event.data.text());
  var msg = event.data.text();
  sendToAll(msg);
});

function sendToAll(msg) {
  clients.matchAll().then(clients => {
    clients.forEach(client => {
      sendTo(client, msg).then(m =>
        console.log("SW: channel response: " + m));
    })
  })
}

function sendTo(client, msg) {
  return new Promise(function(resolve, reject) {
    var chan = new MessageChannel();

    chan.port1.onmessage = function(event) {
      if (event.data.error) {
        reject(event.data.error);
      } else {
        resolve(event.data);
      }
    };

    client.postMessage(msg, [chan.port2]);
  });
}
