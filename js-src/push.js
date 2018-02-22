'use strict';

import * as grpc from "grpc-web-client";
import * as pb_service from "./service_pb_service"
import * as pb from "./service_pb"

var host = window.location.protocol + "//" + window.location.host;

var service = {};
export default service;

service.connect = (app) => {
  service.app = app;

  service.app.ports.getPushState.subscribe(service.onGetPushState);
  service.app.ports.subscribe.subscribe(service.onSubscribe);

  service.app.ports.unsubscribe.subscribe(function() {
    navigator.serviceWorker.ready.then(function(
      serviceWorkerRegistration) {
      serviceWorkerRegistration.pushManager.getSubscription()
        .then(function(subscription) {
          subscription.unsubscribe().then(function() {
            service.send('Unsubscribed');
          });
        })
        .catch(function(e) {
          console.error('Unable to subscribe to push.', e);
        });
    });
  });

  // Handler for messages coming from the service worker
  navigator.serviceWorker.onmessage = function(event) {
    service.app.ports.accountUpdate.send(null);
  };
}

service.send = (msg) => {
  service.app.ports.setPushState.send(msg);
};

service.onGetPushState = () => {
  if (Notification.permission === 'denied') {
    service.send('Denied');
    return;
  }

  navigator.serviceWorker.register('/worker.js');
  navigator.serviceWorker.ready.then(service.onServiceWorkerReady);
};

service.onServiceWorkerReady = (serviceWorkerRegistration) => {
  serviceWorkerRegistration.pushManager.getSubscription()
    .then(function(subscription) {
      if (!subscription) {
        // TODO NeedsPermission?
        service.send('Unsubscribed');
      } else {
        service.send('Subscribed');
      }
    })
    .catch(function(err) {
      service.send('error:' + err);
    });
}

service.onSubscribe = () => {
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
                body.setAuth(new Uint8Array(
                  subscription.getKey(
                    'auth')));
                body.setP256dh(new Uint8Array(
                  subscription.getKey(
                    'p256dh')));
                grpc.grpc.unary(pb_service.PushService.PutSubscription, {
                  request: body,
                  host: host,
                  onEnd: function(res) {
                    service.send('Subscribed');
                  }
                });
              })
          })
          .catch(function(e) {
            console.error('Unable to subscribe to push.',
              e);
          });
      });
    })
    .catch(function(error) {
      console.log('Looks like there was a problem: \n', error);
    });
}
