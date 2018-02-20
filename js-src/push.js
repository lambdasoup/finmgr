var service = require("./service_pb_service");
var pb = require("./service_pb");
var grpc = require("grpc-web-client");

var host = window.location.protocol + "//" + window.location.host;

export function connect(app) {
  var send = app.ports.setPushState.send;

  app.ports.getPushState.subscribe(function() {
    if (!available()) {
      send('NotAvailable');
      return;
    }

    if (Notification.permission === 'denied') {
      send('Denied');
      return;
    }

    navigator.serviceWorker.register('/worker.js');
    navigator.serviceWorker.ready.then(function(serviceWorkerRegistration) {
      serviceWorkerRegistration.pushManager.getSubscription()
        .then(function(subscription) {
          if (!subscription) {
            // TODO NeedsPermission?
            send('Unsubscribed');
          } else {
            send('Subscribed');
          }
        })
        .catch(function(err) {
          send('error:' + err);
        });
    });
  });

  app.ports.subscribe.subscribe(function() {
    fetch('web-push/publicKey').then(function(response) {
        response.arrayBuffer().then(function(buffer) {
          var publicKey = new Uint8Array(buffer);
          navigator.serviceWorker.ready.then(function(
              serviceWorkerRegistration) {
              serviceWorkerRegistration.pushManager.subscribe({
                  userVisibleOnly: true,
                  applicationServerKey: publicKey,
                })
                .then(function(subscription) {
                  var body = new pb.Subscription();
                  body.setEndpoint(subscription.endpoint);
                  body.setAuth(new Uint8Array(subscription.getKey(
                    'auth')));
                  body.setP256dh(new Uint8Array(subscription.getKey(
                    'p256dh')));
                  grpc.grpc.unary(service.PushService.PutSubscription, {
                    request: body,
                    host: host,
                    onEnd: function(res) {
                      send('Subscribed');
                    }
                  });
                })
            })
            .catch(function(e) {
              console.error('Unable to subscribe to push.', e);
            });
        });
      })
      .catch(function(error) {
        console.log('Looks like there was a problem: \n', error);
      });
  });

  app.ports.unsubscribe.subscribe(function() {
    navigator.serviceWorker.ready.then(function(serviceWorkerRegistration) {
      serviceWorkerRegistration.pushManager.getSubscription()
        .then(function(subscription) {
          subscription.unsubscribe().then(function() {
            send('Unsubscribed');
          });
        })
        .catch(function(e) {
          console.error('Unable to subscribe to push.', e);
        });
    });
  });

  // Handler for messages coming from the service worker
  navigator.serviceWorker.onmessage = function(event) {
    app.ports.accountUpdate.send(null);
  };
}

function available() {
  if (!('serviceWorker' in navigator)) {
    return false;
  }
  if (!('PushManager' in window)) {
    return false;
  }
  if (!('showNotification' in ServiceWorkerRegistration.prototype)) {
    return false;
  }
  return true;
}
